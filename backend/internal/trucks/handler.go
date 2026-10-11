package trucks

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/httpx"
)

// Handler xử lý các route HTTP cho xe bồn.
type Handler struct {
	service *Service
}

// NewHandler tạo Handler mới.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List xử lý GET /api/mixer-trucks.
func (h *Handler) List(c *gin.Context) {
	p := httpx.ParsePagination(c)
	result, err := h.service.List(p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "không lấy được danh sách xe bồn")
		return
	}
	httpx.OK(c, result)
}

// Get xử lý GET /api/mixer-trucks/:id.
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	truck, err := h.service.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, truck)
}

// Create xử lý POST /api/mixer-trucks.
func (h *Handler) Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}
	if input.TruckCode == "" || input.PlateNumber == "" {
		httpx.Error(c, http.StatusBadRequest, "mã xe và biển số là bắt buộc")
		return
	}
	created, err := h.service.Create(input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.Created(c, created)
}

// Update xử lý PUT /api/mixer-trucks/:id.
func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}
	if input.TruckCode == "" || input.PlateNumber == "" {
		httpx.Error(c, http.StatusBadRequest, "mã xe và biển số là bắt buộc")
		return
	}
	updated, err := h.service.Update(id, input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, updated)
}

// Delete xử lý DELETE /api/mixer-trucks/:id (xóa mềm).
func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id); err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, gin.H{"deleted": true, "id": id})
}

// Restore xử lý POST /api/mixer-trucks/:id/restore — khôi phục bản ghi đã xóa mềm.
func (h *Handler) Restore(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	truck, err := h.service.Restore(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, truck)
}

// parseID đọc và kiểm tra tham số :id trên URL.
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		httpx.Error(c, http.StatusBadRequest, "id khong hop le")
		return 0, false
	}
	return id, true
}

// InactiveCodeExists là mã lỗi ổn định cho trường hợp mã/biển số trùng bản ghi đã xóa mềm.
const InactiveCodeExists = "INACTIVE_CODE_EXISTS"

// handleError chuyển lỗi nghiệp vụ thành mã HTTP phù hợp.
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.Error(c, http.StatusNotFound, "không tìm thấy xe bồn")
	case errors.Is(err, ErrCodeExists):
		httpx.Error(c, http.StatusConflict, "mã xe hoặc biển số đã tồn tại")
	case errors.Is(err, ErrDriverNotFound):
		httpx.Error(c, http.StatusBadRequest, "tài xế mặc định không tồn tại hoặc đã ngừng hoạt động")
	default:
		if inactive, ok := AsInactiveCodeError(err); ok {
			httpx.ConflictWithCode(
				c,
				"Đã từng tồn tại xe bồn với mã xe hoặc biển số này",
				InactiveCodeExists,
				inactive.Existing,
			)
			return
		}
		httpx.Error(c, http.StatusInternalServerError, "lỗi hệ thống")
	}
}
