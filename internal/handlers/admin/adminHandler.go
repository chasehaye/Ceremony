package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"Ceremony/internal/dtos"
	"Ceremony/internal/models"
	"Ceremony/internal/config"
)

func ListUsers(c *gin.Context, db *gorm.DB) {
	var users []models.User
	if err := db.Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to fetch users"})
		return
	}

	items := make([]UserListItem, len(users))
	for i, u := range users {
		items[i] = UserListItem{
			ID:         u.ID,
			CreatedAt:  u.CreatedAt,
			UpdatedAt:  u.UpdatedAt,
			Name:       u.Name,
			Email:      u.Email,
			IsAdmin:    u.IsAdmin,
			IsVerified: u.IsVerified,
			IsApproved: u.IsApproved,
			IsBanned:   u.IsBanned,
			CanCreate:  u.CanCreate,
			IsSuperAdmin: u.Email == config.App.AdminEmail,
		}
	}

	c.JSON(http.StatusOK, ListUsersResponse{Users: items, Total: len(items)})
}

func ApproveUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Model(&models.User{}).Where("id = ?", id).Update("is_approved", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to approve user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User approved"})
}

func RejectUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Model(&models.User{}).Where("id = ?", id).Update("is_approved", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to reject user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User rejected"})
}

func BanUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Model(&models.User{}).Where("id = ?", id).Update("is_banned", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to ban user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User banned"})
}

func UnbanUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Model(&models.User{}).Where("id = ?", id).Update("is_banned", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to unban user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User unbanned"})
}

func GrantCreate(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Model(&models.User{}).Where("id = ?", id).Update("can_create", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to grant create permission"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create permission granted"})
}

func RevokeCreate(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Model(&models.User{}).Where("id = ?", id).Update("can_create", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to revoke create permission"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create permission revoked"})
}

func GrantAdmin(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Model(&models.User{}).
		Where("id = ?", id).
		Update("is_admin", true).Error; err != nil {

		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{
			Error: "Failed to grant admin role",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin role granted",
	})
}

func RevokeAdmin(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Model(&models.User{}).
		Where("id = ?", id).
		Update("is_admin", false).Error; err != nil {

		c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{
			Error: "Failed to revoke admin role",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin role revoked",
	})
}