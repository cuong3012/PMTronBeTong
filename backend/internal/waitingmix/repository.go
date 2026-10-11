package waitingmix

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi không tìm thấy vị trí cấp phối chờ.
var ErrNotFound = errors.New("không tìm thấy vị trí cấp phối chờ")

const selectColumns = `
	SELECT w.id, w.stt, w.concrete_grade_id, COALESCE(g.grade_code, ''),
	       COALESCE(g.printed_name, ''),
	       w.customer_id, COALESCE(c.customer_name, ''),
	       w.volume_m3, w.m3me, w.status, w.created_at, w.updated_at
	FROM waiting_mixes w
	LEFT JOIN concrete_grades g ON g.id = w.concrete_grade_id
	LEFT JOIN customers c ON c.id = w.customer_id`

// Repository truy vấn bảng waiting_mixes (20 vị trí cố định).
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// List trả về đủ 20 vị trí, sắp xếp theo stt.
func (r *Repository) List() ([]models.WaitingMix, error) {
	rows, err := r.db.Query(selectColumns + " ORDER BY w.stt")
	if err != nil {
		return nil, fmt.Errorf("truy vấn danh sách cấp phối chờ thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.WaitingMix, 0)
	for rows.Next() {
		w, err := scanWaitingMix(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *w)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("đọc danh sách cấp phối chờ thất bại: %w", err)
	}
	return items, nil
}

// GetByID lấy một vị trí theo id.
func (r *Repository) GetByID(id int64) (*models.WaitingMix, error) {
	w, err := scanWaitingMix(r.db.QueryRow(selectColumns+" WHERE w.id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return w, nil
}

// Update cập nhật vị trí cấp phối chờ.
func (r *Repository) Update(id int64, input UpdateInput) (*models.WaitingMix, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	_, err := r.db.Exec(
		`UPDATE waiting_mixes
		 SET concrete_grade_id = ?, customer_id = ?, volume_m3 = ?, m3me = ?, status = ?, updated_at = datetime('now')
		 WHERE id = ?`,
		input.ConcreteGradeID, input.CustomerID, input.VolumeM3, input.M3Me, input.Status, id,
	)
	if err != nil {
		return nil, fmt.Errorf("cập nhật cấp phối chờ id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// GradeActive kiểm tra mác bê tông đang hoạt động.
func (r *Repository) GradeActive(id int64) (bool, error) {
	return r.existsActive("concrete_grades", id)
}

// CustomerActive kiểm tra khách hàng đang hoạt động.
func (r *Repository) CustomerActive(id int64) (bool, error) {
	return r.existsActive("customers", id)
}

func (r *Repository) existsActive(table string, id int64) (bool, error) {
	var count int
	// Bảng lấy từ hằng nội bộ, không phải input người dùng -> an toàn.
	err := r.db.QueryRow(
		fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ? AND is_active = 1", table),
		id,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("kiểm tra tham chiếu %s id=%d thất bại: %w", table, id, err)
	}
	return count > 0, nil
}

func scanWaitingMix(s interface{ Scan(dest ...interface{}) error }) (*models.WaitingMix, error) {
	w := &models.WaitingMix{}
	var gradeID, customerID sql.NullInt64
	if err := s.Scan(
		&w.ID, &w.STT, &gradeID, &w.ConcreteGradeCode, &w.ConcreteGradeName,
		&customerID, &w.CustomerName,
		&w.VolumeM3, &w.M3Me, &w.Status, &w.CreatedAt, &w.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if gradeID.Valid {
		id := gradeID.Int64
		w.ConcreteGradeID = &id
	}
	if customerID.Valid {
		id := customerID.Int64
		w.CustomerID = &id
	}
	return w, nil
}
