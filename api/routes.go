package api

import (
	"smartbudget/db"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.Engine, db *db.DB) {
	userHandler := NewUserHandler(db)
	userHandler.RegisterRoutes(router)
}

// RegisterInsightRoutes registers all insight-related routes
func RegisterInsightRoutes(router *gin.Engine, db *db.DB, aiClient *AIClient) {
	insightHandler := NewInsightHandler(db, aiClient)
	insightHandler.RegisterRoutes(router)
}
