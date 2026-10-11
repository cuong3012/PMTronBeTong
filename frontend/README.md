# Frontend — Hệ thống quản lý trạm trộn bê tông

Ứng dụng React + TypeScript + Vite cho hệ thống quản lý trạm trộn bê tông.

## Công nghệ
- React 18 + TypeScript + Vite
- React Router (điều hướng)
- TanStack Query (quản lý dữ liệu API và cache)
- Zustand (lưu phiên đăng nhập, persist localStorage)
- Ant Design (giao diện nghiệp vụ)
- Axios (gọi HTTP API)

## Cấu trúc thư mục
```text
frontend/
├── src/
│   ├── api/        # client axios + các module gọi API theo nghiệp vụ
│   ├── components/ # layout và component dùng chung (CrudPage, ...)
│   ├── pages/      # các trang (Login, Dashboard, danh mục CRUD)
│   ├── routes/     # ProtectedRoute
│   ├── store/      # Zustand store (xác thực)
│   ├── types/      # kiểu dữ liệu dùng chung
│   ├── utils/      # hàm định dạng
│   ├── App.tsx     # định nghĩa route
│   └── main.tsx    # điểm khởi tạo
├── index.html
├── package.json
└── vite.config.ts
```

## Chạy
Cài dependencies và chạy dev server (mặc định cổng 5173):
```powershell
npm install
npm run dev
```

Trong khi phát triển, Vite proxy các request `/api` sang backend Go tại
`http://localhost:8080`. Hãy chạy backend trước:
```powershell
cd ../backend
go run ./cmd/server
```

## Tài khoản mặc định
Backend seed sẵn tài khoản quản trị:
- Tên đăng nhập: `admin`
- Mật khẩu: `Admin@123`

> Hãy đổi mật khẩu ngay sau lần đăng nhập đầu tiên trên môi trường thực tế.

## Chức năng đã có
- Đăng nhập / đăng xuất, lưu phiên bằng token Bearer.
- Dashboard hiển thị tổng số khách hàng, công trường, tài xế, xe bồn.
- CRUD (thêm/sửa/xóa mềm/tìm kiếm/phân trang) cho:
  - Khách hàng (`/customers`)
  - Công trường (`/construction-sites`)
  - Tài xế (`/drivers`)
  - Xe bồn (`/mixer-trucks`)

## Hợp đồng API (đọc trước khi code)
Xem **[`../docs/API_CONTRACT.md`](../docs/API_CONTRACT.md)** — đặc tả đầy đủ các
endpoint, request/response, kiểu dữ liệu và mã lỗi mà backend cung cấp.

Các module backend **mới** cần Frontend triển khai thêm:
- Cấu hình trạm: `/station` (GET/PATCH, một bản ghi).
- Mác bê tông: `/concrete-grades` (CRUD, xóa mềm, `pg1..pg6` ≤ 4 chữ số thập phân).
- Danh mục thành phần: `/component-catalog` (CRUD, `GET` trả `{ items }` không phân trang).
- Cấp phối chờ: `/waiting-mixes` (GET/PATCH, đúng 20 vị trí).
- Phiếu cân: `/weighing-tickets` (CRUD + `/stats`, kèm `details` khi lấy chi tiết).

## Phân quyền trên giao diện
- Nút **Xóa** chỉ hiển thị với vai trò `admin` hoặc `manager`, khớp với
  `requireAdminOrManager` của backend. Ẩn nút chỉ là hỗ trợ UX; backend vẫn
  kiểm tra quyền trên mọi request.

## Build
```powershell
npm run build
```
Kết quả build nằm trong `dist/`, có thể phục vụ qua web server/CDN hoặc để
backend phục vụ file tĩnh.

## Ghi chú
- API backend dùng tiền tố `/api` (không phải `/api/v1` như bản nháp đặc tả).
- Xóa danh mục là **xóa mềm**: backend chỉ đặt `is_active = 0`.
