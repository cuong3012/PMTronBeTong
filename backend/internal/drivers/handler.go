package drivers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/httpx"
)

// Handler xử lý các route HTTP cho tài xế.
type Handler struct {
	service *Service
}

// NewHandler tạo Handler mới.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List xử lý GET /api/drivers.
func (h *Handler) List(c *gin.Context) {
	p := httpx.ParsePagination(c)
	result, err := h.service.List(p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "không lấy được danh sách tài xế")
		return
	}
	httpx.OK(c, result)
}

// Get xử lý GET /api/drivers/:id.
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	driver, err := h.service.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, driver)
}

// Create xu ly POST /api/drivers (tao moi tai xe).
func (h *Handler) Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}
	if input.DriverCode == "" || input.FullName == "" {
		httpx.Error(c, http.StatusBadRequest, "mã tài xế và họ tên là bắt buộc")
		return
	}
	created, err := h.service.Create(input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.Created(c, created)
}

// Update xử lý PUT /api/drivers/:id.
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
	if input.DriverCode == "" || input.FullName == "" {
		httpx.Error(c, http.StatusBadRequest, "mã tài xế và họ tên là bắt buộc")
		return
	}
	updated, err := h.service.Update(id, input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, updated)
}

// Delete xử lý DELETE /api/drivers/:id (xóa mềm).
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

// Restore xử lý POST /api/drivers/:id/restore — khôi phục bản ghi đã xóa mềm.
func (h *Handler) Restore(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	driver, err := h.service.Restore(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, driver)
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

// InactiveCodeExists là mã lỗi ổn định cho trường hợp mã trùng bản ghi đã xóa mềm.
const InactiveCodeExists = "INACTIVE_CODE_EXISTS"

// handleError chuyển lỗi nghiệp vụ thành mã HTTP phù hợp.
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.Error(c, http.StatusNotFound, "không tìm thấy tài xế")
	case errors.Is(err, ErrCodeExists):
		httpx.Error(c, http.StatusConflict, "mã tài xế đã tồn tại")
	default:
		if inactive, ok := AsInactiveCodeError(err); ok {
			httpx.ConflictWithCode(
				c,
				"Đã từng tồn tại tài xế với mã này",
				InactiveCodeExists,
				inactive.Existing,
			)
			return
		}
		httpx.Error(c, http.StatusInternalServerError, "lỗi hệ thống")
	}
}
