package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "Ceremony/internal/dtos"
    "Ceremony/internal/models"
)

// ListUsers godoc
// @Summary      List All Users
// @Description  Returns all users sorted by creation date, newest first.
// @Tags         admin
// @Produce      json
// @Success      200  {object}  ListUsersResponse
// @Failure      500  {object}  dtos.ServerErrorResponse
// @Router       /api/admin/users [get]
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
        }
    }

    c.JSON(http.StatusOK, ListUsersResponse{
        Users: items,
        Total: len(items),
    })
}

// ApproveUser godoc
// @Summary      Approve User
// @Description  Sets is_approved to true for the specified user.
// @Tags         admin
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  dtos.ServerErrorResponse
// @Router       /api/admin/users/{id}/approve [patch]
func ApproveUser(c *gin.Context, db *gorm.DB) {
    id := c.Param("id")

    if err := db.Model(&models.User{}).Where("id = ?", id).Update("is_approved", true).Error; err != nil {
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to approve user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User approved"})
}

// RejectUser godoc
// @Summary      Reject User
// @Description  Sets is_approved to false for the specified user.
// @Tags         admin
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  dtos.ServerErrorResponse
// @Router       /api/admin/users/{id}/reject [patch]
func RejectUser(c *gin.Context, db *gorm.DB) {
    id := c.Param("id")

    if err := db.Model(&models.User{}).Where("id = ?", id).Update("is_approved", false).Error; err != nil {
        c.JSON(http.StatusInternalServerError, dtos.ServerErrorResponse{Error: "Failed to reject user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User rejected"})
}