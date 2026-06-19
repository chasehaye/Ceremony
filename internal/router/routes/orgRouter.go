package routes

import (
	"Ceremony/internal/handlers/org"
	"Ceremony/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Org(r *gin.Engine, db *gorm.DB) {
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(db))
	{
		protected.POST("/organization", func(c *gin.Context) { org.CreateOrg(c, db) })
		protected.GET("/organization/:slug", func(c *gin.Context) { org.GetOrg(c, db) })
		protected.GET("/organizations", func(c *gin.Context) { org.GetUserOrgs(c, db) })
		protected.DELETE("/organization/:slug", func(c *gin.Context) { org.DeleteOrg(c, db) })
		
		orgGroup := protected.Group("/organization/:slug")
		orgGroup.Use(middleware.OrgMiddleware(db))
		{
			orgGroup.GET("/stats", func(c *gin.Context) { org.GetStats(c, db) })
		}
	}
}