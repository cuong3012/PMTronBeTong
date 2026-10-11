package weighing

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/httpx"
)

// Handler xử lý HTTP cho phiếu cân.
type Handler struct {
	service *Service
}

// NewHandler tạo Handler mới.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List xử lý GET /api/v1/weighing-tickets.
func (h *Handler) List(c *gin.Context) {
	p := httpx.ParsePagination(c)
	result, err := h.service.List(p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "không lấy được danh sách phiếu cân")
		return
	}
	httpx.OK(c, result)
}

// Get xử lý GET /api/v1/weighing-tickets/:id.
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	t, err := h.service.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, t)
}

// Create xử lý POST /api/v1/weighing-tickets.
func (h *Handler) Create(c *gin.Context) {
	var input TicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}
	created, err := h.service.Create(input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.Created(c, created)
}

// Delete xử lý DELETE /api/v1/weighing-tickets/:id (xóa mềm).
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

// Stats xử lý GET /api/v1/weighing-tickets/stats?from=&to=.
func (h *Handler) Stats(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	result, err := h.service.Stats(from, to)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "không lấy được thống kê phiếu cân")
		return
	}
	httpx.OK(c, result)
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		httpx.Error(c, http.StatusBadRequest, "id không hợp lệ")
		return 0, false
	}
	return id, true
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.Error(c, http.StatusNotFound, "không tìm thấy phiếu cân")
	case errors.Is(err, ErrCodeExists):
		httpx.Error(c, http.StatusConflict, "mã phiếu cân đã tồn tại")
	case errors.Is(err, ErrInvalidTicket), errors.Is(err, ErrInvalidDetail):
		httpx.Error(c, http.StatusUnprocessableEntity, err.Error())
	default:
		httpx.Error(c, http.StatusInternalServerError, "lỗi hệ thống")
	}
}
