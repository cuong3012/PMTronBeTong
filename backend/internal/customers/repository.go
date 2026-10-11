package customers

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi không tìm thấy khách hàng.
var ErrNotFound = errors.New("không tìm thấy khách hàng")

// ErrCodeExists trả về khi mã khách hàng đã tồn tại (bản ghi đang hoạt động).
var ErrCodeExists = errors.New("mã khách hàng đã tồn tại")

// InactiveCodeError trả về khi mã trùng với một bản ghi đã bị xóa mềm.
// Mang theo bản ghi cũ để tầng handler trả về cho client, giúp client hỏi
// người dùng có muốn khôi phục bản ghi cũ hay không.
type InactiveCodeError struct {
	Existing *models.Customer
}

// Error triển khai interface error.
func (e *InactiveCodeError) Error() string {
	return "mã khách hàng đã từng tồn tại ở bản ghi đã ngừng hoạt động"
}

// AsInactiveCodeError trích InactiveCodeError từ chuỗi lỗi (nếu có).
func AsInactiveCodeError(err error) (*InactiveCodeError, bool) {
	var target *InactiveCodeError
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

// Repository truy vấn bảng customers.
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// List trả về danh sách khách hàng đang hoạt động (is_active = 1) có phân trang.
// Hỗ trợ tìm kiếm theo mã hoặc tên (không phân biệt hoa/thường).
func (r *Repository) List(p httpx.Pagination) ([]models.Customer, int, error) {
	// Điều kiện tìm kiếm: mã hoặc tên chứa từ khóa.
	where := "WHERE is_active = 1"
	args := []interface{}{}
	if p.Search != "" {
		where += " AND (customer_code LIKE ? OR customer_name LIKE ?)"
		like := "%" + p.Search + "%"
		args = append(args, like, like)
	}

	// Đếm tổng số bản ghi thỏa điều kiện.
	var total int
	countQuery := "SELECT COUNT(*) FROM customers " + where
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("đếm khách hàng thất bại: %w", err)
	}

	// Lấy dữ liệu trang hiện tại.
	listQuery := `
		SELECT id, customer_code, customer_name,
		       COALESCE(address, ''), COALESCE(phone, ''),
		       is_active, created_at, updated_at
		FROM customers ` + where + `
		ORDER BY id DESC
		LIMIT ? OFFSET ?`
	args = append(args, p.Limit(), p.Offset())

	rows, err := r.db.Query(listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("truy vấn danh sách khách hàng thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.Customer, 0)
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("đọc danh sách khách hàng thất bại: %w", err)
	}

	return items, total, nil
}

// GetByID lấy một khách hàng theo id.
func (r *Repository) GetByID(id int64) (*models.Customer, error) {
	const query = `
		SELECT id, customer_code, customer_name,
		       COALESCE(address, ''), COALESCE(phone, ''),
		       is_active, created_at, updated_at
		FROM customers WHERE id = ?`

	c := &models.Customer{}
	var isActive int
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.CustomerCode, &c.CustomerName,
		&c.Address, &c.Phone, &isActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("truy vấn khách hàng id=%d thất bại: %w", id, err)
	}
	c.IsActive = isActive == 1
	return c, nil
}

// Create chèn khách hàng mới và trả về bản ghi vừa tạo.
//
// Nếu mã đã tồn tại:
//   - trùng với bản ghi ĐANG hoạt động  -> ErrCodeExists (không thể tạo).
//   - trùng với bản ghi ĐÃ xóa mềm      -> *InactiveCodeError kèm bản ghi cũ,
//     để client hỏi người dùng có muốn khôi phục bản ghi cũ hay không.
func (r *Repository) Create(input CreateInput) (*models.Customer, error) {
	// Kiểm tra trùng mã (cả bản ghi đã xóa mềm) trước khi chèn.
	existing, err := r.getByCode(input.CustomerCode)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if existing.IsActive {
			return nil, ErrCodeExists
		}
		// Trùng với bản ghi đã xóa mềm: trả lỗi đặc biệt kèm bản ghi cũ.
		return nil, &InactiveCodeError{Existing: existing}
	}

	res, err := r.db.Exec(
		`INSERT INTO customers (customer_code, customer_name, address, phone)
		 VALUES (?, ?, ?, ?)`,
		input.CustomerCode, input.CustomerName, input.Address, input.Phone,
	)
	if err != nil {
		return nil, fmt.Errorf("tạo khách hàng thất bại: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("không lấy được id khách hàng vừa tạo: %w", err)
	}
	return r.GetByID(id)
}

// Update cập nhật thông tin khách hàng theo id.
func (r *Repository) Update(id int64, input UpdateInput) (*models.Customer, error) {
	// Đảm bảo bản ghi tồn tại.
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}

	// Nếu đổi mã, kiểm tra trùng với bản ghi khác.
	exists, err := r.codeExists(input.CustomerCode, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCodeExists
	}

	_, err = r.db.Exec(
		`UPDATE customers
		 SET customer_code = ?, customer_name = ?, address = ?, phone = ?,
		     updated_at = datetime('now')
		 WHERE id = ?`,
		input.CustomerCode, input.CustomerName, input.Address, input.Phone, id,
	)
	if err != nil {
		return nil, fmt.Errorf("cập nhật khách hàng id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// SoftDelete xóa mềm khách hàng: chuyển is_active từ 1 thành 0.
func (r *Repository) SoftDelete(id int64) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	_, err := r.db.Exec(
		`UPDATE customers SET is_active = 0, updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("xóa mềm khách hàng id=%d thất bại: %w", id, err)
	}
	return nil
}

// Restore khôi phục khách hàng đã xóa mềm: chuyển is_active từ 0 thành 1.
// Giữ nguyên thông tin cũ của bản ghi.
func (r *Repository) Restore(id int64) (*models.Customer, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	_, err := r.db.Exec(
		`UPDATE customers SET is_active = 1, updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("khôi phục khách hàng id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// codeExists kiểm tra mã khách hàng đã tồn tại hay chưa.
// excludeID dùng để bỏ qua chính bản ghi đang cập nhật (truyền 0 khi tạo mới).
func (r *Repository) codeExists(code string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM customers WHERE customer_code = ? AND id <> ?`,
		code, excludeID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("kiểm tra mã khách hàng thất bại: %w", err)
	}
	return count > 0, nil
}

// getByCode tìm khách hàng theo mã (không phân biệt đang hoạt động hay không).
// Trả về (nil, nil) nếu không tìm thấy.
func (r *Repository) getByCode(code string) (*models.Customer, error) {
	const query = `
		SELECT id, customer_code, customer_name,
		       COALESCE(address, ''), COALESCE(phone, ''),
		       is_active, created_at, updated_at
		FROM customers WHERE customer_code = ? LIMIT 1`
	c, err := scanCustomer(r.db.QueryRow(query, code))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("truy vấn khách hàng theo mã thất bại: %w", err)
	}
	return c, nil
}

// scanCustomer đọc một dòng kết quả thành models.Customer.
// Nhận interface có hàm Scan để dùng chung cho QueryRow và Rows.
func scanCustomer(s interface {
	Scan(dest ...interface{}) error
}) (*models.Customer, error) {
	c := &models.Customer{}
	var isActive int
	err := s.Scan(
		&c.ID, &c.CustomerCode, &c.CustomerName,
		&c.Address, &c.Phone, &isActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("đọc bản ghi khách hàng thất bại: %w", err)
	}
	c.IsActive = isActive == 1
	return c, nil
}
