package database

import (
	"fmt"
	"log"
)

// SeedSampleResult ghi lại số bản ghi mẫu đã chèn, phục vụ log/test.
type SeedSampleResult struct {
	GradesCreated     int
	ComponentsCreated int
	WaitingCreated    int
	WaitingFilled     int
	TicketsCreated    int
	DetailsCreated    int
}

// SeedEssentialData tạo dữ liệu KHỞI TẠO BẮT BUỘC cho các bảng mới:
//   - component_catalog: 17 vị trí thành phần (TP1..PG6).
//   - waiting_mixes: đủ 20 vị trí cấp phối chờ (stt 1..20), trạng thái trống.
//
// Hàm idempotent (chỉ chèn khi bảng còn trống). Đây là phần cần cho hệ thống
// vận hành, KHÁC với dữ liệu mẫu (SeedSampleData) dùng cho thống kê/demo.
func SeedEssentialData(db *DB) (*SeedSampleResult, error) {
	res := &SeedSampleResult{}
	if err := seedComponentCatalog(db, res); err != nil {
		return nil, err
	}
	if err := seedWaitingMixes(db, res); err != nil {
		return nil, err
	}
	return res, nil
}

// SeedSampleData chèn dữ liệu mẫu để thống kê/demo có dữ liệu:
// mác bê tông, khách hàng, công trường, đơn hàng và phiếu cân
// (weighing_tickets) kèm chi tiết (weighing_ticket_details).
//
// KHÔNG tự động chạy khi khởi động để tránh làm bẩn dữ liệu thật; bật bằng
// biến môi trường SEED_SAMPLE=true (xem cmd/server/main.go) hoặc gọi tường minh.
// Hàm idempotent: chỉ chèn khi bảng tương ứng còn trống.
func SeedSampleData(db *DB) (*SeedSampleResult, error) {
	res := &SeedSampleResult{}

	if err := seedConcreteGrades(db, res); err != nil {
		return nil, err
	}
	if err := seedCustomerAndOrder(db); err != nil {
		return nil, err
	}
	if err := fillWaitingMixesDemo(db, res); err != nil {
		return nil, err
	}
	if err := seedWeighingTickets(db, res); err != nil {
		return nil, err
	}

	return res, nil
}

