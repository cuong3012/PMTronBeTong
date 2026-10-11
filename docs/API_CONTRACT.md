# Hợp đồng API (Backend ↔ Frontend)

Tài liệu này là **nguồn tham chiếu chính thức** để Frontend code theo Backend.
Mọi endpoint, tên trường, kiểu dữ liệu, mã lỗi dưới đây đều **khớp với code
backend hiện tại**. Nếu Frontend cần thay đổi, hãy thống nhất rồi cập nhật file
này trước.

---

## 1. Quy ước chung

### 1.1 Tiền tố và xác thực
- Base URL: `/api` (KHÔNG có `/v1`). Khi dev, Vite proxy `/api` → `http://localhost:8080`.
- Header xác thực: `Authorization: Bearer <token>` cho mọi endpoint trừ `/auth/login` và `/health`.
- `Content-Type: application/json`.

### 1.2 Định dạng lỗi thống nhất
Mọi lỗi trả về một trong hai dạng:

```jsonc
// Dạng cơ bản
{ "error": "thông báo cho người dùng", "message": "chi tiết (tùy chọn)" }

// Dạng mở rộng (chỉ với lỗi trùng mã bản ghi đã xóa mềm)
{
  "error": "Đã từng tồn tại khách hàng với mã này",
  "code": "INACTIVE_CODE_EXISTS",
  "data": { /* bản ghi cũ đã bị xóa mềm */ }
}
```

Mã HTTP dùng chung:
| Mã | Ý nghĩa |
|----|---------|
| 200 | OK |
| 201 | Tạo mới thành công |
| 400 | Dữ liệu/JSON không hợp lệ, `id` không hợp lệ |
| 401 | Chưa đăng nhập / token sai/hết hạn |
| 403 | Không đủ quyền (không phải admin/manager) hoặc tài khoản bị khóa |
| 404 | Không tìm thấy bản ghi |
| 409 | Trùng mã (mã đang hoạt động hoặc đã từng tồn tại) |
| 422 | Dữ liệu hợp lệ về cú pháp nhưng vi phạm nghiệp vụ (ví dụ PG > 4 chữ số thập phân) |
| 500 | Lỗi hệ thống |

### 1.3 Phân trang
Các endpoint danh sách nhận query: `page` (mặc định 1), `per_page` (mặc định 20, tối đa 100), `search`.
Trả về:

```jsonc
{
  "items": [ /* ... */ ],
  "page": 1,
  "per_page": 20,
  "total": 42
}
```

### 1.4 Xóa mềm và khôi phục
- Danh mục (customers, construction-sites, mixer-trucks, drivers, concrete-grades)
  dùng **xóa mềm**: chỉ đặt `is_active = 0`, vẫn `GET` được theo id.
- Khi tạo/đổi mã trùng với bản ghi **đã xóa mềm**, backend trả `409` kèm
  `code: "INACTIVE_CODE_EXISTS"` và bản ghi cũ trong `data` → FE hỏi người dùng
  có muốn khôi phục không.

### 1.5 Quyền
- `requireAuth`: mọi người dùng đã đăng nhập.
- `requireAdminOrManager`: chỉ role `admin` hoặc `manager`.
- Cột "Quyền ghi" dưới đây cho biết role cần để POST/PUT/PATCH/DELETE.

### 1.6 Bảng loại dữ liệu quy ước
- Thời gian: chuỗi `YYYY-MM-DD HH:MM:SS` (UTC, do SQLite `datetime('now')`).
- Tiền/khối lượng: số thực (`float64`). Trường vắng mặt ⇒ `0`.
- Cột `*_id` có thể `null` (đã ghi rõ ở từng model).

---

## 2. Xác thực — `/auth`

### POST `/api/auth/login` (công khai)
Request:
```json
{ "username": "admin", "password": "Admin@123" }
```
Response 200:
```jsonc
{
  "token": "eyJ...",
  "expires_at": "2026-10-12T10:00:00Z",
  "user": {
    "id": 1, "username": "admin", "full_name": "Quản trị viên",
    "role_code": "admin", "role_name": "Quản trị hệ thống",
    "station_id": null
  }
}
```
Lỗi: `400` thiếu field; `401` sai đăng nhập; `403` tài khoản bị khóa.

### GET `/api/auth/me` (auth)
Trả về `UserInfo` như trên.

---

## 3. Cấu hình trạm — `/station`  *(MỚI)*

Trạm là **cấu hình một bản ghi duy nhất**: chỉ đọc và cập nhật, không tạo/xóa.

