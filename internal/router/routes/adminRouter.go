package routes

import (
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"

	"Ceremony/internal/middleware"
    "Ceremony/internal/handlers/admin"
)

func Admin(r *gin.Engine, db *gorm.DB){

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(db))
	protected.Use(middleware.AdminMiddleware(db))
    {
        protected.GET("/admin/users", func(c *gin.Context) {admin.ListUsers(c, db)})
        protected.PATCH("/admin/users/:id/approve", func(c *gin.Context) {admin.ApproveUser(c, db)})
        protected.PATCH("/admin/users/:id/reject", func(c *gin.Context) {admin.RejectUser(c, db)})
    }
}

