package database

import (
	"path/filepath"
	"testing"
)

// openTestDB mở một database SQLite tạm thời trong thư mục test.
// Mỗi test dùng file riêng để không ảnh hưởng lẫn nhau.
func openTestDB(t *testing.T) *DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(path)
	if err != nil {
		t.Fatalf("không mở được database test: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// tableExists kiểm tra một bảng có tồn tại trong database hay không.
func tableExists(t *testing.T, db *DB, name string) bool {
	t.Helper()
	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, name,
	).Scan(&count)
	if err != nil {
		t.Fatalf("kiểm tra bảng %q thất bại: %v", name, err)
	}
	return count > 0
}

// TestMigrateCreatesAllTables chạy migration và xác nhận tất cả bảng được tạo.
func TestMigrateCreatesAllTables(t *testing.T) {
	db := openTestDB(t)

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate thất bại: %v", err)
	}

	expected := []string{
		"stations", "roles", "users", "customers",
		"mixer_trucks", "drivers", "construction_sites",
		"orders", "order_deliveries",
		"audit_logs", "schema_migrations",
	}
	for _, name := range expected {
		if !tableExists(t, db, name) {
			t.Errorf("bảng %q chưa được tạo", name)
		}
	}
}

// TestMigrateIsIdempotent xác nhận chạy migration nhiều lần không gây lỗi
// và không áp dụng trùng version.
func TestMigrateIsIdempotent(t *testing.T) {
	db := openTestDB(t)

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate lần 1 thất bại: %v", err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate lần 2 (idempotent) thất bại: %v", err)
	}

	// Số bản ghi trong schema_migrations phải bằng đúng số migration.
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatalf("đếm schema_migrations thất bại: %v", err)
	}
	if count != len(migrations) {
		t.Errorf("số migration đã ghi phải là %d, nhận được %d", len(migrations), count)
	}
}

// TestCurrentVersion kiểm tra hàm đọc version hiện tại.
func TestCurrentVersion(t *testing.T) {
	db := openTestDB(t)

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate thất bại: %v", err)
	}

	v, err := currentVersion(db)
	if err != nil {
		t.Fatalf("currentVersion thất bại: %v", err)
	}
	if v != len(migrations) {
		t.Errorf("version hiện tại phải là %d, nhận được %d", len(migrations), v)
	}
}

// TestForeignKeyEnabled xác nhận pragma foreign_keys đã được bật.
func TestForeignKeyEnabled(t *testing.T) {
	db := openTestDB(t)

	var fk int
	if err := db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk); err != nil {
		t.Fatalf("đọc PRAGMA foreign_keys thất bại: %v", err)
	}
	if fk != 1 {
		t.Errorf("foreign_keys phải bật (=1), nhận được %d", fk)
	}
}

// TestUniqueConstraint xác nhận ràng buộc UNIQUE được áp dụng (mã khách hàng).
func TestUniqueConstraint(t *testing.T) {
	db := openTestDB(t)
	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate thất bại: %v", err)
	}

	insert := `INSERT INTO customers (customer_code, customer_name) VALUES (?, ?)`
	if _, err := db.Exec(insert, "KH001", "Khách A"); err != nil {
		t.Fatalf("chèn khách hàng đầu tiên thất bại: %v", err)
	}
	// Chèn trùng mã phải lỗi.
	if _, err := db.Exec(insert, "KH001", "Khách B"); err == nil {
		t.Error("mong đợi lỗi UNIQUE khi trùng customer_code, nhưng không có lỗi")
	}
}
