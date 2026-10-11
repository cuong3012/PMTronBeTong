package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config chứa toàn bộ cấu hình runtime của backend.
// Giá trị được đọc từ biến môi trường, nếu không có sẽ dùng mặc định.
type Config struct {
	// HTTP
	Port string

	// Database
	DBPath string

	// JWT
	JWTSecret     string
	JWTExpiryHour int

	// CORS: danh sách origin frontend được phép gọi API.
	// Để "*" trong môi trường phát triển; siết chặt khi triển khai.
	AllowedOrigins []string

	// Môi trường: "development" hoặc "production"
	Env string

	// SeedSample: khi true, nạp thêm dữ liệu mẫu (mác bê tông, khách hàng,
	// công trường, đơn hàng, phiếu cân + chi tiết) để demo/thống kê có số liệu.
	// Mặc định false để không làm bẩn dữ liệu production.
	SeedSample bool
}

// Load đọc cấu hình từ biến môi trường và trả về Config.
func Load() *Config {
	return &Config{
		Port:           getEnv("APP_PORT", "8080"),
		DBPath:         getEnv("DB_PATH", "data/tramtron.db"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-me-in-production"),
		JWTExpiryHour:  getEnvInt("JWT_EXPIRY_HOUR", 12),
		AllowedOrigins: []string{getEnv("CORS_ORIGIN", "http://localhost:5173")},
		Env:            getEnv("APP_ENV", "development"),
		SeedSample:     getEnvBool("SEED_SAMPLE", false),
	}
}

// IsProduction cho biết đang chạy ở môi trường production hay không.
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// JWTExpiry trả về thời hạn token dưới dạng time.Duration.
func (c *Config) JWTExpiry() time.Duration {
	return time.Duration(c.JWTExpiryHour) * time.Hour
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// getEnvBool đọc biến môi trường dạng boolean.
// Chấp nhận "1", "true", "yes", "on" (không phân biệt hoa/thường) là true.
func getEnvBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
