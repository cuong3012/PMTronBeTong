package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"tramtron/backend/internal/models"
)

// ErrInvalidCredentials trả về khi sai username hoặc mật khẩu.
var ErrInvalidCredentials = errors.New("tên đăng nhập hoặc mật khẩu không đúng")

// ErrAccountDisabled trả về khi tài khoản bị khóa (is_active = 0).
var ErrAccountDisabled = errors.New("tài khoản đã bị khóa")

// Service chứa logic xác thực đăng nhập.
type Service struct {
	repo   *Repository
	tokens *TokenManager
}

// NewService khởi tạo Service với repository và token manager.
func NewService(repo *Repository, tokens *TokenManager) *Service {
	return &Service{repo: repo, tokens: tokens}
}

// LoginResult là kết quả trả về sau khi đăng nhập thành công.
type LoginResult struct {
	Token     string    `json:"token"`
	ExpiresAt string    `json:"expires_at"`
	User      *UserInfo `json:"user"`
}

// UserInfo là thông tin người dùng an toàn để trả về cho client (không chứa hash).
type UserInfo struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	RoleCode  string `json:"role_code"`
	RoleName  string `json:"role_name"`
	StationID *int64 `json:"station_id,omitempty"`
}

// Login xác thực username/mật khẩu và trả về token nếu hợp lệ.
//
// Quy trình:
//  1. Tìm người dùng theo username.
//  2. Kiểm tra tài khoản còn hoạt động.
//  3. So khớp mật khẩu với hash bcrypt.
//  4. Sinh JWT và cập nhật last_login.
func (s *Service) Login(username, password string) (*LoginResult, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Không tiết lộ username có tồn tại hay không.
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Kiểm tra tài khoản còn hoạt động.
	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	// So khớp mật khẩu với hash đã lưu.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Sinh token.
	token, expiresAt, err := s.tokens.Generate(user.ID, user.Username, user.RoleCode)
	if err != nil {
		return nil, err
	}

	// Cập nhật thời điểm đăng nhập (không làm hỏng luồng nếu lỗi).
	if err := s.repo.UpdateLastLogin(user.ID); err != nil {
		// Ghi nhận nhưng vẫn trả token; đây là thao tác phụ.
		fmt.Printf("cảnh báo: không cập nhật được last_login cho %s: %v\n", user.Username, err)
	}

	return &LoginResult{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User:      toUserInfo(user),
	}, nil
}

// Me trả về thông tin người dùng theo id (dùng cho API /me).
func (s *Service) Me(userID int64) (*UserInfo, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return toUserInfo(user), nil
}

// toUserInfo chuyển models.User sang UserInfo an toàn để trả về client.
func toUserInfo(u *models.User) *UserInfo {
	return &UserInfo{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		RoleCode:  u.RoleCode,
		RoleName:  u.RoleName,
		StationID: u.StationID,
	}
}
