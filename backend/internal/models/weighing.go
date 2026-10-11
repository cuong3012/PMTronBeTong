package models

// WeighingTicket ánh xạ bảng weighing_tickets (mục 13.1) — phiếu cân.
// Đây là bản ghi kết quả (header/tổng quan), lưu kèm bản sao dữ liệu
// khách hàng/công trường/mác để tra cứu dù dữ liệu gốc đổi về sau.
type WeighingTicket struct {
	ID                        int64   `json:"id"`
	TicketCode                string  `json:"ticket_code"`
	VolumeM3                  float64 `json:"volume_m3"`
	BatchCount                int     `json:"batch_count"`
	M3Me                      float64 `json:"m3me"`
	OrderID                   int64   `json:"order_id"`
	OrderNumber               string  `json:"order_number,omitempty"`
	ConcreteGradeID           int64   `json:"concrete_grade_id"`
	ConcreteGradeCode         string  `json:"concrete_grade_code"`
	ConcreteGradeName         string  `json:"concrete_grade_name"`
	PrintedConcreteGradeName  string  `json:"printed_concrete_grade_name"`
	HumidityValues            string  `json:"humidity_values"`
	StandardPerCube           float64 `json:"standard_per_cube"`
	StandardPerBatch          float64 `json:"standard_per_batch"`
	CustomerName              string  `json:"customer_name"`
	ConstructionSiteAddress   string  `json:"construction_site_address"`
	MaterialName              string  `json:"material_name"`
	RunnerUserID              *int64  `json:"runner_user_id,omitempty"`
	Notes                     string  `json:"notes,omitempty"`
	IsActive                  bool    `json:"is_active"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 string  `json:"updated_at"`

	// Details là danh sách mẻ chi tiết; chỉ điền khi lấy chi tiết 1 phiếu.
	Details []WeighingTicketDetail `json:"details,omitempty"`
}

// WeighingTicketDetail ánh xạ bảng weighing_ticket_details (mục 13.1) — chi tiết mẻ.
// 1 record = 1 mẻ trong 1 phiếu cân, với 16 cột thành phần chính.
type WeighingTicketDetail struct {
	ID                 int64   `json:"id"`
	WeighingTicketID   int64   `json:"weighing_ticket_id"`
	BatchNo            int     `json:"batch_no"`
	BatchVolumeM3      float64 `json:"batch_volume_m3"`
	MixedAt            string  `json:"mixed_at"`
	TP1                float64 `json:"tp1"`
	TP2                float64 `json:"tp2"`
	TP3                float64 `json:"tp3"`
	TP4                float64 `json:"tp4"`
	TP5                float64 `json:"tp5"`
	TP6                float64 `json:"tp6"`
	XI1                float64 `json:"xi1"`
	XI2                float64 `json:"xi2"`
	XI3                float64 `json:"xi3"`
	XI4                float64 `json:"xi4"`
	XI5                float64 `json:"xi5"`
	PG1                float64 `json:"pg1"`
	PG2                float64 `json:"pg2"`
	PG3                float64 `json:"pg3"`
	PG4                float64 `json:"pg4"`
	PG5                float64 `json:"pg5"`
	PG6                float64 `json:"pg6"`
	ManualOverrideJSON string  `json:"manual_override_json,omitempty"`
	Notes              string  `json:"notes,omitempty"`
	IsActive           bool    `json:"is_active"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}
