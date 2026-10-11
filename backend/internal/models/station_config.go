package models

// Station là cấu hình của trạm hiện tại (mục 13.1 — bảng stations).
// Phiên bản đầu chỉ có ĐÚNG MỘT bản ghi; không tạo/xóa qua API, chỉ đọc/cập nhật.
type Station struct {
	ID         int64  `json:"id"`
	StationCode string `json:"station_code"`
	StationName string `json:"station_name"`
	Address    string `json:"address,omitempty"`
	Phone      string `json:"phone,omitempty"`

	// Tên hiển thị cho từng vị trí thành phần TP1–TP6.
	TP1Name string `json:"tp1_name"`
	TP2Name string `json:"tp2_name"`
	TP3Name string `json:"tp3_name"`
	TP4Name string `json:"tp4_name"`
	TP5Name string `json:"tp5_name"`
	TP6Name string `json:"tp6_name"`

	// Tên hiển thị XI1–XI5.
	XI1Name string `json:"xi1_name"`
	XI2Name string `json:"xi2_name"`
	XI3Name string `json:"xi3_name"`
	XI4Name string `json:"xi4_name"`
	XI5Name string `json:"xi5_name"`

	// Tên hiển thị PG1–PG6.
	PG1Name string `json:"pg1_name"`
	PG2Name string `json:"pg2_name"`
	PG3Name string `json:"pg3_name"`
	PG4Name string `json:"pg4_name"`
	PG5Name string `json:"pg5_name"`
	PG6Name string `json:"pg6_name"`

	// Mask đúng 17 ký tự 0/1: thứ tự TP1..TP6, XI1..XI5, PG1..PG6.
	ComponentVisibilityMask string `json:"component_visibility_mask"`
	AccumulateComponents    bool   `json:"accumulate_components"`
	AccumulateCement        bool   `json:"accumulate_cement"`
	BatchMinVolumeM3        float64 `json:"batch_min_volume_m3"`
	BatchMaxVolumeM3        float64 `json:"batch_max_volume_m3"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ConcreteGrade ánh xạ bảng concrete_grades (mục 13.1).
// Mỗi bản ghi là định mức thành phần cho ĐÚNG 1 m³ bê tông (không phải 1 mẻ).
// Thành phần không dùng lưu giá trị 0, không NULL.
//
// DoSut (độ sụt) và CuongDo (cường độ) lưu dạng TEXT vì thực tế thường là
// khoảng hoặc ký hiệu (ví dụ "10-12", "B25", "M300"), không phải số thuần.
type ConcreteGrade struct {
	ID          int64  `json:"id"`
	GradeCode   string `json:"grade_code"`
	PrintedName string `json:"printed_name"`

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

	DoSut   string `json:"do_sut"`
	CuongDo string `json:"cuong_do"`
	Notes   string `json:"notes,omitempty"`

	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ComponentCatalog ánh xạ bảng component_catalog (mục 13.1).
// Đây là metadata tên/định nghĩa vị trí thành phần, không phải dữ liệu sản xuất.
type ComponentCatalog struct {
	ID             int64  `json:"id"`
	SlotName       string `json:"slot_name"`
	HumidityValues string `json:"humidity_values"`
	Unit           string `json:"unit,omitempty"`
	Notes          string `json:"notes,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// WaitingMix ánh xạ bảng waiting_mixes (mục 13.1).
// Đúng 20 vị trí cố định (stt 1..20); chỉ đọc/cập nhật, không tạo/xóa.
// concrete_grade_id/customer_id có thể chưa được gán (vị trí trống).
// Trả kèm tên mác bê tông và tên khách hàng để giao diện hiển thị.
type WaitingMix struct {
	ID                int64   `json:"id"`
	STT               int     `json:"stt"`
	ConcreteGradeID   *int64  `json:"concrete_grade_id,omitempty"`
	ConcreteGradeCode string  `json:"concrete_grade_code,omitempty"`
	ConcreteGradeName string  `json:"concrete_grade_name,omitempty"`
	CustomerID        *int64  `json:"customer_id,omitempty"`
	CustomerName      string  `json:"customer_name,omitempty"`
	VolumeM3          float64 `json:"volume_m3"`
	M3Me              float64 `json:"m3me"`
	Status            int     `json:"status"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}