// seedCustomerAndOrder tạo một khách hàng + công trường + đơn hàng mẫu nếu trống,
// để phiếu cân mẫu có thể tham chiếu tới.
func seedCustomerAndOrder(db *DB) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM customers`).Scan(&count); err != nil {
		return fmt.Errorf("đếm customers thất bại: %w", err)
	}
	if count == 0 {
		if _, err := db.Exec(
			`INSERT INTO customers (customer_code, customer_name, address, phone)
			 VALUES ('KH001', 'Công ty TNHH Xây dựng ABC', 'Hà Nội', '0900000001')`,
		); err != nil {
			return fmt.Errorf("chèn khách hàng mẫu thất bại: %w", err)
		}
	}

	if err := db.QueryRow(`SELECT COUNT(*) FROM construction_sites`).Scan(&count); err != nil {
		return fmt.Errorf("đếm construction_sites thất bại: %w", err)
	}
	if count == 0 {
		if _, err := db.Exec(
			`INSERT INTO construction_sites (site_code, site_name)
			 VALUES ('CT001', 'Dự án Chung cư Vinhomes')`,
		); err != nil {
			return fmt.Errorf("chèn công trường mẫu thất bại: %w", err)
		}
	}

	if err := db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&count); err != nil {
		return fmt.Errorf("đếm orders thất bại: %w", err)
	}
	if count == 0 {
		customerID := firstID(db, "customers")
		siteID := firstID(db, "construction_sites")
		if customerID > 0 && siteID > 0 {
			if _, err := db.Exec(
				`INSERT INTO orders (order_number, customer_id, construction_site_id, ordered_volume_m3)
				 VALUES ('DH001', ?, ?, 50)`,
				customerID, siteID,
			); err != nil {
				return fmt.Errorf("chèn đơn hàng mẫu thất bại: %w", err)
			}
		}
	}
	return nil
}

// seedComponentCatalog tạo 17 vị trí thành phần nếu bảng trống.
func seedComponentCatalog(db *DB, res *SeedSampleResult) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM component_catalog`).Scan(&count); err != nil {
		return fmt.Errorf("đếm component_catalog thất bại: %w", err)
	}
	if count > 0 {
		return nil
	}

	slots := []struct {
		name    string
		unit    string
		humidity string
	}{
		{"TP1", "kg", "12.5,11.8,10.7,9.6,8.9,7.4"},
		{"TP2", "kg", ""},
		{"TP3", "kg", ""},
		{"TP4", "kg", ""},
		{"TP5", "kg", ""},
		{"TP6", "kg", ""},
		{"XI1", "kg", ""},
		{"XI2", "kg", ""},
		{"XI3", "kg", ""},
		{"XI4", "kg", ""},
		{"XI5", "kg", ""},
		{"PG1", "lít", ""},
		{"PG2", "lít", ""},
		{"PG3", "lít", ""},
		{"PG4", "lít", ""},
		{"PG5", "lít", ""},
		{"PG6", "lít", ""},
	}
	for _, s := range slots {
		if _, err := db.Exec(
			`INSERT INTO component_catalog (slot_name, humidity_values, unit) VALUES (?, ?, ?)`,
			s.name, s.humidity, s.unit,
		); err != nil {
			return fmt.Errorf("chèn thành phần %q thất bại: %w", s.name, err)
		}
		res.ComponentsCreated++
	}
	return nil
}

// seedConcreteGrades tạo vài mác bê tông mẫu nếu bảng trống.
func seedConcreteGrades(db *DB, res *SeedSampleResult) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM concrete_grades`).Scan(&count); err != nil {
		return fmt.Errorf("đếm concrete_grades thất bại: %w", err)
	}
	if count > 0 {
		return nil
	}

	// tp1..tp6 = cát/đá, xi1..xi5 = xi măng, pg1..pg6 = phụ gia (mẫu minh họa).
	// do_sut/cuong_do là TEXT (có thể là khoảng hoặc ký hiệu).
	grades := []struct {
		code, name                    string
		tp1, tp2, tp3, tp4, xi1, nuoc float64
		pg1                           float64
		doSut, cuongDo                string
	}{
		{"M250", "Bê tông M250", 750, 1100, 0, 0, 320, 175, 0, "10-12", "M250"},
		{"M300", "Bê tông M300", 720, 1080, 0, 0, 360, 170, 4.5, "10-12", "M300"},
		{"M350", "Bê tông M350", 700, 1060, 0, 0, 400, 165, 5.2, "8-10", "M350"},
	}
	for _, g := range grades {
		if _, err := db.Exec(`
			INSERT INTO concrete_grades
				(grade_code, printed_name, tp1, tp2, tp3, tp4, tp5, tp6,
				 xi1, xi2, xi3, xi4, xi5, nuoc, pg1, do_sut, cuong_do)
			VALUES (?, ?, ?, ?, ?, ?, 0, 0, ?, 0, 0, 0, 0, ?, ?, ?, ?)`,
			g.code, g.name, g.tp1, g.tp2, g.tp3, g.tp4,
			g.xi1, g.nuoc, g.pg1, g.doSut, g.cuongDo,
		); err != nil {
			return fmt.Errorf("chèn mác %q thất bại: %w", g.code, err)
		}
		res.GradesCreated++
	}
	return nil
}

// seedWaitingMixes tạo đủ 20 vị trí cấp phối chờ (stt 1..20, trống) nếu bảng trống.
func seedWaitingMixes(db *DB, res *SeedSampleResult) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM waiting_mixes`).Scan(&count); err != nil {
		return fmt.Errorf("đếm waiting_mixes thất bại: %w", err)
	}
	if count > 0 {
		return nil
	}

	for i := 1; i <= 20; i++ {
		_, err := db.Exec(
			`INSERT INTO waiting_mixes (stt, volume_m3, m3me, status)
			 VALUES (?, 0, 0, 0)`,
			i,
		)
		if err != nil {
			return fmt.Errorf("chèn vị trí cấp phối chờ stt=%d thất bại: %w", i, err)
		}
		res.WaitingCreated++
	}
	return nil
}

