package api

import (
	"net/http"

	"smartbudget/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InsightHandler handles insight-related HTTP requests
type InsightHandler struct {
	db       *db.DB
	aiClient *AIClient
}

// NewInsightHandler creates a new insight handler
func NewInsightHandler(db *db.DB, aiClient *AIClient) *InsightHandler {
	return &InsightHandler{
		db:       db,
		aiClient: aiClient,
	}
}

// RegisterRoutes registers all insight-related routes
func (h *InsightHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/insights", h.CreateInsight)
	router.GET("/insights/:id", h.GetInsight)
	router.GET("/users/:user_id/insights", h.GetUserInsights)
	router.PUT("/insights/:id", h.UpdateInsight)
	router.DELETE("/insights/:id", h.DeleteInsight)
}

// CreateInsight handles insight creation
func (h *InsightHandler) CreateInsight(c *gin.Context) {
	var insight db.Insight
	if err := c.ShouldBindJSON(&insight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate insight type
	if !isValidInsightType(insight.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight type"})
		return
	}

	// Generate UUID for new insight
	insight.ID = uuid.New()

	// Create insight in database
	if err := h.db.CreateInsight(c.Request.Context(), &insight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create insight"})
		return
	}

	c.JSON(http.StatusCreated, insight)
}

// GetInsight handles insight retrieval
func (h *InsightHandler) GetInsight(c *gin.Context) {
	id := c.Param("id")
	insightID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight ID"})
		return
	}

	insight, err := h.db.GetInsight(c.Request.Context(), insightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get insight"})
		return
	}

	if insight == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insight not found"})
		return
	}

	c.JSON(http.StatusOK, insight)
}

// GetUserInsights handles retrieving all insights for a user
func (h *InsightHandler) GetUserInsights(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	insights, err := h.db.GetUserInsights(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user insights"})
		return
	}

	c.JSON(http.StatusOK, insights)
}

// UpdateInsight handles insight updates
func (h *InsightHandler) UpdateInsight(c *gin.Context) {
	id := c.Param("id")
	insightID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight ID"})
		return
	}

	var insight db.Insight
	if err := c.ShouldBindJSON(&insight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate insight type
	if !isValidInsightType(insight.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight type"})
		return
	}

	insight.ID = insightID

	// Update insight in database
	if err := h.db.UpdateInsight(c.Request.Context(), &insight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update insight"})
		return
	}

	c.JSON(http.StatusOK, insight)
}

// DeleteInsight handles insight deletion
func (h *InsightHandler) DeleteInsight(c *gin.Context) {
	id := c.Param("id")
	insightID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight ID"})
		return
	}

	if err := h.db.DeleteInsight(c.Request.Context(), insightID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete insight"})
		return
	}

	c.Status(http.StatusNoContent)
}

// isValidInsightType checks if the insight type is valid
func isValidInsightType(insightType string) bool {
	validTypes := map[string]bool{
		"spending_pattern": true,
		"budget_alert":     true,
		"savings_goal":     true,
		"category_insight": true,
	}
	return validTypes[insightType]
}
