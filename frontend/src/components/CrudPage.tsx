import { useState } from 'react'
import {
  Alert,
  Button,
  Card,
  Input,
  Popconfirm,
  Space,
  Table,
  Typography,
  message,
} from 'antd'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getErrorMessage } from '@/api/client'
import { ADMIN_MANAGER_ROLES, useAuthStore } from '@/store/auth'
import type { PageResult } from '@/types'

// Cấu hình chung cho một trang CRUD danh mục.
//
// Tham số generic T là kiểu bản ghi; TInput là kiểu dữ liệu form gửi lên API.
export interface CrudPageProps<T> {
  /** Tiêu đề trang. */
  title: string
  /** Khóa cache của React Query (ví dụ 'customers'). */
  queryKey: string
  /** Hàm gọi API danh sách có phân trang. */
  list: (params: {
    page: number
    per_page: number
    search?: string
  }) => Promise<PageResult<T>>
  /** Hàm xóa mềm theo id. */
  remove: (id: number) => Promise<void>
  /** Trả về khóa duy nhất của mỗi bản ghi cho Table. */
  rowKey: (record: T) => number
  /** Định nghĩa cột của bảng (không gồm cột hành động). */
  columns: ColumnsType<T>
  /** Mở form tạo mới. */
  onCreate: () => void
  /** Mở form chỉnh sửa cho bản ghi. */
  onEdit: (record: T) => void
  /** Nhãn thực thể dùng trong thông báo (ví dụ 'khách hàng'). */
  entityLabel: string
  /** Có hiển thị nút tạo/xóa hay không (mặc định true). */
  canCreate?: boolean
  /** Cho phép xóa mềm (chỉ admin/manager). */
  canDelete?: boolean
  /** Refetch lại các query khác sau khi thay đổi dữ liệu (tùy chọn). */
  invalidateKeys?: string[][]
}

const PAGE_SIZE = 10

/**
 * Trang CRUD tổng quát: bảng phân trang + tìm kiếm + nút thêm/sửa/xóa.
 * Việc render form được ủy quyền cho component cha qua onCreate/onEdit.
 */
export default function CrudPage<T>({
  title,
  queryKey,
  list,
  remove,
  rowKey,
  columns,
  onCreate,
  onEdit,
  entityLabel,
  canCreate = true,
  canDelete = true,
  invalidateKeys = [],
}: CrudPageProps<T>) {
  const queryClient = useQueryClient()
  const hasRole = useAuthStore((s) => s.hasRole)
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [searchInput, setSearchInput] = useState('')

  // Cho phép xóa chỉ khi được bật và người dùng là admin/manager.
  const deleteAllowed = canDelete && hasRole(...ADMIN_MANAGER_ROLES)

  const query = useQuery({
    queryKey: [queryKey, { page, search }],
    queryFn: () => list({ page, per_page: PAGE_SIZE, search }),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => remove(id),
    onSuccess: () => {
      message.success(`Đã xóa ${entityLabel}`)
      // Làm mới danh sách và các truy vấn liên quan.
      queryClient.invalidateQueries({ queryKey: [queryKey] })
      invalidateKeys.forEach((key) =>
        queryClient.invalidateQueries({ queryKey: key }),
      )
    },
    onError: (err) => message.error(getErrorMessage(err)),
  })

  // Ghép cột hành động vào cuối bảng.
  const allColumns: ColumnsType<T> = [
    ...columns,
    {
      title: 'Hành động',
      key: 'actions',
      width: 180,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button type="link" onClick={() => onEdit(record)}>
            Sửa
          </Button>
          {deleteAllowed && (
            <Popconfirm
              title={`Xóa ${entityLabel}?`}
              description="Bản ghi sẽ được chuyển sang trạng thái ngừng hoạt động."
              okText="Xóa"
              cancelText="Hủy"
              okButtonProps={{ danger: true }}
              onConfirm={() => deleteMutation.mutate(rowKey(record))}
            >
              <Button type="link" danger>
                Xóa
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Typography.Title level={3} style={{ marginTop: 0 }}>
        {title}
      </Typography.Title>

      <Card>
        {query.isError && (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="Không tải được dữ liệu"
            description={getErrorMessage(query.error)}
            action={
              <Button size="small" onClick={() => query.refetch()}>
                Thử lại
              </Button>
            }
          />
        )}

        <Space
          style={{
            marginBottom: 16,
            width: '100%',
            justifyContent: 'space-between',
          }}
          wrap
        >
          <Input.Search
            placeholder={`Tìm ${entityLabel}...`}
            allowClear
            enterButton={<SearchOutlined />}
            style={{ maxWidth: 320 }}
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            onSearch={(value) => {
              setSearch(value.trim())
              setPage(1)
            }}
          />
          {canCreate && (
            <Button type="primary" icon={<PlusOutlined />} onClick={onCreate}>
              Thêm {entityLabel}
            </Button>
          )}
        </Space>

        <Table<T>
          rowKey={rowKey}
          columns={allColumns}
          dataSource={query.data?.items ?? []}
          loading={query.isLoading}
          scroll={{ x: 'max-content' }}
          pagination={{
            current: page,
            pageSize: PAGE_SIZE,
            total: query.data?.total ?? 0,
            showSizeChanger: false,
            onChange: setPage,
            showTotal: (total) => `Tổng ${total} bản ghi`,
          }}
        />
      </Card>
    </div>
  )
}
