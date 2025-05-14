package api

import (
	"net/http"
	"smartbudget/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	db *db.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *db.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

// RegisterRoutes registers all user-related routes
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/users", h.CreateUser)
	router.GET("/users/:id", h.GetUser)
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate UUID for new user
	user.ID = uuid.New()

	// Create user in database
	if err := h.db.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles user retrieval
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.db.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
