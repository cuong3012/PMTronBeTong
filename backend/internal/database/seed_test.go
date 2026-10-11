package database

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// TestSeedDefaults chạy seed trên database mới và kiểm tra dữ liệu được tạo.
func TestSeedDefaults(t *testing.T) {
	db := openTestDB(t)
	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate thất bại: %v", err)
	}

	result, err := SeedDefaults(db)
	if err != nil {
		t.Fatalf("SeedDefaults thất bại: %v", err)
	}

	if result.RolesCreated != len(defaultRoles) {
		t.Errorf("số vai trò tạo phải là %d, nhận được %d", len(defaultRoles), result.RolesCreated)
	}
	if result.StationsCreated != 1 {
		t.Errorf("số trạm tạo phải là 1, nhận được %d", result.StationsCreated)
	}
	if !result.AdminCreated {
		t.Error("tài khoản admin phải được tạo")
	}

	// Kiểm tra đúng số bản ghi trong bảng.
	assertCount(t, db, "roles", len(defaultRoles))
	assertCount(t, db, "stations", 1)
	assertCount(t, db, "users", 1)
}

// TestSeedIsIdempotent xác nhận chạy seed lần hai không tạo thêm bản ghi.
func TestSeedIsIdempotent(t *testing.T) {
	db := openTestDB(t)
	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate thất bại: %v", err)
	}

	if _, err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults lần 1 thất bại: %v", err)
	}

	result2, err := SeedDefaults(db)
	if err != nil {
		t.Fatalf("SeedDefaults lần 2 thất bại: %v", err)
	}
	if result2.RolesCreated != 0 || result2.StationsCreated != 0 || result2.AdminCreated {
		t.Errorf("lần seed thứ hai không được tạo thêm bản ghi, nhận được: %+v", result2)
	}

	// Số bản ghi vẫn không đổi.
	assertCount(t, db, "roles", len(defaultRoles))
	assertCount(t, db, "users", 1)
}

// TestAdminPasswordIsHashed xác nhận mật khẩu admin được lưu dưới dạng băm bcrypt
// và có thể so khớp với mật khẩu gốc "Admin@123".
func TestAdminPasswordIsHashed(t *testing.T) {
	db := openTestDB(t)
	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate thất bại: %v", err)
	}
	if _, err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults thất bại: %v", err)
	}

	var hash string
	if err := db.QueryRow(`SELECT password_hash FROM users WHERE username = ?`, "admin").Scan(&hash); err != nil {
		t.Fatalf("đọc password_hash thất bại: %v", err)
	}

	// Mật khẩu gốc không được lưu trực tiếp.
	if hash == "Admin@123" {
		t.Fatal("mật khẩu admin bị lưu dạng thô, không an toàn")
	}

	// Băm phải khớp với mật khẩu gốc.
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("Admin@123")); err != nil {
		t.Errorf("băm mật khẩu không khớp mật khẩu gốc: %v", err)
	}
}

// assertCount kiểm tra số bản ghi của một bảng.
func assertCount(t *testing.T, db *DB, table string, want int) {
	t.Helper()
	var got int
	// Lưu ý: tên bảng không thể truyền qua placeholder, nhưng ở đây là hằng
	// nội bộ trong test nên an toàn.
	if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&got); err != nil {
		t.Fatalf("đếm bảng %q thất bại: %v", table, err)
	}
	if got != want {
		t.Errorf("bảng %q phải có %d bản ghi, nhận được %d", table, want, got)
	}
}
