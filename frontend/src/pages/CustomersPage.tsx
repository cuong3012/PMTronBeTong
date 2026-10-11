import { useState } from 'react'
import {
  Form,
  Input,
  Modal,
  Tag,
  message,
} from 'antd'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { ColumnsType } from 'antd/es/table'
import CrudPage from '@/components/CrudPage'
import { handleInactiveCodeConflict } from '@/components/useInactiveRestore'
import { customersApi } from '@/api/customers'
import { getErrorMessage } from '@/api/client'
import { formatActive, formatDateTime } from '@/utils/format'
import type { Customer, CustomerInput } from '@/types'

const columns: ColumnsType<Customer> = [
  { title: 'Mã KH', dataIndex: 'customer_code', key: 'customer_code' },
  { title: 'Tên khách hàng', dataIndex: 'customer_name', key: 'customer_name' },
  { title: 'Điện thoại', dataIndex: 'phone', key: 'phone', render: (v) => v || '—' },
  { title: 'Địa chỉ', dataIndex: 'address', key: 'address', render: (v) => v || '—' },
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

// Trang quản lý khách hàng: danh sách + form tạo/sửa.
export default function CustomersPage() {
  const queryClient = useQueryClient()
  const [form] = Form.useForm<CustomerInput>()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<Customer | null>(null)

  const saveMutation = useMutation({
    mutationFn: (values: CustomerInput) =>
      editing
        ? customersApi.update(editing.id, values)
        : customersApi.create(values),
    onSuccess: () => {
      message.success(editing ? 'Đã cập nhật khách hàng' : 'Đã thêm khách hàng')
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      setOpen(false)
    },
    onError: (err) => {
      // Nếu mã trùng bản ghi đã xóa mềm, hỏi người dùng có khôi phục không.
      const handled = handleInactiveCodeConflict({
        err,
        entityLabel: 'khách hàng',
        summarize: (existing) => {
          const c = existing as Customer
          return {
            primary: [c.customer_code, c.customer_name]
              .filter(Boolean)
              .join(' — '),
            details: c.phone ? [`Điện thoại: ${c.phone}`] : undefined,
          }
        },
        restore: (id) => customersApi.restore(id),
        onRestored: () => {
          queryClient.invalidateQueries({ queryKey: ['customers'] })
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

  const openEdit = (record: Customer) => {
    setEditing(record)
    setOpen(true)
  }

  return (
    <>
      <CrudPage<Customer>
        title="Quản lý khách hàng"
        queryKey="customers"
        entityLabel="khách hàng"
        list={(params) => customersApi.list(params)}
        remove={(id) => customersApi.remove(id)}
        rowKey={(r) => r.id}
        columns={columns}
        onCreate={openCreate}
        onEdit={openEdit}
      />

      <Modal
        title={editing ? 'Sửa khách hàng' : 'Thêm khách hàng'}
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={saveMutation.isPending}
        okText={editing ? 'Cập nhật' : 'Thêm mới'}
        cancelText="Hủy"
        destroyOnClose
      >
        <Form<CustomerInput>
          form={form}
          layout="vertical"
          onFinish={(values) => saveMutation.mutate(values)}
          preserve={false}
          initialValues={
            editing
              ? {
                customer_code: editing.customer_code,
                customer_name: editing.customer_name,
                address: editing.address,
                phone: editing.phone,
              }
              : {}
          }
        >
          <Form.Item
            name="customer_code"
            label="Mã khách hàng"
            rules={[
              { required: true, message: 'Vui lòng nhập mã khách hàng' },
              { max: 50, message: 'Tối đa 50 ký tự' },
            ]}
          >
            <Input placeholder="KH001" />
          </Form.Item>
          <Form.Item
            name="customer_name"
            label="Tên khách hàng"
            rules={[{ required: true, message: 'Vui lòng nhập tên khách hàng' }]}
          >
            <Input placeholder="Công ty ABC" />
          </Form.Item>
          <Form.Item name="phone" label="Số điện thoại">
            <Input placeholder="0901234567" />
          </Form.Item>
          <Form.Item name="address" label="Địa chỉ">
            <Input.TextArea rows={2} placeholder="Địa chỉ khách hàng" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
