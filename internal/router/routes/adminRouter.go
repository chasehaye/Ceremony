package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"Ceremony/internal/middleware"
	"Ceremony/internal/handlers/admin"
)

func Admin(r *gin.Engine, db *gorm.DB) {
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(db))
	protected.Use(middleware.AdminMiddleware(db))
	{
		users := protected.Group("/admin/users")
		users.GET("", func(c *gin.Context) { admin.ListUsers(c, db) })

		// All mutations on a specific user are blocked from targeting the super admin.
		mut := users.Group("/:id")
		mut.Use(middleware.ProtectSuperAdminMiddleware(db))
		{
			mut.PATCH("/approve", func(c *gin.Context) { admin.ApproveUser(c, db) })
			mut.PATCH("/reject", func(c *gin.Context) { admin.RejectUser(c, db) })
			mut.PATCH("/ban", func(c *gin.Context) { admin.BanUser(c, db) })
			mut.PATCH("/unban", func(c *gin.Context) { admin.UnbanUser(c, db) })
			mut.PATCH("/grant-create", func(c *gin.Context) { admin.GrantCreate(c, db) })
			mut.PATCH("/revoke-create", func(c *gin.Context) { admin.RevokeCreate(c, db) })

			mut.PATCH("/grant-admin",
				middleware.SuperAdminMiddleware(db),
				func(c *gin.Context) { admin.GrantAdmin(c, db) },
			)
			mut.PATCH("/revoke-admin",
				middleware.SuperAdminMiddleware(db),
				func(c *gin.Context) { admin.RevokeAdmin(c, db) },
			)
		}
	}
}