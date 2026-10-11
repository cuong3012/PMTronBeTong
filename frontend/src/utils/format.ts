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
