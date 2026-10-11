import { useState } from 'react'
import {
  Alert,
  Button,
  Card,
  Input,
  Table,
  Typography,
  message,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { componentCatalogApi } from '@/api/componentCatalog'
import { getErrorMessage } from '@/api/client'
import { ADMIN_MANAGER_ROLES, useAuthStore } from '@/store/auth'
import type { ComponentCatalog } from '@/types'

// Trang danh mục thành phần: CHỈ dùng để cập nhật giá trị độ ẩm của các vị trí
// đã được seed sẵn. Không thêm, không xóa, không sửa các trường khác.
export default function ComponentCatalogPage() {
  const hasRole = useAuthStore((s) => s.hasRole)
  const canEdit = hasRole(...ADMIN_MANAGER_ROLES)

  const query = useQuery({
    queryKey: ['component-catalog'],
    queryFn: () => componentCatalogApi.list(),
  })

  const columns: ColumnsType<ComponentCatalog> = [
    { title: 'Vị trí', dataIndex: 'slot_name', key: 'slot_name', width: 120 },
    {
      title: 'Đơn vị',
      dataIndex: 'unit',
      key: 'unit',
      width: 100,
      render: (v) => v || '—',
    },
    {
      title: 'Giá trị độ ẩm (6 giá trị)',
      dataIndex: 'humidity_values',
      key: 'humidity_values',
      render: (_, record) => (
        <HumidityEditor record={record} canEdit={canEdit} />
      ),
    },
  ]

  return (
    <div>
      <Typography.Title level={3} style={{ marginTop: 0 }}>
        Danh mục thành phần
      </Typography.Title>
      <Typography.Paragraph type="secondary">
        Cập nhật giá trị độ ẩm cho từng vị trí thành phần. Nhập các giá trị ngăn
        cách bởi dấu phẩy (ví dụ: 12.5,11.8,10.7,13.2,9.8,11.0).
      </Typography.Paragraph>

      <Card>
        {query.isError && (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="Không tải được danh mục thành phần"
            description={getErrorMessage(query.error)}
            action={
              <Button size="small" onClick={() => query.refetch()}>
                Thử lại
              </Button>
            }
          />
        )}

        <Table<ComponentCatalog>
          rowKey={(r) => r.id}
          columns={columns}
          dataSource={query.data ?? []}
          loading={query.isLoading}
          pagination={false}
          scroll={{ x: 'max-content' }}
        />
      </Card>
    </div>
  )
}

// Ô nhập độ ẩm cho 1 dòng: sửa field humidity_values và lưu riêng dòng đó.
//
// Khi lưu, gửi lại nguyên slot_name/unit/notes cũ (PUT là cập nhật toàn bộ bản
// ghi và slot_name bắt buộc), chỉ thay humidity_values.
function HumidityEditor({
  record,
  canEdit,
}: {
  record: ComponentCatalog
  canEdit: boolean
}) {
  const queryClient = useQueryClient()
  const [value, setValue] = useState(record.humidity_values ?? '')

  const saveMutation = useMutation({
    mutationFn: (humidityValues: string) =>
      componentCatalogApi.update(record.id, {
        slot_name: record.slot_name,
        humidity_values: humidityValues,
        unit: record.unit ?? '',
        notes: record.notes ?? '',
      }),
    onSuccess: () => {
      message.success(`Đã cập nhật độ ẩm cho ${record.slot_name}`)
      queryClient.invalidateQueries({ queryKey: ['component-catalog'] })
    },
    onError: (err) => message.error(getErrorMessage(err)),
  })

  if (!canEdit) {
    return <Typography.Text>{record.humidity_values || '—'}</Typography.Text>
  }

  return (
    <Input.Search
      value={value}
      placeholder="12.5,11.8,10.7,13.2,9.8,11.0"
      enterButton="Lưu"
      allowClear
      onChange={(e) => setValue(e.target.value)}
      onSearch={(v) => saveMutation.mutate(v.trim())}
      loading={saveMutation.isPending}
      style={{ maxWidth: 420 }}
    />
  )
}