// fillWaitingMixesDemo điền vài vị trí cấp phối chờ đầu tiên (dữ liệu mẫu).
func fillWaitingMixesDemo(db *DB, res *SeedSampleResult) error {
	gradeID := firstID(db, "concrete_grades")
	customerID := firstID(db, "customers")
	if gradeID == 0 || customerID == 0 {
		return nil
	}
	for i := 1; i <= 5; i++ {
		if _, err := db.Exec(
			`UPDATE waiting_mixes SET concrete_grade_id = ?, customer_id = ?, volume_m3 = ?, m3me = ?, status = ?
			 WHERE stt = ?`,
			gradeID, customerID, float64(6+i), 1.5, 1, i,
		); err != nil {
			return fmt.Errorf("điền vị trí cấp phối chờ stt=%d thất bại: %w", i, err)
		}
		res.WaitingFilled++
	}
	return nil
}

// seedWeighingTickets tạo vài phiếu cân mẫu kèm chi tiết nếu bảng trống.
func seedWeighingTickets(db *DB, res *SeedSampleResult) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM weighing_tickets`).Scan(&count); err != nil {
		return fmt.Errorf("đếm weighing_tickets thất bại: %w", err)
	}
	if count > 0 {
		return nil
	}

	// Cần có order/khách để tham chiếu; nếu chưa có thì bỏ qua.
	orderID := firstID(db, "orders")
	customerID := firstID(db, "customers")
	if orderID == 0 || customerID == 0 {
		log.Printf("[seed] bỏ qua phiếu cân mẫu do thiếu order/customer")
		return nil
	}

	var customerName string
	if err := db.QueryRow(`SELECT customer_name FROM customers WHERE id = ?`, customerID).
		Scan(&customerName); err != nil {
		return fmt.Errorf("đọc khách hàng cho phiếu cân mẫu thất bại: %w", err)
	}

	// Lấy danh sách mác bê tông để tạo phiếu đa dạng (khác mác giữa các phiếu).
	gradeID, gradeCode, gradeName := gradeRef(db)

	// spec mô tả một phiếu cân mẫu. Tên/khối lượng mẻ được suy ra từ m3me.
	type spec struct {
		code    string // mã phiếu
		batch   int    // số mẻ
		m3me    float64
		dateMod string // '-N days' cho datetime('now', ...)
		active  int    // 1 = đang hoạt động, 0 = đã xóa mềm
	}
	specs := []spec{
		{"PC001", 3, 2.0, "-6 days", 1},
		{"PC002", 4, 2.0, "-3 days", 1},
		{"PC003", 5, 2.0, "-1 days", 1},
		{"PC004", 3, 2.5, "-0 days", 1},
		{"PC005", 2, 3.0, "-0 days", 1},
		{"PC006", 4, 1.5, "-2 days", 0}, // ví dụ phiếu đã xóa mềm
	}

	for _, s := range specs {
		volume := float64(s.batch) * s.m3me
		standardPerBatch := s.m3me * 2400

		// Chèn header. created_at/updated_at lấy theo mốc ngày tương đối.
		r, err := db.Exec(`
			INSERT INTO weighing_tickets
				(ticket_code, volume_m3, batch_count, m3me, order_id, concrete_grade_id,
				 concrete_grade_code, concrete_grade_name, printed_concrete_grade_name,
				 humidity_values, standard_per_cube, standard_per_batch,
				 customer_name, construction_site_address, material_name,
				 is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			        datetime('now', ?), datetime('now', ?))`,
			s.code, volume, s.batch, s.m3me, orderID, gradeID,
			gradeCode, gradeName, gradeName,
			"12.5,11.8,10.7,9.6,8.9,7.4", 2400, standardPerBatch,
			customerName, "Công trường mẫu", "Cát vàng",
			s.active, s.dateMod, s.dateMod,
		)
		if err != nil {
			return fmt.Errorf("tạo phiếu cân mẫu %q thất bại: %w", s.code, err)
		}
		ticketID, _ := r.LastInsertId()
		res.TicketsCreated++

		// Tạo các mẻ chi tiết với 16 thành phần (theo định mức/m³ của mác).
		for b := 1; b <= s.batch; b++ {
			d := buildBatchDetail(s.m3me)
			_, err := db.Exec(`
				INSERT INTO weighing_ticket_details
					(weighing_ticket_id, batch_no, batch_volume_m3, mixed_at,
					 tp1, tp2, tp3, tp4, tp5, tp6,
					 xi1, xi2, xi3, xi4, xi5,
					 pg1, pg2, pg3, pg4, pg5, pg6,
					 is_active, created_at, updated_at)
				VALUES (?, ?, ?, datetime('now', ?),
				        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				        ?, datetime('now', ?), datetime('now', ?))`,
				ticketID, b, s.m3me, s.dateMod,
				d.tp1, d.tp2, d.tp3, d.tp4, d.tp5, d.tp6,
				d.xi1, d.xi2, d.xi3, d.xi4, d.xi5,
				d.pg1, d.pg2, d.pg3, d.pg4, d.pg5, d.pg6,
				s.active, s.dateMod, s.dateMod,
			)
			if err != nil {
				return fmt.Errorf("tạo chi tiết mẻ %d của phiếu %q thất bại: %w", b, s.code, err)
			}
			res.DetailsCreated++
		}
	}
	return nil
}

// batchComponents là khối lượng 16 thành phần của một mẻ (kg/lít).
type batchComponents struct {
	tp1, tp2, tp3, tp4, tp5, tp6 float64
	xi1, xi2, xi3, xi4, xi5      float64
	pg1, pg2, pg3, pg4, pg5, pg6 float64
}

// buildBatchDetail tính khối lượng từng thành phần cho một mẻ có khối lượng m3me.
// Định mức minh họa cho 1 m³: cát 750, đá 1100, xi măng 320; có thêm độ ẩm nhỏ.
func buildBatchDetail(m3me float64) batchComponents {
	return batchComponents{
		tp1: 750 * m3me,  // cát vàng
		tp2: 1100 * m3me, // đá dăm
		xi1: 320 * m3me,  // xi măng PCB40
		pg1: 2.5 * m3me,  // phụ gia
	}
}

// gradeRef lấy một mác bê tông đang hoạt động để gắn vào phiếu cân mẫu.
// Trả (0, "", "") nếu chưa có mác nào.
func gradeRef(db *DB) (int64, string, string) {
	var id int64
	var code, name string
	row := db.QueryRow(
		`SELECT id, grade_code, printed_name FROM concrete_grades
		 WHERE is_active = 1 ORDER BY id LIMIT 1`,
	)
	if err := row.Scan(&id, &code, &name); err != nil {
		return 0, "", ""
	}
	return id, code, name
}

// firstID trả về id nhỏ nhất của bảng, 0 nếu bảng trống.
// Tên bảng do code cung cấp (không phải input người dùng) nên an toàn.
func firstID(db *DB, table string) int64 {
	var id int64
	if err := db.QueryRow("SELECT COALESCE(MIN(id), 0) FROM " + table).Scan(&id); err != nil {
		return 0
	}
	return id
}