### GET `/api/station` (auth)
```jsonc
{
  "id": 1,
  "station_code": "TRAM01",
  "station_name": "Trạm trộn số 1",
  "address": "Khu công nghiệp",
  "phone": "",
  "tp1_name": "Cát vàng", "tp2_name": "", "tp3_name": "", "tp4_name": "", "tp5_name": "", "tp6_name": "",
  "xi1_name": "Xi măng PCB40", "xi2_name": "", "xi3_name": "", "xi4_name": "", "xi5_name": "",
  "pg1_name": "Phụ gia", "pg2_name": "", "pg3_name": "", "pg4_name": "", "pg5_name": "", "pg6_name": "",
  "component_visibility_mask": "10000000000000000",
  "accumulate_components": true,
  "accumulate_cement": true,
  "batch_min_volume_m3": 0,
  "batch_max_volume_m3": 0,
  "created_at": "2026-10-11 04:34:38",
  "updated_at": "2026-10-11 04:34:38"
}
```
- `component_visibility_mask`: **đúng 17 ký tự** `0/1`, thứ tự `TP1..TP6, XI1..XI5, PG1..PG6`. `1` = hiển thị.
- Lỗi: `404` chưa có cấu hình.

### PATCH `/api/station` (admin/manager)
Body = cùng các trường ở trên (bỏ `id`, `created_at`, `updated_at`).
```json
{
  "station_code": "TRAM01", "station_name": "Trạm 1",
  "tp1_name": "Cát vàng", "xi1_name": "PCB40", "pg1_name": "Phụ gia",
  "component_visibility_mask": "10000000000000000",
  "accumulate_components": true, "accumulate_cement": true,
  "batch_min_volume_m3": 1.5, "batch_max_volume_m3": 6
}
```
- `station_code`, `station_name`, `component_visibility_mask` là **bắt buộc**, không rỗng.
- Mask sai định dạng ⇒ `422`; max < min ⇒ `422`; mã/tên rỗng ⇒ `422`.

---

## 4. Mác bê tông / Cấp phối — `/concrete-grades`  *(MỚI)*

Định mức thành phần cho **đúng 1 m³** bê tông. Thành phần không dùng lưu `0`.

### GET `/api/concrete-grades` (auth) — phân trang + `search` theo mã/tên
```jsonc
{
  "items": [
    {
      "id": 1, "grade_code": "M250", "printed_name": "Bê tông M250",
      "tp1": 750, "tp2": 1100, "tp3": 0, "tp4": 0, "tp5": 0, "tp6": 0,
      "xi1": 320, "xi2": 0, "xi3": 0, "xi4": 0, "xi5": 0,
      "nuoc": 175,
      "pg1": 0, "pg2": 0, "pg3": 0, "pg4": 0, "pg5": 0, "pg6": 0,
      "do_sut": "10-12", "cuong_do": "B25",
      "notes": "",
      "is_active": true,
      "created_at": "2026-10-11 04:34:38", "updated_at": "2026-10-11 04:34:38"
    }
  ],
  "page": 1, "per_page": 20, "total": 3
}
```

### GET `/api/concrete-grades/:id` (auth) — một bản ghi ở trên.

### POST `/api/concrete-grades` (admin/manager)
```json
{
  "grade_code": "M300", "printed_name": "Bê tông M300",
  "tp1": 720, "tp2": 1080, "xi1": 360, "nuoc": 170, "pg1": 4.5,
  "do_sut": "10-12", "cuong_do": "B25", "notes": ""
}
```
- `grade_code`, `printed_name` bắt buộc, không rỗng ⇒ nếu rỗng `422`.
- **`do_sut`, `cuong_do` là TEXT** (không phải số) — có thể là khoảng/ký hiệu như
  `"10-12"`, `"B25"`, `"M300"`. Tối đa 100 ký tự mỗi trường.
- **`pg1..pg6` tối đa 4 chữ số thập phân** ⇒ vượt `422`.
- Trùng mã đang hoạt động ⇒ `409`.
- Trùng mã đã xóa mềm ⇒ `409` + `code: INACTIVE_CODE_EXISTS` + `data` = bản ghi cũ.
  *(Lưu ý: module grades hiện chưa có endpoint restore.)*

### PUT `/api/concrete-grades/:id` (admin/manager) — body như POST.

### DELETE `/api/concrete-grades/:id` (admin/manager) → `{ "deleted": true, "id": 1 }`

---

## 5. Danh mục thành phần — `/component-catalog`  *(MỚI)*

Metadata tên/định nghĩa vị trí thành phần (không phải dữ liệu sản xuất).
Backend seed sẵn **17 bản ghi** (`TP1..TP6, XI1..XI5, PG1..PG6`).

