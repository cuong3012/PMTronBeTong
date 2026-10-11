import { apiClient } from './client'
import type { LoginResult, UserInfo } from '@/types'

// API xác thực (khớp backend/internal/auth).
export const authApi = {
  /** Đăng nhập; trả token và thông tin người dùng. */
  async login(username: string, password: string): Promise<LoginResult> {
    const { data } = await apiClient.post<LoginResult>('/auth/login', {
      username,
      password,
    })
    return data
  },

  /** Lấy thông tin người dùng hiện tại từ token. */
  async me(): Promise<UserInfo> {
    const { data } = await apiClient.get<UserInfo>('/auth/me')
    return data
  },
}
