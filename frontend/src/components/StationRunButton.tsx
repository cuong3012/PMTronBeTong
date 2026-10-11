import { Button, Tooltip } from 'antd'
import { PlayCircleOutlined, AppstoreOutlined } from '@ant-design/icons'
import { useLocation, useNavigate } from 'react-router-dom'

// Đường dẫn của trang chạy trạm trộn.
const RUN_PATH = '/run'

/**
 * Nút nổi ở góc dưới bên trái để chuyển nhanh giữa khu quản lý và trang chạy
 * trạm trộn — không tải lại trang (điều hướng SPA bằng React Router).
 *
 * - Khi KHÔNG ở trang chạy trạm: nút đỏ "Chạy trạm trộn" -> điều hướng tới /run.
 * - Khi ĐANG ở trang chạy trạm: nút đổi thành "Về quản lý" -> quay về /.
 */
export default function StationRunButton() {
  const navigate = useNavigate()
  const location = useLocation()

  const isRunPage = location.pathname.startsWith(RUN_PATH)

  const handleClick = () => {
    navigate(isRunPage ? '/' : RUN_PATH)
  }

  return (
    <Tooltip
      title={isRunPage ? 'Quay về trang quản lý' : 'Mở màn hình chạy trạm trộn'}
      placement="right"
    >
      <Button
        type="primary"
        danger
        size="large"
        icon={isRunPage ? <AppstoreOutlined /> : <PlayCircleOutlined />}
        onClick={handleClick}
        style={{
          position: 'fixed',
          left: 24,
          bottom: 24,
          zIndex: 1000,
          height: 52,
          paddingInline: 20,
          fontSize: 16,
          fontWeight: 600,
          boxShadow: '0 6px 16px rgba(0,0,0,0.25)',
        }}
      >
        {isRunPage ? 'Về quản lý' : 'Chạy trạm trộn'}
      </Button>
    </Tooltip>
  )
}
