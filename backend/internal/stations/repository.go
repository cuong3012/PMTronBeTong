package stations

import (
	"database/sql"
	"errors"
	"fmt"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/models"
)

// ErrNotFound trả về khi chưa có bản ghi cấu hình trạm.
var ErrNotFound = errors.New("chưa có cấu hình trạm")

const selectColumns = `
	SELECT id, station_code, station_name, COALESCE(address, ''), COALESCE(phone, ''),
	       tp1_name, tp2_name, tp3_name, tp4_name, tp5_name, tp6_name,
	       xi1_name, xi2_name, xi3_name, xi4_name, xi5_name,
	       pg1_name, pg2_name, pg3_name, pg4_name, pg5_name, pg6_name,
	       component_visibility_mask, accumulate_components, accumulate_cement,
	       batch_min_volume_m3, batch_max_volume_m3, created_at, updated_at
	FROM stations`

// Repository truy vấn cấu hình trạm (đúng một bản ghi).
type Repository struct {
	db *database.DB
}

// NewRepository tạo Repository mới.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Get đọc cấu hình trạm. Vì chỉ có 1 bản ghi, lấy bản ghi nhỏ nhất theo id.
func (r *Repository) Get() (*models.Station, error) {
	row := r.db.QueryRow(selectColumns + " ORDER BY id LIMIT 1")
	s, err := scanStation(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Update cập nhật cấu hình trạm duy nhất. Giữ nguyên id/created_at.
func (r *Repository) Update(input UpdateInput) (*models.Station, error) {
	current, err := r.Get()
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(`
		UPDATE stations SET
			station_code = ?, station_name = ?, address = ?, phone = ?,
			tp1_name = ?, tp2_name = ?, tp3_name = ?, tp4_name = ?, tp5_name = ?, tp6_name = ?,
			xi1_name = ?, xi2_name = ?, xi3_name = ?, xi4_name = ?, xi5_name = ?,
			pg1_name = ?, pg2_name = ?, pg3_name = ?, pg4_name = ?, pg5_name = ?, pg6_name = ?,
			component_visibility_mask = ?, accumulate_components = ?, accumulate_cement = ?,
			batch_min_volume_m3 = ?, batch_max_volume_m3 = ?,
			updated_at = datetime('now')
		WHERE id = ?`,
		input.StationCode, input.StationName, input.Address, input.Phone,
		input.TP1Name, input.TP2Name, input.TP3Name, input.TP4Name, input.TP5Name, input.TP6Name,
		input.XI1Name, input.XI2Name, input.XI3Name, input.XI4Name, input.XI5Name,
		input.PG1Name, input.PG2Name, input.PG3Name, input.PG4Name, input.PG5Name, input.PG6Name,
		input.ComponentVisibilityMask, boolToInt(input.AccumulateComponents), boolToInt(input.AccumulateCement),
		input.BatchMinVolumeM3, input.BatchMaxVolumeM3,
		current.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("cập nhật cấu hình trạm thất bại: %w", err)
	}
	return r.Get()
}

func scanStation(s interface{ Scan(dest ...interface{}) error }) (*models.Station, error) {
	st := &models.Station{}
	var accumComponents, accumCement int
	if err := s.Scan(
		&st.ID, &st.StationCode, &st.StationName, &st.Address, &st.Phone,
		&st.TP1Name, &st.TP2Name, &st.TP3Name, &st.TP4Name, &st.TP5Name, &st.TP6Name,
		&st.XI1Name, &st.XI2Name, &st.XI3Name, &st.XI4Name, &st.XI5Name,
		&st.PG1Name, &st.PG2Name, &st.PG3Name, &st.PG4Name, &st.PG5Name, &st.PG6Name,
		&st.ComponentVisibilityMask, &accumComponents, &accumCement,
		&st.BatchMinVolumeM3, &st.BatchMaxVolumeM3, &st.CreatedAt, &st.UpdatedAt,
	); err != nil {
		return nil, err
	}
	st.AccumulateComponents = accumComponents == 1
	st.AccumulateCement = accumCement == 1
	return st, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
