package org

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"Ceremony/internal/config"
	"Ceremony/internal/dtos"
	"Ceremony/internal/models"
)

var validInviteRoles = map[string]bool{"member": true, "admin": true}
var errInviteConsumed = errors.New("invite already consumed")

func generateInviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// CreateInvite — owner/admin only. Runs under OrgMiddleware, so orgID/orgRole are set.
func CreateInvite(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	userID := c.MustGet("userID").(uint)
	role := c.MustGet("orgRole").(string)

	if role != "owner" && role != "admin" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "You do not have permission to invite members"})
		return
	}

	var input CreateInviteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{Error: err.Error()})
		return
	}

	invitedRole := strings.ToLower(strings.TrimSpace(input.Role))
	if invitedRole == "" {
		invitedRole = "member"
	}
	if !validInviteRoles[invitedRole] {
		c.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{Error: "Invalid role"})
		return
	}
	// Admins can only invite members; only owners can mint admins.
	if invitedRole == "admin" && role != "owner" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "Only the owner can invite admins"})
		return
	}

	token, err := generateInviteToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to generate invite"})
		return
	}

	invite := models.OrganizationInvite{
		OrganizationID: orgID,
		InvitedByID:    userID,
		Email:          strings.ToLower(strings.TrimSpace(input.Email)),
		Role:           invitedRole,
		Token:          token,
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour),
	}
	if err := db.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to create invite"})
		return
	}

	c.JSON(http.StatusCreated, InviteResponse{
		ID:        invite.ID,
		Token:     invite.Token,
		Email:     invite.Email,
		Role:      invite.Role,
		ExpiresAt: invite.ExpiresAt,
		Link:      config.App.FrontendURL + "/invite/" + invite.Token,
	})
}

// GetInvite — public preview so the invitee sees what they're joining before logging in.
func GetInvite(c *gin.Context, db *gorm.DB) {
	token := c.Param("token")

	var invite models.OrganizationInvite
	if err := db.Preload("Organization").Where("token = ?", token).First(&invite).Error; err != nil {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Invite not found"})
		return
	}
	if invite.Used || time.Now().After(invite.ExpiresAt) {
		c.JSON(http.StatusGone, dtos.ForbiddenResponse{Error: "This invite has expired or already been used"})
		return
	}

	c.JSON(http.StatusOK, InvitePreviewResponse{
		OrganizationName: invite.Organization.Name,
		OrganizationSlug: invite.Organization.Slug,
		Role:             invite.Role,
		Email:            invite.Email,
	})
}

// AcceptInvite — requires login (but not membership). Claims the invite + creates membership atomically.
func AcceptInvite(c *gin.Context, db *gorm.DB) {
	userID := c.MustGet("userID").(uint)
	token := c.Param("token")

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse{Error: "User not found"})
		return
	}

	var invite models.OrganizationInvite
	if err := db.Where("token = ?", token).First(&invite).Error; err != nil {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Invite not found"})
		return
	}
	if invite.Used || time.Now().After(invite.ExpiresAt) {
		c.JSON(http.StatusGone, dtos.ForbiddenResponse{Error: "This invite has expired or already been used"})
		return
	}

	// Email-bound invites can only be accepted by the matching account.
	if invite.Email != "" && !strings.EqualFold(invite.Email, user.Email) {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "This invite was issued to a different email address"})
		return
	}

	var existing int64
	if err := db.Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", invite.OrganizationID, userID).
		Count(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to check membership"})
		return
	}
	if existing > 0 {
		c.JSON(http.StatusConflict, dtos.ForbiddenResponse{Error: "You are already a member of this organization"})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// Atomically claim the invite: only one accept can flip used false->true.
		res := tx.Model(&models.OrganizationInvite{}).
			Where("id = ? AND used = ?", invite.ID, false).
			Update("used", true)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errInviteConsumed
		}
		return tx.Create(&models.OrganizationMember{
			OrganizationID: invite.OrganizationID,
			UserID:         userID,
			Role:           invite.Role,
		}).Error
	})
	if errors.Is(err, errInviteConsumed) {
		c.JSON(http.StatusGone, dtos.ForbiddenResponse{Error: "This invite has already been used"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to accept invite"})
		return
	}

	var org models.Organization
	db.Select("name", "slug").First(&org, invite.OrganizationID)

	c.JSON(http.StatusOK, AcceptInviteResponse{
		Message:          "You have joined the organization",
		OrganizationName: org.Name,
		OrganizationSlug: org.Slug,
		Role:             invite.Role,
	})
}

// ListInvites — owner/admin: outstanding (unused) invites.
func ListInvites(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	role := c.MustGet("orgRole").(string)

	if role != "owner" && role != "admin" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "You do not have permission to view invites"})
		return
	}

	var invites []models.OrganizationInvite
	if err := db.Where("organization_id = ? AND used = ?", orgID, false).
		Order("created_at DESC").Find(&invites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to fetch invites"})
		return
	}

	items := make([]InviteResponse, 0, len(invites))
	for _, inv := range invites {
		items = append(items, InviteResponse{
			ID:        inv.ID,
			Token:     inv.Token,
			Email:     inv.Email,
			Role:      inv.Role,
			ExpiresAt: inv.ExpiresAt,
			Link:      config.App.FrontendURL + "/invite/" + inv.Token,
		})
	}
	c.JSON(http.StatusOK, ListInvitesResponse{Invites: items})
}

// RevokeInvite — owner/admin: delete an outstanding invite.
func RevokeInvite(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	role := c.MustGet("orgRole").(string)
	id := c.Param("id")

	if role != "owner" && role != "admin" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "You do not have permission to revoke invites"})
		return
	}

	res := db.Where("id = ? AND organization_id = ?", id, orgID).Delete(&models.OrganizationInvite{})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to revoke invite"})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Invite not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Invite revoked"})
}