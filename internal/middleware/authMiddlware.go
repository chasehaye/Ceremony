package middleware

import (
	"Ceremony/internal/crypt"
    "Ceremony/internal/dtos"
    "net/http"
	"github.com/gin-gonic/gin"
    "Ceremony/internal/models"
    "gorm.io/gorm"
)


func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString, err := c.Cookie("token")
        
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{
				Error: "Authentication required",
			})
            return
        }
        token, err := crypt.ValidateJWT(tokenString)
        if err != nil || token == nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{
				Error: "Invalid or expired session",
			})
            return
        }
        uid, err := crypt.GetUserIDFromJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{
				Error: "Failed to identify user from session",
			})
			return
		}

        var user models.User
		if err := db.First(&user, uid).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{
				Error: "User not found",
			})
			return
		}

        if user.IsBanned {
			c.AbortWithStatusJSON(http.StatusForbidden, dtos.ForbiddenResponse{
				Error: "Your account has been banned",
			})
			return
		}

        c.Set("userID", uint(uid))
        c.Next()
    }
}