package stations

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/httpx"
)

// Handler xử lý HTTP cho cấu hình trạm.
type Handler struct {
	service *Service
}

// NewHandler tạo Handler mới.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Get xử lý GET /api/v1/station.
func (h *Handler) Get(c *gin.Context) {
	st, err := h.service.Get()
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, st)
}

// Update xử lý PATCH /api/v1/station.
func (h *Handler) Update(c *gin.Context) {
	var input UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}
	st, err := h.service.Update(input)
	if err != nil {
		handleError(c, err)
		return
	}
	httpx.OK(c, st)
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.Error(c, http.StatusNotFound, "chưa có cấu hình trạm")
	case errors.Is(err, ErrInvalidMask), errors.Is(err, ErrInvalidBatchRange), errors.Is(err, ErrInvalidStationCode):
		httpx.Error(c, http.StatusUnprocessableEntity, err.Error())
	default:
		httpx.Error(c, http.StatusInternalServerError, "lỗi hệ thống")
	}
}
