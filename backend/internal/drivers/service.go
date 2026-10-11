package drivers

import (
	"strings"

	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// CreateInput là dữ liệu đầu vào khi tạo tài xế.
type CreateInput struct {
	DriverCode string `json:"driver_code" binding:"required,max=50"`
	FullName   string `json:"full_name" binding:"required,max=200"`
}

// UpdateInput là dữ liệu đầu vào khi cập nhật tài xế.
type UpdateInput struct {
	DriverCode string `json:"driver_code" binding:"required,max=50"`
	FullName   string `json:"full_name" binding:"required,max=200"`
}

// Service chứa logic nghiệp vụ cho tài xế.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List trả về danh sách tài xế có phân trang.
func (s *Service) List(p httpx.Pagination) (httpx.PageResult, error) {
	items, total, err := s.repo.List(p)
	if err != nil {
		return httpx.PageResult{}, err
	}
	return httpx.NewPageResult(items, p.Page, p.PerPage, total), nil
}

// Get lấy tài xế theo id.
func (s *Service) Get(id int64) (*models.Driver, error) {
	return s.repo.GetByID(id)
}

// Create chèn tài xế mới, cắt khoảng trắng thừa ở mã và tên.
func (s *Service) Create(input CreateInput) (*models.Driver, error) {
	input.DriverCode = strings.TrimSpace(input.DriverCode)
	input.FullName = strings.TrimSpace(input.FullName)
	return s.repo.Create(input)
}

// Update chỉnh sửa tài xế theo id.
func (s *Service) Update(id int64, input UpdateInput) (*models.Driver, error) {
	input.DriverCode = strings.TrimSpace(input.DriverCode)
	input.FullName = strings.TrimSpace(input.FullName)
	return s.repo.Update(id, input)
}

// Delete xóa mềm tài xế theo id.
func (s *Service) Delete(id int64) error {
	return s.repo.SoftDelete(id)
}

// Restore khôi phục tài xế đã xóa mềm theo id.
func (s *Service) Restore(id int64) (*models.Driver, error) {
	return s.repo.Restore(id)
}
