package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
	"gorm.io/gorm"

    "Ceremony/internal/dtos"
    "Ceremony/internal/models"
)

func AdminMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        uidInterface, exists := c.Get("userID")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{
                Error: "Authentication required",
            })
            return
        }

        uid := uidInterface.(uint)

        var user models.User
        if err := db.First(&user, uid).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{
                Error: "User not found",
            })
            return
        }

        if !user.IsAdmin {
            c.AbortWithStatusJSON(http.StatusForbidden, dtos.ForbiddenResponse{
                Error: "Admin access required",
            })
            return
        }

        c.Next()
    }
}