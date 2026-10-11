import axios, { AxiosError } from 'axios'
import { useAuthStore } from '@/store/auth'

// Instance axios dùng chung cho toàn bộ API.
//
// Backend Go phục vụ API dưới tiền tố /api. Khi phát triển, Vite proxy
// chuyển /api sang http://localhost:8080 (xem vite.config.ts).
export const apiClient = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor: tự gắn Bearer token nếu đã đăng nhập.
apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor: bắt lỗi 401 để đăng xuất và điều hướng về login.
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      // Token hết hạn/không hợp lệ: xóa phiên hiện tại.
      useAuthStore.getState().logout()
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

/** Định dạng lỗi thống nhất từ backend: { error, message }. */
interface ApiErrorBody {
  error?: string
  message?: string
  /** Mã lỗi ổn định (chỉ có ở một số lỗi đặc biệt, ví dụ INACTIVE_CODE_EXISTS). */
  code?: string
  /** Dữ liệu bổ sung kèm theo lỗi (ví dụ bản ghi cũ đã bị xóa mềm). */
  data?: unknown
}

/**
 * Mã lỗi backend trả về khi mã/biển số trùng với một bản ghi đã bị xóa mềm.
 * FE dùng mã này để hỏi người dùng có muốn khôi phục bản ghi cũ hay không.
 */
export const INACTIVE_CODE_EXISTS = 'INACTIVE_CODE_EXISTS'

/** Kết quả trích xuất khi gặp lỗi trùng mã với bản ghi đã xóa mềm. */
export interface InactiveCodeConflict<T> {
  /** Bản ghi cũ đã bị xóa mềm (để hiển thị cho người dùng đối chiếu). */
  existing: T
  /** Thông báo lỗi từ backend. */
  message: string
}

/**
 * Phát hiện lỗi "trùng mã với bản ghi đã xóa mềm" từ error của axios.
 * Trả về bản ghi cũ kèm thông báo nếu đúng loại lỗi, ngược lại trả null.
 */
export function parseInactiveCodeConflict<T>(
  err: unknown,
): InactiveCodeConflict<T> | null {
  if (!axios.isAxiosError(err)) return null
  const body = err.response?.data as ApiErrorBody | undefined
  if (body?.code !== INACTIVE_CODE_EXISTS) return null
  return {
    existing: body.data as T,
    message: body.error ?? 'Đã từng tồn tại bản ghi với mã này',
  }
}

/**
 * Kiểm tra error có phải do KHÔNG kết nối được backend hay không.
 *
 * Bao gồm: backend chưa chạy/đã tắt (connection refused), sai địa chỉ,
 * mất mạng, lỗi DNS, hoặc request bị timeout/hủy.
 *
 * Khi `error.response` không tồn tại => request không nhận được phản hồi HTTP
 * nào từ server, tức là không kết nối được tới backend.
 */
export function isNetworkError(err: unknown): boolean {
  return axios.isAxiosError(err) && !err.response
}

/** Thông báo chuẩn khi không kết nối được backend. */
export const BACKEND_UNREACHABLE_MESSAGE =
  'Không kết nối được tới máy chủ (backend). Vui lòng kiểm tra backend đã chạy chưa và thử lại.'

/**
 * Trích thông báo lỗi dễ đọc từ error của axios.
 * Ưu tiên message chi tiết, sau đó tới error, cuối cùng là thông báo mặc định.
 * Lỗi mạng (không kết nối được backend) trả thông báo rõ ràng riêng.
 */
export function getErrorMessage(err: unknown): string {
  if (axios.isAxiosError(err)) {
    // Không có response => không kết nối được backend.
    if (!err.response) {
      return BACKEND_UNREACHABLE_MESSAGE
    }
    const body = err.response.data as ApiErrorBody | undefined
    if (body?.message) return body.message
    if (body?.error) return body.error
    if (err.message) return err.message
  }
  if (err instanceof Error) return err.message
  return 'Đã xảy ra lỗi không xác định'
}
