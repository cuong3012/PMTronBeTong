package drivers

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi không tìm thấy tài xế.
var ErrNotFound = errors.New("không tìm thấy tài xế")

// ErrCodeExists trả về khi mã tài xế đã tồn tại (bản ghi đang hoạt động).
var ErrCodeExists = errors.New("mã tài xế đã tồn tại")

// InactiveCodeError trả về khi mã trùng với một bản ghi đã bị xóa mềm.
// Mang theo bản ghi cũ để tầng handler trả về cho client.
type InactiveCodeError struct {
	Existing *models.Driver
}

// Error triển khai interface error.
func (e *InactiveCodeError) Error() string {
	return "mã tài xế đã từng tồn tại ở bản ghi đã ngừng hoạt động"
}

// AsInactiveCodeError trích InactiveCodeError từ chuỗi lỗi (nếu có).
func AsInactiveCodeError(err error) (*InactiveCodeError, bool) {
	var target *InactiveCodeError
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

// Repository truy vấn bảng drivers.
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// List trả về danh sách tài xế đang hoạt động, có phân trang + tìm kiếm.
func (r *Repository) List(p httpx.Pagination) ([]models.Driver, int, error) {
	where := "WHERE is_active = 1"
	args := []interface{}{}
	if p.Search != "" {
		where += " AND (driver_code LIKE ? OR full_name LIKE ?)"
		like := "%" + p.Search + "%"
		args = append(args, like, like)
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM drivers "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("đếm tài xế thất bại: %w", err)
	}

	query := `
		SELECT id, driver_code, full_name, is_active, created_at, updated_at
		FROM drivers ` + where + `
		ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, p.Limit(), p.Offset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("truy vấn danh sách tài xế thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.Driver, 0)
	for rows.Next() {
		var d models.Driver
		var isActive int
		if err := rows.Scan(&d.ID, &d.DriverCode, &d.FullName, &isActive, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("đọc bản ghi tài xế thất bại: %w", err)
		}
		d.IsActive = isActive == 1
		items = append(items, d)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("đọc danh sách tài xế thất bại: %w", err)
	}
	return items, total, nil
}

// GetByID lấy tài xế theo id.
func (r *Repository) GetByID(id int64) (*models.Driver, error) {
	const query = `
		SELECT id, driver_code, full_name, is_active, created_at, updated_at
		FROM drivers WHERE id = ?`

	d := &models.Driver{}
	var isActive int
	err := r.db.QueryRow(query, id).Scan(
		&d.ID, &d.DriverCode, &d.FullName, &isActive, &d.CreatedAt, &d.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("truy vấn tài xế id=%d thất bại: %w", id, err)
	}
	d.IsActive = isActive == 1
	return d, nil
}

// Create chèn tài xế mới.
//
// Nếu mã đã tồn tại:
//   - trùng bản ghi ĐANG hoạt động -> ErrCodeExists.
//   - trùng bản ghi ĐÃ xóa mềm     -> *InactiveCodeError kèm bản ghi cũ.
func (r *Repository) Create(input CreateInput) (*models.Driver, error) {
	existing, err := r.getByCode(input.DriverCode)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if existing.IsActive {
			return nil, ErrCodeExists
		}
		return nil, &InactiveCodeError{Existing: existing}
	}

	res, err := r.db.Exec(
		`INSERT INTO drivers (driver_code, full_name) VALUES (?, ?)`,
		input.DriverCode, input.FullName,
	)
	if err != nil {
		return nil, fmt.Errorf("tạo tài xế thất bại: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("không lấy được id tài xế: %w", err)
	}
	return r.GetByID(id)
}

// Update cập nhật tài xế theo id.
func (r *Repository) Update(id int64, input UpdateInput) (*models.Driver, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	exists, err := r.codeExists(input.DriverCode, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCodeExists
	}
	_, err = r.db.Exec(
		`UPDATE drivers SET driver_code = ?, full_name = ?, updated_at = datetime('now') WHERE id = ?`,
		input.DriverCode, input.FullName, id,
	)
	if err != nil {
		return nil, fmt.Errorf("cap nhat tai xe id=%d that bai: %w", id, err)
	}
	return r.GetByID(id)
}

// SoftDelete xóa mềm tài xế.
func (r *Repository) SoftDelete(id int64) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	_, err := r.db.Exec(
		`UPDATE drivers SET is_active = 0, updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("xóa mềm tài xế id=%d thất bại: %w", id, err)
	}
	return nil
}

// Restore khôi phục tài xế đã xóa mềm: chuyển is_active từ 0 thành 1.
func (r *Repository) Restore(id int64) (*models.Driver, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	_, err := r.db.Exec(
		`UPDATE drivers SET is_active = 1, updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("khôi phục tài xế id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// codeExists kiểm tra trùng mã tài xế (bỏ qua excludeID).
func (r *Repository) codeExists(code string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM drivers WHERE driver_code = ? AND id <> ?`,
		code, excludeID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("kiểm tra mã tài xế thất bại: %w", err)
	}
	return count > 0, nil
}

// getByCode tìm tài xế theo mã (bất kể trạng thái). (nil, nil) nếu không có.
func (r *Repository) getByCode(code string) (*models.Driver, error) {
	const query = `
		SELECT id, driver_code, full_name, is_active, created_at, updated_at
		FROM drivers WHERE driver_code = ? LIMIT 1`
	d := &models.Driver{}
	var isActive int
	err := r.db.QueryRow(query, code).Scan(
		&d.ID, &d.DriverCode, &d.FullName, &isActive, &d.CreatedAt, &d.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("truy vấn tài xế theo mã thất bại: %w", err)
	}
	d.IsActive = isActive == 1
	return d, nil
}
