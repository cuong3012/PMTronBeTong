# Đặc tả yêu cầu phần mềm - Trạm trộn bê tông

## 1. Thông tin dự án
- Tên phần mềm: Hệ thống quản lý trạm trộn bê tông
- Người yêu cầu: Ban quản lý / bộ phận vận hành trạm
- Người phát triển: Đội ngũ phát triển phần mềm
- Ngày tạo: 2026-10-10
- Trạng thái: Bản nháp chuyển đổi từ hệ thống C# sang React + Go

## 2. Mục tiêu tổng quan
Hệ thống này nhằm hỗ trợ quản lý toàn bộ quy trình trạm trộn bê tông từ đơn hàng, nguyên vật liệu, công thức, sản xuất, kiểm soát chất lượng đến xuất hàng và báo cáo. Mục tiêu chính là:
- Tự động hóa quy trình thao tác nghiệp vụ
- Giảm sai sót do nhập liệu thủ công
- Theo dõi tồn kho, khối lượng bê tông, thời gian và hiệu suất sản xuất
- Tăng khả năng báo cáo và ra quyết định cho quản lý
- Cung cấp nền tảng hiện đại, dễ mở rộng trên web

## 3. Bối cảnh và vấn đề hiện tại
- Hiện tại, quy trình quản lý trạm trộn bê tông có thể đang triển khai theo dạng Excel, phần mềm cũ hoặc thao tác thủ công.
- Các vấn đề chính gồm:
  - Nhập đơn hàng và công thức theo cách thủ công, dễ sai lệch
  - Khó theo dõi tồn kho nguyên vật liệu theo thời gian thực
  - Khó kiểm soát số lượng bê tông đã sản xuất, giao hàng và tồn kho cuối ca
  - Báo cáo và tổng hợp dữ liệu mất nhiều thời gian, dễ sai sót
  - Quyền truy cập và phân quyền chưa rõ ràng
- Điều này dẫn đến:
  - Tăng chi phí vận hành và nhân sự
  - Tỷ lệ sai đơn, sai công thức hoặc sai khối lượng cao hơn mong muốn
  - Khó phản hồi nhanh khi xảy ra sự cố trong sản xuất

## 4. Người dùng và vai trò
### 4.1 Người dùng chính
- Admin hệ thống
- Quản lý trạm / cấp quản lý
- Điều phối viên đơn hàng
- Nhân viên kho nguyên liệu
- Nhân viên trạm trộn / vận hành máy
- Nhân viên kiểm soát chất lượng
- Kỹ thuật viên / kỹ sư công thức
- Khách hàng / đối tác (nếu cần quyền xem báo cáo hoặc đơn hàng)

### 4.2 Vai trò theo chức năng
- Admin: quản lý người dùng, phân quyền, cấu hình hệ thống
- Quản lý: theo dõi sản xuất, báo cáo, quyết định điều phối
- Điều phối viên: tạo đơn hàng, cập nhật tiến độ, kiểm tra lịch giao
- Kho: quản lý nhập kho, xuất kho, tồn kho nguyên vật liệu
- Vận hành: theo dõi máy móc, công thức, số liệu sản xuất từng mẻ
- Kỹ thuật: cấu hình công thức, kiểm tra độ tương thích nguyên liệu
- QC: ghi nhận kiểm tra chất lượng, lưu kết quả mẫu thử

## 5. Yêu cầu chức năng
### 5.1 Chức năng chung
1. Hệ thống phải hỗ trợ đăng nhập và xác thực người dùng.
2. Hệ thống phải phân quyền theo vai trò rõ ràng.
3. Hệ thống phải lưu trữ dữ liệu và lịch sử thay đổi theo thời gian.
4. Hệ thống phải hiển thị dữ liệu theo dạng bảng, biểu đồ và báo cáo trực quan.
5. Hệ thống phải có nhật ký hoạt động (audit log) cho các thao tác quan trọng.

### 5.2 Chức năng nghiệp vụ chính
- Quản lý khách hàng và thông tin đối tác
- Quản lý đơn hàng và lịch giao hàng
- Quản lý nguyên vật liệu và hóa đơn nhập kho
- Quản lý công thức trộn bê tông theo từng loại và mác bê tông
- Quản lý ca sản xuất / mẻ trộn / quy trình vận hành
- Theo dõi khối lượng bê tông sản xuất theo thời gian thực
- Quản lý máy móc, thiết bị và trạng thái hoạt động
- Quản lý kiểm tra chất lượng đầu ra
- Quản lý tồn kho, xuất kho và sử dụng nguyên vật liệu
- Quản lý báo cáo tổng hợp cho quản lý

### 5.3 Quy trình nghiệp vụ
- Bước 1: Tạo thông tin khách hàng và đơn hàng
- Bước 2: Xác định công thức bê tông phù hợp với yêu cầu kỹ thuật
- Bước 3: Kiểm tra tồn kho nguyên liệu và chuẩn bị nguyên vật liệu cho mẻ
- Bước 4: Khởi tạo ca sản xuất và nhập thông tin mẻ trộn
- Bước 5: Theo dõi thời gian, số lượng, máy móc và chất lượng trong quá trình sản xuất
- Bước 6: Ghi nhận kết quả kiểm tra chất lượng và xác nhận sản phẩm
- Bước 7: Xuất hàng hoặc giao hàng và cập nhật tồn kho cuối ca
- Bước 8: Tổng hợp báo cáo theo ngày / tuần / tháng

### 5.4 Yêu cầu dữ liệu
- Dữ liệu cần lưu:
  - Khách hàng, đối tác, địa chỉ, thông tin liên hệ
  - Đơn hàng, loại bê tông, khối lượng, thời gian giao
  - Nguyên vật liệu, đơn vị tính, giá thành, tồn kho
  - Công thức trộn, tỷ lệ thành phần, mức độ pha trộn
  - Mẻ sản xuất, máy móc, thời gian bắt đầu/kết thúc
  - Kết quả kiểm tra chất lượng và mẫu thử
  - Tài khoản người dùng, phân quyền, lịch sử thao tác
- Ràng buộc dữ liệu:
  - Trường bắt buộc không được null
  - Mã đơn hàng, mã nguyên liệu, mã công thức phải duy nhất
  - Tồn kho không được âm khi phát sinh xuất kho
  - Số liệu sản xuất phải khớp giữa mẻ trộn và xuất hàng

