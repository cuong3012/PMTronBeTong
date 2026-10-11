package httpx

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Pagination chứa tham số phân trang đọc từ query string.
type Pagination struct {
	Page    int    // trang hiện tại, bắt đầu từ 1
	PerPage int    // số bản ghi mỗi trang
	Search  string // từ khóa tìm kiếm (tùy module diễn giải)
}

// Offset tính vị trí bắt đầu cho câu SQL LIMIT/OFFSET.
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// Limit trả về số bản ghi cần lấy (chính là PerPage).
func (p Pagination) Limit() int {
	return p.PerPage
}

// ParsePagination đọc page, per_page và search từ query string.
//
// Giá trị mặc định: page=1, per_page=20. Giới hạn per_page tối đa 100
// để tránh truy vấn quá lớn làm treo hệ thống.
func ParsePagination(c *gin.Context) Pagination {
	page := atoiDefault(c.Query("page"), 1)
	if page < 1 {
		page = 1
	}

	perPage := atoiDefault(c.Query("per_page"), 20)
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	return Pagination{
		Page:    page,
		PerPage: perPage,
		Search:  strings.TrimSpace(c.Query("search")),
	}
}

// PageResult là cấu trúc trả về cho danh sách có phân trang.
type PageResult struct {
	Items   interface{} `json:"items"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
	Total   int         `json:"total"`
}

// NewPageResult tạo kết quả phân trang từ danh sách và số lượng.
func NewPageResult(items interface{}, page, perPage, total int) PageResult {
	return PageResult{Items: items, Page: page, PerPage: perPage, Total: total}
}

// atoiDefault chuyển chuỗi sang int, trả về giá trị mặc định nếu lỗi.
func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
