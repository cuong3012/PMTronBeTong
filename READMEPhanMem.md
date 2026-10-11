# Hệ thống quản lý trạm trộn bê tông IMI

Ứng dụng quản lý trạm trộn bê tông gồm hai phần độc lập:

- `frontend/` — Ứng dụng React + TypeScript + Vite
- `backend/` — Go Backend API (Gin, SQLite, JWT)

Xem chi tiết yêu cầu tại [`SOFTWARE_REQUIREMENTS.md`](./SOFTWARE_REQUIREMENTS.md).

## Kiến trúc

```text
01PMTramTron/
├── frontend/          # React + TypeScript (Vite)
├── backend/           # Go API + SQLite
├── SOFTWARE_REQUIREMENTS.md
└── README.md
```

Frontend chỉ gọi Go Backend qua HTTP API (`/api/...`); trình duyệt không kết
nối trực tiếp tới SQLite hay PLC.

## Chạy backend

Cần cài Go (xem `backend/go.mod`). Từ thư mục `backend/`:

```powershell
go run ./cmd/server
```

Backend lắng nghe ở `http://localhost:8080`, tự chạy migration và seed dữ liệu
(tài khoản admin mặc định: `admin` / `Admin@123`).

## Chạy frontend

Cần Node.js. Từ thư mục `frontend/`:

```powershell
npm install
npm run dev
```

Frontend chạy ở `http://localhost:5173` và proxy `/api` sang backend.

## Chức năng hiện có

- **Đăng nhập / đăng xuất** bằng token Bearer (`POST /api/auth/login`).
- **Dashboard** tổng quan số lượng khách hàng, công trường, tài xế, xe bồn.
- **CRUD danh mục** (thêm / sửa / xóa mềm / tìm kiếm / phân trang):
  - Khách hàng
  - Công trường
  - Tài xế
  - Xe bồn (kèm chọn tài xế mặc định)

## API backend (tiền tố `/api`)

| Nhóm | Endpoint |
|---|---|
| Xác thực | `POST /auth/login`, `GET /auth/me` |
| Khách hàng | `GET/POST /customers`, `GET/PUT/DELETE /customers/{id}` |
| Công trường | `GET/POST /construction-sites`, `GET/PUT/DELETE /construction-sites/{id}` |
| Tài xế | `GET/POST /drivers`, `GET/PUT/DELETE /drivers/{id}` |
| Xe bồn | `GET/POST /mixer-trucks`, `GET/PUT/DELETE /mixer-trucks/{id}` |

- Danh sách trả `{ items, page, per_page, total }`; hỗ trợ `page`, `per_page`, `search`.
- Lỗi trả `{ error, message }`.
- Xóa danh mục là **xóa mềm** và chỉ dành cho vai trò `admin`/`manager`.

## Phân quyền

Ẩn nút Xóa trên giao diện chỉ là hỗ trợ UX; backend vẫn kiểm tra quyền trên mọi
request (middleware `RequireAuth` + `RequireRole`).
