package trucks

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi không tìm thấy xe bồn.
var ErrNotFound = errors.New("không tìm thấy xe bồn")

// ErrCodeExists trả về khi mã xe hoặc biển số đã tồn tại.
var ErrCodeExists = errors.New("mã xe hoặc biển số đã tồn tại")

// ErrDriverNotFound trả về khi tài xế mặc định không tồn tại hoặc đã ngừng hoạt động.
var ErrDriverNotFound = errors.New("tài xế mặc định không tồn tại hoặc đã ngừng hoạt động")

// InactiveCodeError trả về khi mã xe hoặc biển số trùng với một bản ghi đã
// bị xóa mềm. Mang theo bản ghi cũ để tầng handler trả về cho client.
type InactiveCodeError struct {
	Existing *models.MixerTruck
}

// Error triển khai interface error.
func (e *InactiveCodeError) Error() string {
	return "mã xe hoặc biển số đã từng tồn tại ở bản ghi đã ngừng hoạt động"
}

// AsInactiveCodeError trích InactiveCodeError từ chuỗi lỗi (nếu có).
func AsInactiveCodeError(err error) (*InactiveCodeError, bool) {
	var target *InactiveCodeError
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

// selectColumns là danh sách cột chung cho các truy vấn, kèm JOIN tên tài xế mặc định.
const selectColumns = `
	SELECT t.id, t.truck_code, t.plate_number,
	       t.default_driver_id, COALESCE(d.full_name, ''),
	       t.is_active, t.created_at, t.updated_at
	FROM mixer_trucks t
	LEFT JOIN drivers d ON d.id = t.default_driver_id`

// Repository truy vấn bảng mixer_trucks.
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// List trả về danh sách xe bồn đang hoạt động, có phân trang + tìm kiếm.
func (r *Repository) List(p httpx.Pagination) ([]models.MixerTruck, int, error) {
	where := "WHERE t.is_active = 1"
	args := []interface{}{}
	if p.Search != "" {
		where += " AND (t.truck_code LIKE ? OR t.plate_number LIKE ?)"
		like := "%" + p.Search + "%"
		args = append(args, like, like)
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM mixer_trucks t "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("đếm xe bồn thất bại: %w", err)
	}

	query := selectColumns + " " + where + " ORDER BY t.id DESC LIMIT ? OFFSET ?"
	args = append(args, p.Limit(), p.Offset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("truy vấn danh sách xe bồn thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.MixerTruck, 0)
	for rows.Next() {
		t, err := scanTruck(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("đọc danh sách xe bồn thất bại: %w", err)
	}
	return items, total, nil
}

// GetByID lấy xe bồn theo id.
func (r *Repository) GetByID(id int64) (*models.MixerTruck, error) {
	query := selectColumns + " WHERE t.id = ?"
	t, err := scanTruck(r.db.QueryRow(query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

// Create chèn xe bồn mới.
//
// Nếu mã xe hoặc biển số đã tồn tại:
//   - trùng bản ghi ĐANG hoạt động -> ErrCodeExists.
//   - trùng bản ghi ĐÃ xóa mềm     -> *InactiveCodeError kèm bản ghi cũ.
func (r *Repository) Create(input CreateInput) (*models.MixerTruck, error) {
	existing, err := r.getByCodeOrPlate(input.TruckCode, input.PlateNumber)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if existing.IsActive {
			return nil, ErrCodeExists
		}
		return nil, &InactiveCodeError{Existing: existing}
	}
	if err := r.validateDriver(input.DefaultDriverID); err != nil {
		return nil, err
	}

	res, err := r.db.Exec(
		`INSERT INTO mixer_trucks (truck_code, plate_number, default_driver_id)
		 VALUES (?, ?, ?)`,
		input.TruckCode, input.PlateNumber, input.DefaultDriverID,
	)
	if err != nil {
		return nil, fmt.Errorf("tạo xe bồn thất bại: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("không lấy được id xe bồn: %w", err)
	}
	return r.GetByID(id)
}

// Update cập nhật xe bồn theo id.
func (r *Repository) Update(id int64, input UpdateInput) (*models.MixerTruck, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	exists, err := r.codeOrPlateExists(input.TruckCode, input.PlateNumber, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCodeExists
	}
	if err := r.validateDriver(input.DefaultDriverID); err != nil {
		return nil, err
	}

	_, err = r.db.Exec(
		`UPDATE mixer_trucks
		 SET truck_code = ?, plate_number = ?, default_driver_id = ?, updated_at = datetime('now')
		 WHERE id = ?`,
		input.TruckCode, input.PlateNumber, input.DefaultDriverID, id,
	)
	if err != nil {
		return nil, fmt.Errorf("cap nhat xe bon id=%d that bai: %w", id, err)
	}
	return r.GetByID(id)
}

// SoftDelete xóa mềm xe bồn.
func (r *Repository) SoftDelete(id int64) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	_, err := r.db.Exec(
		`UPDATE mixer_trucks SET is_active = 0, updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("xóa mềm xe bồn id=%d thất bại: %w", id, err)
	}
	return nil
}

// Restore khôi phục xe bồn đã xóa mềm: chuyển is_active từ 0 thành 1.
func (r *Repository) Restore(id int64) (*models.MixerTruck, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	_, err := r.db.Exec(
		`UPDATE mixer_trucks SET is_active = 1, updated_at = datetime('now') WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("khôi phục xe bồn id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// codeOrPlateExists kiểm tra trùng mã xe hoặc biển số (bỏ qua excludeID).
func (r *Repository) codeOrPlateExists(truckCode, plateNumber string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM mixer_trucks
		 WHERE (truck_code = ? OR plate_number = ?) AND id <> ?`,
		truckCode, plateNumber, excludeID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("kiểm tra mã/biển số xe thất bại: %w", err)
	}
	return count > 0, nil
}

// getByCodeOrPlate tìm xe bồn theo mã xe hoặc biển số (bất kể trạng thái).
// Trả về (nil, nil) nếu không tìm thấy.
func (r *Repository) getByCodeOrPlate(truckCode, plateNumber string) (*models.MixerTruck, error) {
	query := selectColumns + " WHERE (t.truck_code = ? OR t.plate_number = ?) LIMIT 1"
	t, err := scanTruck(r.db.QueryRow(query, truckCode, plateNumber))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("truy vấn xe bồn theo mã/biển số thất bại: %w", err)
	}
	return t, nil
}

// validateDriver đảm bảo tài xế mặc định (nếu có) tồn tại và đang hoạt động.
// Truyền nil để bỏ trống tài xế.
func (r *Repository) validateDriver(driverID *int64) error {
	if driverID == nil {
		return nil
	}
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM drivers WHERE id = ? AND is_active = 1`,
		*driverID,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("kiểm tra tài xế mặc định thất bại: %w", err)
	}
	if count == 0 {
		return ErrDriverNotFound
	}
	return nil
}

// scanTruck đọc một dòng kết quả thành models.MixerTruck.
func scanTruck(s interface {
	Scan(dest ...interface{}) error
}) (*models.MixerTruck, error) {
	t := &models.MixerTruck{}
	var isActive int
	var driverID sql.NullInt64
	if err := s.Scan(
		&t.ID, &t.TruckCode, &t.PlateNumber,
		&driverID, &t.DefaultDriverName,
		&isActive, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if driverID.Valid {
		id := driverID.Int64
		t.DefaultDriverID = &id
	}
	t.IsActive = isActive == 1
	return t, nil
}
