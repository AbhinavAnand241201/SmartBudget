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

func setupTestUserRouter(t *testing.T) (*gin.Engine, *db.DB, *db.User) {
	gin.SetMode(gin.TestMode)
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	database, err := db.NewDB(dbURL)
	if err != nil {
		t.Skipf("Could not connect to test DB: %v", err)
	}
	user := &db.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}
	_ = database.CreateUser(nil, user)
	handler := NewUserHandler(database)
	r := gin.Default()
	handler.RegisterRoutes(r)
	return r, database, user
}

func TestUserHandler(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		r, database, _ := setupTestUserRouter(t)
		defer database.Close()
		user := &db.User{
			ID:    uuid.New(),
			Email: "newuser@example.com",
			Name:  "New User",
		}
		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var response db.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, response.Email)
		assert.Equal(t, user.Name, response.Name)
	})

	t.Run("GetUser", func(t *testing.T) {
		r, database, user := setupTestUserRouter(t)
		defer database.Close()
		req, _ := http.NewRequest("GET", "/users/"+user.ID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response db.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Email, response.Email)
		assert.Equal(t, user.Name, response.Name)
	})
}
