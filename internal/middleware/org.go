package middleware

import (
	"Ceremony/internal/dtos"
	"Ceremony/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OrgMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		slug := c.Param("slug")

		var org models.Organization
		if err := db.Where("slug = ?", slug).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, dtos.NotFoundErrorResponse{
					Error: "Organization not found",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, dtos.ServerErrorResponse{
				Error: "Failed to fetch organization",
			})
			return
		}

		var member models.OrganizationMember
		if err := db.Where("organization_id = ? AND user_id = ?", org.ID, userID).First(&member).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dtos.ForbiddenResponse{
				Error: "You are not a member of this organization",
			})
			return
		}

		c.Set("orgID", org.ID)
		c.Set("orgRole", member.Role)
		c.Next()
	}
}