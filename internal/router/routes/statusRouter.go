package routes

import (
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"

	"Ceremony/internal/handlers"
)

func Status(r *gin.Engine, db *gorm.DB){

	guest := r.Group("/api")
    {
        guest.GET("/status", func(c *gin.Context) { handlers.Ping(c, db) })
    }
}
