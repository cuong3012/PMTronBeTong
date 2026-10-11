import { Layout, Menu, Typography, Dropdown, Avatar, Space, Button } from 'antd'
import {
  DashboardOutlined,
  TeamOutlined,
  EnvironmentOutlined,
  IdcardOutlined,
  CarOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import BackendStatusBanner from '@/components/BackendStatusBanner'
import StationRunButton from '@/components/StationRunButton'

const { Header, Sider, Content } = Layout

// Menu điều hướng chính; key khớp đường dẫn route.
const menuItems = [
  { key: '/', icon: <DashboardOutlined />, label: 'Tổng quan' },
  { key: '/customers', icon: <TeamOutlined />, label: 'Khách hàng' },
  {
    key: '/construction-sites',
    icon: <EnvironmentOutlined />,
    label: 'Công trường',
  },
  { key: '/drivers', icon: <IdcardOutlined />, label: 'Tài xế' },
  { key: '/mixer-trucks', icon: <CarOutlined />, label: 'Xe bồn' },
]

// Layout chính sau khi đăng nhập: sidebar + header + nội dung.
export default function AppLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)

  // Xác định mục menu đang chọn dựa trên path (khớp cả route con).
  const selectedKey =
    menuItems.find((m) => m.key !== '/' && location.pathname.startsWith(m.key))
      ?.key ?? '/'

  const handleLogout = () => {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark" breakpoint="lg" collapsedWidth="0">
        <div
          style={{
            color: '#fff',
            fontWeight: 700,
            textAlign: 'center',
            padding: '16px 8px',
            fontSize: 15,
          }}
        >
          TRẠM TRỘN BÊ TÔNG
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>

      <Layout>
        <Header
          style={{
            background: '#fff',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            paddingInline: 24,
          }}
        >
          <Typography.Title level={4} style={{ margin: 0 }}>
            Hệ thống quản lý trạm trộn
          </Typography.Title>

          <Dropdown
            menu={{
              items: [
                {
                  key: 'logout',
                  icon: <LogoutOutlined />,
                  label: 'Đăng xuất',
                  onClick: handleLogout,
                },
              ],
            }}
          >
            <Space style={{ cursor: 'pointer' }}>
              <Avatar icon={<UserOutlined />} />
              <span>
                {user?.full_name}
                {user?.role_name ? ` (${user.role_name})` : ''}
              </span>
              <Button type="text" icon={<LogoutOutlined />} />
            </Space>
          </Dropdown>
        </Header>

        <Content className="app-content">
          <BackendStatusBanner />
          <Outlet />
        </Content>
      </Layout>

      {/* Nút nổi góc dưới trái: chuyển nhanh tới màn hình chạy trạm trộn. */}
      <StationRunButton />
    </Layout>
  )
}