### GET `/api/component-catalog` (auth) — query tùy chọn `search` theo `slot_name`
Trả về **không phân trang**:
```jsonc
{
  "items": [
    {
      "id": 1, "slot_name": "TP1",
      "humidity_values": "12.5,11.8,10.7,9.6,8.9,7.4",
      "unit": "kg", "notes": "",
      "created_at": "...", "updated_at": "..."
    }
  ]
}
```

### GET `/api/component-catalog/:id` (auth) — một bản ghi.

### POST `/api/component-catalog` (admin/manager)
```json
{ "slot_name": "TP1", "humidity_values": "12.5,11.8", "unit": "kg", "notes": "" }
```
- `slot_name` bắt buộc, không rỗng ⇒ `422`; trùng ⇒ `409`.

### PUT `/api/component-catalog/:id` (admin/manager) — body như POST.
### DELETE `/api/component-catalog/:id` (admin/manager) → **xóa cứng** `{ "deleted": true, "id": 1 }`.

---

## 6. Cấp phối chờ — `/waiting-mixes`  *(MỚI)*

Đúng **20 vị trí cố định** (`stt` 1..20), seed sẵn khi khởi động. Chỉ đọc/cập nhật;
KHÔNG tạo/xóa qua API.

### GET `/api/waiting-mixes` (auth) — trả đủ 20 vị trí (không phân trang)
```jsonc
{
  "items": [
    {
      "id": 1, "stt": 1,
      "concrete_grade_id": 1,          // có thể vắng/null nếu vị trí trống
      "concrete_grade_code": "M250",
      "concrete_grade_name": "Bê tông M250",
      "customer_id": 1,                // có thể vắng/null nếu vị trí trống
      "customer_name": "Công ty ABC",
      "volume_m3": 7, "m3me": 1.5, "status": 1,
      "created_at": "...", "updated_at": "..."
    }
    // ... đủ 20 phần tử
  ]
}
```

### GET `/api/waiting-mixes/:id` (auth) — một vị trí.

### PATCH `/api/waiting-mixes/:id` (auth)
```json
{ "concrete_grade_id": 1, "customer_id": 1, "volume_m3": 7, "m3me": 1.5, "status": 1 }
```
- `concrete_grade_id`, `customer_id` bắt buộc (> 0).
- `status` chỉ nhận `0` hoặc `1` ⇒ sai `422`.
- `volume_m3 > 0` và `m3me >= 0` ⇒ sai `422`.
- `concrete_grade_id` / `customer_id` phải là bản ghi **đang hoạt động** ⇒ sai `422`.

---

## 7. Phiếu cân — `/weighing-tickets`  *(MỚI)*

Phiếu cân (= `thongtinphieucan`) và chi tiết mẻ (= `chitietphieucan`).
Là **bản ghi kết quả**: lưu kèm bản sao mác/khách/công trường tại thời điểm cân.

### GET `/api/weighing-tickets` (auth) — phân trang + `search` theo mã/khách/mác
```jsonc
{
  "items": [
    {
      "id": 1, "ticket_code": "PC001",
      "volume_m3": 6, "batch_count": 3, "m3me": 2,
      "order_id": 1, "order_number": "DH001",
      "concrete_grade_id": 1, "concrete_grade_code": "M250",
      "concrete_grade_name": "Bê tông M250",
      "printed_concrete_grade_name": "Bê tông M250",
      "humidity_values": "12.5,11.8,10.7,9.6,8.9,7.4",
      "standard_per_cube": 2400, "standard_per_batch": 4800,
      "customer_name": "Công ty ABC",
      "construction_site_address": "Công trường mẫu",
      "material_name": "Cát vàng",
      "runner_user_id": null,
      "notes": "",
      "is_active": true,
      "created_at": "...", "updated_at": "..."
    }
  ],
  "page": 1, "per_page": 20, "total": 4
}
```
> Danh sách **không** kèm `details`; gọi `GET /:id` để lấy chi tiết.

### GET `/api/weighing-tickets/:id` (auth) — **kèm `details`**
```jsonc
{
  /* ...các trường như trên... */
  "details": [
    {
      "id": 1, "weighing_ticket_id": 1, "batch_no": 1,
      "batch_volume_m3": 2, "mixed_at": "2026-10-11 08:00:00",
      "tp1": 750, "tp2": 1100, "tp3": 0, "tp4": 0, "tp5": 0, "tp6": 0,
      "xi1": 320, "xi2": 0, "xi3": 0, "xi4": 0, "xi5": 0,
      "pg1": 0, "pg2": 0, "pg3": 0, "pg4": 0, "pg5": 0, "pg6": 0,
      "manual_override_json": "",
      "notes": "",
      "is_active": true, "created_at": "...", "updated_at": "..."
    }
  ]
}
```

