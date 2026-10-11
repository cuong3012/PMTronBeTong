import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'

// Bọc các route yêu cầu đăng nhập; chưa có token thì chuyển về /login.
export default function ProtectedRoute({
  children,
}: {
  children: React.ReactNode
}) {
  const token = useAuthStore((s) => s.token)
  const location = useLocation()

  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  return <>{children}</>
}
