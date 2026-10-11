package grades

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/httpx"
)

// Handler xử lý HTTP cho mác bê tông.
type Handler struct {
	service *Service
}

// NewHandler tạo Handler mới.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List xử lý GET /api/v1/concrete-grades.
func (h *Handler) List(c *gin.Context) {
	p := httpx.ParsePagination(c)
	result, err := h.service.List(p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "không lấy được danh sách mác bê tông")
		return
	}
	httpx.OK(c, result)
}

// Get xử lý GET /api/v1/concrete-grades/:id.
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	g, err := h.service.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, g)
}

// Create xử lý POST /api/v1/concrete-grades.
func (h *Handler) Create(c *gin.Context) {
	var input GradeInput
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

// Update xử lý PUT /api/v1/concrete-grades/:id.
func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input GradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}
	g, err := h.service.Update(id, input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, g)
}

// Delete xử lý DELETE /api/v1/concrete-grades/:id (xóa mềm).
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
		httpx.Error(c, http.StatusNotFound, "không tìm thấy mác bê tông")
	case errors.Is(err, ErrCodeExists):
		httpx.Error(c, http.StatusConflict, "mã mác bê tông đã tồn tại")
	case errors.Is(err, ErrInvalidCode), errors.Is(err, ErrPgTooManyDecimals):
		httpx.Error(c, http.StatusUnprocessableEntity, err.Error())
	default:
		httpx.Error(c, http.StatusInternalServerError, "lỗi hệ thống")
	}
}
