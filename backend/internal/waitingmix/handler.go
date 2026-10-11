package waitingmix

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/httpx"
)

// Handler xử lý HTTP cho cấp phối chờ.
type Handler struct {
	service *Service
}

// NewHandler tạo Handler mới.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List xử lý GET /api/v1/waiting-mixes (trả đủ 20 vị trí).
func (h *Handler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "không lấy được danh sách cấp phối chờ")
		return
	}
	httpx.OK(c, gin.H{"items": items})
}

// Get xử lý GET /api/v1/waiting-mixes/:id.
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.service.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, item)
}

// Update xử lý PATCH /api/v1/waiting-mixes/:id.
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
	item, err := h.service.Update(id, input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, item)
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
		httpx.Error(c, http.StatusNotFound, "không tìm thấy vị trí cấp phối chờ")
	case errors.Is(err, ErrInvalidStatus), errors.Is(err, ErrInvalidVolume),
		errors.Is(err, ErrInvalidGrade), errors.Is(err, ErrInvalidCustomer):
		httpx.Error(c, http.StatusUnprocessableEntity, err.Error())
	default:
		httpx.Error(c, http.StatusInternalServerError, "lỗi hệ thống")
	}
}
