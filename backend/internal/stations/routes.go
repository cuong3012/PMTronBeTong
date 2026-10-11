package stations

import (
	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/database"
)

// Register gắn route cấu hình trạm (một bản ghi) vào nhóm router.
// Cập nhật chỉ dành cho admin/manager.
func Register(api *gin.RouterGroup, db *database.DB, requireAuth, requireAdminOrManager gin.HandlerFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	group := api.Group("/station")
	group.Use(requireAuth)

	group.GET("", h.Get)
	group.PATCH("", requireAdminOrManager, h.Update)
}
