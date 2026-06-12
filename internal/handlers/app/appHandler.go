package app

import (
    "log"
    "net/http"
    "errors"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "Ceremony/internal/crypt"
    "Ceremony/internal/dtos"
    "Ceremony/internal/models"
)

func CreateApp(c *gin.Context, db *gorm.DB) {
    orgID := c.MustGet("orgID").(uint)

    var input CreateAppInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{
            Error: "Invalid input",
            Details: map[string]string{
                "name": "Name is required",
            },
        })
        return
    }

    var slug string
    for range 5 {
        s, err := crypt.GenerateSlug()
        if err != nil {
            c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to generate app slug"})
            return
        }
        var existing models.App
        if errors.Is(db.Where("slug = ?", s).First(&existing).Error, gorm.ErrRecordNotFound) {
            slug = s
            break
        }
    }

    if slug == "" {
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to generate unique app slug"})
        return
    }

    apiKey, err := crypt.GenerateToken()
    if err != nil {
        log.Printf("Failed to generate API key: %v", err)
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to generate API key"})
        return
    }

    app := models.App{
        OrganizationID: orgID,
        Name:           input.Name,
        Slug:           slug,
        Description:    input.Description,
        APIKey:         apiKey,
    }

    if err := db.Create(&app).Error; err != nil {
        log.Printf("Failed to create app: %v", err)
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to create app"})
        return
    }

    c.JSON(http.StatusCreated, AppResponse{
        ID:          app.ID,
        Slug:        app.Slug,
        Name:        app.Name,
        Description: app.Description,
        APIKey:      app.APIKey,
        IsActive:    app.IsActive,
        CreatedAt:   app.CreatedAt.String(),
    })
}

func ListApps(c *gin.Context, db *gorm.DB) {
    orgID := c.MustGet("orgID").(uint)

    var apps []models.App
    if err := db.Where("organization_id = ?", orgID).Order("created_at DESC").Find(&apps).Error; err != nil {
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to fetch apps"})
        return
    }

    items := make([]AppResponse, len(apps))
    for i, a := range apps {
        items[i] = AppResponse{
            ID:          a.ID,
            Slug:        a.Slug,
            Name:        a.Name,
            Description: a.Description,
            APIKey:      a.APIKey,
            IsActive:    a.IsActive,
            CreatedAt:   a.CreatedAt.String(),
        }
    }

    c.JSON(http.StatusOK, ListAppsResponse{Apps: items, Total: len(items)})
}

func GetApp(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	slug := c.Param("appSlug")

	var app models.App
	result := db.Where("slug = ? AND organization_id = ?", slug, orgID).First(&app)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "App not found"})
		return
	}

	c.JSON(http.StatusOK, AppResponse{
		ID:          app.ID,
		Slug:        app.Slug,
		Name:        app.Name,
		Description: app.Description,
		APIKey:      app.APIKey,
		IsActive:    app.IsActive,
		CreatedAt:   app.CreatedAt.String(),
	})
}

func DeleteApp(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	slug := c.Param("appSlug")

	result := db.Where("slug = ? AND organization_id = ?", slug, orgID).Delete(&models.App{})
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "App not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "App deleted"})
}

func RotateKey(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	slug := c.Param("appSlug")

	var app models.App
	result := db.Where("slug = ? AND organization_id = ?", slug, orgID).First(&app)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "App not found"})
		return
	}

	newKey, err := crypt.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to generate new key"})
		return
	}

	if err := db.Model(&app).Update("api_key", newKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to rotate key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": newKey})
}