package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"smartbudget/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/plaid/plaid-go/v10/plaid"
)

type PlaidClient interface {
	TransactionsGet(ctx interface{}) PlaidTransactionsGetRequest
}

type PlaidTransactionsGetRequest interface {
	TransactionsGetRequest(req interface{}) PlaidTransactionsGetRequest
	Execute() (struct {
		Transactions []struct {
			Amount   float64
			Category []string
			Name     string
			Date     string
		}
	}, *http.Response, error)
}

// FinancialDataHandler handles financial data operations
type FinancialDataHandler struct {
	userDB         db.UserDB
	transactionDB  db.TransactionDB
	costOfLivingDB db.CostOfLivingDB
	plaidClient    PlaidClient
	httpGet        func(url string) (*http.Response, error)
	numbeoAPIKey   string
}

// NewFinancialDataHandler creates a new FinancialDataHandler
func NewFinancialDataHandler(
	userDB db.UserDB,
	transactionDB db.TransactionDB,
	costOfLivingDB db.CostOfLivingDB,
	plaidClient PlaidClient,
	httpGet func(url string) (*http.Response, error),
	numbeoAPIKey string,
) *FinancialDataHandler {
	return &FinancialDataHandler{
		userDB:         userDB,
		transactionDB:  transactionDB,
		costOfLivingDB: costOfLivingDB,
		plaidClient:    plaidClient,
		httpGet:        httpGet,
		numbeoAPIKey:   numbeoAPIKey,
	}
}

// RegisterRoutes registers the financial data routes
func (h *FinancialDataHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/sync-transactions", h.SyncTransactions)
	router.POST("/fetch-cost-of-living", h.FetchCostOfLiving)
}

// SyncTransactionsRequest represents a request to sync transactions
type SyncTransactionsRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
}

// SyncTransactionsResponse represents a response from syncing transactions
type SyncTransactionsResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// SyncTransactions syncs transactions from Plaid
func (h *FinancialDataHandler) SyncTransactions(c *gin.Context) {
	var req SyncTransactionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user exists
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userDB.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get transactions from Plaid
	startDate := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	transactionsGetRequest := h.plaidClient.TransactionsGet(nil).TransactionsGetRequest(plaid.NewTransactionsGetRequest(
		req.AccessToken,
		startDate,
		endDate,
	))

	transactionsGetResponse, _, err := transactionsGetRequest.Execute()
	if err != nil {
		if err.Error() == "Invalid token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions from Plaid"})
		return
	}

	// Store transactions
	count := 0
	for _, t := range transactionsGetResponse.Transactions {
		amount := float64(t.Amount)
		if t.Amount < 0 {
			amount = -amount
		}

		transaction := &db.Transaction{
			UserID:      userID,
			Amount:      amount,
			Category:    t.Category[0],
			Description: t.Name,
			Date:        t.Date,
		}

		err := h.transactionDB.CreateTransaction(c.Request.Context(), transaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store transaction"})
			return
		}
		count++
	}

	c.JSON(http.StatusOK, SyncTransactionsResponse{
		Status: "success",
		Count:  count,
	})
}

// FetchCostOfLivingRequest represents a request to fetch cost of living data
type FetchCostOfLivingRequest struct {
	City string `json:"city" binding:"required"`
}

// FetchCostOfLivingResponse represents a response from fetching cost of living data
type FetchCostOfLivingResponse struct {
	City   string          `json:"city"`
	Prices json.RawMessage `json:"prices"`
}

// FetchCostOfLiving fetches cost of living data from Numbeo
func (h *FinancialDataHandler) FetchCostOfLiving(c *gin.Context) {
	var req FetchCostOfLivingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch data from Numbeo
	url := fmt.Sprintf("https://www.numbeo.com/api/city_prices?api_key=%s&city=%s", h.numbeoAPIKey, req.City)
	resp, err := h.httpGet(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cost of living data"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cost of living data"})
		return
	}

	var prices json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cost of living data"})
		return
	}

	// Store in database
	err = h.costOfLivingDB.UpsertCostOfLiving(c.Request.Context(), req.City, prices)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store cost of living data"})
		return
	}

	c.JSON(http.StatusOK, FetchCostOfLivingResponse{
		City:   req.City,
		Prices: prices,
	})
}
