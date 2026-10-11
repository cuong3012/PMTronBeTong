package models

// Role ánh xạ bảng roles (mục 13.1).
type Role struct {
	ID          int64  `json:"id"`
	RoleCode    string `json:"role_code"`
	RoleName    string `json:"role_name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// User ánh xạ bảng users (mục 13.1).
//
// LƯU Ý: trường PasswordHash có json:"-" để KHÔNG bao giờ trả về client.
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone,omitempty"`
	RoleID       int64  `json:"role_id"`
	StationID    *int64 `json:"station_id,omitempty"` // nil nếu là tài khoản toàn hệ thống
	IsActive     bool   `json:"is_active"`
	LastLoginAt  string `json:"last_login_at,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`

	// Các trường tiện ích khi JOIN với bảng roles (không lưu trong bảng users).
	RoleCode string `json:"role_code,omitempty"`
	RoleName string `json:"role_name,omitempty"`
}
