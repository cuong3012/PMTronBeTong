package weighing

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi không tìm thấy phiếu cân.
var ErrNotFound = errors.New("không tìm thấy phiếu cân")

// ErrCodeExists trả về khi mã phiếu cân đã tồn tại.
var ErrCodeExists = errors.New("mã phiếu cân đã tồn tại")

const ticketColumns = `
	SELECT w.id, w.ticket_code, w.volume_m3, w.batch_count, w.m3me,
	       w.order_id, COALESCE(o.order_number, ''),
	       w.concrete_grade_id, w.concrete_grade_code, w.concrete_grade_name,
	       w.printed_concrete_grade_name, w.humidity_values,
	       w.standard_per_cube, w.standard_per_batch,
	       w.customer_name, w.construction_site_address, w.material_name,
	       w.runner_user_id, COALESCE(w.notes, ''), w.is_active,
	       w.created_at, w.updated_at
	FROM weighing_tickets w
	LEFT JOIN orders o ON o.id = w.order_id`

const detailColumns = `
	SELECT id, weighing_ticket_id, batch_no, batch_volume_m3, mixed_at,
	       tp1, tp2, tp3, tp4, tp5, tp6,
	       xi1, xi2, xi3, xi4, xi5,
	       pg1, pg2, pg3, pg4, pg5, pg6,
	       COALESCE(manual_override_json, ''), COALESCE(notes, ''),
	       is_active, created_at, updated_at
	FROM weighing_ticket_details`

