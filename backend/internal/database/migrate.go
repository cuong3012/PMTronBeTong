package database

import (
	"database/sql"
	"fmt"
	"log"
)

// migration đại diện cho một bước thay đổi schema có phiên bản.
type migration struct {
	version int    // số phiên bản, tăng dần và không đổi sau khi đã phát hành
	name    string // mô tả ngắn gọn nội dung migration
	sql     string // câu lệnh SQL thực thi (có thể chứa nhiều câu lệnh)
}

// migrations là danh sách TẤT CẢ migration theo thứ tự tăng dần của version.
//
// QUY TẮC QUAN TRỌNG:
//   - Không bao giờ sửa nội dung một migration đã phát hành, vì database đã
//     chạy nó sẽ không chạy lại. Muốn thay đổi -> thêm migration version mới.
//   - Version phải là duy nhất và tăng dần.
var migrations = []migration{
	{
		version: 1,
		name:    "tạo các bảng nền tảng: trạm, vai trò, người dùng, khách hàng",
		sql: `
		-- Bảng trạm trộn
		CREATE TABLE IF NOT EXISTS stations (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			station_code  TEXT    NOT NULL UNIQUE,
			station_name  TEXT    NOT NULL,
			address       TEXT,
			phone         TEXT,
			is_active     INTEGER NOT NULL DEFAULT 1,
			created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
		);

		-- Bảng vai trò người dùng
		CREATE TABLE IF NOT EXISTS roles (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			role_code   TEXT    NOT NULL UNIQUE,
			role_name   TEXT    NOT NULL,
			description TEXT,
			created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
		);

		-- Bảng tài khoản người dùng
		CREATE TABLE IF NOT EXISTS users (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			username      TEXT    NOT NULL UNIQUE,
			password_hash TEXT    NOT NULL,
			full_name     TEXT    NOT NULL,
			phone         TEXT,
			role_id       INTEGER NOT NULL REFERENCES roles(id),
			station_id    INTEGER REFERENCES stations(id),
			is_active     INTEGER NOT NULL DEFAULT 1,
			last_login_at TEXT,
			created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
		);

		-- Bảng khách hàng
		CREATE TABLE IF NOT EXISTS customers (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			customer_code TEXT    NOT NULL UNIQUE,
			customer_name TEXT    NOT NULL,
			address       TEXT,
			phone         TEXT,
			is_active     INTEGER NOT NULL DEFAULT 1,
			created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
		);
		`,
	},
	{
		version: 2,
		name:    "tạo bảng danh mục vận tải và công trường",
		sql: `
		-- Bảng xe bồn
		CREATE TABLE IF NOT EXISTS mixer_trucks (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			truck_code   TEXT    NOT NULL UNIQUE,
			plate_number TEXT    NOT NULL UNIQUE,
			is_active    INTEGER NOT NULL DEFAULT 1,
			created_at   TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at   TEXT    NOT NULL DEFAULT (datetime('now'))
		);

		-- Bảng tài xế (không bắt buộc có tài khoản đăng nhập)
		CREATE TABLE IF NOT EXISTS drivers (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			driver_code TEXT    NOT NULL UNIQUE,
			full_name   TEXT    NOT NULL,
			is_active   INTEGER NOT NULL DEFAULT 1,
			created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
		);

		-- Bảng công trường / địa điểm thi công
		CREATE TABLE IF NOT EXISTS construction_sites (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			site_code  TEXT    NOT NULL UNIQUE,
			site_name  TEXT    NOT NULL,
			is_active  INTEGER NOT NULL DEFAULT 1,
			created_at TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT    NOT NULL DEFAULT (datetime('now'))
		);
		`,
	},
	{
		version: 3,
		name:    "tạo bảng đơn hàng và chuyến giao hàng",
		sql: `
		-- Bảng đơn hàng
		CREATE TABLE IF NOT EXISTS orders (
			id                   INTEGER PRIMARY KEY AUTOINCREMENT,
			order_number         TEXT    NOT NULL UNIQUE,
			customer_id          INTEGER NOT NULL REFERENCES customers(id),
			construction_site_id INTEGER NOT NULL REFERENCES construction_sites(id),
			ordered_volume_m3    REAL    NOT NULL,
			cumulative_volume_m3 REAL    NOT NULL DEFAULT 0,
			is_active            INTEGER NOT NULL DEFAULT 1,
			notes                TEXT,
			created_at           TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at           TEXT    NOT NULL DEFAULT (datetime('now'))
		);

		-- Bảng chuyến giao hàng (một đơn có thể nhiều chuyến)
		CREATE TABLE IF NOT EXISTS order_deliveries (
			id                   INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id             INTEGER NOT NULL REFERENCES orders(id),
			trip_no              INTEGER NOT NULL,
			truck_id             INTEGER REFERENCES mixer_trucks(id),
			driver_id            INTEGER REFERENCES drivers(id),
			planned_volume_m3    REAL    NOT NULL,
			loaded_volume_m3     REAL,
			delivered_volume_m3  REAL,
			status               TEXT    NOT NULL DEFAULT 'planned',
			departed_at          TEXT,
			delivered_at         TEXT,
			notes                TEXT,
			created_at           TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at           TEXT    NOT NULL DEFAULT (datetime('now')),
			-- Mỗi đơn không được có hai chuyến cùng số thứ tự
			UNIQUE (order_id, trip_no)
		);

		-- Chỉ mục giúp tra cứu chuyến giao theo đơn nhanh hơn
		CREATE INDEX IF NOT EXISTS idx_order_deliveries_order_id
			ON order_deliveries(order_id);
		`,
	},
	{
		version: 4,
		name:    "tạo bảng audit log cho các thao tác quan trọng",
		sql: `
		-- Nhật ký hoạt động (mục 5.1 yêu cầu audit log)
		CREATE TABLE IF NOT EXISTS audit_logs (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id    INTEGER REFERENCES users(id),
			action     TEXT    NOT NULL,
			entity     TEXT    NOT NULL,
			entity_id  INTEGER,
			detail     TEXT,
			created_at TEXT    NOT NULL DEFAULT (datetime('now'))
		);

		CREATE INDEX IF NOT EXISTS idx_audit_logs_entity
			ON audit_logs(entity, entity_id);
		`,
	},
	{
		version: 5,
		name:    "bổ sung tài xế mặc định cho xe bồn",
		sql: `
		-- Tài xế mặc định của xe (mục 13.1). Giao diện tự điền khi điều phối,
		-- nhưng có thể để trống nếu xe chưa được gán tài xế.
		-- Lưu ý: SQLite không cho thêm ràng buộc FOREIGN KEY qua ALTER TABLE,
		-- nên tính hợp lệ của default_driver_id được kiểm tra ở tầng ứng dụng.
		ALTER TABLE mixer_trucks ADD COLUMN default_driver_id INTEGER;

		CREATE INDEX IF NOT EXISTS idx_mixer_trucks_default_driver
			ON mixer_trucks(default_driver_id);
		`,
	},
}

