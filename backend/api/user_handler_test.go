package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"smartbudget/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	handler := NewUserHandler()
	handler.RegisterRoutes(r)

	t.Run("CreateUser", func(t *testing.T) {
		// Create test user
		user := db.NewUser("test@example.com")
		body, _ := json.Marshal(user)

		// Create request
		req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Perform request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code)
		var response db.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, response.Email)
	})

	t.Run("GetUser", func(t *testing.T) {
		// Create request
		req, _ := http.NewRequest("GET", "/api/users/123", nil)
		w := httptest.NewRecorder()

		// Perform request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "123", response["id"])
	})
}
