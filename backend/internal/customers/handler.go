package customers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/httpx"
)

// Handler xử lý các route HTTP cho khách hàng.
type Handler struct {
	service *Service
}

// NewHandler tạo Handler mới.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List xử lý GET /api/customers — danh sách khách hàng có phân trang.
func (h *Handler) List(c *gin.Context) {
	p := httpx.ParsePagination(c)
	result, err := h.service.List(p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "không lấy được danh sách khách hàng")
		return
	}
	httpx.OK(c, result)
}

// Get xử lý GET /api/customers/:id — chi tiết một khách hàng.
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	customer, err := h.service.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, customer)
}

// Create xử lý POST /api/customers — tạo khách hàng mới.
func (h *Handler) Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}
	if input.CustomerCode == "" || input.CustomerName == "" {
		httpx.Error(c, http.StatusBadRequest, "mã và tên khách hàng là bắt buộc")
		return
	}
	customer, err := h.service.Create(input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.Created(c, customer)
}

// Update xử lý PUT /api/customers/:id — cập nhật khách hàng.
func (h *Handler) Update(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var input UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}
	if input.CustomerCode == "" || input.CustomerName == "" {
		httpx.Error(c, http.StatusBadRequest, "mã và tên khách hàng là bắt buộc")
		return
	}
	customer, err := h.service.Update(id, input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, customer)
}

// Delete xử lý DELETE /api/customers/:id — xóa mềm khách hàng.
func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id); err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, gin.H{"deleted": true, "id": id})
}

// Restore xử lý POST /api/customers/:id/restore — khôi phục khách hàng đã xóa mềm.
func (h *Handler) Restore(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	customer, err := h.service.Restore(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, customer)
}

// parseIDParam đọc và kiểm tra tham số :id trên URL.
// Trả về (id, true) nếu hợp lệ; nếu không, tự trả lỗi 400 và trả (0, false).
func parseIDParam(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		httpx.Error(c, http.StatusBadRequest, "id không hợp lệ")
		return 0, false
	}
	return id, true
}

// InactiveCodeExists là mã lỗi ổn định mà client dùng để nhận biết mã trùng
// với bản ghi đã xóa mềm (cần hỏi người dùng có muốn khôi phục hay không).
const InactiveCodeExists = "INACTIVE_CODE_EXISTS"

// handleError chuyển lỗi nghiệp vụ thành mã HTTP phù hợp.
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.Error(c, http.StatusNotFound, "không tìm thấy khách hàng")
	case errors.Is(err, ErrCodeExists):
		httpx.Error(c, http.StatusConflict, "mã khách hàng đã tồn tại")
	default:
		// Trùng mã với bản ghi đã xóa mềm: trả kèm bản ghi cũ cho client.
		if inactive, ok := AsInactiveCodeError(err); ok {
			httpx.ConflictWithCode(
				c,
				"Đã từng tồn tại khách hàng với mã này",
				InactiveCodeExists,
				inactive.Existing,
			)
			return
		}
		httpx.Error(c, http.StatusInternalServerError, "lỗi hệ thống")
	}
}
