package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken trả về khi token không hợp lệ, hết hạn hoặc sai chữ ký.
var ErrInvalidToken = errors.New("token không hợp lệ")

// Claims là nội dung (payload) nhúng trong JWT.
// Kế thừa RegisteredClaims để có exp/iat, đồng thời thêm thông tin nghiệp vụ.
type Claims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"username"`
	RoleCode string `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager tạo và xác thực JWT.
type TokenManager struct {
	secret []byte
	expiry time.Duration
}

// NewTokenManager khởi tạo TokenManager với khóa ký và thời hạn token.
func NewTokenManager(secret string, expiry time.Duration) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
		expiry: expiry,
	}
}

// Generate tạo một JWT mới cho người dùng với các thông tin cơ bản.
// Token được ký bằng HS256.
func (m *TokenManager) Generate(userID int64, username, roleCode string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.expiry)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleCode: roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ký JWT thất bại: %w", err)
	}
	return signed, expiresAt, nil
}

// Parse xác thực chuỗi token và trả về Claims nếu hợp lệ.
func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Chỉ chấp nhận phương thức ký HMAC để tránh tấn công đổi thuật toán.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("phương thức ký không hợp lệ: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
