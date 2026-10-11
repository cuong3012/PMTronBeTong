package grades

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi không tìm thấy mác bê tông.
var ErrNotFound = errors.New("không tìm thấy mác bê tông")

// ErrCodeExists trả về khi mã mác đã tồn tại ở bản ghi đang hoạt động.
var ErrCodeExists = errors.New("mã mác bê tông đã tồn tại")

const selectColumns = `
	SELECT id, grade_code, printed_name,
	       tp1, tp2, tp3, tp4, tp5, tp6,
	       xi1, xi2, xi3, xi4, xi5, nuoc,
	       pg1, pg2, pg3, pg4, pg5, pg6,
	       do_sut, cuong_do, COALESCE(notes, ''),
	       is_active, created_at, updated_at
	FROM concrete_grades`

// Repository truy vấn bảng concrete_grades.
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// List trả về danh sách mác đang hoạt động, phân trang + tìm kiếm.
func (r *Repository) List(p httpx.Pagination) ([]models.ConcreteGrade, int, error) {
	where := "WHERE is_active = 1"
	args := []interface{}{}
	if p.Search != "" {
		where += " AND (grade_code LIKE ? OR printed_name LIKE ?)"
		like := "%" + p.Search + "%"
		args = append(args, like, like)
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM concrete_grades "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("đếm mác bê tông thất bại: %w", err)
	}

	query := selectColumns + " " + where + " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, p.Limit(), p.Offset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("truy vấn danh sách mác bê tông thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.ConcreteGrade, 0)
	for rows.Next() {
		g, err := scanGrade(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *g)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("đọc danh sách mác bê tông thất bại: %w", err)
	}
	return items, total, nil
}

// GetByID lấy mác theo id.
func (r *Repository) GetByID(id int64) (*models.ConcreteGrade, error) {
	g, err := scanGrade(r.db.QueryRow(selectColumns+" WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Create chèn mác mới.
func (r *Repository) Create(input GradeInput) (*models.ConcreteGrade, error) {
	exists, err := r.codeExists(input.GradeCode, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCodeExists
	}

	res, err := r.db.Exec(`
		INSERT INTO concrete_grades
			(grade_code, printed_name,
			 tp1, tp2, tp3, tp4, tp5, tp6,
			 xi1, xi2, xi3, xi4, xi5, nuoc,
			 pg1, pg2, pg3, pg4, pg5, pg6,
			 do_sut, cuong_do, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.GradeCode, input.PrintedName,
		input.TP1, input.TP2, input.TP3, input.TP4, input.TP5, input.TP6,
		input.XI1, input.XI2, input.XI3, input.XI4, input.XI5, input.Nuoc,
		input.PG1, input.PG2, input.PG3, input.PG4, input.PG5, input.PG6,
		input.DoSut, input.CuongDo, input.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("tạo mác bê tông thất bại: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("không lấy được id mác bê tông: %w", err)
	}
	return r.GetByID(id)
}

// Update cập nhật mác theo id.
func (r *Repository) Update(id int64, input GradeInput) (*models.ConcreteGrade, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, err
	}
	exists, err := r.codeExists(input.GradeCode, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCodeExists
	}

	_, err = r.db.Exec(`
		UPDATE concrete_grades SET
			grade_code = ?, printed_name = ?,
			tp1 = ?, tp2 = ?, tp3 = ?, tp4 = ?, tp5 = ?, tp6 = ?,
			xi1 = ?, xi2 = ?, xi3 = ?, xi4 = ?, xi5 = ?, nuoc = ?,
			pg1 = ?, pg2 = ?, pg3 = ?, pg4 = ?, pg5 = ?, pg6 = ?,
			do_sut = ?, cuong_do = ?, notes = ?, updated_at = datetime('now')
		WHERE id = ?`,
		input.GradeCode, input.PrintedName,
		input.TP1, input.TP2, input.TP3, input.TP4, input.TP5, input.TP6,
		input.XI1, input.XI2, input.XI3, input.XI4, input.XI5, input.Nuoc,
		input.PG1, input.PG2, input.PG3, input.PG4, input.PG5, input.PG6,
		input.DoSut, input.CuongDo, input.Notes, id,
	)
	if err != nil {
		return nil, fmt.Errorf("cập nhật mác bê tông id=%d thất bại: %w", id, err)
	}
	return r.GetByID(id)
}

// SoftDelete xóa mềm mác.
func (r *Repository) SoftDelete(id int64) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	_, err := r.db.Exec(
		`UPDATE concrete_grades SET is_active = 0, updated_at = datetime('now') WHERE id = ?`, id,
	)
	if err != nil {
		return fmt.Errorf("xóa mềm mác bê tông id=%d thất bại: %w", id, err)
	}
	return nil
}

func (r *Repository) codeExists(code string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM concrete_grades WHERE grade_code = ? AND id <> ?`,
		code, excludeID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("kiểm tra mã mác bê tông thất bại: %w", err)
	}
	return count > 0, nil
}

func scanGrade(s interface{ Scan(dest ...interface{}) error }) (*models.ConcreteGrade, error) {
	g := &models.ConcreteGrade{}
	var isActive int
	if err := s.Scan(
		&g.ID, &g.GradeCode, &g.PrintedName,
		&g.TP1, &g.TP2, &g.TP3, &g.TP4, &g.TP5, &g.TP6,
		&g.XI1, &g.XI2, &g.XI3, &g.XI4, &g.XI5, &g.Nuoc,
		&g.PG1, &g.PG2, &g.PG3, &g.PG4, &g.PG5, &g.PG6,
		&g.DoSut, &g.CuongDo, &g.Notes,
		&isActive, &g.CreatedAt, &g.UpdatedAt,
	); err != nil {
		return nil, err
	}
	g.IsActive = isActive == 1
	return g, nil
}