## 6. Yêu cầu phi chức năng
### 6.1 Hiệu năng
- Thời gian phản hồi các màn hình chính dưới 2 giây đối với dữ liệu thông thường
- Số lượng người dùng đồng thời cần được xác nhận theo thực tế; SQLite phù hợp nhất khi một Backend API phục vụ một trạm với mức ghi dữ liệu vừa phải
- Nhiều màn hình có thể đọc dữ liệu đồng thời, nhưng các thao tác ghi vào SQLite được tuần tự hóa; cần kiểm thử tải theo số người dùng và tần suất ghi thực tế
- Bảng báo cáo lớn phải xử lý dữ liệu theo ngày / tháng mà không làm treo giao diện

### 6.2 Bảo mật
- Mật khẩu phải được mã hóa bằng thuật toán an toàn
- Có phân quyền theo vai trò người dùng
- Có log thao tác cho các chức năng quan trọng như tạo đơn hàng, sửa công thức, nhập kho, xuất kho
- Có cơ chế session / token bảo mật cho API

### 6.3 Tính khả dụng
- Hệ thống cần hoạt động ổn định trong giờ làm việc
- Dữ liệu cần có khôi phục và backup định kỳ
- Có cảnh báo khi xảy ra lỗi nghiêm trọng hoặc API không khả dụng

### 6.4 Tính mở rộng
- Mở rộng thêm mô-đun cho nhiều trạm, nhiều kho, nhiều công thức mới
- Dễ tích hợp với phần mềm kế toán, CRM hoặc ERP sau này
- API rõ ràng, module hóa để phát triển tiếp trong tương lai

### 6.5 Tương thích
- Hỗ trợ trình duyệt Chrome, Edge, Firefox mới nhất
- Giao diện web đáp ứng trên máy tính bàn và màn hình lớn
- Hệ thống có thể triển khai trên máy chủ Linux hoặc Windows

## 7. Giao diện và trải nghiệm người dùng
- Giao diện chính phải tối ưu cho việc xử lý nhanh trong môi trường công nghiệp
- Màn hình dashboard nên hiển thị các chỉ số quan trọng như tồn kho, đơn hàng hôm nay, sản lượng, thực hiện ca
- Form nhập liệu cần rõ ràng, dễ kiểm tra, có validation phù hợp
- Có thông báo lỗi và xác nhận thao tác rõ ràng
- Có bộ lọc, tìm kiếm và sắp xếp theo nhiều tiêu chí

## 8. Điều kiện chấp nhận
Hệ thống được coi là hoàn thành khi:
1. Người dùng có thể tạo và quản lý đơn hàng, nguyên vật liệu, công thức và mẻ sản xuất.
2. Tồn kho và sản lượng được cập nhật nhất quán qua các quy trình nhập/xuất/hủy/tra cứu.
3. Báo cáo theo ngày / tuần / tháng có thể tạo ra đúng và nhanh chóng.
4. Phân quyền người dùng hoạt động đúng, không ai vượt quyền thao tác.
5. Có kiểm thử cơ bản cho các chức năng chính và ít nhất 1 luồng nghiệp vụ hoàn chỉnh hoạt động ổn định.

## 9. Ràng buộc và giới hạn
- Chi phí đầu tư cần hợp lý nhưng phải đảm bảo hiệu quả lâu dài
- Thời gian triển khai ưu tiên cho các tính năng cốt lõi trước, sau đó mới mở rộng chức năng phụ
- Nền tảng mới cần giữ được các quy trình nghiệp vụ đã có từ hệ thống cũ, tránh mất dữ liệu hoặc luồng làm việc cũ
- Dữ liệu từ hệ thống cũ cần được đối chiếu và cleanse trước khi migrate

## 10. Dự kiến phạm vi và ngoài phạm vi
### 10.1 Trong phạm vi
- Quản lý đơn hàng
- Quản lý nguyên vật liệu và tồn kho
- Quản lý công thức và mẻ trộn
- Quản lý máy móc, ca sản xuất và chất lượng
- Báo cáo quản lý
- Phân quyền, tài khoản và audit log

### 10.2 Ngoài phạm vi
- Tự động điều khiển máy móc qua PLC chưa nằm trong giai đoạn đầu nếu chưa được xác nhận yêu cầu; cần khảo sát PLC, giao thức Modbus và quy trình an toàn trước khi triển khai
- Hệ thống kế toán tài chính đầy đủ
- Mô-đun ERP tổng hợp toàn công ty
- Chức năng mobile phức tạp chưa cần triển khai ngay

## 11. Câu hỏi/điểm cần làm rõ
- Hệ thống cũ hiện đang lưu dữ liệu theo cấu trúc nào?
- Có cần cấp quyền xem/chỉnh sửa cho nhiều phòng ban không?
- Mức độ tích hợp với máy móc, PLC hay cảm biến là bao nhiêu?
- Có phải theo dõi nhiều trạm / nhiều kho / nhiều loại bê tông không?
- Có cần báo cáo theo mốc ngày / ca / máy / khách hàng không?

## 12. Kiến nghị công nghệ đề xuất
### 12.1 Frontend
- React + TypeScript + Vite; mã nguồn frontend nằm riêng trong thư mục `frontend/`
- React Router, React Hook Form, Zustand hoặc Redux Toolkit
- Ant Design hoặc Material UI để xây dựng giao diện nghiệp vụ nhanh
- TanStack Query để quản lý dữ liệu API và cache
- Frontend chỉ gọi Go Backend qua HTTP/HTTPS API; không kết nối trực tiếp tới SQLite hoặc PLC

### 12.2 Backend
- Go (Golang); mã nguồn backend nằm riêng trong thư mục `backend/`
- Gin hoặc Echo cho HTTP API
- Cấu trúc theo module: auth, users, customers, orders, materials, formulas, production, quality, reports
- Dùng middleware cho logging, auth, validation, xử lý lỗi
- Hỗ trợ unit test và integration test
- Backend Go có thể vừa cung cấp API/quản lý dữ liệu vừa giao tiếp PLC Modbus trong cùng một ứng dụng; phần PLC được tổ chức thành module riêng bên trong `backend/`
- Module PLC chạy nền để polling trạng thái/giá trị và xử lý yêu cầu điều khiển đơn giản; không chạy tác vụ polling bên trong luồng xử lý request HTTP
- Chỉ cân nhắc tách module PLC thành tiến trình/service riêng (PLC Gateway) nếu sau này cần triển khai, khởi động, giám sát hoặc mở rộng độc lập

