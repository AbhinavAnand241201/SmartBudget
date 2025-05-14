package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"smartbudget/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestRouterAndUser(t *testing.T) (*gin.Engine, uuid.UUID, *db.DB) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	user := &db.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	database, err := db.NewDB(dbURL)
	if err != nil {
		t.Skipf("Could not connect to test DB: %v", err)
	}
	_ = database.CreateUser(nil, user)
	handler := NewInsightHandler(database, nil)
	r := gin.Default()
	handler.RegisterRoutes(r)
	return r, userID, database
}

func TestInsightHandler(t *testing.T) {
	t.Run("CreateInsight", func(t *testing.T) {
		r, userID, database := setupTestRouterAndUser(t)
		defer database.Close()
		insight := &db.Insight{
			ID:          uuid.New(),
			UserID:      userID,
			Type:        "spending_pattern",
			Title:       "Test Title",
			Description: "Test Description",
			Data:        "{}",
		}
		body, _ := json.Marshal(insight)
		req, _ := http.NewRequest("POST", "/insights", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var response db.Insight
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, insight.UserID, response.UserID)
		assert.Equal(t, insight.Type, response.Type)
		assert.Equal(t, insight.Title, response.Title)
		assert.Equal(t, insight.Description, response.Description)
		assert.Equal(t, insight.Data, response.Data)
	})

	t.Run("CreateInsightInvalidType", func(t *testing.T) {
		r, userID, database := setupTestRouterAndUser(t)
		defer database.Close()
		insight := &db.Insight{
			UserID:      userID,
			Type:        "invalid",
			Title:       "Test Title",
			Description: "Test Description",
			Data:        "{}",
		}
		body, _ := json.Marshal(insight)
		req, _ := http.NewRequest("POST", "/insights", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid insight type", response["error"])
	})

	t.Run("GetUserInsights", func(t *testing.T) {
		r, userID, database := setupTestRouterAndUser(t)
		defer database.Close()
		insight := &db.Insight{
			ID:          uuid.New(),
			UserID:      userID,
			Type:        "spending_pattern",
			Title:       "Test Title",
			Description: "Test Description",
			Data:        "{}",
		}
		body, _ := json.Marshal(insight)
		createReq, _ := http.NewRequest("POST", "/insights", bytes.NewBuffer(body))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		r.ServeHTTP(createW, createReq)
		assert.Equal(t, http.StatusCreated, createW.Code)
		var created db.Insight
		err := json.Unmarshal(createW.Body.Bytes(), &created)
		assert.NoError(t, err)
		req, _ := http.NewRequest("GET", "/users/"+userID.String()+"/insights", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response []db.Insight
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response)
		found := false
		for _, resp := range response {
			if resp.ID == created.ID {
				found = true
				assert.Equal(t, created.UserID, resp.UserID)
				assert.Equal(t, created.Type, resp.Type)
				assert.Equal(t, created.Title, resp.Title)
				assert.Equal(t, created.Description, resp.Description)
				assert.Equal(t, created.Data, resp.Data)
			}
		}
		assert.True(t, found)
	})

	t.Run("GetUserInsightsInvalidID", func(t *testing.T) {
		r, _, database := setupTestRouterAndUser(t)
		defer database.Close()
		req, _ := http.NewRequest("GET", "/users/invalid/insights", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid user ID", response["error"])
	})
}
