package waitingmix

import (
	"errors"

	"tramtron/backend/internal/models"
)

// ErrInvalidStatus trả về khi status không phải 0/1.
var ErrInvalidStatus = errors.New("status chỉ nhận giá trị 0 hoặc 1")

// ErrInvalidVolume trả về khi khối lượng không hợp lệ.
var ErrInvalidVolume = errors.New("khối lượng phải là số dương và m3me không âm")

// ErrInvalidGrade trả về khi mác bê tông không tồn tại hoặc ngừng hoạt động.
var ErrInvalidGrade = errors.New("mác bê tông không tồn tại hoặc đã ngừng hoạt động")

// ErrInvalidCustomer trả về khi khách hàng không tồn tại hoặc ngừng hoạt động.
var ErrInvalidCustomer = errors.New("khách hàng không tồn tại hoặc đã ngừng hoạt động")

// UpdateInput là dữ liệu cập nhật một vị trí cấp phối chờ.
type UpdateInput struct {
	ConcreteGradeID int64   `json:"concrete_grade_id" binding:"required"`
	CustomerID      int64   `json:"customer_id" binding:"required"`
	VolumeM3        float64 `json:"volume_m3"`
	M3Me            float64 `json:"m3me"`
	Status          int     `json:"status"`
}

// Service chứa logic nghiệp vụ cấp phối chờ.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List trả về đủ 20 vị trí cấp phối chờ.
func (s *Service) List() ([]models.WaitingMix, error) {
	return s.repo.List()
}

// Get lấy một vị trí theo id.
func (s *Service) Get(id int64) (*models.WaitingMix, error) {
	return s.repo.GetByID(id)
}

// Update cập nhật một vị trí sau khi kiểm tra nghiệp vụ.
func (s *Service) Update(id int64, input UpdateInput) (*models.WaitingMix, error) {
	if input.Status != 0 && input.Status != 1 {
		return nil, ErrInvalidStatus
	}
	if input.VolumeM3 <= 0 || input.M3Me < 0 {
		return nil, ErrInvalidVolume
	}
	if input.ConcreteGradeID <= 0 {
		return nil, ErrInvalidGrade
	}
	ok, err := s.repo.GradeActive(input.ConcreteGradeID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidGrade
	}
	if input.CustomerID <= 0 {
		return nil, ErrInvalidCustomer
	}
	ok, err = s.repo.CustomerActive(input.CustomerID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCustomer
	}
	return s.repo.Update(id, input)
}
