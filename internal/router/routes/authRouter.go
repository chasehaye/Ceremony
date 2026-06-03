package routes

import (
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"

	"Ceremony/internal/middleware"
    "Ceremony/internal/handlers/auth"
)

func Auth(r *gin.Engine, db *gorm.DB){

	guest := r.Group("/api")
    {
        guest.POST("/user/register", func(c *gin.Context) {auth.RegisterUser(c, db)})
        guest.POST("/user/login", func(c *gin.Context) {auth.LoginUser(c, db)})
        guest.POST("/auth/forgot-password", func(c *gin.Context) {auth.ForgotPassword(c, db)})
        guest.POST("/auth/reset-password/:token", func(c *gin.Context) {auth.ResetPassword(c, db)})
    }




	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/me", func(c *gin.Context) {auth.GetMe(c, db)})
        protected.POST("/user/logout", func(c *gin.Context) {auth.LogOut(c, db)})
    }
}
