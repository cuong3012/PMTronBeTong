package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	// Driver SQLite pure-Go (không cần CGO/C compiler).
	// Đăng ký driver với tên "sqlite".
	_ "modernc.org/sqlite"
)

// DB là kiểu bao (wrapper) quanh *sql.DB để bổ sung tiện ích nếu cần về sau.
type DB struct {
	*sql.DB
}

// Open mở kết nối tới file SQLite tại đường dẫn path.
//
// Hàm thực hiện các bước:
//  1. Tạo thư mục chứa file database nếu chưa tồn tại.
//  2. Mở kết nối với các tham số DSN quan trọng:
//     - _pragma=foreign_keys(1): BẬT kiểm tra khóa ngoại (SQLite mặc định TẮT).
//     - _pragma=journal_mode(WAL): dùng Write-Ahead Logging để tăng khả năng
//       đọc/ghi đồng thời. Đây là cấu hình khuyến nghị cho ứng dụng web.
//     - _pragma=busy_timeout(5000): khi có lock ghi, chờ tối đa 5 giây rồi mới
//       báo lỗi busy, tránh lỗi "database is locked" tức thời.
//     - _pragma=synchronous(NORMAL): cân bằng giữa an toàn dữ liệu và tốc độ,
//       phù hợp khi đã bật WAL.
//  3. Kiểm tra kết nối bằng ping để phát hiện lỗi sớm.
func Open(path string) (*DB, error) {
	// Bước 1: đảm bảo thư mục cha của file database tồn tại.
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("tạo thư mục database %q thất bại: %w", dir, err)
		}
	}

	// Bước 2: ghép DSN với các pragma cần thiết.
	// Mỗi _pragma sẽ được driver thực thi ngay khi mở kết nối.
	dsn := path +
		"?_pragma=foreign_keys(1)" +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=synchronous(NORMAL)"

	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("mở kết nối SQLite thất bại: %w", err)
	}

	// SQLite (với driver pure-Go) không hỗ trợ tốt việc mở nhiều kết nối ghi
	// song song vào cùng một file qua nhiều connection. Giới hạn 1 kết nối
	// ghi tại một thời điểm để tránh lỗi "database is locked".
	// Vì đây là ứng dụng một trạm, mức ghi vừa phải nên cấu hình này phù hợp.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	// Bước 3: ping để xác nhận kết nối hoạt động thực sự.
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("ping SQLite thất bại: %w", err)
	}

	return &DB{DB: sqlDB}, nil
}

// Close đóng kết nối database và giải phóng tài nguyên.
func (d *DB) Close() error {
	if d == nil || d.DB == nil {
		return nil
	}
	return d.DB.Close()
}
