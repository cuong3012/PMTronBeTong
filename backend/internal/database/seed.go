package database

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Vai trò khởi đầu của hệ thống (mục 13.1 — bảng roles).
const (
	RoleAdmin      = "admin"      // Quản trị hệ thống
	RoleManager    = "manager"    // Quản lý trạm
	RoleDispatcher = "dispatcher" // Điều phối viên đơn hàng
	RoleOperator   = "operator"   // Nhân viên vận hành
)

// seedRole mô tả một vai trò cần chèn.
type seedRole struct {
	code        string
	name        string
	description string
}

// defaultRoles là danh sách vai trò mặc định, được tạo khi database trống.
var defaultRoles = []seedRole{
	{RoleAdmin, "Quản trị hệ thống", "Quản lý người dùng, phân quyền và cấu hình hệ thống"},
	{RoleManager, "Quản lý trạm", "Theo dõi sản xuất, báo cáo và điều phối"},
	{RoleDispatcher, "Điều phối viên", "Tạo đơn hàng, cập nhật tiến độ, kiểm tra lịch giao"},
	{RoleOperator, "Nhân viên vận hành", "Theo dõi máy móc, công thức và số liệu sản xuất"},
}

// Seed kết quả trả về sau khi chèn dữ liệu khởi đầu, phục vụ log và test.
type Seed struct {
	RolesCreated     int
	StationsCreated  int
	AdminCreated     bool
	AdminUsername    string
	AdminPassword    string // chỉ dùng để log ở môi trường phát triển
}

// SeedDefaults chèn dữ liệu khởi đầu (vai trò, trạm, tài khoản admin) nếu chưa có.
//
// Hàm này an toàn khi chạy lại nhiều lần (idempotent):
//   - Vai trò: dùng INSERT ... ON CONFLICT(role_code) DO NOTHING.
//   - Trạm: dùng ON CONFLICT(station_code) DO NOTHING.
//   - Tài khoản admin: chỉ tạo nếu chưa có username "admin".
//
// Thông tin admin mặc định:
//   - username: admin
//   - password: Admin@123 (CẦN đổi ngay sau lần đăng nhập đầu tiên)
func SeedDefaults(db *DB) (*Seed, error) {
	result := &Seed{}

	// 1. Chèn các vai trò mặc định.
	for _, r := range defaultRoles {
		res, err := db.Exec(
			`INSERT INTO roles (role_code, role_name, description)
			 VALUES (?, ?, ?)
			 ON CONFLICT(role_code) DO NOTHING`,
			r.code, r.name, r.description,
		)
		if err != nil {
			return nil, fmt.Errorf("chèn vai trò %q thất bại: %w", r.code, err)
		}
		if n, _ := res.RowsAffected(); n > 0 {
			result.RolesCreated++
		}
	}

	// 2. Chèn trạm mẫu nếu chưa có.
	res, err := db.Exec(
		`INSERT INTO stations (station_code, station_name, address)
		 VALUES (?, ?, ?)
		 ON CONFLICT(station_code) DO NOTHING`,
		"TRAM01", "Trạm trộn số 1", "Khu công nghiệp",
	)
	if err != nil {
		return nil, fmt.Errorf("chèn trạm mẫu thất bại: %w", err)
	}
	if n, _ := res.RowsAffected(); n > 0 {
		result.StationsCreated++
	}

	// 3. Chèn tài khoản admin nếu chưa tồn tại.
	var exists int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, "admin").Scan(&exists); err != nil {
		return nil, fmt.Errorf("kiểm tra tài khoản admin thất bại: %w", err)
	}

	if exists == 0 {
		const adminUsername = "admin"
		const adminPassword = "Admin@123"

		hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("băm mật khẩu admin thất bại: %w", err)
		}

		// Lấy id vai trò admin vừa chèn.
		var roleID int64
		if err := db.QueryRow(`SELECT id FROM roles WHERE role_code = ?`, RoleAdmin).Scan(&roleID); err != nil {
			return nil, fmt.Errorf("không tìm thấy vai trò admin: %w", err)
		}

		if _, err := db.Exec(
			`INSERT INTO users (username, password_hash, full_name, role_id, station_id)
			 VALUES (?, ?, ?, ?, NULL)`,
			adminUsername, string(hash), "Quản trị viên", roleID,
		); err != nil {
			return nil, fmt.Errorf("tạo tài khoản admin thất bại: %w", err)
		}

		result.AdminCreated = true
		result.AdminUsername = adminUsername
		result.AdminPassword = adminPassword
	}

	return result, nil
}

// EnsureSchemaAndSeed là tiện ích gộp: chạy migration rồi seed dữ liệu mặc định.
// Được gọi khi khởi động server để đảm bảo database luôn sẵn sàng.
func EnsureSchemaAndSeed(db *DB) error {
	if err := Migrate(db); err != nil {
		return err
	}

	result, err := SeedDefaults(db)
	if err != nil {
		return err
	}

	log.Printf("[seed] vai trò mới: %d, trạm mới: %d, tạo admin: %v",
		result.RolesCreated, result.StationsCreated, result.AdminCreated)
	if result.AdminCreated {
		log.Printf("[seed] tài khoản admin mặc định: %s / %s (hãy đổi mật khẩu sau khi đăng nhập)",
			result.AdminUsername, result.AdminPassword)
	}
	return nil
}