// Repository truy vấn phiếu cân và chi tiết.
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// List trả về danh sách phiếu cân đang hoạt động, phân trang + tìm kiếm.
func (r *Repository) List(p httpx.Pagination) ([]models.WeighingTicket, int, error) {
	where := "WHERE w.is_active = 1"
	args := []interface{}{}
	if p.Search != "" {
		where += " AND (w.ticket_code LIKE ? OR w.customer_name LIKE ? OR w.concrete_grade_code LIKE ?)"
		like := "%" + p.Search + "%"
		args = append(args, like, like, like)
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM weighing_tickets w "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("đếm phiếu cân thất bại: %w", err)
	}

	query := ticketColumns + " " + where + " ORDER BY w.id DESC LIMIT ? OFFSET ?"
	args = append(args, p.Limit(), p.Offset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("truy vấn danh sách phiếu cân thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.WeighingTicket, 0)
	for rows.Next() {
		t, err := scanTicket(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("đọc danh sách phiếu cân thất bại: %w", err)
	}
	return items, total, nil
}

// GetByID lấy phiếu cân kèm danh sách mẻ chi tiết.
func (r *Repository) GetByID(id int64) (*models.WeighingTicket, error) {
	t, err := scanTicket(r.db.QueryRow(ticketColumns+" WHERE w.id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	details, err := r.listDetails(id)
	if err != nil {
		return nil, err
	}
	t.Details = details
	return t, nil
}

// listDetails lấy các mẻ của một phiếu cân.
func (r *Repository) listDetails(ticketID int64) ([]models.WeighingTicketDetail, error) {
	rows, err := r.db.Query(detailColumns+" WHERE weighing_ticket_id = ? ORDER BY batch_no", ticketID)
	if err != nil {
		return nil, fmt.Errorf("truy vấn chi tiết phiếu cân thất bại: %w", err)
	}
	defer rows.Close()

	items := make([]models.WeighingTicketDetail, 0)
	for rows.Next() {
		d, err := scanDetail(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("đọc chi tiết phiếu cân thất bại: %w", err)
	}
	return items, nil
}

// CreateTicket chèn phiếu cân và toàn bộ mẻ chi tiết trong MỘT transaction.
func (r *Repository) CreateTicket(input TicketInput) (*models.WeighingTicket, error) {
	exists, err := r.codeExists(input.TicketCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCodeExists
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("bắt đầu transaction tạo phiếu cân thất bại: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(`
		INSERT INTO weighing_tickets
			(ticket_code, volume_m3, batch_count, m3me, order_id, concrete_grade_id,
			 concrete_grade_code, concrete_grade_name, printed_concrete_grade_name,
			 humidity_values, standard_per_cube, standard_per_batch,
			 customer_name, construction_site_address, material_name, runner_user_id, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.TicketCode, input.VolumeM3, input.BatchCount, input.M3Me, input.OrderID,
		input.ConcreteGradeID, input.ConcreteGradeCode, input.ConcreteGradeName,
		input.PrintedConcreteGradeName, input.HumidityValues,
		input.StandardPerCube, input.StandardPerBatch,
		input.CustomerName, input.ConstructionSiteAddress, input.MaterialName,
		input.RunnerUserID, input.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("tạo phiếu cân thất bại: %w", err)
	}
	ticketID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("không lấy được id phiếu cân: %w", err)
	}

	for _, d := range input.Details {
		_, err := tx.Exec(`
			INSERT INTO weighing_ticket_details
				(weighing_ticket_id, batch_no, batch_volume_m3, mixed_at,
				 tp1, tp2, tp3, tp4, tp5, tp6,
				 xi1, xi2, xi3, xi4, xi5,
				 pg1, pg2, pg3, pg4, pg5, pg6,
				 manual_override_json, notes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			ticketID, d.BatchNo, d.BatchVolumeM3, d.MixedAt,
			d.TP1, d.TP2, d.TP3, d.TP4, d.TP5, d.TP6,
			d.XI1, d.XI2, d.XI3, d.XI4, d.XI5,
			d.PG1, d.PG2, d.PG3, d.PG4, d.PG5, d.PG6,
			nullIfEmpty(d.ManualOverrideJSON), d.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("tạo chi tiết mẻ %d thất bại: %w", d.BatchNo, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tạo phiếu cân thất bại: %w", err)
	}
	return r.GetByID(ticketID)
}

// SoftDeleteTicket xóa mềm phiếu cân và các mẻ chi tiết.
func (r *Repository) SoftDeleteTicket(id int64) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("bắt đầu transaction xóa phiếu cân thất bại: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(
		`UPDATE weighing_tickets SET is_active = 0, updated_at = datetime('now') WHERE id = ?`, id,
	); err != nil {
		return fmt.Errorf("xóa mềm phiếu cân id=%d thất bại: %w", id, err)
	}
	if _, err := tx.Exec(
		`UPDATE weighing_ticket_details SET is_active = 0, updated_at = datetime('now') WHERE weighing_ticket_id = ?`, id,
	); err != nil {
		return fmt.Errorf("xóa mềm chi tiết phiếu cân id=%d thất bại: %w", id, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit xóa phiếu cân thất bại: %w", err)
	}
	return nil
}

func (r *Repository) codeExists(code string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM weighing_tickets WHERE ticket_code = ?`, code,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("kiểm tra mã phiếu cân thất bại: %w", err)
	}
	return count > 0, nil
}

func scanTicket(s interface{ Scan(dest ...interface{}) error }) (*models.WeighingTicket, error) {
	t := &models.WeighingTicket{}
	var isActive int
	var runnerID sql.NullInt64
	if err := s.Scan(
		&t.ID, &t.TicketCode, &t.VolumeM3, &t.BatchCount, &t.M3Me,
		&t.OrderID, &t.OrderNumber,
		&t.ConcreteGradeID, &t.ConcreteGradeCode, &t.ConcreteGradeName,
		&t.PrintedConcreteGradeName, &t.HumidityValues,
		&t.StandardPerCube, &t.StandardPerBatch,
		&t.CustomerName, &t.ConstructionSiteAddress, &t.MaterialName,
		&runnerID, &t.Notes, &isActive,
		&t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if runnerID.Valid {
		id := runnerID.Int64
		t.RunnerUserID = &id
	}
	t.IsActive = isActive == 1
	return t, nil
}

func scanDetail(s interface{ Scan(dest ...interface{}) error }) (*models.WeighingTicketDetail, error) {
	d := &models.WeighingTicketDetail{}
	var isActive int
	if err := s.Scan(
		&d.ID, &d.WeighingTicketID, &d.BatchNo, &d.BatchVolumeM3, &d.MixedAt,
		&d.TP1, &d.TP2, &d.TP3, &d.TP4, &d.TP5, &d.TP6,
		&d.XI1, &d.XI2, &d.XI3, &d.XI4, &d.XI5,
		&d.PG1, &d.PG2, &d.PG3, &d.PG4, &d.PG5, &d.PG6,
		&d.ManualOverrideJSON, &d.Notes,
		&isActive, &d.CreatedAt, &d.UpdatedAt,
	); err != nil {
		return nil, err
	}
	d.IsActive = isActive == 1
	return d, nil
}

func nullIfEmpty(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

// StatsRow là một dòng thống kê tổng hợp theo mác bê tông.
type StatsRow struct {
	ConcreteGradeCode string  `json:"concrete_grade_code"`
	ConcreteGradeName string  `json:"concrete_grade_name"`
	TicketCount       int     `json:"ticket_count"`
	TotalVolumeM3     float64 `json:"total_volume_m3"`
	TotalBatchCount   int     `json:"total_batch_count"`
}

// StatsByGrade tổng hợp phiếu cân theo mác bê tông (cho báo cáo).
// from/to là mốc 'YYYY-MM-DD'; rỗng nghĩa là không giới hạn.
func (r *Repository) StatsByGrade(from, to string) ([]StatsRow, error) {
	where := "WHERE is_active = 1"
	args := []interface{}{}
	if from != "" {
		where += " AND date(created_at) >= date(?)"
		args = append(args, from)
	}
	if to != "" {
		where += " AND date(created_at) <= date(?)"
		args = append(args, to)
	}

	query := `
		SELECT concrete_grade_code, concrete_grade_name,
		       COUNT(*) AS ticket_count,
		       COALESCE(SUM(volume_m3), 0) AS total_volume,
		       COALESCE(SUM(batch_count), 0) AS total_batches
		FROM weighing_tickets ` + where + `
		GROUP BY concrete_grade_code, concrete_grade_name
		ORDER BY total_volume DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("truy vấn thống kê theo mác thất bại: %w", err)
	}
	defer rows.Close()

	out := make([]StatsRow, 0)
	for rows.Next() {
		var s StatsRow
		if err := rows.Scan(&s.ConcreteGradeCode, &s.ConcreteGradeName, &s.TicketCount, &s.TotalVolumeM3, &s.TotalBatchCount); err != nil {
			return nil, fmt.Errorf("đọc thống kê theo mác thất bại: %w", err)
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("đọc thống kê theo mác thất bại: %w", err)
	}
	return out, nil
}

// DailyRow là một dòng thống kê theo ngày.
type DailyRow struct {
	Day           string  `json:"day"`
	TicketCount   int     `json:"ticket_count"`
	TotalVolumeM3 float64 `json:"total_volume_m3"`
}

// StatsDaily tổng hợp sản lượng theo ngày (cho biểu đồ).
func (r *Repository) StatsDaily(from, to string) ([]DailyRow, error) {
	where := "WHERE is_active = 1"
	args := []interface{}{}
	if from != "" {
		where += " AND date(created_at) >= date(?)"
		args = append(args, from)
	}
	if to != "" {
		where += " AND date(created_at) <= date(?)"
		args = append(args, to)
	}

	query := `
		SELECT date(created_at) AS day,
		       COUNT(*) AS ticket_count,
		       COALESCE(SUM(volume_m3), 0) AS total_volume
		FROM weighing_tickets ` + where + `
		GROUP BY date(created_at)
		ORDER BY day ASC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("truy vấn thống kê theo ngày thất bại: %w", err)
	}
	defer rows.Close()

	out := make([]DailyRow, 0)
	for rows.Next() {
		var d DailyRow
		if err := rows.Scan(&d.Day, &d.TicketCount, &d.TotalVolumeM3); err != nil {
			return nil, fmt.Errorf("đọc thống kê theo ngày thất bại: %w", err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("đọc thống kê theo ngày thất bại: %w", err)
	}
	return out, nil
}
