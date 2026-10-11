package models

// ConstructionSite ánh xạ bảng construction_sites (mục 13.1).
// Xóa là xóa mềm (is_active: 1 -> 0) để giữ lịch sử đơn hàng.
type ConstructionSite struct {
	ID        int64  `json:"id"`
	SiteCode  string `json:"site_code"`
	SiteName  string `json:"site_name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// MixerTruck ánh xạ bảng mixer_trucks (mục 13.1).
// Xóa là xóa mềm; xe đã dùng trong chuyến giao vẫn được giữ để bảo toàn lịch sử.
//
// DefaultDriverID là tài xế mặc định của xe (có thể rỗng). Đây chỉ là gợi ý
// để giao diện tự điền khi điều phối chuyến, không ràng buộc tài xế thực tế
// của mọi chuyến. DefaultDriverName là tên tài xế mặc định (đọc kèm qua JOIN).
type MixerTruck struct {
	ID                int64  `json:"id"`
	TruckCode         string `json:"truck_code"`
	PlateNumber       string `json:"plate_number"`
	DefaultDriverID   *int64 `json:"default_driver_id,omitempty"`
	DefaultDriverName string `json:"default_driver_name,omitempty"`
	IsActive          bool   `json:"is_active"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// Driver ánh xạ bảng drivers (mục 13.1).
// Tài xế không bắt buộc có tài khoản đăng nhập; xóa là xóa mềm.
type Driver struct {
	ID         int64  `json:"id"`
	DriverCode string `json:"driver_code"`
	FullName   string `json:"full_name"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
