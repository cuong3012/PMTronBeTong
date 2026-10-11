# Backend — Hệ thống quản lý trạm trộn bê tông

Go Backend API cho hệ thống quản lý trạm trộn bê tông.

## Công nghệ
- Go 1.27
- Gin (HTTP API)
- SQLite (modernc.org/sqlite — pure Go, không cần CGO)
- JWT cho xác thực
- bcrypt cho mật khẩu

## Cấu trúc thư mục
```text
backend/
├── cmd/server/          # entrypoint
├── internal/
│   ├── config/          # cấu hình
│   ├── database/        # kết nối DB, migrations, seed
│   ├── middleware/      # auth, logging, error handling
│   ├── models/          # struct ánh xạ bảng
│   ├── auth/            # login, JWT
│   ├── customers/       # CRUD khách hàng
│   ├── sites/           # CRUD công trường
│   ├── trucks/          # CRUD xe bồn
│   ├── drivers/         # CRUD tài xế
│   ├── stations/        # cấu hình trạm (một bản ghi)
│   ├── grades/          # CRUD mác bê tông / cấp phối
│   ├── components/      # danh mục thành phần cấp phối
│   ├── waitingmix/      # cấp phối chờ (20 vị trí)
│   ├── weighing/        # phiếu cân + chi tiết + thống kê
│   ├── orders/          # CRUD đơn hàng
│   └── ...              # các module khác
├── data/                # file SQLite (không commit)
└── go.mod
```

## Chạy
```powershell
go run ./cmd/server
```

### Nạp dữ liệu mẫu (tùy chọn)
Mặc định backend chỉ nạp dữ liệu tối thiểu (roles, trạm, admin, danh mục thành
phần, 20 vị trí cấp phối chờ). Để nạp thêm **dữ liệu mẫu** (mác bê tông, khách
hàng, công trường, đơn hàng và **phiếu cân + chi tiết phiếu cân**), bật biến
môi trường `SEED_SAMPLE`:

```powershell
$env:SEED_SAMPLE = "true"
go run ./cmd/server
```

Chỉ chèn khi bảng còn trống (idempotent), nên chạy lại nhiều lần vẫn an toàn.

## Biến môi trường chính
| Biến | Mặc định | Ý nghĩa |
|------|----------|---------|
| `APP_PORT` | `8080` | Cổng HTTP |
| `DB_PATH` | `data/tramtron.db` | Đường dẫn file SQLite |
| `JWT_SECRET` | `dev-secret-...` | Khóa ký JWT (đổi khi production) |
| `JWT_EXPIRY_HOUR` | `12` | Hạn token (giờ) |
| `CORS_ORIGIN` | `http://localhost:5173` | Origin frontend được phép |
| `APP_ENV` | `development` | `development`/`production` |
| `SEED_SAMPLE` | `false` | Bật nạp dữ liệu mẫu khi khởi động |

## Hợp đồng API
Đặc tả đầy đủ các endpoint backend cung cấp (để Frontend code theo):
xem **[`../docs/API_CONTRACT.md`](../docs/API_CONTRACT.md)**.
Cập nhật tài liệu này mỗi khi thay đổi route/trường dữ liệu.
