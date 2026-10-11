package drivers

import (
	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/database"
)

// Register gắn các route tài xế vào nhóm router cho trước.
// Tất cả yêu cầu xác thực; thao tác xóa chỉ dành cho admin/manager.
func Register(api *gin.RouterGroup, db *database.DB, requireAuth, requireAdminOrManager gin.HandlerFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	group := api.Group("/drivers")
	group.Use(requireAuth)

	group.GET("", h.List)
	group.GET("/:id", h.Get)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", requireAdminOrManager, h.Delete)
	// Khôi phục bản ghi đã xóa mềm: chỉ admin/manager được phép.
	group.POST("/:id/restore", requireAdminOrManager, h.Restore)
}
