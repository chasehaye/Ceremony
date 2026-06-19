package org

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"Ceremony/internal/dtos"
	"Ceremony/internal/models"
)

// ListMembers — any member of the org may view the roster.
func ListMembers(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)

	var members []models.OrganizationMember
	if err := db.Preload("User").
		Where("organization_id = ?", orgID).
		Order("role, created_at").
		Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to fetch members"})
		return
	}

	items := make([]MemberResponse, 0, len(members))
	for _, m := range members {
		items = append(items, MemberResponse{
			UserID: m.UserID,
			Name:   m.User.Name,
			Email:  m.User.Email,
			Role:   m.Role,
		})
	}
	c.JSON(http.StatusOK, ListMembersResponse{Members: items})
}

// ChangeMemberRole — owner only. Promotes/demotes between admin and member.
func ChangeMemberRole(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	role := c.MustGet("orgRole").(string)
	targetID := c.Param("userID")

	if role != "owner" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "Only the owner can change roles"})
		return
	}

	var input ChangeRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{Error: err.Error()})
		return
	}
	newRole := strings.ToLower(strings.TrimSpace(input.Role))
	if newRole != "admin" && newRole != "member" {
		c.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{Error: "Role must be admin or member"})
		return
	}

	var target models.OrganizationMember
	if err := db.Where("organization_id = ? AND user_id = ?", orgID, targetID).First(&target).Error; err != nil {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Member not found"})
		return
	}
	if target.Role == "owner" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "The owner's role cannot be changed"})
		return
	}

	if err := db.Model(&models.OrganizationMember{}).
		Where("id = ?", target.ID).
		Update("role", newRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to update role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role updated", "role": newRole})
}

// RemoveMember — owner/admin. Hard delete. Owner is protected; admins can only remove members.
func RemoveMember(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	callerRole := c.MustGet("orgRole").(string)
	targetID := c.Param("userID")

	if callerRole != "owner" && callerRole != "admin" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "You do not have permission to remove members"})
		return
	}

	var target models.OrganizationMember
	if err := db.Where("organization_id = ? AND user_id = ?", orgID, targetID).First(&target).Error; err != nil {
		c.JSON(http.StatusNotFound, dtos.NotFoundErrorResponse{Error: "Member not found"})
		return
	}
	if target.Role == "owner" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "The owner cannot be removed"})
		return
	}
	if callerRole == "admin" && target.Role != "member" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "Admins can only remove members"})
		return
	}

	if err := db.Where("id = ?", target.ID).Delete(&models.OrganizationMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to remove member"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// LeaveOrg — any member leaves themselves. The owner must transfer or delete first.
func LeaveOrg(c *gin.Context, db *gorm.DB) {
	orgID := c.MustGet("orgID").(uint)
	callerID := c.MustGet("userID").(uint)
	role := c.MustGet("orgRole").(string)

	if role == "owner" {
		c.JSON(http.StatusForbidden, dtos.ForbiddenResponse{Error: "The owner must transfer ownership or delete the organization before leaving"})
		return
	}

	if err := db.Where("organization_id = ? AND user_id = ?", orgID, callerID).
		Delete(&models.OrganizationMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to leave organization"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "You have left the organization"})
}