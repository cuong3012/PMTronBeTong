package components

import (
	"errors"
	"strings"

	"tramtron/backend/internal/models"
)

// ErrInvalidSlot trả về khi tên vị trí rỗng.
var ErrInvalidSlot = errors.New("tên vị trí thành phần là bắt buộc")

// CatalogInput là dữ liệu tạo/cập nhật thành phần.
type CatalogInput struct {
	SlotName       string `json:"slot_name" binding:"required,max=100"`
	HumidityValues string `json:"humidity_values"`
	Unit           string `json:"unit" binding:"max=30"`
	Notes          string `json:"notes" binding:"max=300"`
}

// Service chứa logic nghiệp vụ danh mục thành phần.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List trả về danh sách thành phần.
func (s *Service) List(search string) ([]models.ComponentCatalog, error) {
	return s.repo.List(strings.TrimSpace(search))
}

// Get lấy thành phần theo id.
func (s *Service) Get(id int64) (*models.ComponentCatalog, error) {
	return s.repo.GetByID(id)
}

// Create tạo thành phần mới.
func (s *Service) Create(input CatalogInput) (*models.ComponentCatalog, error) {
	input.SlotName = strings.TrimSpace(input.SlotName)
	if input.SlotName == "" {
		return nil, ErrInvalidSlot
	}
	return s.repo.Create(input)
}

// Update cập nhật thành phần.
func (s *Service) Update(id int64, input CatalogInput) (*models.ComponentCatalog, error) {
	input.SlotName = strings.TrimSpace(input.SlotName)
	if input.SlotName == "" {
		return nil, ErrInvalidSlot
	}
	return s.repo.Update(id, input)
}

// Delete xóa thành phần.
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
