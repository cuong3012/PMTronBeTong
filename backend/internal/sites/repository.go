package sites

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi không tìm thấy công trường.
var ErrNotFound = errors.New("không tìm thấy công trường")

// ErrCodeExists trả về khi mã công trường đã tồn tại (bản ghi đang hoạt động).
var ErrCodeExists = errors.New("mã công trường đã tồn tại")

// InactiveCodeError trả về khi mã trùng với một bản ghi đã bị xóa mềm.
// Mang theo bản ghi cũ để tầng handler trả về cho client.
type InactiveCodeError struct {
	Existing *models.ConstructionSite
}

// Error triển khai interface error.
func (e *InactiveCodeError) Error() string {
	return "mã công trường đã từng tồn tại ở bản ghi đã ngừng hoạt động"
}

// AsInactiveCodeError trích InactiveCodeError từ chuỗi lỗi (nếu có).
func AsInactiveCodeError(err error) (*InactiveCodeError, bool) {
	var target *InactiveCodeError
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

// Repository truy vấn bảng construction_sites.
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// List trả về danh sách công trường đang hoạt động, có phân trang + tìm kiếm.
func (r *Repository) List(p httpx.Pagination) ([]models.ConstructionSite, int, error) {
	where := "WHERE is_active = 1"
	args := []interface{}{}
	if p.Search != "" {
		where += " AND (site_code LIKE ? OR site_name LIKE ?)"
		like := "%" + p.Search + "%"
		args = append(args, like, like)
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM construction_sites "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("đếm công trường thất bại: %w", err)
	}

	query := `
		SELECT id, site_code, site_name, is_active, created_at, updated_at
		FROM construction_sites ` + where + `
		ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, p.Limit(), p.Offset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("truy vấn danh sách công trường thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.ConstructionSite, 0)
	for rows.Next() {
		var s models.ConstructionSite
		var isActive int
		if err := rows.Scan(&s.ID, &s.SiteCode, &s.SiteName, &isActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("đọc bản ghi công trường thất bại: %w", err)
		}
		s.IsActive = isActive == 1
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("đọc danh sách công trường thất bại: %w", err)
	}
	return items, total, nil
}

// GetByID lấy công trường theo id.
func (r *Repository) GetByID(id int64) (*models.ConstructionSite, error) {
	const query = `
		SELECT id, site_code, site_name, is_active, created_at, updated_at
		FROM construction_sites WHERE id = ?`

	s := &models.ConstructionSite{}
	var isActive int
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.SiteCode, &s.SiteName, &isActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("truy vấn công trường id=%d thất bại: %w", id, err)
	}
	s.IsActive = isActive == 1
	return s, nil
}

// Create chèn công trường mới.
//
// Nếu mã đã tồn tại:
//   - trùng bản ghi ĐANG hoạt động -> ErrCodeExists.
//   - trùng bản ghi ĐÃ xóa mềm     -> *InactiveCodeError kèm bản ghi cũ.
func (r *Repository) Create(input CreateInput) (*models.ConstructionSite, error) {
	existing, err := r.getByCode(input.SiteCode)
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
		`INSERT INTO construction_sites (site_code, site_name) VALUES (?, ?)`,
		input.SiteCode, input.SiteName,
	)
	if err != nil {
		return nil, fmt.Errorf("tạo công trường thất bại: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("không lấy được id công trường: %w", err)
	}
	return r.GetByID(id)
}

// Update cập nhật công trường theo id.
func (r *Repository) Update(id int64, input UpdateInput) (*models.ConstructionSite, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	exists, err := r.codeExists(input.SiteCode, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCodeExists
	}
	_, err = r.db.Exec(
		`UPDATE construction_sites SET site_code = ?, site_name = ?, updated_at = datetime('now') WHERE id = ?`,
		input.SiteCode, input.SiteName, id,
	)
	if err != nil {
		return nil, fmt.Errorf("cập nhật công trường id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// SoftDelete xóa mềm công trường.
func (r *Repository) SoftDelete(id int64) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	_, err := r.db.Exec(
		`UPDATE construction_sites SET is_active = 0, updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("xóa mềm công trường id=%d thất bại: %w", id, err)
	}
	return nil
}

// Restore khôi phục công trường đã xóa mềm: chuyển is_active từ 0 thành 1.
func (r *Repository) Restore(id int64) (*models.ConstructionSite, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	_, err := r.db.Exec(
		`UPDATE construction_sites SET is_active = 1, updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("khôi phục công trường id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// codeExists kiểm tra trùng mã công trường (bỏ qua excludeID).
func (r *Repository) codeExists(code string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM construction_sites WHERE site_code = ? AND id <> ?`,
		code, excludeID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("kiểm tra mã công trường thất bại: %w", err)
	}
	return count > 0, nil
}

// getByCode tìm công trường theo mã (bất kể trạng thái). (nil, nil) nếu không có.
func (r *Repository) getByCode(code string) (*models.ConstructionSite, error) {
	const query = `
		SELECT id, site_code, site_name, is_active, created_at, updated_at
		FROM construction_sites WHERE site_code = ? LIMIT 1`
	s := &models.ConstructionSite{}
	var isActive int
	err := r.db.QueryRow(query, code).Scan(
		&s.ID, &s.SiteCode, &s.SiteName, &isActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("truy vấn công trường theo mã thất bại: %w", err)
	}
	s.IsActive = isActive == 1
	return s, nil
}
