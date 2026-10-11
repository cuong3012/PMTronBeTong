import { useEffect, useState } from 'react'
import { Alert, Button, Space } from 'antd'
import { ReloadOutlined } from '@ant-design/icons'
import { BACKEND_UNREACHABLE_MESSAGE } from '@/api/client'
import { checkBackendHealth } from '@/api/health'

/**
 * Banner cảnh báo khi không kết nối được backend.
 *
 * Tự thăm dò endpoint /health định kỳ. Chỉ hiển thị khi backend không phản hồi
 * (backend tắt, mất mạng, sai địa chỉ). Khi backend hoạt động lại thì tự ẩn.
 */
export default function BackendStatusBanner() {
  const [online, setOnline] = useState<boolean | null>(null)

  const probe = async () => {
    setOnline(await checkBackendHealth())
  }

  useEffect(() => {
    probe()
    const timer = setInterval(probe, 10_000)
    return () => clearInterval(timer)
  }, [])

  // Chưa kiểm tra xong hoặc backend vẫn ổn: không hiển thị gì.
  if (online !== false) return null

  return (
    <Alert
      type="error"
      showIcon
      banner
      message="Mất kết nối tới máy chủ"
      description={
        <Space>
          <span>{BACKEND_UNREACHABLE_MESSAGE}</span>
          <Button size="small" icon={<ReloadOutlined />} onClick={probe}>
            Kiểm tra lại
          </Button>
        </Space>
      }
      style={{ marginBottom: 16 }}
    />
  )
}
