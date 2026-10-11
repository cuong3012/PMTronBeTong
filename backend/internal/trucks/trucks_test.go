package trucks

import (
	"errors"
	"path/filepath"
	"testing"

	"tramtron/backend/internal/database"
	"tramtron/backend/internal/httpx"
)

// setup trả về Service và DB để test có thể chèn dữ liệu phụ trợ (tài xế).
func setup(t *testing.T) (*Service, *database.DB) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "trucks.db")
	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("mo database test that bai: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := database.EnsureSchemaAndSeed(db); err != nil {
		t.Fatalf("seed database test that bai: %v", err)
	}
	return NewService(NewRepository(db)), db
}

// insertDriver chèn tài xế trực tiếp và trả về id.
func insertDriver(t *testing.T, db *database.DB, code, name string, active int) int64 {
	t.Helper()
	res, err := db.Exec(
		`INSERT INTO drivers (driver_code, full_name, is_active) VALUES (?, ?, ?)`,
		code, name, active,
	)
	if err != nil {
		t.Fatalf("chen tai xe that bai: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("lay id tai xe that bai: %v", err)
	}
	return id
}

func TestCreateNormalizesPlate(t *testing.T) {
	svc, _ := setup(t)
	created, err := svc.Create(CreateInput{TruckCode: "XE001", PlateNumber: "29c-12345"})
	if err != nil {
		t.Fatalf("Create that bai: %v", err)
	}
	if created.PlateNumber != "29C-12345" {
		t.Errorf("bien so phai duoc viet hoa, nhan duoc %q", created.PlateNumber)
	}
}

func TestDuplicateTruckCode(t *testing.T) {
	svc, _ := setup(t)
	if _, err := svc.Create(CreateInput{TruckCode: "XE001", PlateNumber: "29C-11111"}); err != nil {
		t.Fatalf("Create lan 1 that bai: %v", err)
	}
	_, err := svc.Create(CreateInput{TruckCode: "XE001", PlateNumber: "29C-22222"})
	if !errors.Is(err, ErrCodeExists) {
		t.Errorf("mong doi ErrCodeExists khi trung ma xe, nhan duoc %v", err)
	}
}

func TestDuplicatePlateNumber(t *testing.T) {
	svc, _ := setup(t)
	first, err := svc.Create(CreateInput{TruckCode: "XE001", PlateNumber: "29C-11111"})
	if err != nil || first.ID == 0 {
		t.Fatalf("Create ban ghi dau tien that bai: %v", err)
	}
	_, err = svc.Create(CreateInput{TruckCode: "XE002", PlateNumber: "29C-11111"})
	if !errors.Is(err, ErrCodeExists) {
		t.Errorf("mong doi ErrCodeExists khi trung bien so, nhan duoc %v", err)
	}
}

func TestSoftDeleteHidesFromList(t *testing.T) {
	svc, _ := setup(t)
	created, _ := svc.Create(CreateInput{TruckCode: "XE001", PlateNumber: "29C-11111"})
	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("Delete that bai: %v", err)
	}
	result, _ := svc.List(httpx.Pagination{Page: 1, PerPage: 20})
	if result.Total != 0 {
		t.Errorf("sau xoa mem phai rong, nhan duoc %d", result.Total)
	}
}

func TestCreateWithDefaultDriver(t *testing.T) {
	svc, db := setup(t)
	driverID := insertDriver(t, db, "TX001", "Nguyen Van A", 1)

	created, err := svc.Create(CreateInput{
		TruckCode:       "XE001",
		PlateNumber:     "29C-11111",
		DefaultDriverID: &driverID,
	})
	if err != nil {
		t.Fatalf("Create voi tai xe mac dinh that bai: %v", err)
	}
	if created.DefaultDriverID == nil || *created.DefaultDriverID != driverID {
		t.Errorf("tai xe mac dinh phai la %d, nhan duoc %v", driverID, created.DefaultDriverID)
	}
	if created.DefaultDriverName != "Nguyen Van A" {
		t.Errorf("ten tai xe mac dinh sai, nhan duoc %q", created.DefaultDriverName)
	}
}

func TestCreateWithUnknownDriver(t *testing.T) {
	svc, _ := setup(t)
	missing := int64(9999)
	_, err := svc.Create(CreateInput{
		TruckCode:       "XE001",
		PlateNumber:     "29C-11111",
		DefaultDriverID: &missing,
	})
	if !errors.Is(err, ErrDriverNotFound) {
		t.Errorf("mong doi ErrDriverNotFound, nhan duoc %v", err)
	}
}

func TestCreateWithInactiveDriver(t *testing.T) {
	svc, db := setup(t)
	driverID := insertDriver(t, db, "TX001", "Da Nghi", 0)

	_, err := svc.Create(CreateInput{
		TruckCode:       "XE001",
		PlateNumber:     "29C-11111",
		DefaultDriverID: &driverID,
	})
	if !errors.Is(err, ErrDriverNotFound) {
		t.Errorf("tai xe ngung hoat dong phai bi tu choi, nhan duoc %v", err)
	}
}

func TestUpdateDefaultDriver(t *testing.T) {
	svc, db := setup(t)
	driverID := insertDriver(t, db, "TX001", "Nguyen Van A", 1)

	created, err := svc.Create(CreateInput{TruckCode: "XE001", PlateNumber: "29C-11111"})
	if err != nil || created.ID == 0 {
		t.Fatalf("Create ban dau that bai: %v", err)
	}
	if created.DefaultDriverID != nil {
		t.Fatal("xe moi tao chua gan tai xe phai co DefaultDriverID nil")
	}

	updated, err := svc.Update(created.ID, UpdateInput{
		TruckCode:       "XE001",
		PlateNumber:     "29C-11111",
		DefaultDriverID: &driverID,
	})
	if err != nil {
		t.Fatalf("Update gan tai xe that bai: %v", err)
	}
	if updated.DefaultDriverID == nil || *updated.DefaultDriverID != driverID {
		t.Errorf("sau update tai xe mac dinh phai la %d, nhan duoc %v", driverID, updated.DefaultDriverID)
	}
}