// Migrate chạy tất cả migration chưa được áp dụng lên database.
//
// Cơ chế:
//  1. Tạo bảng schema_migrations để ghi nhớ version nào đã chạy.
//  2. Đọc version lớn nhất hiện có trong database.
//  3. Với mỗi migration có version lớn hơn version đã chạy:
//     - Chạy trong một transaction để đảm bảo nguyên tử.
//     - Ghi bản ghi version vào schema_migrations.
//
// Trả về lỗi nếu bất kỳ bước nào thất bại.
func Migrate(db *DB) error {
	// Bước 1: bảng theo dõi các migration đã áp dụng.
	const createTable = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    INTEGER PRIMARY KEY,
			name       TEXT    NOT NULL,
			applied_at TEXT    NOT NULL DEFAULT (datetime('now'))
		);`
	if _, err := db.Exec(createTable); err != nil {
		return fmt.Errorf("tạo bảng schema_migrations thất bại: %w", err)
	}

	// Bước 2: xác định version lớn nhất đã chạy.
	current, err := currentVersion(db)
	if err != nil {
		return err
	}

	applied := 0
	// Bước 3: áp dụng tuần tự các migration còn thiếu.
	for _, m := range migrations {
		if m.version <= current {
			continue // đã chạy trước đó, bỏ qua
		}
		if err := applyMigration(db, m); err != nil {
			return err
		}
		log.Printf("[migrate] đã áp dụng version %d: %s", m.version, m.name)
		applied++
	}

	if applied == 0 {
		log.Printf("[migrate] database đã ở version mới nhất (%d), không cần chạy thêm", current)
	} else {
		log.Printf("[migrate] hoàn tất, đã áp dụng %d migration", applied)
	}
	return nil
}

// currentVersion trả về version migration lớn nhất đã áp dụng.
// Nếu chưa có migration nào, trả về 0.
func currentVersion(db *DB) (int, error) {
	var version sql.NullInt64
	row := db.QueryRow(`SELECT MAX(version) FROM schema_migrations`)
	if err := row.Scan(&version); err != nil {
		return 0, fmt.Errorf("đọc version hiện tại thất bại: %w", err)
	}
	if !version.Valid {
		return 0, nil
	}
	return int(version.Int64), nil
}

// applyMigration thực thi một migration và ghi lại lịch sử trong cùng transaction.
func applyMigration(db *DB, m migration) error {
	// Bắt đầu transaction: nếu SQL lỗi, toàn bộ thay đổi được rollback.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("bắt đầu transaction cho migration %d thất bại: %w", m.version, err)
	}

	// Rollback phòng trường hợp thoát sớm (khi đã Commit thành công thì
	// Rollback sau đó là vô hại, trả về lỗi ErrTxDone và bị bỏ qua).
	defer func() { _ = tx.Rollback() }()

	// Chạy nội dung SQL của migration.
	if _, err := tx.Exec(m.sql); err != nil {
		return fmt.Errorf("thực thi migration %d (%s) thất bại: %w", m.version, m.name, err)
	}

	// Ghi lại version đã áp dụng.
	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (version, name) VALUES (?, ?)`,
		m.version, m.name,
	); err != nil {
		return fmt.Errorf("ghi lịch sử migration %d thất bại: %w", m.version, err)
	}

	// Xác nhận transaction.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %d thất bại: %w", m.version, err)
	}
	return nil
}
