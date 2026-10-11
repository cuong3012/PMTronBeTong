package customers

import (
	"errors"
	"path/filepath"
	"testing"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
)

// setupService tạo database tạm + migrate + seed, trả về Service để test.
func setupService(t *testing.T) *Service {
	t.Helper()
	path := filepath.Join(t.TempDir(), "customers.db")
	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("mở database test thất bại: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := database.EnsureSchemaAndSeed(db); err != nil {
		t.Fatalf("seed database test thất bại: %v", err)
	}

	return NewService(NewRepository(db))
}

// TestCreateAndGet tạo khách hàng rồi lấy lại theo id.
func TestCreateAndGet(t *testing.T) {
	svc := setupService(t)

	created, err := svc.Create(CreateInput{
		CustomerCode: "KH001",
		CustomerName: "Công ty Xây dựng ABC",
		Address:      "Hà Nội",
		Phone:        "0900000001",
	})
	if err != nil {
		t.Fatalf("Create thất bại: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("id khách hàng phải được gán")
	}
	if !created.IsActive {
		t.Error("khách hàng mới phải có is_active = true")
	}

	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatalf("Get thất bại: %v", err)
	}
	if got.CustomerCode != "KH001" {
		t.Errorf("mã khách hàng phải là KH001, nhận được %q", got.CustomerCode)
	}
}

// TestCreateDuplicateCode xác nhận không cho trùng mã khách hàng.
func TestCreateDuplicateCode(t *testing.T) {
	svc := setupService(t)

	in := CreateInput{CustomerCode: "KH001", CustomerName: "Khách A"}
	if _, err := svc.Create(in); err != nil {
		t.Fatalf("Create lần 1 thất bại: %v", err)
	}

	_, err := svc.Create(CreateInput{CustomerCode: "KH001", CustomerName: "Khách B"})
	if !errors.Is(err, ErrCodeExists) {
		t.Errorf("mong đợi ErrCodeExists, nhận được %v", err)
	}
}

// TestUpdate xác nhận cập nhật thay đổi đúng dữ liệu.
func TestUpdate(t *testing.T) {
	svc := setupService(t)

	created, err := svc.Create(CreateInput{CustomerCode: "KH001", CustomerName: "Khách A"})
	if err != nil {
		t.Fatalf("Create thất bại: %v", err)
	}

	updated, err := svc.Update(created.ID, UpdateInput{
		CustomerCode: "KH001",
		CustomerName: "Khách A (đổi tên)",
		Phone:        "0912345678",
	})
	if err != nil {
		t.Fatalf("Update thất bại: %v", err)
	}
	if updated.CustomerName != "Khách A (đổi tên)" {
		t.Errorf("tên chưa được cập nhật, nhận được %q", updated.CustomerName)
	}
	if updated.Phone != "0912345678" {
		t.Errorf("số điện thoại chưa được cập nhật, nhận được %q", updated.Phone)
	}
}

// TestSoftDelete xác nhận xóa là xóa mềm: bản ghi vẫn tồn tại nhưng ẩn khỏi danh sách.
func TestSoftDelete(t *testing.T) {
	svc := setupService(t)

	created, err := svc.Create(CreateInput{CustomerCode: "KH001", CustomerName: "Khách A"})
	if err != nil {
		t.Fatalf("Create thất bại: %v", err)
	}

	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("Delete thất bại: %v", err)
	}

	// Bản ghi vẫn lấy được theo id nhưng is_active = false.
	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatalf("Get sau xóa mềm thất bại: %v", err)
	}
	if got.IsActive {
		t.Error("sau xóa mềm, is_active phải là false")
	}

	// Danh sách thông thường không còn khách hàng này.
	result, err := svc.List(httpx.Pagination{Page: 1, PerPage: 20})
	if err != nil {
		t.Fatalf("List thất bại: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("danh sách phải rỗng sau xóa mềm, nhận được total = %d", result.Total)
	}
}

// TestListSearchAndPaging kiểm tra tìm kiếm và phân trang.
func TestListSearchAndPaging(t *testing.T) {
	svc := setupService(t)

	for _, c := range []CreateInput{
		{CustomerCode: "KH001", CustomerName: "Công ty ABC"},
		{CustomerCode: "KH002", CustomerName: "Công ty XYZ"},
		{CustomerCode: "KH003", CustomerName: "ABC Xây dựng"},
	} {
		if _, err := svc.Create(c); err != nil {
			t.Fatalf("Create thất bại: %v", err)
		}
	}

	// Tìm kiếm "ABC" phải khớp 2 bản ghi.
	all, err := svc.List(httpx.Pagination{Page: 1, PerPage: 20, Search: "ABC"})
	if err != nil {
		t.Fatalf("List all thất bại: %v", err)
	}
	if all.Total != 2 {
		t.Errorf("tìm 'ABC' phải khớp 2 bản ghi, nhận được %d", all.Total)
	}

	// Phân trang: mỗi trang 2 bản ghi.
	page1, err := svc.List(httpx.Pagination{Page: 1, PerPage: 2})
	if err != nil {
		t.Fatalf("List page 1 thất bại: %v", err)
	}
	if page1.Total != 3 {
		t.Errorf("tổng phải là 3, nhận được %d", page1.Total)
	}
	items, _ := page1.Items.([]interface{})
	if len(items) > 2 {
		t.Errorf("trang 1 tối đa 2 bản ghi, nhận được %d", len(items))
	}
}

// TestGetNotFound xác nhận lấy id không tồn tại trả lỗi.
func TestGetNotFound(t *testing.T) {
	svc := setupService(t)
	if _, err := svc.Get(9999); !errors.Is(err, ErrNotFound) {
		t.Errorf("mong đợi ErrNotFound, nhận được %v", err)
	}
}