### POST `/api/weighing-tickets` (admin/manager)
Tạo phiếu + toàn bộ mẻ trong **một transaction**.
```jsonc
{
  "ticket_code": "PC999",
  "volume_m3": 4, "batch_count": 2, "m3me": 2,
  "order_id": 1,
  "concrete_grade_id": 1,
  "concrete_grade_code": "M250",
  "concrete_grade_name": "Bê tông M250",
  "printed_concrete_grade_name": "Bê tông M250",
  "humidity_values": "12.5,11.8",
  "standard_per_cube": 2400, "standard_per_batch": 4800,
  "customer_name": "Công ty ABC",
  "construction_site_address": "Công trường mẫu",
  "material_name": "Cát vàng",
  "runner_user_id": null,
  "notes": "",
  "details": [
    { "batch_no": 1, "batch_volume_m3": 2, "mixed_at": "2026-10-11T08:00:00Z",
      "tp1": 750, "tp2": 1100, "xi1": 320 },
    { "batch_no": 2, "batch_volume_m3": 2, "mixed_at": "2026-10-11T08:10:00Z",
      "tp1": 750, "tp2": 1100, "xi1": 320 }
  ]
}
```
Bắt buộc: `ticket_code` (không rỗng), `order_id`, `concrete_grade_id`,
`concrete_grade_code`, `concrete_grade_name`, `printed_concrete_grade_name`, `customer_name`.
Mỗi phần tử `details` cần `batch_no > 0`, `batch_volume_m3 > 0`, `mixed_at` không rỗng.
- `volume_m3 <= 0` hoặc `batch_count <= 0` ⇒ `422`.
- Chi tiết thiếu/sai ⇒ `422`; trùng `ticket_code` ⇒ `409`.
- Trả về **201** kèm phiếu đầy đủ (đã có `details`).

### DELETE `/api/weighing-tickets/:id` (admin/manager)
Xóa mềm phiếu **và** các mẻ của nó → `{ "deleted": true, "id": 5 }`.

### GET `/api/weighing-tickets/stats` (auth) — query tùy chọn `from`, `to` (`YYYY-MM-DD`)
```jsonc
{
  "by_grade": [
    { "concrete_grade_code": "M250", "concrete_grade_name": "Bê tông M250",
      "ticket_count": 4, "total_volume_m3": 31.5, "total_batch_count": 15 }
  ],
  "daily": [
    { "day": "2026-10-05", "ticket_count": 1, "total_volume_m3": 6 }
  ]
}
```

---

## 8. Danh mục cũ (đã có, giữ nguyên)

| Module | Đường dẫn | Quyền ghi | Ghi chú |
|--------|-----------|-----------|---------|
| Khách hàng | `/customers` | admin/manager | xóa mềm + `POST /:id/restore` |
| Công trường | `/construction-sites` | admin/manager | xóa mềm + `/:id/restore` |
| Xe bồn | `/mixer-trucks` | admin/manager | xóa mềm + `/:id/restore` |
| Tài xế | `/drivers` | admin/manager | xóa mềm + `/:id/restore` |

Mẫu chuẩn cho cả 4 module (ví dụ `/customers`):
- `GET /customers?page&per_page&search` → `PageResult<Customer>`
- `GET /customers/:id` → `Customer`
- `POST /customers` → 201 `Customer`
- `PUT /customers/:id` → `Customer`
- `DELETE /customers/:id` → `{ "deleted": true, "id" }`
- `POST /customers/:id/restore` → `Customer`

Trường `Customer`: `id, customer_code, customer_name, address, phone, is_active, created_at, updated_at`.
Trường `ConstructionSite`: `id, site_code, site_name, is_active, created_at, updated_at`.
Trường `MixerTruck`: `id, truck_code, plate_number, default_driver_id (nullable), default_driver_name, is_active, created_at, updated_at`.
Trường `Driver`: `id, driver_code, full_name, is_active, created_at, updated_at`.

---

## 9. Bảng route tổng hợp

