package weighing

import (
	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/database"
)

// Register gắn route phiếu cân vào nhóm router.
// Tạo/xóa phiếu cân yêu cầu quyền admin/manager.
func Register(api *gin.RouterGroup, db *database.DB, requireAuth, requireAdminOrManager gin.HandlerFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	group := api.Group("/weighing-tickets")
	group.Use(requireAuth)

	group.GET("", h.List)
	group.GET("/stats", h.Stats)
	group.GET("/:id", h.Get)
	group.POST("", requireAdminOrManager, h.Create)
	group.DELETE("/:id", requireAdminOrManager, h.Delete)
}
