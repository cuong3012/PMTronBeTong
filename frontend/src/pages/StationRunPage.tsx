import { Card, Typography, Result } from 'antd'
import StationRunButton from '@/components/StationRunButton'

// Trang "Chạy trạm trộn" — khu vực vận hành sản xuất.
//
// Trang hiển thị TOÀN MÀN HÌNH (không nằm trong AppLayout nên không có sidebar),
// giúp người vận hành có tối đa không gian. Dùng nút nổi đỏ ở góc dưới trái
// (StationRunButton) để quay lại trang quản lý mà không tải lại trang.
//
// Hiện là khung sườn (placeholder): phần điều khiển PLC, mẻ trộn và theo dõi
// sản lượng theo thời gian thực sẽ được triển khai sau theo đặc tả.
export default function StationRunPage() {
  return (
    <div
      style={{
        minHeight: '100vh',
        background: '#0b1220',
        color: '#fff',
        padding: 24,
      }}
    >
      <Typography.Title level={3} style={{ marginTop: 0, color: '#fff' }}>
        Chạy trạm trộn
      </Typography.Title>
      <Typography.Paragraph style={{ color: 'rgba(255,255,255,0.65)' }}>
        Khu vực vận hành trạm trộn bê tông theo thời gian thực.
      </Typography.Paragraph>

      <Card>
        <Result
          status="info"
          title="Tính năng đang được phát triển"
          subTitle="Màn hình điều khiển mẻ trộn, theo dõi sản lượng và kết nối PLC sẽ được bổ sung ở giai đoạn tiếp theo."
        />
      </Card>

      {/* Nút nổi đỏ ở góc dưới trái để quay về trang quản lý. */}
      <StationRunButton />
    </div>
  )
}
