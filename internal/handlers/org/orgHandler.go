package org

import (
	"Ceremony/internal/dtos"
	"Ceremony/internal/models"
	"Ceremony/internal/crypt"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrg(c *gin.Context, db *gorm.DB) {
	userID := c.MustGet("userID").(uint)

	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{Error: err.Error()})
		return
	}

	var slug string

	for range 5 {
		s, err := crypt.GenerateSlug()
		if err != nil {
			c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{
				Error: "Failed to generate organization slug",
			})
			return
		}

		var count int64
		if err := db.Model(&models.Organization{}).
			Where("slug = ?", s).
			Count(&count).Error; err != nil {

			c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{
				Error: "Failed to check slug uniqueness",
			})
			return
		}

		if count == 0 {
			slug = s
			break
		}
	}

	if slug == "" {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to generate unique slug"})
		return
	}

	org := models.Organization{
		Name: input.Name,
		Slug: slug,
	}

	if err := db.Create(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to create organization"})
		return
	}

	member := models.OrganizationMember{
		OrganizationID: org.ID,
		UserID:         userID,
		Role:           "owner",
	}

	if err := db.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to assign organization owner"})
		return
	}

	c.JSON(http.StatusCreated, CreateOrgResponse{
		Organization: OrgResponse{
			Name: org.Name,
			Slug: org.Slug,
		},
	})
}

func GetOrg(c *gin.Context, db *gorm.DB) {
	userID := c.MustGet("userID").(uint)
	slug := c.Param("slug")

	var org models.Organization
	if err := db.Where("slug = ?", slug).First(&org).Error; err != nil {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Organization not found"})
		return
	}

	var member models.OrganizationMember
	if err := db.Where("organization_id = ? AND user_id = ?", org.ID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "Access denied"})
		return
	}

	c.JSON(http.StatusOK, OrgWithRoleResponse{
		Organization: OrgResponse{
			Name: org.Name,
			Slug: org.Slug,
		},
		Role: member.Role,
	})
}

func GetUserOrgs(c *gin.Context, db *gorm.DB) {
	userID := c.MustGet("userID").(uint)

	var memberships []models.OrganizationMember
	if err := db.Preload("Organization").
		Joins("JOIN organizations ON organizations.id = organization_members.organization_id AND organizations.deleted_at IS NULL").
		Where("user_id = ?", userID).
		Order("organization_members.created_at DESC").
		Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to fetch organizations"})
		return
	}

	result := make([]OrgWithRoleResponse, 0, len(memberships))
	for _, m := range memberships {
		if m.Organization.ID == 0 {
			continue
		}
		result = append(result, OrgWithRoleResponse{
			Organization: OrgResponse{
				Name: m.Organization.Name,
				Slug: m.Organization.Slug,
			},
			Role: m.Role,
		})
	}

	c.JSON(http.StatusOK, UserOrgsResponse{Organizations: result})
}

func DeleteOrg(c *gin.Context, db *gorm.DB) {
	userID := c.MustGet("userID").(uint)
	slug := c.Param("slug")

	var org models.Organization
	if err := db.Where("slug = ?", slug).First(&org).Error; err != nil {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Organization not found"})
		return
	}

	var member models.OrganizationMember
	if err := db.Where("organization_id = ? AND user_id = ? AND role = ?", org.ID, userID, "owner").First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "Only the owner can delete this organization"})
		return
	}

	if err := db.Delete(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, DeleteOrgResponse{Message: "Organization deleted"})
}

func GetStats(c *gin.Context, db *gorm.DB) {
    orgID := c.MustGet("orgID").(uint)

    var totalEmails int64
    db.Model(&models.EmailLog{}).Where("organization_id = ?", orgID).Count(&totalEmails)

    var sentEmails int64
    db.Model(&models.EmailLog{}).Where("organization_id = ? AND status = ?", orgID, "sent").Count(&sentEmails)

    var failedEmails int64
    db.Model(&models.EmailLog{}).Where("organization_id = ? AND status = ?", orgID, "failed").Count(&failedEmails)

    var pendingEmails int64
    db.Model(&models.EmailLog{}).Where("organization_id = ? AND status = ?", orgID, "pending").Count(&pendingEmails)

    var totalTemplates int64
    db.Model(&models.EmailTemplate{}).Where("organization_id = ?", orgID).Count(&totalTemplates)

    var activeApps int64
    db.Model(&models.App{}).Where("organization_id = ? AND is_active = ?", orgID, true).Count(&activeApps)

    var recentLogs []models.EmailLog
    db.Where("organization_id = ?", orgID).
        Order("created_at DESC").
        Limit(5).
        Find(&recentLogs)

    recentItems := make([]RecentLogResponse, len(recentLogs))
    for i, log := range recentLogs {
        recentItems[i] = RecentLogResponse{
            ToEmail:   log.ToEmail,
            Subject:   log.Subject,
            Status:    log.Status,
            CreatedAt: log.CreatedAt,
        }
    }

    c.JSON(http.StatusOK, StatsResponse{
        TotalEmails:    totalEmails,
        SentEmails:     sentEmails,
        FailedEmails:   failedEmails,
        PendingEmails:  pendingEmails,
        TotalTemplates: totalTemplates,
        ActiveApps:     activeApps,
        RecentLogs:     recentItems,
    })
}