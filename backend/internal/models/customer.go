package models

// Customer ánh xạ bảng customers (mục 13.1).
//
// Lưu ý nghiệp vụ: xóa khách hàng là XÓA MỀM — backend chỉ chuyển
// is_active từ 1 thành 0, không xóa vật lý dòng dữ liệu.
type Customer struct {
	ID           int64  `json:"id"`
	CustomerCode string `json:"customer_code"`
	CustomerName string `json:"customer_name"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
