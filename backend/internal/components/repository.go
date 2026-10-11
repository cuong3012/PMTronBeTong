package components

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi không tìm thấy thành phần.
var ErrNotFound = errors.New("không tìm thấy thành phần")

// ErrSlotExists trả về khi slot_name đã tồn tại.
var ErrSlotExists = errors.New("tên vị trí thành phần đã tồn tại")

const selectColumns = `
	SELECT id, slot_name, humidity_values, COALESCE(unit, ''), COALESCE(notes, ''),
	       created_at, updated_at
	FROM component_catalog`

// Repository truy vấn bảng component_catalog.
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// List trả về danh sách thành phần (không phân trang, số lượng nhỏ).
func (r *Repository) List(search string) ([]models.ComponentCatalog, error) {
	query := selectColumns
	args := []interface{}{}
	if search != "" {
		query += " WHERE slot_name LIKE ?"
		args = append(args, "%"+search+"%")
	}
	query += " ORDER BY id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("truy vấn danh sách thành phần thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.ComponentCatalog, 0)
	for rows.Next() {
		c, err := scanComponent(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("đọc danh sách thành phần thất bại: %w", err)
	}
	return items, nil
}

// GetByID lấy thành phần theo id.
func (r *Repository) GetByID(id int64) (*models.ComponentCatalog, error) {
	c, err := scanComponent(r.db.QueryRow(selectColumns+" WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Create chèn thành phần mới.
func (r *Repository) Create(input CatalogInput) (*models.ComponentCatalog, error) {
	exists, err := r.slotExists(input.SlotName, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlotExists
	}
	res, err := r.db.Exec(
		`INSERT INTO component_catalog (slot_name, humidity_values, unit, notes)
		 VALUES (?, ?, ?, ?)`,
		input.SlotName, input.HumidityValues, input.Unit, input.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("tạo thành phần thất bại: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("không lấy được id thành phần: %w", err)
	}
	return r.GetByID(id)
}

// Update cập nhật thành phần theo id.
func (r *Repository) Update(id int64, input CatalogInput) (*models.ComponentCatalog, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	exists, err := r.slotExists(input.SlotName, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlotExists
	}
	_, err = r.db.Exec(
		`UPDATE component_catalog
		 SET slot_name = ?, humidity_values = ?, unit = ?, notes = ?, updated_at = datetime('now')
		 WHERE id = ?`,
		input.SlotName, input.HumidityValues, input.Unit, input.Notes, id,
	)
	if err != nil {
		return nil, fmt.Errorf("cập nhật thành phần id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// Delete xóa thành phần theo id (xóa cứng — bảng metadata).
func (r *Repository) Delete(id int64) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	_, err := r.db.Exec(`DELETE FROM component_catalog WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("xóa thành phần id=%d thất bại: %w", id, err)
	}
	return nil
}

func (r *Repository) slotExists(slot string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM component_catalog WHERE slot_name = ? AND id <> ?`,
		slot, excludeID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("kiểm tra tên vị trí thành phần thất bại: %w", err)
	}
	return count > 0, nil
}

func scanComponent(s interface{ Scan(dest ...interface{}) error }) (*models.ComponentCatalog, error) {
	c := &models.ComponentCatalog{}
	if err := s.Scan(
		&c.ID, &c.SlotName, &c.HumidityValues, &c.Unit, &c.Notes,
		&c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return c, nil
}