### 12.2.1 Cấu trúc thư mục repository dự kiến
```text
01PMTramTron/
├── frontend/          # Ứng dụng React + TypeScript
├── backend/           # Go Backend API và các module nghiệp vụ
├── SOFTWARE_REQUIREMENTS.md
└── README.md
```
- `frontend/` và `backend/` được phát triển, cài dependencies và build độc lập; frontend giao tiếp với backend qua API đã thống nhất.
- Database SQLite và cấu hình truy cập database thuộc trách nhiệm của `backend/`; không đặt file database trong thư mục frontend hoặc gửi file database tới trình duyệt.
- Trong giai đoạn đầu, một ứng dụng Go trong `backend/` chịu trách nhiệm đồng thời cho HTTP API, nghiệp vụ, SQLite và module PLC chạy nền; không cần triển khai hai backend riêng.
- Khi triển khai, có thể phục vụ file build của frontend từ web server/CDN hoặc để backend phục vụ file tĩnh; cách triển khai không làm thay đổi việc tách mã nguồn hai thư mục.
- Nếu sau này cần PLC Gateway chạy độc lập, có thể thêm thư mục/thành phần riêng (ví dụ `plc-gateway/`); trước mắt phần Modbus có thể là module bên trong `backend/`.

### 12.3 Database
- Đề xuất SQLite cho phiên bản đầu để giảm công cài đặt và vận hành; Backend Go kết nối trực tiếp tới một file database đặt trên ổ đĩa cục bộ của máy chủ
- Bật foreign key, dùng transaction cho thao tác nghiệp vụ nhiều bảng, cấu hình WAL và busy timeout phù hợp; dùng migration có phiên bản để quản lý thay đổi schema
- Không đặt file SQLite trên thư mục mạng dùng chung, không cho nhiều bản Backend độc lập ghi cùng một file, và không để React truy cập trực tiếp file database
- Thiết kế schema relational rõ ràng với khóa ngoại giữa đơn hàng, công thức, mẻ trộn, nguyên vật liệu và tồn kho
- Redis chưa cần ở giai đoạn đầu; chỉ bổ sung nếu xuất hiện nhu cầu cache hoặc xử lý hàng đợi cụ thể
- Có thể chuyển sang PostgreSQL khi cần nhiều Backend instance, nhiều trạm cùng ghi vào một hệ thống, tần suất ghi đồng thời cao, HA/replication, hoặc khi kiểm thử tải cho thấy SQLite không đáp ứng

### 12.4 Kiến trúc đề xuất
- Frontend: React gọi Go Backend API qua HTTP/HTTPS; trình duyệt không kết nối trực tiếp tới PLC
- Go Backend: một ứng dụng xử lý đăng nhập, phân quyền, nghiệp vụ, SQLite và API cho giao diện; trong cùng ứng dụng, module PLC chạy nền kết nối Modbus TCP/RTU trong mạng nội bộ để đọc trạng thái/giá trị và thực hiện lệnh đơn giản đã được kiểm tra
- Module PLC phải xử lý timeout, mất kết nối/reconnect, kiểm tra dữ liệu đọc được và tránh lưu trùng một sự kiện/mẻ; các thao tác ghi PLC phải được kiểm tra quyền và điều kiện nghiệp vụ ở backend
- Chỉ tách PLC Gateway thành tiến trình riêng nếu yêu cầu vận hành sau này cần cô lập hoặc mở rộng độc lập; khi tách, dùng giao tiếp nội bộ có kiểm soát và không mở gateway công khai ra Internet
- Database giai đoạn đầu: SQLite là nguồn dữ liệu nghiệp vụ chính, đặt trên ổ đĩa cục bộ của máy chạy Backend; chỉ Backend chịu trách nhiệm đọc/ghi database
- Nếu PLC collector sau này tách thành tiến trình riêng, nó gửi dữ liệu cho Backend API hoặc một giao diện nội bộ có kiểm soát; tránh để nhiều tiến trình cùng ghi trực tiếp SQLite nếu không thực sự cần
- An toàn vận hành: PLC vẫn chịu trách nhiệm interlock và logic điều khiển thời gian thực; backend/gateway không thay thế các cơ chế an toàn tại PLC
- Triển khai: giai đoạn đầu chạy một tiến trình Go cho HTTP API, nghiệp vụ, SQLite và module PLC chạy nền; chỉ tách PLC Gateway khi có yêu cầu vận hành riêng
- Authentication: dùng access token dạng Bearer cho API; nếu sau này tách PLC Gateway thì bổ sung xác thực dịch vụ riêng
- Deployment: Docker Compose phù hợp giai đoạn đầu; chỉ cân nhắc Kubernetes khi quy mô vận hành thực sự cần

### 12.5 Khuyến nghị quyết định số lượng backend
- Giai đoạn đầu dùng một Go Backend API trong `backend/` để phục vụ React, quản lý SQLite và chạy module PLC nền; không cần hai backend riêng.
- Module PLC dùng kết nối nội bộ để đọc trạng thái và xử lý điều khiển đơn giản; nếu yêu cầu triển khai/độ sẵn sàng thay đổi, có thể tách module này thành PLC Gateway sau.
- Dù chạy chung hay tách riêng, chỉ backend/module PLC được phép giao tiếp trực tiếp với PLC; React chỉ gọi API.

### 12.6 Quy ước API giữa Backend và Frontend
API này là hợp đồng để Backend Go triển khai và Frontend React tích hợp. Các endpoint dưới đây là đặc tả dự kiến, chưa phải mã nguồn đã triển khai.

#### Quy ước chung
- Base path: `/api/v1`.
- Giao thức: HTTP trong phát triển nội bộ; HTTPS khi triển khai qua mạng.
- Request/response dùng JSON UTF-8; tên trường JSON dùng `snake_case` thống nhất với tên trường mô tả ở phần database.
- Ngày giờ dùng chuỗi ISO 8601 có múi giờ, lưu/truyền theo UTC, ví dụ `2026-10-10T08:30:00Z`. Ngày không có giờ dùng `YYYY-MM-DD`.
- ID là số nguyên; các khối lượng m³ truyền dưới dạng số, không dùng chuỗi đã định dạng có dấu phẩy.
- API biểu diễn cờ `is_active` bằng JSON boolean (`true`/`false`); Backend chuyển đổi phù hợp với SQLite (`1`/`0`).
- Frontend không gửi `id`, `created_at`, `updated_at` khi tạo bản ghi. Backend sinh ID và thời gian; Backend cập nhật `updated_at` khi sửa.
- Các trường tham chiếu như `customer_id`, `construction_site_id`, `truck_id` phải là ID hợp lệ và bản ghi được tham chiếu phải đang hoạt động tại thời điểm tạo/cập nhật.
- Khi tạo/cập nhật mác bê tông, các trường định lượng thành phần là bắt buộc; nếu frontend không có giá trị cho một thành phần thì gửi `0` (Backend cũng có thể áp dụng mặc định `0` khi trường bị bỏ qua). `pg1`–`pg6` không được có quá 4 chữ số sau dấu thập phân.

