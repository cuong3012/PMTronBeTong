package sites

import (
	"errors"
	"path/filepath"
	"testing"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
)

func setupService(t *testing.T) *Service {
	t.Helper()
	path := filepath.Join(t.TempDir(), "sites.db")
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
	created, err := svc.Create(CreateInput{SiteCode: "CT001", SiteName: "Cong truong ABC"})
	if err != nil {
		t.Fatalf("Create that bai: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("id cong truong phai duoc gan sau khi tao")
	}
	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatalf("Get that bai: %v", err)
	}
	if got.SiteName != "Cong truong ABC" {
		t.Errorf("ten cong truong sai, nhan duoc %q", got.SiteName)
	}
}

func TestDuplicateCode(t *testing.T) {
	svc := setupService(t)
	if _, err := svc.Create(CreateInput{SiteCode: "CT001", SiteName: "A"}); err != nil {
		t.Fatalf("Create lan 1 that bai: %v", err)
	}
	_, err := svc.Create(CreateInput{SiteCode: "CT001", SiteName: "B"})
	if !errors.Is(err, ErrCodeExists) {
		t.Errorf("mong doi ErrCodeExists, nhan duoc %v", err)
	}
}

func TestSoftDeleteHidesFromList(t *testing.T) {
	svc := setupService(t)
	created, _ := svc.Create(CreateInput{SiteCode: "CT001", SiteName: "A"})
	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("Delete that bai: %v", err)
	}
	result, err := svc.List(httpx.Pagination{Page: 1, PerPage: 20})
	if err != nil {
		t.Fatalf("List that bai: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("sau xoa mem phai rong, nhan duoc %d", result.Total)
	}
}