| Method | Path | Quyền | Mục |
|--------|------|-------|-----|
| POST | `/api/auth/login` | công khai | 2 |
| GET | `/api/auth/me` | auth | 2 |
| GET | `/api/station` | auth | 3 |
| PATCH | `/api/station` | admin/manager | 3 |
| GET | `/api/concrete-grades` | auth | 4 |
| GET | `/api/concrete-grades/:id` | auth | 4 |
| POST | `/api/concrete-grades` | admin/manager | 4 |
| PUT | `/api/concrete-grades/:id` | admin/manager | 4 |
| DELETE | `/api/concrete-grades/:id` | admin/manager | 4 |
| GET | `/api/component-catalog` | auth | 5 |
| GET | `/api/component-catalog/:id` | auth | 5 |
| POST | `/api/component-catalog` | admin/manager | 5 |
| PUT | `/api/component-catalog/:id` | admin/manager | 5 |
| DELETE | `/api/component-catalog/:id` | admin/manager | 5 |
| GET | `/api/waiting-mixes` | auth | 6 |
| GET | `/api/waiting-mixes/:id` | auth | 6 |
| PATCH | `/api/waiting-mixes/:id` | auth | 6 |
| GET | `/api/weighing-tickets` | auth | 7 |
| GET | `/api/weighing-tickets/stats` | auth | 7 |
| GET | `/api/weighing-tickets/:id` | auth | 7 |
| POST | `/api/weighing-tickets` | admin/manager | 7 |
| DELETE | `/api/weighing-tickets/:id` | admin/manager | 7 |

---

## 10. Gợi ý triển khai Frontend

Tạo các file API module tương ứng trong `frontend/src/api/` theo mẫu `sites.ts`:
`station.ts`, `concreteGrades.ts`, `componentCatalog.ts`, `waitingMixes.ts`, `weighingTickets.ts`,
và bổ sung interface vào `frontend/src/types/index.ts`:

```ts
// types/index.ts (bổ sung)
export interface Station { id: number; station_code: string; /* ...17 tên + mask... */ }
export interface ConcreteGrade { /* tp1..tp6, xi1..xi5, nuoc, pg1..pg6: number; do_sut, cuong_do: string; ... */ }
export interface ComponentCatalog { id: number; slot_name: string; humidity_values: string; unit?: string; notes?: string; /* ... */ }
export interface WaitingMix { id: number; stt: number; concrete_grade_id?: number | null; customer_id?: number | null; /* ... */ }
export interface WeighingTicket { /* header... */ details?: WeighingTicketDetail[] }
export interface WeighingTicketDetail { /* batch_no, batch_volume_m3, mixed_at, 16 thành phần... */ }
```

Lưu ý khi code:
- `GET /component-catalog` và `GET /waiting-mixes` trả về `{ items }`, **không** có phân trang
  (khác với `PageResult`). Đừng tái sử dụng nhầm `CrudPage`.
- `concrete_grade_id`/`customer_id` của `WaitingMix` và `runner_user_id` của phiếu cân có thể thiếu → khai báo `| null` hoặc optional.
- Khi gặp `409` + `code === 'INACTIVE_CODE_EXISTS'`, dùng `parseInactiveCodeConflict` có sẵn trong `api/client.ts`.
- Nút ghi/xóa chỉ hiển thị cho `admin`/`manager` (khớp `requireAdminOrManager`).
- **`do_sut`/`cuong_do` của mác bê tông là `string`** (dùng `Input` text, KHÔNG dùng `InputNumber`).

---

## 11. Dữ liệu seed khi khởi động

### 11.1 Bắt buộc (luôn chạy)
4 roles, 1 trạm (`TRAM01`), tài khoản `admin`, **17 `component_catalog`**, **20 `waiting_mixes`**.

### 11.2 Dữ liệu mẫu (tùy chọn — bật bằng `SEED_SAMPLE=true`)
Không tự chạy, để tránh làm bẩn dữ liệu production. Khi khởi động với
`SEED_SAMPLE=true`, backend nạp thêm (idempotent — chỉ chèn khi bảng trống):
- 3 mác bê tông: `M250`, `M300`, `M350`.
- 1 khách hàng (`KH001`), 1 công trường (`CT001`), 1 đơn (`DH001`).
- Điền 5 vị trí `waiting_mixes` đầu.
- **6 phiếu cân** (`PC001`..`PC006`) kèm **21 chi tiết mẻ** trải trên nhiều ngày
  (để thống kê có số liệu). Trong đó `PC006` đặt `is_active=0` làm ví dụ
  phiếu đã xóa mềm (sẽ ẩn khỏi danh sách/stats).

Cách chạy:

```powershell
# PowerShell
$env:SEED_SAMPLE = "true"
go run ./cmd/server
```

```bash
# bash
SEED_SAMPLE=true go run ./cmd/server
```

> Lưu ý: `SEED_SAMPLE` chỉ có tác dụng khi bảng `weighing_tickets` (và các bảng
> liên quan) còn trống. Nếu đã có dữ liệu, backend bỏ qua để không ghi đè.

> Cập nhật tài liệu này mỗi khi backend đổi route/trường. Đây là hợp đồng giữa hai phía.
