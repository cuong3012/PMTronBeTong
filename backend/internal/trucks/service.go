package trucks

import (
	"strings"

	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// CreateInput là dữ liệu đầu vào khi tạo xe bồn.
type CreateInput struct {
	TruckCode       string `json:"truck_code" binding:"required,max=50"`
	PlateNumber     string `json:"plate_number" binding:"required,max=30"`
	DefaultDriverID *int64 `json:"default_driver_id"` // tùy chọn, có thể bỏ trống
}

// UpdateInput là dữ liệu đầu vào khi cập nhật xe bồn.
type UpdateInput struct {
	TruckCode       string `json:"truck_code" binding:"required,max=50"`
	PlateNumber     string `json:"plate_number" binding:"required,max=30"`
	DefaultDriverID *int64 `json:"default_driver_id"` // tùy chọn, có thể bỏ trống
}

// Service chứa logic nghiệp vụ cho xe bồn.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List trả về danh sách xe bồn có phân trang.
func (s *Service) List(p httpx.Pagination) (httpx.PageResult, error) {
	items, total, err := s.repo.List(p)
	if err != nil {
		return httpx.PageResult{}, err
	}
	return httpx.NewPageResult(items, p.Page, p.PerPage, total), nil
}

// Get lấy xe bồn theo id.
func (s *Service) Get(id int64) (*models.MixerTruck, error) {
	return s.repo.GetByID(id)
}

// Create chèn xe bồn mới, biển số được chuẩn hóa thành chữ hoa.
func (s *Service) Create(input CreateInput) (*models.MixerTruck, error) {
	input.TruckCode = strings.TrimSpace(input.TruckCode)
	input.PlateNumber = strings.ToUpper(strings.TrimSpace(input.PlateNumber))
	return s.repo.Create(input)
}

// Update cập nhật xe bồn theo id.
func (s *Service) Update(id int64, input UpdateInput) (*models.MixerTruck, error) {
	input.TruckCode = strings.TrimSpace(input.TruckCode)
	input.PlateNumber = strings.ToUpper(strings.TrimSpace(input.PlateNumber))
	return s.repo.Update(id, input)
}

// Delete xóa mềm xe bồn theo id.
func (s *Service) Delete(id int64) error {
	return s.repo.SoftDelete(id)
}

// Restore khôi phục xe bồn đã xóa mềm theo id.
func (s *Service) Restore(id int64) (*models.MixerTruck, error) {
	return s.repo.Restore(id)
}