#### Xác thực
- `POST /api/v1/auth/login`: nhận `{ "username": "...", "password": "..." }`; trả access token và thông tin người dùng cơ bản.
- `POST /api/v1/auth/logout`: kết thúc phiên/token theo cơ chế backend triển khai.
- `GET /api/v1/auth/me`: trả thông tin tài khoản hiện tại.
- Các API nghiệp vụ yêu cầu header `Authorization: Bearer <access_token>`, trừ API đăng nhập.
- Backend kiểm tra quyền trên mọi request; ẩn nút ở frontend chỉ là hỗ trợ UX, không thay thế phân quyền phía server.

#### Định dạng response
- API đọc một bản ghi trả đối tượng JSON trực tiếp, ví dụ `{ "id": 12, "customer_code": "KH001" }`.
- API danh sách trả dạng phân trang:
```json
{
  "items": [],
  "page": 1,
  "page_size": 20,
  "total": 0
}
```
- Danh sách hỗ trợ `page`, `page_size`, `search`, `is_active`; giới hạn `page_size` tối đa 100. Backend áp dụng mặc định `is_active=true`, trừ khi request rõ ràng yêu cầu xem bản ghi đã xóa mềm và người dùng có quyền.
- Lỗi trả cấu trúc thống nhất:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Dữ liệu không hợp lệ",
    "details": [
      { "field": "customer_name", "message": "Trường này là bắt buộc" }
    ]
  }
}
```
- Dùng HTTP status phù hợp: `200` đọc/cập nhật thành công, `201` tạo thành công, `204` xóa mềm thành công, `400` request sai, `401` chưa xác thực, `403` không đủ quyền, `404` không tìm thấy, `409` trùng mã hoặc xung đột dữ liệu, `422` dữ liệu nghiệp vụ không hợp lệ, `500` lỗi hệ thống.

#### Endpoint danh mục và tài khoản
Các nhóm danh mục hỗ trợ `GET` danh sách, `GET /{id}`, `POST` tạo, `PATCH /{id}` cập nhật và `DELETE /{id}` xóa mềm:
- `/users` — tài khoản; tạo/cập nhật tài khoản cần được giới hạn theo quyền
- `/customers` — khách hàng
- `/mixer-trucks` — xe bồn
- `/drivers` — tài xế
- `/construction-sites` — công trường
- `/concrete-grades` — mác bê tông/công thức cấp phối; hỗ trợ đọc, tạo, cập nhật và xóa mềm
- `/waiting-mixes` — 20 vị trí cố định của danh sách cấp phối chờ; đọc và cập nhật vị trí, không tạo/xóa vị trí qua API

`/stations` là cấu hình của đúng một trạm trong phiên bản hiện tại, không phải danh mục nhiều trạm: `GET /api/v1/station` đọc cấu hình duy nhất; `PATCH /api/v1/station` cập nhật cấu hình đó. Không hỗ trợ tạo hoặc xóa trạm qua API. Backend khởi tạo một bản ghi cấu hình khi cài đặt lần đầu; lần cập nhật sau chỉ sửa bản ghi này và cập nhật `updated_at`.

Riêng `/roles` giai đoạn đầu chỉ cần `GET /roles` để frontend tải danh sách vai trò. Không cho xóa vật lý vai trò đang được tài khoản sử dụng. Nếu cần đặt/đổi mật khẩu, dùng endpoint riêng có quyền kiểm soát, không nhận `password_hash` từ frontend.

`DELETE` trên danh mục/tài khoản là xóa mềm: Backend đặt `is_active=false` (hoặc `0` trong SQLite), không xóa dòng dữ liệu. Nếu bản ghi đã được tham chiếu, Backend vẫn giữ tham chiếu lịch sử và ngăn không cho chọn lại bản ghi không hoạt động.

API xe bồn trả `default_driver_id`; có thể trả thêm object `default_driver` gồm `id`, `driver_code`, `full_name` để Frontend hiển thị ngay tài xế mặc định khi chọn xe. Khi lưu chuyến giao, Frontend gửi `driver_id` được chọn cuối cùng; Backend không tự ghi đè lựa chọn tài xế của chuyến bằng cấu hình mặc định.

#### Endpoint đơn hàng và chuyến giao
- `/orders`: `GET` danh sách, `POST` tạo đơn.
- `/orders/{id}`: `GET` chi tiết, `PATCH` cập nhật thông tin được phép, `DELETE` xóa mềm.
- `/orders/{id}/deliveries`: `GET` danh sách chuyến của đơn, `POST` tạo chuyến mới.
- `/deliveries/{id}`: `GET` chi tiết chuyến, `PATCH` cập nhật điều phối/trạng thái/số liệu chuyến theo quyền và quy tắc nghiệp vụ.

Tạo/cập nhật đơn nhận các trường nghiệp vụ như `order_number`, `customer_id`, `construction_site_id`, `concrete_grade_id`, `ordered_volume_m3`, `notes`. Backend kiểm tra số đơn duy nhất và các tham chiếu khách hàng/công trường/mác bê tông đang hoạt động.

`cumulative_volume_m3` là giá trị do Backend tính/cập nhật dựa trên quy trình dữ liệu đã xác nhận; Frontend chỉ đọc, không được gửi giá trị này để ghi trực tiếp. Khi tạo/cập nhật chuyến, Backend kiểm tra tổng khối lượng dự kiến không vượt khối lượng đặt; các thao tác cập nhật đơn và chuyến liên quan phải dùng transaction phù hợp.

- `GET /api/v1/waiting-mixes`: trả đủ 20 vị trí cấp phối chờ, sắp xếp theo `stt`.
- `PATCH /api/v1/waiting-mixes/{id}`: cập nhật thông tin của một vị trí; chỉ nhận các trường nghiệp vụ được phép.
- Backend bảo đảm luôn có đúng 20 vị trí, `stt` từ 1 đến 20 và không trùng; Frontend không tạo hoặc xóa vị trí.
- `status` chỉ nhận số nguyên `0` hoặc `1`. Quy tắc cụ thể về thời điểm một vị trí chuyển từ `1` (vừa trộn) về `0` (cũ) cần được thống nhất theo luồng thao tác/PLC trước khi triển khai.

#### Hành vi Frontend dựa trên API
- Khi tải form đơn/chuyến/cấp phối chờ, Frontend gọi API danh sách khách hàng, công trường, mác bê tông, xe, tài xế đang hoạt động.
- Khi chọn xe bồn, Frontend lấy `default_driver`/`default_driver_id` từ dữ liệu xe và tự điền tài xế; người dùng có thể chọn tài xế khác trước khi gửi request.
- Frontend hiển thị lỗi validation theo `details[].field`, thông báo lỗi chung theo `message`, và không tự giả lập thành công khi API lỗi.
- Frontend làm mới cache/dữ liệu liên quan sau khi tạo, cập nhật hoặc xóa mềm thành công.

## 13. Đề xuất thiết kế cơ sở dữ liệu ban đầu
### 13.1 Các bảng trong phạm vi mô tả hiện tại
Đây là mô tả yêu cầu dữ liệu ở mức thiết kế; chưa tạo file SQL hoặc migration. Tên bảng/cột là tên gợi ý và có thể điều chỉnh khi đối chiếu database, quy trình phần mềm C# hiện tại.

#### `stations` — Trạm trộn
- Mục đích: lưu thông tin và cấu hình vận hành của trạm hiện tại. Phiên bản đầu chỉ có đúng một bản ghi; các thay đổi cập nhật bản ghi này, không tạo thêm trạm.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa bản ghi cấu hình duy nhất; bản ghi phải luôn có cùng một ID sau khi khởi tạo |
| `station_code` | TEXT | Có | Mã trạm |
| `station_name` | TEXT | Có | Tên trạm |
| `address` | TEXT | Không | Địa chỉ trạm |
| `phone` | TEXT | Không | Số điện thoại liên hệ |
| `tp1_name`, `tp2_name`, `tp3_name`, `tp4_name`, `tp5_name`, `tp6_name` | TEXT | Có | Tên hiển thị cho từng vị trí thành phần TP1–TP6 |
| `xi1_name`, `xi2_name`, `xi3_name`, `xi4_name`, `xi5_name` | TEXT | Có | Tên hiển thị cho từng vị trí xi măng XI1–XI5 |
| `pg1_name`, `pg2_name`, `pg3_name`, `pg4_name`, `pg5_name`, `pg6_name` | TEXT | Có | Tên hiển thị cho từng vị trí phụ gia PG1–PG6 |
| `component_visibility_mask` | TEXT | Có | Chuỗi đúng 17 ký tự `0`/`1`, xác định vị trí thành phần nào được hiển thị trên phần mềm; vị trí đầu tiên TP1 mặc định bật (`1`), các vị trí còn lại cấu hình theo trạm |
| `accumulate_components` | INTEGER | Có | Cờ cộng dồn thành phần cốt liệu: `1` bật, `0` tắt; mặc định theo cấu hình trạm |
| `accumulate_cement` | INTEGER | Có | Cờ cộng dồn xi măng: `1` bật, `0` tắt; mặc định theo cấu hình trạm |
| `batch_min_volume_m3` | REAL | Có | Khối lượng mẻ trộn tối thiểu (m³) |
| `batch_max_volume_m3` | REAL | Có | Khối lượng mẻ trộn tối đa (m³), phải lớn hơn hoặc bằng khối lượng tối thiểu |
| `created_at` | TEXT | Có | Thời điểm tạo cấu hình |
| `updated_at` | TEXT | Có | Thời điểm cập nhật cấu hình gần nhất |

- Thứ tự 17 ký tự trong `component_visibility_mask` cố định là `TP1, TP2, TP3, TP4, TP5, TP6, XI1, XI2, XI3, XI4, XI5, PG1, PG2, PG3, PG4, PG5, PG6`; ký tự đầu tiên ứng với TP1. `1` nghĩa là hiển thị, `0` nghĩa là ẩn. Ví dụ `10001001000000000` nghĩa là chỉ TP1, TP5 và XI2 được hiển thị.
- Backend phải kiểm tra mask có đúng 17 ký tự và chỉ chứa `0`/`1`. Frontend dùng mask để hiển thị/ẩn các trường thành phần tương ứng.
- Ý nghĩa chi tiết của “cộng dồn thành phần” và “cộng dồn xi” cần được xác nhận theo logic vận hành/PLC hiện tại; ở mức cấu hình, các trường trên được xem là cờ bật/tắt.
- Khi cập nhật, Backend sửa bản ghi duy nhất này và giữ nguyên `id`/`created_at`; không cung cấp thao tác xóa mềm vì đây là cấu hình bắt buộc của hệ thống.

#### `users` — Tài khoản người dùng
- Mục đích: lưu tài khoản đăng nhập và thông tin người thao tác hệ thống.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `username` | TEXT | Có | Tên đăng nhập; không trùng |
| `password_hash` | TEXT | Có | Mật khẩu đã băm an toàn; không lưu mật khẩu gốc |
| `full_name` | TEXT | Có | Họ và tên người dùng |
| `phone` | TEXT | Không | Số điện thoại |
| `role_id` | INTEGER | Có | Vai trò của người dùng, tham chiếu `roles.id` |
| `station_id` | INTEGER | Không | Trạm được phân công, tham chiếu `stations.id`; để trống nếu là tài khoản toàn hệ thống |
| `is_active` | INTEGER | Có | `1` đang hoạt động, `0` bị khóa/ngừng hoạt động; mặc định `1` |
| `last_login_at` | TEXT | Không | Thời điểm đăng nhập gần nhất |
| `created_at` | TEXT | Có | Thời điểm tạo tài khoản |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |

#### `roles` — Vai trò
- Mục đích: xác định nhóm quyền của tài khoản.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `role_code` | TEXT | Có | Mã vai trò; không trùng, ví dụ `admin`, `dispatcher` |
| `role_name` | TEXT | Có | Tên vai trò hiển thị |
| `description` | TEXT | Không | Mô tả phạm vi công việc của vai trò |
| `created_at` | TEXT | Có | Thời điểm tạo bản ghi |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |

- Vai trò khởi đầu có thể gồm quản trị hệ thống, quản lý, điều phối và nhân viên vận hành. Quyền chi tiết theo từng chức năng có thể bổ sung sau nếu cần.

#### `concrete_grades` — Mác bê tông / cấp phối
- Mục đích: lưu một mác bê tông cùng định mức thành phần để sản xuất đúng **1 m³ bê tông**. Mỗi bản ghi là một công thức/định mức cho một khối bê tông, không phải một mẻ sản xuất thực tế.
- Các trường thành phần `tp1`–`tp6`, `xi1`–`xi5`, `nuoc`, `pg1`–`pg6` là trường bắt buộc và lưu khối lượng/định lượng nguyên liệu dùng cho 1 m³ theo đơn vị vận hành đã thống nhất. Nếu thành phần không dùng trong cấp phối thì giá trị mặc định là `0`, không để trống/NULL.
- Riêng `pg1`–`pg6` hỗ trợ tối đa 4 chữ số sau dấu thập phân. Backend cần từ chối giá trị có nhiều hơn 4 chữ số thập phân thay vì tự làm tròn gây thay đổi định lượng.
- Đề xuất ban đầu là kg/m³ cho cát, đá và xi măng; đơn vị phụ gia cần xác nhận theo cách PLC/công thức hiện có (kg, lít hoặc đơn vị khác).
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `grade_code` | TEXT | Có | Mã mác bê tông; không trùng |
| `printed_name` | TEXT | Có | Tên mác bê tông hiển thị/in trên phiếu |
| `tp1` | REAL | Có | Định lượng thành phần TP1 cho 1 m³ (nếu không dùng mặc định `0`; ví dụ cát/đá 1, cần ánh xạ tên vật liệu chính xác) |
| `tp2` | REAL | Có | Định lượng TP2 cho 1 m³; nếu không dùng mặc định `0` |
| `tp3` | REAL | Có | Định lượng TP3 cho 1 m³; nếu không dùng mặc định `0` |
| `tp4` | REAL | Có | Định lượng TP4 cho 1 m³; nếu không dùng mặc định `0` |
| `tp5` | REAL | Có | Định lượng TP5 cho 1 m³; nếu không dùng mặc định `0` |
| `tp6` | REAL | Có | Định lượng TP6 cho 1 m³; nếu không dùng mặc định `0` |
| `xi1` | REAL | Có | Định lượng xi măng loại 1 cho 1 m³; nếu không dùng mặc định `0` |
| `xi2` | REAL | Có | Định lượng xi măng loại 2 cho 1 m³; nếu không dùng mặc định `0` |
| `xi3` | REAL | Có | Định lượng xi măng loại 3 cho 1 m³; nếu không dùng mặc định `0` |
| `xi4` | REAL | Có | Định lượng xi măng loại 4 cho 1 m³; nếu không dùng mặc định `0` |
| `xi5` | REAL | Có | Định lượng xi măng loại 5 cho 1 m³; nếu không dùng mặc định `0` |
| `nuoc` | REAL | Có | Định lượng nước cho 1 m³; nếu không dùng mặc định `0` |
| `pg1` | REAL | Có | Định lượng phụ gia loại 1 cho 1 m³, tối đa 4 chữ số thập phân; nếu không dùng mặc định `0` |
| `pg2` | REAL | Có | Định lượng phụ gia loại 2 cho 1 m³, tối đa 4 chữ số thập phân; nếu không dùng mặc định `0` |
| `pg3` | REAL | Có | Định lượng phụ gia loại 3 cho 1 m³, tối đa 4 chữ số thập phân; nếu không dùng mặc định `0` |
| `pg4` | REAL | Có | Định lượng phụ gia loại 4 cho 1 m³, tối đa 4 chữ số thập phân; nếu không dùng mặc định `0` |
| `pg5` | REAL | Có | Định lượng phụ gia loại 5 cho 1 m³, tối đa 4 chữ số thập phân; nếu không dùng mặc định `0` |
| `pg6` | REAL | Có | Định lượng phụ gia loại 6 cho 1 m³, tối đa 4 chữ số thập phân; nếu không dùng mặc định `0` |
| `do_sut` | REAL | Có | Độ sụt yêu cầu; đơn vị cần thống nhất (thường ghi cm hoặc mm) |
| `cuong_do` | REAL | Có | Cường độ bê tông; đơn vị/cách ghi theo tiêu chuẩn trạm (ví dụ MPa hoặc cấp/mác) |
| `notes` | TEXT | Không | Ghi chú về cấp phối |
| `created_at` | TEXT | Có | Thời điểm tạo bản ghi |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |
| `is_active` | INTEGER | Có | Cờ xóa mềm: `1` đang sử dụng, `0` đã xóa mềm; mặc định `1` |

- TP1–TP6, XI1–XI5 và PG1–PG6 là các vị trí thành phần như cách gọi hiện tại; trước khi triển khai cần xác nhận mỗi vị trí tương ứng chính xác với vật liệu nào và đơn vị gì. Nếu danh sách vật liệu thay đổi thường xuyên, có thể chuẩn hóa thành bảng thành phần cấp phối riêng ở giai đoạn thiết kế database.
- Không xóa vật lý mác/cấp phối đã được đơn hàng hoặc sản xuất sử dụng; chuyển `is_active` về `0` và giữ bản ghi để bảo toàn lịch sử. Khi sửa cấp phối đã dùng, cần cân nhắc tạo phiên bản mới để các đơn hàng/mẻ cũ tiếp tục tham chiếu đúng công thức đã áp dụng.

#### `waiting_mixes` — Cấp phối chờ
- Mục đích: lưu tối đa 20 vị trí cấp phối chờ để người dùng chọn/điều phối cho sản xuất. Mỗi dòng là một vị trí cố định có số thứ tự; không phải lịch sử các mẻ đã trộn.
- Hệ thống khởi tạo đủ 20 dòng với `stt` từ 1 đến 20. Khi vị trí không dùng, vẫn giữ dòng đó và cập nhật thông tin/trạng thái thay vì xóa dòng hoặc thay đổi số thứ tự.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ của vị trí |
| `stt` | INTEGER | Có | Số thứ tự vị trí, duy nhất trong khoảng 1–20 |
| `concrete_grade_id` | INTEGER | Có | Mác/cấp phối, tham chiếu `concrete_grades.id` |
| `customer_id` | INTEGER | Có | Khách hàng, tham chiếu `customers.id` |
| `volume_m3` | REAL | Có | Khối lượng cấp phối/đơn chờ, tính bằng m³ |
| `m3me` | REAL | Có | Khối lượng m³ cho một mẻ; ý nghĩa chính xác cần xác nhận theo cách gọi tại trạm |
| `status` | INTEGER | Có | Chỉ nhận `0` hoặc `1`: `1` là vừa trộn, `0` là cũ rồi |

- Trường `concrete_grade_id` và `customer_id` lưu ID tham chiếu; API có thể trả thêm thông tin tên/mã tương ứng để giao diện hiển thị.
- Khi chọn vị trí cấp phối chờ, Backend phải kiểm tra mác bê tông và khách hàng đang hoạt động, khối lượng là số dương, `m3me` không âm, và `status` chỉ là `0` hoặc `1`.
- Quy tắc tự động chuyển `status` từ `1` sang `0` (ví dụ khi một vị trí khác vừa trộn) chưa được xác định; cần chốt trước khi nối PLC/triển khai nghiệp vụ này.

#### `customers` — Khách hàng
- Mục đích: lưu thông tin khách hàng để lập và tra cứu đơn hàng.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `customer_code` | TEXT | Có | Mã khách hàng dùng để tra cứu; không trùng |
| `customer_name` | TEXT | Có | Tên cá nhân hoặc tên công ty |
| `address` | TEXT | Không | Địa chỉ khách hàng |
| `phone` | TEXT | Không | Số điện thoại liên hệ |
| `created_at` | TEXT | Có | Thời điểm tạo bản ghi |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |
| `is_active` | INTEGER | Có | Cờ xóa mềm: `1` là bản ghi đang sử dụng, `0` là đã xóa mềm; mặc định `1` |

- Khi xóa khách hàng, Backend cập nhật `is_active` từ `1` thành `0`, không xóa vật lý dòng dữ liệu. Danh sách khách hàng thông thường chỉ hiển thị bản ghi `is_active = 1`; dữ liệu đã tham chiếu bởi đơn hàng vẫn được giữ để bảo toàn lịch sử.
- Kiểu dữ liệu và tên trường hiện là đề xuất; sẽ chốt khi bắt đầu tạo database và đối chiếu dữ liệu thực tế.

#### `mixer_trucks` — Xe bồn
- Mục đích: quản lý xe dùng để giao bê tông.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `truck_code` | TEXT | Có | Mã xe nội bộ; không trùng |
| `plate_number` | TEXT | Có | Biển số xe; không trùng |
| `default_driver_id` | INTEGER | Không | Tài xế mặc định của xe, tham chiếu `drivers.id`; có thể để trống nếu xe chưa được gán tài xế |
| `created_at` | TEXT | Có | Thời điểm tạo bản ghi |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |
| `is_active` | INTEGER | Có | Cờ xóa mềm: `1` là xe đang sử dụng, `0` là đã xóa mềm; mặc định `1` |

- Khi xóa xe bồn, Backend cập nhật `is_active` từ `1` thành `0`, không xóa vật lý dòng dữ liệu. Xe đã được dùng trong chuyến giao vẫn được giữ để bảo toàn lịch sử.
- `default_driver_id` là tài xế gợi ý mặc định của xe, không có nghĩa xe luôn do tài xế đó lái trong mọi chuyến.

#### `drivers` — Tài xế
- Mục đích: lưu hồ sơ tài xế được phân công chuyến giao; tài xế không bắt buộc phải có tài khoản đăng nhập.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `driver_code` | TEXT | Có | Mã tài xế nội bộ; không trùng |
| `full_name` | TEXT | Có | Họ và tên tài xế |
| `created_at` | TEXT | Có | Thời điểm tạo hồ sơ |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |
| `is_active` | INTEGER | Có | Cờ xóa mềm: `1` là tài xế đang sử dụng, `0` là đã xóa mềm; mặc định `1` |

- Khi xóa tài xế, Backend cập nhật `is_active` từ `1` thành `0`, không xóa vật lý dòng dữ liệu. Tài xế đã được gán vào chuyến giao vẫn được giữ để bảo toàn lịch sử.

#### `construction_sites` — Công trường
- Mục đích: lưu danh sách công trường/địa điểm thi công để chọn khi lập đơn hàng.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `site_code` | TEXT | Có | Mã công trường; không trùng |
| `site_name` | TEXT | Có | Tên công trường/dự án |
| `created_at` | TEXT | Có | Thời điểm tạo bản ghi |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |
| `is_active` | INTEGER | Có | Cờ xóa mềm: `1` là công trường đang sử dụng, `0` là đã xóa mềm; mặc định `1` |

- Khi xóa công trường, Backend cập nhật `is_active` từ `1` thành `0`, không xóa vật lý dữ liệu để giữ lịch sử đơn hàng.

#### `orders` — Đơn hàng
- Mục đích: lưu yêu cầu đặt bê tông của khách hàng.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `order_number` | TEXT | Có | Số/mã đơn hàng; không trùng |
| `customer_id` | INTEGER | Có | Khách hàng, tham chiếu `customers.id` |
| `construction_site_id` | INTEGER | Có | Công trường, tham chiếu `construction_sites.id` |
| `concrete_grade_id` | INTEGER | Có | Mác/cấp phối bê tông, tham chiếu `concrete_grades.id` |
| `ordered_volume_m3` | REAL | Có | Khối lượng khách đặt, tính bằng m³ |
| `cumulative_volume_m3` | REAL | Có | Khối lượng đã thực hiện/giao lũy kế cho đơn, tính bằng m³; khởi tạo bằng `0` |
| `is_active` | INTEGER | Có | Cờ xóa mềm: `1` là đơn đang hoạt động, `0` là đã xóa mềm; mặc định `1` |
| `created_at` | TEXT | Có | Thời điểm tạo bản ghi |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |
| `notes` | TEXT | Không | Ghi chú nghiệp vụ |

- `cumulative_volume_m3` nên được Backend tính/cập nhật từ dữ liệu thực tế của các chuyến giao hoặc mẻ sản xuất đã xác nhận, không nhập tùy ý từ giao diện. Cần chốt quy trình nào là nguồn chính thức trước khi triển khai.
- Khi xóa mềm đơn hàng, Backend cập nhật `is_active` từ `1` thành `0`, không xóa vật lý; các chuyến giao và dữ liệu liên quan vẫn được giữ để tra cứu lịch sử.

#### `order_deliveries` — Chuyến giao hàng
- Mục đích: tách các chuyến giao khỏi đơn hàng để hỗ trợ giao một đơn thành nhiều chuyến hoặc nhiều xe.
- Các trường dự kiến:

| Tên trường | Kiểu dữ liệu dự kiến | Bắt buộc | Ý nghĩa |
|---|---|---:|---|
| `id` | INTEGER | Có | Khóa chính nội bộ, tự tăng |
| `order_id` | INTEGER | Có | Đơn hàng, tham chiếu `orders.id` |
| `trip_no` | INTEGER | Có | Số thứ tự chuyến trong đơn |
| `truck_id` | INTEGER | Không | Xe bồn được điều phối, tham chiếu `mixer_trucks.id` |
| `driver_id` | INTEGER | Không | Tài xế được điều phối, tham chiếu `drivers.id` |
| `planned_volume_m3` | REAL | Có | Khối lượng dự kiến chở ở chuyến này |
| `loaded_volume_m3` | REAL | Không | Khối lượng thực tế đã nạp lên xe |
| `delivered_volume_m3` | REAL | Không | Khối lượng thực tế đã giao |
| `status` | TEXT | Có | Trạng thái chuyến; mặc định `planned` (đã lên kế hoạch) |
| `departed_at` | TEXT | Không | Thời điểm xe rời trạm |
| `delivered_at` | TEXT | Không | Thời điểm hoàn tất giao |
| `notes` | TEXT | Không | Ghi chú/chênh lệch/sự cố của chuyến |
| `created_at` | TEXT | Có | Thời điểm tạo chuyến |
| `updated_at` | TEXT | Có | Thời điểm cập nhật gần nhất |

- Xe và tài xế có thể chưa được gán khi chuyến mới được tạo và đang chờ điều phối.
- Khi người dùng chọn xe bồn trên giao diện, hệ thống tự điền tài xế mặc định của xe vào `driver_id`; người dùng được phép đổi tài xế trước khi lưu chuyến.
- `order_deliveries.driver_id` lưu tài xế thực tế được chọn cho chuyến đó. Nếu cấu hình tài xế mặc định của xe thay đổi về sau, không thay đổi tài xế đã lưu trên các chuyến cũ.

### 13.2 Quan hệ chính
- Một trạm cấu hình duy nhất có thể được tham chiếu bởi người dùng, xe bồn và tài xế; phiên bản đầu không quản lý danh mục nhiều trạm vận hành.
- Một khách hàng có thể có nhiều đơn hàng; mỗi đơn hàng thuộc một khách hàng.
- Một công trường có thể được tham chiếu bởi nhiều đơn hàng; mỗi đơn hàng thuộc một công trường.
- Một mác/cấp phối bê tông có thể được dùng trong nhiều đơn hàng; mỗi đơn hàng tham chiếu một mác/cấp phối.
- Mỗi vị trí trong `waiting_mixes` tham chiếu một mác bê tông và một khách hàng.
- Một đơn hàng có thể có nhiều chuyến giao; mỗi chuyến thuộc đúng một đơn hàng.
- Mỗi chuyến có thể được gán một xe bồn và một tài xế.
- Một xe bồn có thể có một tài xế mặc định để giao diện tự điền khi điều phối; tài xế chuyến giao vẫn được phép khác tài xế mặc định.
- Mỗi tài khoản có một vai trò; tài khoản quản trị toàn hệ thống có thể không gắn cố định với một trạm.
- Mỗi trường `*_id` nêu trong mô tả là khóa ngoại tới bảng tương ứng; Backend cần kiểm tra các bản ghi được tham chiếu còn hợp lệ và đang hoạt động khi phân công hoặc tạo nghiệp vụ.

### 13.3 Quy tắc dữ liệu và nghiệp vụ dự kiến
- Số đơn hàng, mã khách hàng, mã trạm, mã xe, mã tài xế, mã công trường và mã mác bê tông phải duy nhất.
- Không xóa vật lý dữ liệu đã được dùng trong nghiệp vụ; chuyển sang trạng thái ngừng hoạt động để giữ lịch sử.
- Tổng khối lượng dự kiến của các chuyến giao không được vượt quá khối lượng đơn hàng; Backend cần kiểm tra khi tạo hoặc sửa chuyến.
- Thời điểm lưu trong database theo ISO 8601 UTC; giao diện hiển thị theo múi giờ địa phương.
- Các thao tác tạo/sửa nhiều bản ghi liên quan cần được thực hiện nguyên tử bằng transaction.

### 13.4 Phạm vi thiết kế tiếp theo
Các bảng nguyên vật liệu, sản xuất/mẻ trộn, kiểm soát chất lượng, audit log và dữ liệu PLC sẽ mô tả sau khi người dùng cung cấp thêm quy trình và cấu trúc database hiện tại. Chỉ tạo SQL/migration khi bắt đầu triển khai phần đó.

## 14. Lịch sử cập nhật
| Phiên bản | Ngày | Mô tả |
|----------|------|-------|
| 0.1 | 2026-10-10 | Khởi tạo bản nháp đầu tiên |
| 0.2 | 2026-10-10 | Chuyển sang đặc tả phù hợp với hệ thống trạm trộn bê tông và kiến trúc React + Go |
| 0.3 | 2026-10-10 | Bổ sung khuyến nghị Backend API và PLC Gateway cho tích hợp Modbus |
| 0.4 | 2026-10-10 | Đề xuất SQLite cho giai đoạn đầu và nêu điều kiện chuyển sang PostgreSQL |
| 0.5 | 2026-10-10 | Mô tả yêu cầu các bảng tài khoản, đơn hàng, khách hàng, xe bồn, tài xế và trạm; chưa tạo SQL |
| 0.6 | 2026-10-10 | Bổ sung tên trường, kiểu dữ liệu dự kiến, tính bắt buộc và ý nghĩa cho từng bảng trong phạm vi hiện tại |
| 0.7 | 2026-10-10 | Bổ sung bảng công trường và liên kết công trường với đơn hàng |
| 0.8 | 2026-10-10 | Điều chỉnh các trường bảng đơn hàng theo yêu cầu: khối lượng đặt, lũy kế, cờ xóa mềm và ghi chú |
| 0.9 | 2026-10-10 | Quy định tách mã nguồn React frontend và Go backend thành hai thư mục riêng |
| 0.10 | 2026-10-10 | Bổ sung quy ước API Backend-Frontend, endpoint, xác thực, phân trang, lỗi và hành vi form |
| 0.11 | 2026-10-10 | Bổ sung bảng mác bê tông/cấp phối theo định lượng cho 1 m³ và liên kết với đơn hàng/API |
| 0.12 | 2026-10-10 | Quy định trường thành phần cấp phối bắt buộc mặc định 0 và phụ gia tối đa 4 chữ số thập phân |
| 0.13 | 2026-10-10 | Làm rõ thành phần không sử dụng được lưu mặc định bằng 0; API từ chối phụ gia vượt quá 4 chữ số thập phân |
| 0.14 | 2026-10-10 | Chuyển bảng trạm thành cấu hình một bản ghi, bổ sung tên/hiển thị thành phần, cờ cộng dồn và giới hạn mẻ |
| 0.15 | 2026-10-10 | Làm rõ mặt nạ hiển thị 17 vị trí thành phần và mặc định bật TP1 |
| 0.16 | 2026-10-10 | Bổ sung bảng cấp phối chờ gồm 20 vị trí cố định và API đọc/cập nhật |

## 15. Ghi chú
Đây là bản đặc tả đầu tiên theo hướng hệ thống thực tế của trạm trộn bê tông. Tài liệu có thể tiếp tục được cập nhật theo cấu trúc database, quy trình operation, và yêu cầu nghiệp vụ chi tiết mà bạn đang có sẵn.
