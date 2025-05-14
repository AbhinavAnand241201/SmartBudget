package api

import (
	"net/http"

	"smartbudget/db"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	// Add database connection or service here later
}

// NewUserHandler creates a new UserHandler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/api/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Add database operation
	c.JSON(http.StatusCreated, user)
}

// GetUser handles getting a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	// TODO: Add database operation
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// UpdateUser handles updating a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Add database operation
	c.JSON(http.StatusOK, gin.H{"id": id, "user": user})
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	// TODO: Add database operation
	c.JSON(http.StatusOK, gin.H{"id": id})
}
