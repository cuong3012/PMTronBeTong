package grades

import (
	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/database"
)

// Register gắn route mác bê tông vào nhóm router.
// Tạo/sửa/xóa yêu cầu quyền admin/manager (cấu hình công thức).
func Register(api *gin.RouterGroup, db *database.DB, requireAuth, requireAdminOrManager gin.HandlerFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	group := api.Group("/concrete-grades")
	group.Use(requireAuth)

	group.GET("", h.List)
	group.GET("/:id", h.Get)
	group.POST("", requireAdminOrManager, h.Create)
	group.PUT("/:id", requireAdminOrManager, h.Update)
	group.DELETE("/:id", requireAdminOrManager, h.Delete)
}
