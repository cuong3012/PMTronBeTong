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
│   ├── orders/          # CRUD đơn hàng
│   └── ...              # các module khác
├── data/                # file SQLite (không commit)
└── go.mod
```

## Chạy
```powershell
go run ./cmd/server
```
