package router

import (
	"log"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "Ceremony/internal/router/routes"
	"Ceremony/internal/config"
    "Ceremony/internal/middleware"


    _ "Ceremony/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(db *gorm.DB) *gin.Engine {
    if config.IsProduction() {
        log.Println("Running in Production mode")
        gin.SetMode(gin.ReleaseMode)
    } else {
        log.Println("Running in Development mode")
        gin.SetMode(gin.DebugMode)
    }

    r := gin.New()
    r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
    r.Use(gin.Recovery())
    if !config.IsProduction() {
        r.Use(gin.Logger())
    }

    r.Use(middleware.CORSMiddleware())

    routes.Status(r, db)
    routes.Auth(r, db)
    routes.Admin(r, db)
    routes.App(r, db)
    routes.MailLog(r, db)
    routes.Org(r, db)


    if !config.IsProduction() {
        r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    }

    return r
}

// public  // open to anyone, API key auth
// guest   // frontend, no token needed  
// auth    // frontend, token required
// admin   // frontend, token + admin