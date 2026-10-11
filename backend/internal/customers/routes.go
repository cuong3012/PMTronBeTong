package customers

import (
	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/database"
)

// Register gắn các route khách hàng vào nhóm router cho trước.
//
// Các route đều yêu cầu xác thực (requireAuth được truyền từ ngoài vào).
// Xóa khách hàng chỉ dành cho admin và manager.
func Register(api *gin.RouterGroup, db *database.DB, requireAuth, requireAdminOrManager gin.HandlerFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	group := api.Group("/customers")
	group.Use(requireAuth)

	group.GET("", h.List)
	group.GET("/:id", h.Get)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", requireAdminOrManager, h.Delete)
	// Khôi phục bản ghi đã xóa mềm: chỉ admin/manager được phép.
	group.POST("/:id/restore", requireAdminOrManager, h.Restore)
}
