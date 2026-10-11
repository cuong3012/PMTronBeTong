import { Card, Col, Row, Statistic, Typography, Empty } from 'antd'
import {
  TeamOutlined,
  EnvironmentOutlined,
  IdcardOutlined,
  CarOutlined,
} from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/store/auth'
import { customersApi } from '@/api/customers'
import { sitesApi } from '@/api/sites'
import { driversApi } from '@/api/drivers'
import { trucksApi } from '@/api/trucks'

// Dashboard hiển thị các chỉ số nhanh lấy từ các API danh mục.
export default function DashboardPage() {
  const user = useAuthStore((s) => s.user)

  // Lấy tổng số bản ghi của mỗi danh mục (chỉ cần 1 bản ghi/trang).
  const customers = useQuery({
    queryKey: ['customers', { page: 1, per_page: 1 }],
    queryFn: () => customersApi.list({ page: 1, per_page: 1 }),
  })
  const sites = useQuery({
    queryKey: ['construction-sites', { page: 1, per_page: 1 }],
    queryFn: () => sitesApi.list({ page: 1, per_page: 1 }),
  })
  const drivers = useQuery({
    queryKey: ['drivers', { page: 1, per_page: 1 }],
    queryFn: () => driversApi.list({ page: 1, per_page: 1 }),
  })
  const trucks = useQuery({
    queryKey: ['mixer-trucks', { page: 1, per_page: 1 }],
    queryFn: () => trucksApi.list({ page: 1, per_page: 1 }),
  })

  return (
    <div>
      <Typography.Title level={3}>
        Xin chào, {user?.full_name ?? 'người dùng'}
      </Typography.Title>
      <Typography.Paragraph type="secondary">
        Tổng quan nhanh các danh mục đang quản lý trong hệ thống.
      </Typography.Paragraph>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Khách hàng"
              value={customers.data?.total ?? 0}
              prefix={<TeamOutlined />}
              loading={customers.isLoading}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Công trường"
              value={sites.data?.total ?? 0}
              prefix={<EnvironmentOutlined />}
              loading={sites.isLoading}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Tài xế"
              value={drivers.data?.total ?? 0}
              prefix={<IdcardOutlined />}
              loading={drivers.isLoading}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Xe bồn"
              value={trucks.data?.total ?? 0}
              prefix={<CarOutlined />}
              loading={trucks.isLoading}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ marginTop: 24 }}>
        <Empty
          description="Các mô-đun đơn hàng, công thức, sản xuất và báo cáo sẽ được bổ sung sau"
        />
      </Card>
    </div>
  )
}
