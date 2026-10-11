package auth

import (
	"testing"
	"time"
)

// TestGenerateAndParse xác nhận token sinh ra được parse lại đúng thông tin.
func TestGenerateAndParse(t *testing.T) {
	tm := NewTokenManager("test-secret", time.Hour)

	token, exp, err := tm.Generate(42, "admin", "admin")
	if err != nil {
		t.Fatalf("Generate thất bại: %v", err)
	}
	if token == "" {
		t.Fatal("token rỗng")
	}
	if exp.Before(time.Now()) {
		t.Error("thời điểm hết hạn phải ở tương lai")
	}

	claims, err := tm.Parse(token)
	if err != nil {
		t.Fatalf("Parse thất bại: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("UserID phải là 42, nhận được %d", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("Username phải là admin, nhận được %q", claims.Username)
	}
	if claims.RoleCode != "admin" {
		t.Errorf("RoleCode phải là admin, nhận được %q", claims.RoleCode)
	}
}

// TestParseWrongSecret xác nhận token ký bằng khóa khác không hợp lệ.
func TestParseWrongSecret(t *testing.T) {
	tm1 := NewTokenManager("secret-1", time.Hour)
	tm2 := NewTokenManager("secret-2", time.Hour)

	token, _, err := tm1.Generate(1, "u", "operator")
	if err != nil {
		t.Fatalf("Generate thất bại: %v", err)
	}

	if _, err := tm2.Parse(token); err == nil {
		t.Error("mong đợi lỗi khi parse token với khóa khác, nhưng không có lỗi")
	}
}

// TestParseExpired xác nhận token hết hạn bị từ chối.
func TestParseExpired(t *testing.T) {
	// Thời hạn âm để token hết hạn ngay.
	tm := NewTokenManager("secret", -time.Minute)

	token, _, err := tm.Generate(1, "u", "operator")
	if err != nil {
		t.Fatalf("Generate thất bại: %v", err)
	}

	if _, err := tm.Parse(token); err == nil {
		t.Error("mong đợi lỗi khi parse token hết hạn, nhưng không có lỗi")
	}
}

// TestParseGarbage xác nhận chuỗi không phải token sẽ bị từ chối.
func TestParseGarbage(t *testing.T) {
	tm := NewTokenManager("secret", time.Hour)
	if _, err := tm.Parse("khong-phai-token"); err == nil {
		t.Error("mong đợi lỗi khi parse chuỗi không hợp lệ, nhưng không có lỗi")
	}
}
