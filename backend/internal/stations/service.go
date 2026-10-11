package stations

import (
	"errors"
	"strings"

	"tramtron/backend/internal/models"
)

// ErrInvalidMask trả về khi component_visibility_mask không đúng định dạng.
var ErrInvalidMask = errors.New("component_visibility_mask phải gồm đúng 17 ký tự 0/1")

// ErrInvalidBatchRange trả về khi khoảng khối lượng mẻ không hợp lệ.
var ErrInvalidBatchRange = errors.New("khối lượng mẻ tối đa phải lớn hơn hoặc bằng tối thiểu")

// ErrInvalidStationCode trả về khi mã/tên trạm rỗng.
var ErrInvalidStationCode = errors.New("mã và tên trạm là bắt buộc")

// UpdateInput là dữ liệu cập nhật cấu hình trạm.
type UpdateInput struct {
	StationCode string `json:"station_code" binding:"required,max=50"`
	StationName string `json:"station_name" binding:"required,max=200"`
	Address     string `json:"address" binding:"max=300"`
	Phone       string `json:"phone" binding:"max=30"`

	TP1Name string `json:"tp1_name" binding:"max=100"`
	TP2Name string `json:"tp2_name" binding:"max=100"`
	TP3Name string `json:"tp3_name" binding:"max=100"`
	TP4Name string `json:"tp4_name" binding:"max=100"`
	TP5Name string `json:"tp5_name" binding:"max=100"`
	TP6Name string `json:"tp6_name" binding:"max=100"`
	XI1Name string `json:"xi1_name" binding:"max=100"`
	XI2Name string `json:"xi2_name" binding:"max=100"`
	XI3Name string `json:"xi3_name" binding:"max=100"`
	XI4Name string `json:"xi4_name" binding:"max=100"`
	XI5Name string `json:"xi5_name" binding:"max=100"`
	PG1Name string `json:"pg1_name" binding:"max=100"`
	PG2Name string `json:"pg2_name" binding:"max=100"`
	PG3Name string `json:"pg3_name" binding:"max=100"`
	PG4Name string `json:"pg4_name" binding:"max=100"`
	PG5Name string `json:"pg5_name" binding:"max=100"`
	PG6Name string `json:"pg6_name" binding:"max=100"`

	ComponentVisibilityMask string  `json:"component_visibility_mask" binding:"required"`
	AccumulateComponents    bool    `json:"accumulate_components"`
	AccumulateCement        bool    `json:"accumulate_cement"`
	BatchMinVolumeM3        float64 `json:"batch_min_volume_m3"`
	BatchMaxVolumeM3        float64 `json:"batch_max_volume_m3"`
}

// Service chứa logic cấu hình trạm.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Get đọc cấu hình trạm.
func (s *Service) Get() (*models.Station, error) {
	return s.repo.Get()
}

// Update cập nhật cấu hình trạm sau khi kiểm tra nghiệp vụ.
func (s *Service) Update(input UpdateInput) (*models.Station, error) {
	input.StationCode = strings.TrimSpace(input.StationCode)
	input.StationName = strings.TrimSpace(input.StationName)
	if input.StationCode == "" || input.StationName == "" {
		return nil, ErrInvalidStationCode
	}
	if err := ValidateVisibilityMask(input.ComponentVisibilityMask); err != nil {
		return nil, err
	}
	if input.BatchMaxVolumeM3 < input.BatchMinVolumeM3 {
		return nil, ErrInvalidBatchRange
	}
	return s.repo.Update(input)
}

// ValidateVisibilityMask kiểm tra mask đúng 17 ký tự 0/1.
func ValidateVisibilityMask(mask string) error {
	if len(mask) != 17 {
		return ErrInvalidMask
	}
	for _, ch := range mask {
		if ch != '0' && ch != '1' {
			return ErrInvalidMask
		}
	}
	return nil
}
