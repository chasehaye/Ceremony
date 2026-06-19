package middleware

import (
    "Ceremony/internal/dtos"
    "net/http"
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "Ceremony/internal/models"
    "Ceremony/internal/config"
    "errors"
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

func SuperAdminMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        uidVal, exists := c.Get("userID")
        if !exists {
            c.JSON(401, gin.H{"error": "missing user context"})
            c.Abort()
            return
        }

        uid := uidVal.(uint)

        var user models.User
        if err := db.Select("email").First(&user, uid).Error; err != nil {
            c.JSON(401, gin.H{"error": "user not found"})
            c.Abort()
            return
        }

        if user.Email != config.App.AdminEmail {
            c.JSON(403, gin.H{"error": "super admin access required"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func ProtectSuperAdminMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.Next()
			return
		}

		var user models.User
		if err := db.Select("email").Where("id = ?", id).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, dtos.ServerErrorResponse{
					Error: "Failed to load user",
				})
			}
			return
		}

		if user.Email == config.App.AdminEmail {
			c.AbortWithStatusJSON(http.StatusForbidden, dtos.ForbiddenResponse{
				Error: "This admin account cannot be modified",
			})
			return
		}

		c.Next()
	}
}