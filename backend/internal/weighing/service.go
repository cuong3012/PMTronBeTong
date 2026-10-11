package weighing

import (
	"errors"
	"strings"

	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// ErrInvalidTicket trả về khi dữ liệu phiếu cân không hợp lệ.
var ErrInvalidTicket = errors.New("mã phiếu cân, số m³ và số mẻ là bắt buộc")

// ErrInvalidDetail trả về khi chi tiết mẻ không hợp lệ.
var ErrInvalidDetail = errors.New("mỗi mẻ cần có số mẻ và khối lượng m3 hợp lệ")

// DetailInput là dữ liệu một mẻ chi tiết khi tạo phiếu cân.
type DetailInput struct {
	BatchNo            int     `json:"batch_no"`
	BatchVolumeM3      float64 `json:"batch_volume_m3"`
	MixedAt            string  `json:"mixed_at" binding:"required"`
	TP1                float64 `json:"tp1"`
	TP2                float64 `json:"tp2"`
	TP3                float64 `json:"tp3"`
	TP4                float64 `json:"tp4"`
	TP5                float64 `json:"tp5"`
	TP6                float64 `json:"tp6"`
	XI1                float64 `json:"xi1"`
	XI2                float64 `json:"xi2"`
	XI3                float64 `json:"xi3"`
	XI4                float64 `json:"xi4"`
	XI5                float64 `json:"xi5"`
	PG1                float64 `json:"pg1"`
	PG2                float64 `json:"pg2"`
	PG3                float64 `json:"pg3"`
	PG4                float64 `json:"pg4"`
	PG5                float64 `json:"pg5"`
	PG6                float64 `json:"pg6"`
	ManualOverrideJSON string  `json:"manual_override_json"`
	Notes              string  `json:"notes" binding:"max=500"`
}

// TicketInput là dữ liệu tạo phiếu cân kèm các mẻ.
type TicketInput struct {
	TicketCode                string        `json:"ticket_code" binding:"required,max=50"`
	VolumeM3                  float64       `json:"volume_m3"`
	BatchCount                int           `json:"batch_count"`
	M3Me                      float64       `json:"m3me"`
	OrderID                   int64         `json:"order_id" binding:"required"`
	ConcreteGradeID           int64         `json:"concrete_grade_id" binding:"required"`
	ConcreteGradeCode         string        `json:"concrete_grade_code" binding:"required,max=50"`
	ConcreteGradeName         string        `json:"concrete_grade_name" binding:"required,max=200"`
	PrintedConcreteGradeName  string        `json:"printed_concrete_grade_name" binding:"required,max=200"`
	HumidityValues            string        `json:"humidity_values"`
	StandardPerCube           float64       `json:"standard_per_cube"`
	StandardPerBatch          float64       `json:"standard_per_batch"`
	CustomerName              string        `json:"customer_name" binding:"required,max=200"`
	ConstructionSiteAddress   string        `json:"construction_site_address" binding:"max=300"`
	MaterialName              string        `json:"material_name" binding:"max=200"`
	RunnerUserID              *int64        `json:"runner_user_id"`
	Notes                     string        `json:"notes" binding:"max=500"`
	Details                   []DetailInput `json:"details"`
}

// Service chứa logic nghiệp vụ phiếu cân.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List trả về danh sách phiếu cân phân trang.
func (s *Service) List(p httpx.Pagination) (httpx.PageResult, error) {
	items, total, err := s.repo.List(p)
	if err != nil {
		return httpx.PageResult{}, err
	}
	return httpx.NewPageResult(items, p.Page, p.PerPage, total), nil
}

// Get lấy phiếu cân kèm chi tiết.
func (s *Service) Get(id int64) (*models.WeighingTicket, error) {
	return s.repo.GetByID(id)
}

// Create tạo phiếu cân mới (kèm chi tiết, trong một transaction).
func (s *Service) Create(input TicketInput) (*models.WeighingTicket, error) {
	input.TicketCode = strings.TrimSpace(input.TicketCode)
	if input.TicketCode == "" || input.VolumeM3 <= 0 || input.BatchCount <= 0 {
		return nil, ErrInvalidTicket
	}
	for _, d := range input.Details {
		if d.BatchNo <= 0 || d.BatchVolumeM3 <= 0 || strings.TrimSpace(d.MixedAt) == "" {
			return nil, ErrInvalidDetail
		}
	}
	return s.repo.CreateTicket(input)
}

// Delete xóa mềm phiếu cân (kèm chi tiết).
func (s *Service) Delete(id int64) error {
	return s.repo.SoftDeleteTicket(id)
}

// Stats trả về thống kê tổng hợp theo mác và theo ngày.
func (s *Service) Stats(from, to string) (map[string]interface{}, error) {
	byGrade, err := s.repo.StatsByGrade(from, to)
	if err != nil {
		return nil, err
	}
	daily, err := s.repo.StatsDaily(from, to)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"by_grade": byGrade,
		"daily":    daily,
	}, nil
}
