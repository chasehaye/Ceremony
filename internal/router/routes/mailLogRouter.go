package routes

import (
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"

	"Ceremony/internal/middleware"
    "Ceremony/internal/handlers/mailLog"
)

func MailLog(r *gin.Engine, db *gorm.DB){



    public := r.Group("/api")
    public.Use(middleware.APIKeyMiddleware(db))
    {
        public.POST("/mail/send", func(c *gin.Context) { mailLog.Send(c, db) })
    }

    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware(db))
    {
        org := protected.Group("/organization/:slug")
        org.Use(middleware.OrgMiddleware(db))
        {
            org.GET("/mail/logs", func(c *gin.Context) { mailLog.Logs(c, db) })
        }
    }
}
