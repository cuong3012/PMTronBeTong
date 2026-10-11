import { useState } from 'react'
import { Form, Input, Modal, Tag, message } from 'antd'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { ColumnsType } from 'antd/es/table'
import CrudPage from '@/components/CrudPage'
import { handleInactiveCodeConflict } from '@/components/useInactiveRestore'
import { driversApi } from '@/api/drivers'
import { getErrorMessage } from '@/api/client'
import { formatActive, formatDateTime } from '@/utils/format'
import type { Driver, DriverInput } from '@/types'

const columns: ColumnsType<Driver> = [
  { title: 'Mã tài xế', dataIndex: 'driver_code', key: 'driver_code' },
  { title: 'Họ và tên', dataIndex: 'full_name', key: 'full_name' },
  {
    title: 'Trạng thái',
    dataIndex: 'is_active',
    key: 'is_active',
    render: (v: boolean) => (
      <Tag color={v ? 'green' : 'default'}>{formatActive(v)}</Tag>
    ),
  },
  {
    title: 'Cập nhật',
    dataIndex: 'updated_at',
    key: 'updated_at',
    render: (v: string) => formatDateTime(v),
  },
]

// Trang quản lý tài xế.
export default function DriversPage() {
  const queryClient = useQueryClient()
  const [form] = Form.useForm<DriverInput>()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<Driver | null>(null)

  const saveMutation = useMutation({
    mutationFn: (values: DriverInput) =>
      editing ? driversApi.update(editing.id, values) : driversApi.create(values),
    onSuccess: () => {
      message.success(editing ? 'Đã cập nhật tài xế' : 'Đã thêm tài xế')
      queryClient.invalidateQueries({ queryKey: ['drivers'] })
      setOpen(false)
    },
    onError: (err) => {
      // Nếu mã trùng bản ghi đã xóa mềm, hỏi người dùng có khôi phục không.
      const handled = handleInactiveCodeConflict({
        err,
        entityLabel: 'tài xế',
        summarize: (existing) => {
          const d = existing as Driver
          return {
            primary: [d.driver_code, d.full_name].filter(Boolean).join(' — '),
          }
        },
        restore: (id) => driversApi.restore(id),
        onRestored: () => {
          queryClient.invalidateQueries({ queryKey: ['drivers'] })
          setOpen(false)
        },
      })
      if (!handled) message.error(getErrorMessage(err))
    },
  })

  const openCreate = () => {
    setEditing(null)
    setOpen(true)
  }

  const openEdit = (record: Driver) => {
    setEditing(record)
    setOpen(true)
  }

  return (
    <>
      <CrudPage<Driver>
        title="Quản lý tài xế"
        queryKey="drivers"
        entityLabel="tài xế"
        list={(params) => driversApi.list(params)}
        remove={(id) => driversApi.remove(id)}
        rowKey={(r) => r.id}
        columns={columns}
        onCreate={openCreate}
        onEdit={openEdit}
      />

      <Modal
        title={editing ? 'Sửa tài xế' : 'Thêm tài xế'}
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={saveMutation.isPending}
        okText={editing ? 'Cập nhật' : 'Thêm mới'}
        cancelText="Hủy"
        destroyOnClose
      >
        <Form<DriverInput>
          form={form}
          layout="vertical"
          onFinish={(values) => saveMutation.mutate(values)}
          preserve={false}
          initialValues={
            editing
              ? {
                driver_code: editing.driver_code,
                full_name: editing.full_name,
              }
              : {}
          }
        >
          <Form.Item
            name="driver_code"
            label="Mã tài xế"
            rules={[
              { required: true, message: 'Vui lòng nhập mã tài xế' },
              { max: 50, message: 'Tối đa 50 ký tự' },
            ]}
          >
            <Input placeholder="TX001" />
          </Form.Item>
          <Form.Item
            name="full_name"
            label="Họ và tên"
            rules={[{ required: true, message: 'Vui lòng nhập họ tên' }]}
          >
            <Input placeholder="Nguyễn Văn A" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
