package routes

import (
	"Ceremony/internal/handlers/template"
	"Ceremony/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Template(r *gin.Engine, db *gorm.DB) {
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(db))
	{
		org := protected.Group("/organization/:slug")
		org.Use(middleware.OrgMiddleware(db))
		{
			org.POST("/templates", func(c *gin.Context) { template.Createtemplate(c, db) })
			org.GET("/templates", func(c *gin.Context) { template.Listtemplates(c, db) })
			org.GET("/templates/:templateSlug", func(c *gin.Context) { template.Gettemplate(c, db) })
			org.DELETE("/templates/:templateSlug", func(c *gin.Context) { template.Deletetemplate(c, db) })
		}
	}
}