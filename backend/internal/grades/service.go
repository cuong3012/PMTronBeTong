package grades

import (
	"fmt"
	"strings"

	"tramtron/backend/internal/httpx"
	"tramtron/backend/internal/models"
)

// ErrInvalidCode trả về khi mã/tên mác rỗng.
var ErrInvalidCode = fmt.Errorf("mã và tên mác bê tông là bắt buộc")

// ErrPgTooManyDecimals trả về khi phụ gia PG1–PG6 có quá 4 chữ số thập phân.
var ErrPgTooManyDecimals = fmt.Errorf("phụ gia pg1–pg6 chỉ được có tối đa 4 chữ số thập phân")

// GradeInput là dữ liệu tạo/cập nhật mác bê tông.
type GradeInput struct {
	GradeCode   string `json:"grade_code" binding:"required,max=50"`
	PrintedName string `json:"printed_name" binding:"required,max=200"`

	TP1  float64 `json:"tp1"`
	TP2  float64 `json:"tp2"`
	TP3  float64 `json:"tp3"`
	TP4  float64 `json:"tp4"`
	TP5  float64 `json:"tp5"`
	TP6  float64 `json:"tp6"`
	XI1  float64 `json:"xi1"`
	XI2  float64 `json:"xi2"`
	XI3  float64 `json:"xi3"`
	XI4  float64 `json:"xi4"`
	XI5  float64 `json:"xi5"`
	Nuoc float64 `json:"nuoc"`
	PG1  float64 `json:"pg1"`
	PG2  float64 `json:"pg2"`
	PG3  float64 `json:"pg3"`
	PG4  float64 `json:"pg4"`
	PG5  float64 `json:"pg5"`
	PG6  float64 `json:"pg6"`

	DoSut   string `json:"do_sut" binding:"max=100"`
	CuongDo string `json:"cuong_do" binding:"max=100"`
	Notes   string `json:"notes" binding:"max=500"`
}

// Service chứa logic nghiệp vụ mác bê tông.
type Service struct {
	repo *Repository
}

// NewService tạo Service mới.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List trả về danh sách mác bê tông phân trang.
func (s *Service) List(p httpx.Pagination) (httpx.PageResult, error) {
	items, total, err := s.repo.List(p)
	if err != nil {
		return httpx.PageResult{}, err
	}
	return httpx.NewPageResult(items, p.Page, p.PerPage, total), nil
}

// Get lấy mác theo id.
func (s *Service) Get(id int64) (*models.ConcreteGrade, error) {
	return s.repo.GetByID(id)
}

// Create tạo mác mới.
func (s *Service) Create(input GradeInput) (*models.ConcreteGrade, error) {
	if err := validate(input); err != nil {
		return nil, err
	}
	input.GradeCode = strings.TrimSpace(input.GradeCode)
	input.PrintedName = strings.TrimSpace(input.PrintedName)
	return s.repo.Create(input)
}

// Update cập nhật mác theo id.
func (s *Service) Update(id int64, input GradeInput) (*models.ConcreteGrade, error) {
	if err := validate(input); err != nil {
		return nil, err
	}
	input.GradeCode = strings.TrimSpace(input.GradeCode)
	input.PrintedName = strings.TrimSpace(input.PrintedName)
	return s.repo.Update(id, input)
}

// Delete xóa mềm mác theo id.
func (s *Service) Delete(id int64) error {
	return s.repo.SoftDelete(id)
}

// validate kiểm tra nghiệp vụ cho mác bê tông.
func validate(input GradeInput) error {
	if strings.TrimSpace(input.GradeCode) == "" || strings.TrimSpace(input.PrintedName) == "" {
		return ErrInvalidCode
	}
	for _, v := range []float64{input.PG1, input.PG2, input.PG3, input.PG4, input.PG5, input.PG6} {
		if decimals(v) > 4 {
			return ErrPgTooManyDecimals
		}
	}
	return nil
}

// decimals đếm số chữ số sau dấu thập phân của một số thực.
func decimals(v float64) int {
	s := fmt.Sprintf("%v", v)
	if i := strings.IndexByte(s, '.'); i >= 0 {
		// Bỏ ký tự khoa học nếu có (ví dụ 1e-05) để đếm an toàn.
		frac := s[i+1:]
		if e := strings.IndexAny(frac, "eE"); e >= 0 {
			frac = frac[:e]
		}
		frac = strings.TrimRight(frac, "0")
		return len(frac)
	}
	return 0
}
