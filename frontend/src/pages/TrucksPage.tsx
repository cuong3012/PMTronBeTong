import { useState } from 'react'
import { Form, Input, Modal, Select, Tag, message } from 'antd'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { ColumnsType } from 'antd/es/table'
import CrudPage from '@/components/CrudPage'
import { handleInactiveCodeConflict } from '@/components/useInactiveRestore'
import { trucksApi } from '@/api/trucks'
import { driversApi } from '@/api/drivers'
import { getErrorMessage } from '@/api/client'
import { formatActive, formatDateTime } from '@/utils/format'
import type { Driver, MixerTruck, MixerTruckInput } from '@/types'

const columns: ColumnsType<MixerTruck> = [
  { title: 'Mã xe', dataIndex: 'truck_code', key: 'truck_code' },
  { title: 'Biển số', dataIndex: 'plate_number', key: 'plate_number' },
  {
    title: 'Tài xế mặc định',
    dataIndex: 'default_driver_name',
    key: 'default_driver_name',
    render: (v?: string) => v || '—',
  },
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

// Trang quản lý xe bồn. Form cho chọn tài xế mặc định từ danh sách tài xế.
export default function TrucksPage() {
  const queryClient = useQueryClient()
  const [form] = Form.useForm<MixerTruckInput>()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<MixerTruck | null>(null)

  // Danh sách tài xế cho dropdown (lấy tối đa 100 bản ghi đang hoạt động).
  const drivers = useQuery({
    queryKey: ['drivers', { page: 1, per_page: 100 }],
    queryFn: () => driversApi.list({ page: 1, per_page: 100 }),
  })
  const driverOptions = (drivers.data?.items ?? []).map((d: Driver) => ({
    label: `${d.driver_code} — ${d.full_name}`,
    value: d.id,
  }))

  const saveMutation = useMutation({
    mutationFn: (values: MixerTruckInput) =>
      editing ? trucksApi.update(editing.id, values) : trucksApi.create(values),
    onSuccess: () => {
      message.success(editing ? 'Đã cập nhật xe bồn' : 'Đã thêm xe bồn')
      queryClient.invalidateQueries({ queryKey: ['mixer-trucks'] })
      setOpen(false)
    },
    onError: (err) => {
      // Nếu mã xe/biển số trùng bản ghi đã xóa mềm, hỏi người dùng có khôi phục không.
      const handled = handleInactiveCodeConflict({
        err,
        entityLabel: 'xe bồn',
        summarize: (existing) => {
          const t = existing as MixerTruck
          return {
            primary: [t.truck_code, t.plate_number].filter(Boolean).join(' — '),
            details: t.default_driver_name
              ? [`Tài xế mặc định: ${t.default_driver_name}`]
              : undefined,
          }
        },
        restore: (id) => trucksApi.restore(id),
        onRestored: () => {
          queryClient.invalidateQueries({ queryKey: ['mixer-trucks'] })
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

  const openEdit = (record: MixerTruck) => {
    setEditing(record)
    setOpen(true)
  }

  return (
    <>
      <CrudPage<MixerTruck>
        title="Quản lý xe bồn"
        queryKey="mixer-trucks"
        entityLabel="xe bồn"
        list={(params) => trucksApi.list(params)}
        remove={(id) => trucksApi.remove(id)}
        rowKey={(r) => r.id}
        columns={columns}
        onCreate={openCreate}
        onEdit={openEdit}
        invalidateKeys={[['drivers']]}
      />

      <Modal
        title={editing ? 'Sửa xe bồn' : 'Thêm xe bồn'}
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={saveMutation.isPending}
        okText={editing ? 'Cập nhật' : 'Thêm mới'}
        cancelText="Hủy"
        destroyOnClose
      >
        <Form<MixerTruckInput>
          form={form}
          layout="vertical"
          preserve={false}
          initialValues={
            editing
              ? {
                truck_code: editing.truck_code,
                plate_number: editing.plate_number,
                default_driver_id: editing.default_driver_id ?? null,
              }
              : {}
          }
          onFinish={(values) =>
            saveMutation.mutate({
              truck_code: values.truck_code,
              plate_number: values.plate_number,
              // Chuẩn hóa null khi không chọn tài xế.
              default_driver_id: values.default_driver_id ?? null,
            })
          }
        >
          <Form.Item
            name="truck_code"
            label="Mã xe"
            rules={[
              { required: true, message: 'Vui lòng nhập mã xe' },
              { max: 50, message: 'Tối đa 50 ký tự' },
            ]}
          >
            <Input placeholder="XE001" />
          </Form.Item>
          <Form.Item
            name="plate_number"
            label="Biển số xe"
            rules={[
              { required: true, message: 'Vui lòng nhập biển số' },
              { max: 30, message: 'Tối đa 30 ký tự' },
            ]}
          >
            <Input placeholder="51C-12345" />
          </Form.Item>
          <Form.Item name="default_driver_id" label="Tài xế mặc định">
            <Select
              allowClear
              placeholder="Chọn tài xế (tùy chọn)"
              options={driverOptions}
              loading={drivers.isLoading}
              showSearch
              optionFilterProp="label"
            />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
