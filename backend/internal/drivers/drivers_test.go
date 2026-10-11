package drivers

import (
	"errors"
	"path/filepath"
	"testing"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
)

func setupService(t *testing.T) *Service {
	t.Helper()
	path := filepath.Join(t.TempDir(), "drivers.db")
	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("mo database test that bai: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := database.EnsureSchemaAndSeed(db); err != nil {
		t.Fatalf("seed database test that bai: %v", err)
	}
	return NewService(NewRepository(db))
}

func TestCreateAndGet(t *testing.T) {
	svc := setupService(t)
	created, err := svc.Create(CreateInput{DriverCode: "TX001", FullName: "Nguyen Van A"})
	if err != nil {
		t.Fatalf("Create that bai: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("id tai xe phai duoc gan sau khi tao")
	}
	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatalf("Get that bai: %v", err)
	}
	if got.FullName != "Nguyen Van A" {
		t.Errorf("ho ten tai xe sai, nhan duoc %q", got.FullName)
	}
}

func TestDuplicateDriverCode(t *testing.T) {
	svc := setupService(t)
	first, err := svc.Create(CreateInput{DriverCode: "TX001", FullName: "A"})
	if err != nil || first.ID == 0 {
		t.Fatalf("Create ban ghi dau tien that bai: %v", err)
	}
	_, err = svc.Create(CreateInput{DriverCode: "TX001", FullName: "B"})
	if !errors.Is(err, ErrCodeExists) {
		t.Errorf("mong doi ErrCodeExists, nhan duoc %v", err)
	}
}

func TestSoftDeleteHidesFromList(t *testing.T) {
	svc := setupService(t)
	created, _ := svc.Create(CreateInput{DriverCode: "TX001", FullName: "A"})
	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("Delete that bai: %v", err)
	}
	result, _ := svc.List(httpx.Pagination{Page: 1, PerPage: 20})
	if result.Total != 0 {
		t.Errorf("sau xoa mem phai rong, nhan duoc %d", result.Total)
	}
}
