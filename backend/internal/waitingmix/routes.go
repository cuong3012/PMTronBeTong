package waitingmix

import (
	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/database"
)

// Register gắn route cấp phối chờ vào nhóm router.
// Chỉ đọc và cập nhật; không tạo/xóa vị trí qua API.
func Register(api *gin.RouterGroup, db *database.DB, requireAuth gin.HandlerFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	group := api.Group("/waiting-mixes")
	group.Use(requireAuth)

	group.GET("", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
}
