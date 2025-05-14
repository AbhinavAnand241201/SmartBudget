package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"smartbudget/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestRouterAndUser() (*gin.Engine, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	user := &db.User{
		ID:    userID,
		Email: "test@example.com",
	}
	mockDB := db.NewMockDB()
	_ = mockDB.CreateUser(nil, user)
	// Patch handler to use this mockDB
	handler := &InsightHandler{db: mockDB}
	r := gin.Default()
	handler.RegisterRoutes(r)
	return r, userID
}

func TestInsightHandler(t *testing.T) {
	t.Run("CreateInsight", func(t *testing.T) {
		r, userID := setupTestRouterAndUser()
		insight := &db.Insight{
			ID:      uuid.New(),
			UserID:  userID,
			Type:    db.Alert,
			Content: "Test alert",
		}
		body, _ := json.Marshal(insight)
		req, _ := http.NewRequest("POST", "/api/insights", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var response db.Insight
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, insight.UserID, response.UserID)
		assert.Equal(t, insight.Type, response.Type)
		assert.Equal(t, insight.Content, response.Content)
	})

	t.Run("CreateInsightInvalidType", func(t *testing.T) {
		r, userID := setupTestRouterAndUser()
		insight := &db.Insight{
			UserID:  userID,
			Type:    "invalid",
			Content: "Test alert",
		}
		body, _ := json.Marshal(insight)
		req, _ := http.NewRequest("POST", "/api/insights", bytes.NewBuffer(body))
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
		r, userID := setupTestRouterAndUser()
		insight := &db.Insight{
			ID:      uuid.New(),
			UserID:  userID,
			Type:    db.Alert,
			Content: "Test alert",
		}
		body, _ := json.Marshal(insight)
		createReq, _ := http.NewRequest("POST", "/api/insights", bytes.NewBuffer(body))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		r.ServeHTTP(createW, createReq)
		assert.Equal(t, http.StatusCreated, createW.Code)
		var created db.Insight
		err := json.Unmarshal(createW.Body.Bytes(), &created)
		assert.NoError(t, err)
		req, _ := http.NewRequest("GET", "/api/insights/user/"+userID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response []db.Insight
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, created.ID, response[0].ID)
		assert.Equal(t, created.UserID, response[0].UserID)
		assert.Equal(t, created.Type, response[0].Type)
		assert.Equal(t, created.Content, response[0].Content)
	})

	t.Run("GetUserInsightsInvalidID", func(t *testing.T) {
		r, _ := setupTestRouterAndUser()
		req, _ := http.NewRequest("GET", "/api/insights/user/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid user ID", response["error"])
	})
}
