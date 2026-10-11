package auth

import (
	"path/filepath"
	"testing"
	"time"

	"tramtron/backend/internal/database"
)

// setupService tạo database tạm, migrate + seed, rồi trả về Service để test.
func setupService(t *testing.T) *Service {
	t.Helper()
	path := filepath.Join(t.TempDir(), "auth.db")
	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("mở database test thất bại: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := database.EnsureSchemaAndSeed(db); err != nil {
		t.Fatalf("seed database test thất bại: %v", err)
	}

	repo := NewRepository(db)
	tm := NewTokenManager("test-secret", time.Hour)
	return NewService(repo, tm)
}

// TestLoginSuccess xác nhận đăng nhập đúng trả về token và thông tin người dùng.
func TestLoginSuccess(t *testing.T) {
	svc := setupService(t)

	result, err := svc.Login("admin", "Admin@123")
	if err != nil {
		t.Fatalf("đăng nhập admin thất bại: %v", err)
	}
	if result.Token == "" {
		t.Error("token phải không rỗng")
	}
	if result.User == nil {
		t.Fatal("thông tin người dùng phải có")
	}
	if result.User.Username != "admin" {
		t.Errorf("username phải là admin, nhận được %q", result.User.Username)
	}
	if result.User.RoleCode != "admin" {
		t.Errorf("role phải là admin, nhận được %q", result.User.RoleCode)
	}
}

// TestLoginWrongPassword xác nhận sai mật khẩu bị từ chối.
func TestLoginWrongPassword(t *testing.T) {
	svc := setupService(t)

	_, err := svc.Login("admin", "sai-mat-khau")
	if err != ErrInvalidCredentials {
		t.Errorf("mong đợi ErrInvalidCredentials, nhận được %v", err)
	}
}

// TestLoginUnknownUser xác nhận username không tồn tại trả về lỗi chung.
func TestLoginUnknownUser(t *testing.T) {
	svc := setupService(t)

	_, err := svc.Login("khong-ton-tai", "bat-ky")
	if err != ErrInvalidCredentials {
		t.Errorf("mong đợi ErrInvalidCredentials, nhận được %v", err)
	}
}

// TestMe xác nhận lấy được thông tin người dùng theo id.
func TestMe(t *testing.T) {
	svc := setupService(t)

	result, err := svc.Login("admin", "Admin@123")
	if err != nil {
		t.Fatalf("đăng nhập thất bại: %v", err)
	}

	info, err := svc.Me(result.User.ID)
	if err != nil {
		t.Fatalf("Me thất bại: %v", err)
	}
	if info.Username != "admin" {
		t.Errorf("username phải là admin, nhận được %q", info.Username)
	}
}
