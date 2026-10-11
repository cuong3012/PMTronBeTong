package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/models"
)

// ErrUserNotFound trả về khi không tìm thấy người dùng theo username.
var ErrUserNotFound = errors.New("không tìm thấy người dùng")

// Repository truy vấn dữ liệu người dùng từ database.
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới với kết nối database cho trước.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// GetByUsername tìm người dùng theo username (JOIN với roles để lấy mã/tên vai trò).
// Trả về ErrUserNotFound nếu không có bản ghi nào.
func (r *Repository) GetByUsername(username string) (*models.User, error) {
	const query = `
		SELECT u.id, u.username, u.password_hash, u.full_name,
		       COALESCE(u.phone, ''), u.role_id, u.station_id, u.is_active,
		       COALESCE(u.last_login_at, ''), u.created_at, u.updated_at,
		       r.role_code, r.role_name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.username = ?`

	u := &models.User{}
	var isActive int
	var stationID sql.NullInt64

	err := r.db.QueryRow(query, username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.FullName,
		&u.Phone, &u.RoleID, &stationID, &isActive,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
		&u.RoleCode, &u.RoleName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("truy vấn người dùng %q thất bại: %w", username, err)
	}

	u.IsActive = isActive == 1
	if stationID.Valid {
		id := stationID.Int64
		u.StationID = &id
	}
	return u, nil
}

// GetByID tìm người dùng theo id (dùng cho middleware khi cần nạp lại thông tin).
func (r *Repository) GetByID(id int64) (*models.User, error) {
	const query = `
		SELECT u.id, u.username, u.password_hash, u.full_name,
		       COALESCE(u.phone, ''), u.role_id, u.station_id, u.is_active,
		       COALESCE(u.last_login_at, ''), u.created_at, u.updated_at,
		       r.role_code, r.role_name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = ?`

	u := &models.User{}
	var isActive int
	var stationID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.FullName,
		&u.Phone, &u.RoleID, &stationID, &isActive,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
		&u.RoleCode, &u.RoleName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("truy vấn người dùng id=%d thất bại: %w", id, err)
	}

	u.IsActive = isActive == 1
	if stationID.Valid {
		sid := stationID.Int64
		u.StationID = &sid
	}
	return u, nil
}

// UpdateLastLogin cập nhật thời điểm đăng nhập gần nhất cho người dùng.
func (r *Repository) UpdateLastLogin(id int64) error {
	_, err := r.db.Exec(
		`UPDATE users SET last_login_at = datetime('now'), updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("cập nhật last_login cho user id=%d thất bại: %w", id, err)
	}
	return nil
}
