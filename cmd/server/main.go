// @title           Log Relay API
// @version         1.0
// @description     API Server for Log Relay Project.
// @host            localhost:8080
// @BasePath        /
package main

import (
    "log"
    "os"

	"github.com/joho/godotenv"
	"Ceremony/internal/database"
	"Ceremony/internal/config"
	"Ceremony/internal/router"
    "Ceremony/internal/models"
)

func main() {
    _ = godotenv.Load() 
    config.CheckRequiredEnvVarsAndLoad()
    db, err := database.ConnectToDB(config.App.DBHost, config.App.DBUser, config.App.DBPass, "ceremonyDB", config.App.DBPort)
    if err != nil {
        log.Fatalln("Failed to connect to database:", err)
    }
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatalln("Failed to get generic database object:", err)
    }
    if err := sqlDB.Ping(); err != nil {
        log.Fatalln("Database unreachable:", err)
    }
    defer sqlDB.Close()
    log.Printf("****************************************************************************")
    log.Println("----Database connection verified")
    err = db.AutoMigrate(
        &models.User{},
        &models.PasswordReset{},
        &models.EmailVerification{},
        &models.Organization{},
        &models.OrganizationMember{},
        &models.Domain{},
        &models.App{},
        &models.EmailTemplate{},
        &models.EmailLog{},
    )
    if err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
    log.Println("---Database migration successful")
    log.Println("-----------------------Connected")
    log.Printf("****************************************************************************")
    r := router.Setup(db)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("****************************************************************************")
    log.Printf("Server starting on port %s in %s mode...", config.App.Port, config.App.Env)
    log.Printf("****************************************************************************")
    r.Run(":" + port)
}
