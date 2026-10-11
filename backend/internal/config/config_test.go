package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	// Đảm bảo các biến môi trường không ảnh hưởng tới test này.
	t.Setenv("APP_PORT", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("JWT_SECRET", "")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Port mặc định phải là 8080, nhận được %q", cfg.Port)
	}
	if cfg.DBPath != "data/tramtron.db" {
		t.Errorf("DBPath mặc định phải là data/tramtron.db, nhận được %q", cfg.DBPath)
	}
	if cfg.JWTExpiryHour != 12 {
		t.Errorf("JWTExpiryHour mặc định phải là 12, nhận được %d", cfg.JWTExpiryHour)
	}
	if cfg.IsProduction() {
		t.Error("Mặc định không phải production")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	t.Setenv("JWT_SECRET", "xyz")
	t.Setenv("JWT_EXPIRY_HOUR", "1")
	t.Setenv("APP_ENV", "production")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Port phải là 9090, nhận được %q", cfg.Port)
	}
	if cfg.JWTSecret != "xyz" {
		t.Errorf("JWTSecret phải là xyz, nhận được %q", cfg.JWTSecret)
	}
	if cfg.JWTExpiry().Hours() != 1 {
		t.Errorf("JWTExpiry phải là 1 giờ, nhận được %v", cfg.JWTExpiry())
	}
	if !cfg.IsProduction() {
		t.Error("Phải nhận diện production")
	}
}
