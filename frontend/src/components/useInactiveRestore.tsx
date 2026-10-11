import { Modal, Typography } from 'antd'
import { ExclamationCircleFilled } from '@ant-design/icons'
import { getErrorMessage, parseInactiveCodeConflict } from '@/api/client'
import { message } from 'antd'

/** Thông tin tóm tắt bản ghi cũ để hiển thị trong hộp thoại xác nhận. */
export interface InactiveRecordSummary {
  /** Dòng mô tả chính, ví dụ "KH001 — Công ty ABC". */
  primary: string
  /** Các dòng thông tin phụ (tùy chọn). */
  details?: string[]
}

/**
 * Hàm khôi phục bản ghi theo id (gọi API restore).
 * Mỗi trang truyền vào hàm restore tương ứng của module mình.
 */
type RestoreFn = (id: number) => Promise<unknown>

/** Lấy id của bản ghi cũ từ dữ liệu backend trả về (các model đều có field id). */
function extractId(existing: unknown): number | null {
  if (existing && typeof existing === 'object' && 'id' in existing) {
    const id = (existing as { id: unknown }).id
    return typeof id === 'number' ? id : null
  }
  return null
}

/**
 * Xử lý lỗi khi lưu (tạo mới) bản ghi.
 *
 * Nếu lỗi là "trùng mã với bản ghi đã xóa mềm", hiển thị hộp thoại xác nhận:
 * - OK  -> gọi restoreFn để khôi phục bản ghi cũ, rồi thông báo cho người dùng
 *          kiểm tra lại bản ghi cũ (giữ nguyên thông tin cũ, không áp dữ liệu mới).
 * - Hủy -> không làm gì.
 *
 * Trả về `true` nếu đã xử lý (là lỗi trùng bản ghi xóa mềm), `false` nếu là
 * lỗi khác (để trang tự xử lý tiếp, ví dụ hiện message.error).
 */
export function handleInactiveCodeConflict(options: {
  err: unknown
  entityLabel: string
  /** Tạo phần tóm tắt bản ghi cũ từ dữ liệu backend trả về. */
  summarize: (existing: unknown) => InactiveRecordSummary
  /** Hàm gọi API khôi phục theo id. */
  restore: RestoreFn
  /** Callback sau khi khôi phục thành công (thường là invalidate danh sách). */
  onRestored?: () => void
}): boolean {
  const { err, entityLabel, summarize, restore, onRestored } = options

  const conflict = parseInactiveCodeConflict<unknown>(err)
  if (!conflict) return false

  const existingId = extractId(conflict.existing)
  const summary = summarize(conflict.existing)

  Modal.confirm({
    title: 'Đã từng tồn tại bản ghi này',
    icon: <ExclamationCircleFilled />,
    okText: 'Khôi phục',
    cancelText: 'Hủy',
    content: (
      <div>
        <Typography.Paragraph style={{ marginBottom: 8 }}>
          Mã bạn nhập đã từng thuộc một {entityLabel} đang ở trạng thái{' '}
          <b>ngừng hoạt động</b>:
        </Typography.Paragraph>
        <Typography.Paragraph style={{ marginBottom: 8 }}>
          <b>{summary.primary}</b>
          {summary.details?.map((line, i) => (
            <div key={i}>{line}</div>
          ))}
        </Typography.Paragraph>
        <Typography.Paragraph style={{ marginBottom: 0 }}>
          Bạn có muốn khôi phục trạng thái hoạt động không? Dữ liệu sẽ được phục
          hồi theo bản ghi cũ (không dùng thông tin vừa nhập).
        </Typography.Paragraph>
      </div>
    ),
    onOk: async () => {
      if (existingId === null) {
        message.error('Không xác định được bản ghi cũ để khôi phục')
        return
      }
      try {
        await restore(existingId)
        message.success(
          `Đã khôi phục ${entityLabel} cũ. Vui lòng kiểm tra và cập nhật lại thông tin nếu cần.`,
          // Giữ thông báo đủ lâu để người dùng đọc.
          6,
        )
        onRestored?.()
      } catch (e) {
        message.error(getErrorMessage(e))
      }
    },
  })

  return true
}
