import axios from 'axios'

/**
 * Kiểm tra backend có đang chạy và phản hồi được hay không.
 *
 * Backend Go định nghĩa endpoint công khai GET /health ở gốc (không nằm dưới
 * tiền tố /api), nên dùng axios trực tiếp thay vì `apiClient` để tránh baseURL.
 *
 * Trả về `true` nếu nhận được phản hồi HTTP thành công, `false` nếu không kết
 * nối được backend hoặc backend phản hồi lỗi.
 *
 * Không bao giờ ném lỗi ra ngoài để phù hợp cho việc thăm dò định kỳ.
 */
export async function checkBackendHealth(timeoutMs = 3000): Promise<boolean> {
  try {
    const res = await axios.get('/health', { timeout: timeoutMs })
    return res.status >= 200 && res.status < 300
  } catch {
    return false
  }
}
