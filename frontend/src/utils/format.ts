import dayjs from 'dayjs'

/** Hiển thị thời điểm ISO 8601 theo giờ địa phương, rỗng nếu không có. */
export function formatDateTime(value?: string): string {
  if (!value) return '—'
  const d = dayjs(value)
  return d.isValid() ? d.format('DD/MM/YYYY HH:mm') : value
}

/** Nhãn trạng thái hoạt động dựa trên cờ is_active. */
export function formatActive(isActive: boolean): string {
  return isActive ? 'Đang hoạt động' : 'Ngừng hoạt động'
}

/**
 * Hiển thị số với dấu chấm làm dấu thập phân (ví dụ 1.3), KHÔNG dùng dấu phẩy
 * kiểu Việt Nam. Dùng locale 'en-US' và bỏ các số 0 thập phân dư.
 * Trả '—' khi giá trị không hợp lệ.
 */
export function formatNumber(value?: number | null, maxDecimals = 2): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '—'
  return new Intl.NumberFormat('en-US', {
    maximumFractionDigits: maxDecimals,
  }).format(value)
}
