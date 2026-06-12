package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "Ceremony/internal/dtos"
    "Ceremony/internal/models"
)

func APIKeyMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.GetHeader("X-API-Key")
        if key == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{Error: "API key required"})
            return
        }

        var app models.App
        if err := db.Where("api_key = ? AND is_active = ?", key, true).First(&app).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{Error: "Invalid or inactive API key"})
            return
        }

        c.Set("appID", app.ID)
        // c.Set("userID", app.UserID)
        c.Next()
    }
}