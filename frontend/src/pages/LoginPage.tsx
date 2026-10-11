import { useEffect, useState } from 'react'
import { Button, Card, Form, Input, Typography, Alert, Space, message } from 'antd'
import { LockOutlined, UserOutlined, ReloadOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { authApi } from '@/api/auth'
import {
  BACKEND_UNREACHABLE_MESSAGE,
  getErrorMessage,
  isNetworkError,
} from '@/api/client'
import { useAuthStore } from '@/store/auth'
import { checkBackendHealth } from '@/api/health'

interface LoginForm {
  username: string
  password: string
}

// Trang đăng nhập: gọi POST /api/auth/login và lưu phiên vào store.
export default function LoginPage() {
  const navigate = useNavigate()
  const setSession = useAuthStore((s) => s.setSession)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  // Cờ cho biết backend có đang kết nối được hay không (null = chưa kiểm tra).
  const [backendOnline, setBackendOnline] = useState<boolean | null>(null)

  // Kiểm tra sức khỏe backend ngay khi mở trang và định kỳ.
  // Giúp cảnh báo sớm cho người dùng trước khi họ bấm đăng nhập.
  const probeBackend = async () => {
    const ok = await checkBackendHealth()
    setBackendOnline(ok)
    return ok
  }

  useEffect(() => {
    probeBackend()
    const timer = setInterval(probeBackend, 10_000)
    return () => clearInterval(timer)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const onFinish = async (values: LoginForm) => {
    setLoading(true)
    setError(null)
    try {
      const result = await authApi.login(values.username, values.password)
      setSession(result.token, result.user)
      setBackendOnline(true)
      message.success('Đăng nhập thành công')
      navigate('/', { replace: true })
    } catch (err) {
      // Nếu không kết nối được backend, cập nhật cờ để hiện cảnh báo.
      if (isNetworkError(err)) {
        setBackendOnline(false)
      }
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #1677ff 0%, #0958d9 100%)',
      }}
    >
      <Card style={{ width: 400, boxShadow: '0 8px 24px rgba(0,0,0,0.2)' }}>
        <Typography.Title level={3} style={{ textAlign: 'center' }}>
          Đăng nhập
        </Typography.Title>
        <Typography.Paragraph style={{ textAlign: 'center', color: '#888' }}>
          Hệ thống quản lý trạm trộn bê tông
        </Typography.Paragraph>

        {backendOnline === false && (
          <Alert
            type="warning"
            showIcon
            style={{ marginBottom: 16 }}
            message="Không kết nối được máy chủ"
            description={
              <Space direction="vertical" size={4}>
                <span>{BACKEND_UNREACHABLE_MESSAGE}</span>
                <Button
                  size="small"
                  icon={<ReloadOutlined />}
                  onClick={probeBackend}
                >
                  Kiểm tra lại
                </Button>
              </Space>
            }
          />
        )}

        {error && (
          <Alert
            type="error"
            message={error}
            showIcon
            style={{ marginBottom: 16 }}
          />
        )}

        <Form<LoginForm>
          layout="vertical"
          onFinish={onFinish}
          autoComplete="off"
          requiredMark={false}
        >
          <Form.Item
            name="username"
            label="Tên đăng nhập"
            rules={[{ required: true, message: 'Vui lòng nhập tên đăng nhập' }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="admin"
              size="large"
            />
          </Form.Item>

          <Form.Item
            name="password"
            label="Mật khẩu"
            rules={[{ required: true, message: 'Vui lòng nhập mật khẩu' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="••••••••"
              size="large"
            />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0 }}>
            <Button
              type="primary"
              htmlType="submit"
              block
              size="large"
              loading={loading}
            >
              Đăng nhập
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}
