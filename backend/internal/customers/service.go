package customers

import (
	"strings"

	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// CreateInput là dữ liệu đầu vào khi tạo khách hàng.
type CreateInput struct {
	CustomerCode string `json:"customer_code" binding:"required,max=50"`
	CustomerName string `json:"customer_name" binding:"required,max=200"`
	Address      string `json:"address" binding:"max=300"`
	Phone        string `json:"phone" binding:"max=30"`
}

// UpdateInput là dữ liệu đầu vào khi cập nhật khách hàng.
type UpdateInput struct {
	CustomerCode string `json:"customer_code" binding:"required,max=50"`
	CustomerName string `json:"customer_name" binding:"required,max=200"`
	Address      string `json:"address" binding:"max=300"`
	Phone        string `json:"phone" binding:"max=30"`
}

// Service chứa logic nghiệp vụ cho khách hàng.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List trả về danh sách khách hàng có phân trang.
func (s *Service) List(p httpx.Pagination) (httpx.PageResult, error) {
	items, total, err := s.repo.List(p)
	if err != nil {
		return httpx.PageResult{}, err
	}
	return httpx.NewPageResult(items, p.Page, p.PerPage, total), nil
}

// Get lấy khách hàng theo id.
func (s *Service) Get(id int64) (*models.Customer, error) {
	return s.repo.GetByID(id)
}

// Create tạo khách hàng mới sau khi chuẩn hóa dữ liệu.
func (s *Service) Create(input CreateInput) (*models.Customer, error) {
	input.CustomerCode = strings.TrimSpace(input.CustomerCode)
	input.CustomerName = strings.TrimSpace(input.CustomerName)
	input.Address = strings.TrimSpace(input.Address)
	input.Phone = strings.TrimSpace(input.Phone)

	return s.repo.Create(input)
}

// Update cập nhật khách hàng theo id.
func (s *Service) Update(id int64, input UpdateInput) (*models.Customer, error) {
	input.CustomerCode = strings.TrimSpace(input.CustomerCode)
	input.CustomerName = strings.TrimSpace(input.CustomerName)
	input.Address = strings.TrimSpace(input.Address)
	input.Phone = strings.TrimSpace(input.Phone)

	return s.repo.Update(id, input)
}

// Delete xóa mềm khách hàng theo id.
func (s *Service) Delete(id int64) error {
	return s.repo.SoftDelete(id)
}

// Restore khôi phục khách hàng đã xóa mềm theo id (chuyển về trạng thái hoạt động).
func (s *Service) Restore(id int64) (*models.Customer, error) {
	return s.repo.Restore(id)
}
