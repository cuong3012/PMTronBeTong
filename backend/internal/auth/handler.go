package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/httpx"
)

// gin context keys — giữ đồng bộ với middleware.CtxUserID.
// Định nghĩa lại ở đây (không import middleware) để tránh vòng lặp import.
const (
	ctxKeyUserID = "auth_user_id"
	ctxKeyRole   = "auth_role_code"
)

// loginRequest là body của yêu cầu đăng nhập.
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Handler xử lý các route liên quan đến xác thực.
type Handler struct {
	service *Service
}

// NewHandler tạo Handler với service cho trước.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Login xử lý POST /api/auth/login.
// Nhận username/password, trả về token và thông tin người dùng.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorWithDetail(c, http.StatusBadRequest, "dữ liệu không hợp lệ", err.Error())
		return
	}

	result, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			httpx.Error(c, http.StatusUnauthorized, "tên đăng nhập hoặc mật khẩu không đúng")
		case errors.Is(err, ErrAccountDisabled):
			httpx.Error(c, http.StatusForbidden, "tài khoản đã bị khóa")
		default:
			httpx.Error(c, http.StatusInternalServerError, "lỗi hệ thống khi đăng nhập")
		}
		return
	}

	httpx.OK(c, result)
}

// Me xử lý GET /api/auth/me (yêu cầu đã xác thực).
// Trả về thông tin người dùng hiện tại từ token.
func (h *Handler) Me(c *gin.Context) {
	v, _ := c.Get(ctxKeyUserID)
	userID, _ := v.(int64)

	info, err := h.service.Me(userID)
	if err != nil {
		httpx.Error(c, http.StatusNotFound, "không tìm thấy người dùng")
		return
	}
	httpx.OK(c, info)
}

// RegisterRoutes gắn các route auth vào router nhóm /api.
// login là route công khai; me yêu cầu middleware xác thực.
func (h *Handler) RegisterRoutes(api *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	api.POST("/auth/login", h.Login)
	api.GET("/auth/me", requireAuth, h.Me)
}
