package template

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "Ceremony/internal/dtos"
    "Ceremony/internal/models"
)

func Createtemplate(c *gin.Context, db *gorm.DB) {
    orgID := c.MustGet("orgID").(uint)

    var input CreateTemplateInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{
            Error: "Invalid input",
            Details: map[string]string{
                "name": "Name is required",
            },
        })
        return
    }

    tmpl := models.EmailTemplate{
        OrganizationID: orgID,
        Name:           input.Name,
        Subject:        input.Subject,
        Body:           input.Body,
        Type:           input.Type,
        IsActive:       true,
    }

    if err := db.Create(&tmpl).Error; err != nil {
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to create template"})
        return
    }

    c.JSON(http.StatusCreated, toTemplateResponse(tmpl))
}

func Listtemplates(c *gin.Context, db *gorm.DB) {
    orgID := c.MustGet("orgID").(uint)

    var templates []models.EmailTemplate
    if err := db.Where("organization_id = ?", orgID).Order("created_at DESC").Find(&templates).Error; err != nil {
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to fetch templates"})
        return
    }

    items := make([]TemplateResponse, len(templates))
    for i, t := range templates {
        items[i] = toTemplateResponse(t)
    }

    c.JSON(http.StatusOK, ListTemplatesResponse{Templates: items, Total: len(items)})
}

func Gettemplate(c *gin.Context, db *gorm.DB) {
    orgID := c.MustGet("orgID").(uint)
    id := c.Param("templateSlug")

    var tmpl models.EmailTemplate
    if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&tmpl).Error; err != nil {
        c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Template not found"})
        return
    }

    c.JSON(http.StatusOK, toTemplateResponse(tmpl))
}

func Deletetemplate(c *gin.Context, db *gorm.DB) {
    orgID := c.MustGet("orgID").(uint)
    id := c.Param("templateSlug")

    result := db.Where("id = ? AND organization_id = ?", id, orgID).Delete(&models.EmailTemplate{})
    if result.Error != nil || result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Template not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Template deleted"})
}

func toTemplateResponse(t models.EmailTemplate) TemplateResponse {
    return TemplateResponse{
        ID:        t.ID,
        Name:      t.Name,
        Subject:   t.Subject,
        Body:      t.Body,
        Type:      t.Type,
        IsActive:  t.IsActive,
        CreatedAt: t.CreatedAt.String(),
    }
}