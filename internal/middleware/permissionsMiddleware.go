package middleware

import (
    "Ceremony/internal/dtos"
    "net/http"
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "Ceremony/internal/models"
)

func CanCreateMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        canCreate := c.MustGet("canCreate").(bool)
        if !canCreate {
            c.AbortWithStatusJSON(http.StatusForbidden, dtos.ForbiddenResponse{
                Error: "You do not have permission to create resources",
            })
            return
        }
        c.Next()
    }
}

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