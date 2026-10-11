package sites

import (
	"strings"

	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// CreateInput là dữ liệu đầu vào khi tạo công trường.
type CreateInput struct {
	SiteCode string `json:"site_code" binding:"required,max=50"`
	SiteName string `json:"site_name" binding:"required,max=200"`
}

// UpdateInput là dữ liệu đầu vào khi cập nhật công trường.
type UpdateInput struct {
	SiteCode string `json:"site_code" binding:"required,max=50"`
	SiteName string `json:"site_name" binding:"required,max=200"`
}

// Service chứa logic nghiệp vụ cho công trường.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List trả về danh sách công trường có phân trang.
func (s *Service) List(p httpx.Pagination) (httpx.PageResult, error) {
	items, total, err := s.repo.List(p)
	if err != nil {
		return httpx.PageResult{}, err
	}
	return httpx.NewPageResult(items, p.Page, p.PerPage, total), nil
}

// Get lấy công trường theo id.
func (s *Service) Get(id int64) (*models.ConstructionSite, error) {
	return s.repo.GetByID(id)
}

// Create tạo mới công trường sau khi chuẩn hóa khoảng trắng.
func (s *Service) Create(input CreateInput) (*models.ConstructionSite, error) {
	input.SiteCode = strings.TrimSpace(input.SiteCode)
	input.SiteName = strings.TrimSpace(input.SiteName)
	return s.repo.Create(input)
}

// Update cập nhật công trường theo id.
func (s *Service) Update(id int64, input UpdateInput) (*models.ConstructionSite, error) {
	input.SiteCode = strings.TrimSpace(input.SiteCode)
	input.SiteName = strings.TrimSpace(input.SiteName)
	return s.repo.Update(id, input)
}

// Delete xóa mềm công trường theo id.
func (s *Service) Delete(id int64) error {
	return s.repo.SoftDelete(id)
}

// Restore khôi phục công trường đã xóa mềm theo id.
func (s *Service) Restore(id int64) (*models.ConstructionSite, error) {
	return s.repo.Restore(id)
}
