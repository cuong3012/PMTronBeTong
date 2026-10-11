package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"tramtron/backend/internal/auth"
	"tramtron/backend/internal/httpx"
)

// ContextKey dùng để đặt/đọc giá trị trong gin.Context một cách an toàn.
const (
	// CtxUserID là khóa lưu user id trong context sau khi xác thực.
	CtxUserID = "auth_user_id"
	// CtxUsername là khóa lưu username trong context.
	CtxUsername = "auth_username"
	// CtxRoleCode là khóa lưu mã vai trò trong context.
	CtxRoleCode = "auth_role_code"
)

// RequireAuth là middleware yêu cầu token JWT hợp lệ.
// Token đọc từ header "Authorization: Bearer <token>".
// Nếu hợp lệ, thông tin người dùng được lưu vào context cho các handler sau.
func RequireAuth(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			httpx.Error(c, http.StatusUnauthorized, "thiếu header Authorization")
			return
		}

		// Yêu cầu định dạng "Bearer <token>".
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpx.Error(c, http.StatusUnauthorized, "định dạng Authorization không hợp lệ, cần 'Bearer <token>'")
			return
		}

		claims, err := tokens.Parse(parts[1])
		if err != nil {
			httpx.Error(c, http.StatusUnauthorized, "token không hợp lệ hoặc đã hết hạn")
			return
		}

		// Lưu thông tin vào context.
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRoleCode, claims.RoleCode)

		c.Next()
	}
}

// RequireRole là middleware yêu cầu người dùng có một trong các vai trò cho trước.
// Phải dùng SAU RequireAuth để đảm bảo role đã có trong context.
func RequireRole(roles ...string) gin.HandlerFunc {
	// Chuyển danh sách vai trò sang map để tra cứu O(1).
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role, exists := c.Get(CtxRoleCode)
		if !exists {
			httpx.Error(c, http.StatusUnauthorized, "chưa xác thực")
			return
		}
		if _, ok := allowed[role.(string)]; !ok {
			httpx.Error(c, http.StatusForbidden, "bạn không có quyền thực hiện thao tác này")
			return
		}
		c.Next()
	}
}

// CurrentUserID đọc user id từ context (trả về 0 nếu chưa xác thực).
func CurrentUserID(c *gin.Context) int64 {
	if v, ok := c.Get(CtxUserID); ok {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// CurrentRoleCode đọc mã vai trò từ context.
func CurrentRoleCode(c *gin.Context) string {
	if v, ok := c.Get(CtxRoleCode); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
