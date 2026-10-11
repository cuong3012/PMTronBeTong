import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { UserInfo } from '@/types'

// Trạng thái xác thực toàn cục.
// Dùng persist để giữ token/người dùng qua các lần tải lại trang (localStorage).
interface AuthState {
  token: string | null
  user: UserInfo | null
  /** Lưu phiên sau khi đăng nhập thành công. */
  setSession: (token: string, user: UserInfo) => void
  /** Xóa phiên (đăng xuất hoặc token hết hạn). */
  logout: () => void
  /** Kiểm tra người dùng hiện tại có thuộc một trong các vai trò cho trước. */
  hasRole: (...roles: string[]) => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      setSession: (token, user) => set({ token, user }),
      logout: () => set({ token: null, user: null }),
      hasRole: (...roles) => {
        const role = get().user?.role_code
        return !!role && roles.includes(role)
      },
    }),
    { name: 'tramtron-auth' },
  ),
)

/** Vai trò được phép xóa mềm trên backend (xem middleware.RequireRole). */
export const ADMIN_MANAGER_ROLES = ['admin', 'manager']
