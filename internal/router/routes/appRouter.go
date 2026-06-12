package routes

import (
	"Ceremony/internal/handlers/app"
	"Ceremony/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func App(r *gin.Engine, db *gorm.DB) {
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		org := protected.Group("/organization/:slug")
		org.Use(middleware.OrgMiddleware(db))
		{
			org.POST("/apps", func(c *gin.Context) { app.CreateApp(c, db) })
			org.GET("/apps", func(c *gin.Context) { app.ListApps(c, db) })
			org.GET("/apps/:appSlug", func(c *gin.Context) { app.GetApp(c, db) })
			org.DELETE("/apps/:appSlug", func(c *gin.Context) { app.DeleteApp(c, db) })
			org.POST("/apps/:appSlug/rotate-key", func(c *gin.Context) { app.RotateKey(c, db) })
		}
	}
}