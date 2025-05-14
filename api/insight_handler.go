package api

import (
	"net/http"

	"smartbudget/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InsightHandler handles insight-related HTTP requests
type InsightHandler struct {
	db *db.MockDB
}

// NewInsightHandler creates a new insight handler
func NewInsightHandler() *InsightHandler {
	return &InsightHandler{
		db: db.NewMockDB(),
	}
}

// RegisterRoutes registers the insight routes
func (h *InsightHandler) RegisterRoutes(r *gin.Engine) {
	insights := r.Group("/api/insights")
	{
		insights.POST("", h.CreateInsight)
		insights.GET("/:id", h.GetInsight)
		insights.GET("/user/:user_id", h.GetUserInsights)
		insights.PUT("/:id", h.UpdateInsight)
		insights.DELETE("/:id", h.DeleteInsight)
	}
}

// CreateInsight creates a new insight
func (h *InsightHandler) CreateInsight(c *gin.Context) {
	var insight db.Insight
	if err := c.ShouldBindJSON(&insight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidInsightType(insight.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight type"})
		return
	}

	createdInsight, err := h.db.CreateInsight(nil, &insight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdInsight)
}

// GetInsight gets an insight by ID
func (h *InsightHandler) GetInsight(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight ID"})
		return
	}

	insight, err := h.db.GetInsight(nil, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insight not found"})
		return
	}

	c.JSON(http.StatusOK, insight)
}

// GetUserInsights gets all insights for a user
func (h *InsightHandler) GetUserInsights(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	insights, err := h.db.GetUserInsights(nil, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insights)
}

// UpdateInsight updates an insight
func (h *InsightHandler) UpdateInsight(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight ID"})
		return
	}

	var insight db.Insight
	if err := c.ShouldBindJSON(&insight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidInsightType(insight.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight type"})
		return
	}

	insight.ID = id
	updatedInsight, err := h.db.UpdateInsight(nil, &insight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedInsight)
}

// DeleteInsight deletes an insight
func (h *InsightHandler) DeleteInsight(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight ID"})
		return
	}

	err = h.db.DeleteInsight(nil, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// isValidInsightType checks if the insight type is valid
func isValidInsightType(insightType db.InsightType) bool {
	switch insightType {
	case db.Alert, db.Tip, db.Prediction:
		return true
	default:
		return false
	}
}
