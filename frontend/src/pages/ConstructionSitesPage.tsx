import { useState } from 'react'
import { Form, Input, Modal, Tag, message } from 'antd'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { ColumnsType } from 'antd/es/table'
import CrudPage from '@/components/CrudPage'
import { handleInactiveCodeConflict } from '@/components/useInactiveRestore'
import { sitesApi } from '@/api/sites'
import { getErrorMessage } from '@/api/client'
import { formatActive, formatDateTime } from '@/utils/format'
import type { ConstructionSite, ConstructionSiteInput } from '@/types'

const columns: ColumnsType<ConstructionSite> = [
  { title: 'Mã công trường', dataIndex: 'site_code', key: 'site_code' },
  { title: 'Tên công trường', dataIndex: 'site_name', key: 'site_name' },
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

// Trang quản lý công trường.
export default function ConstructionSitesPage() {
  const queryClient = useQueryClient()
  const [form] = Form.useForm<ConstructionSiteInput>()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<ConstructionSite | null>(null)

  const saveMutation = useMutation({
    mutationFn: (values: ConstructionSiteInput) =>
      editing ? sitesApi.update(editing.id, values) : sitesApi.create(values),
    onSuccess: () => {
      message.success(editing ? 'Đã cập nhật công trường' : 'Đã thêm công trường')
      queryClient.invalidateQueries({ queryKey: ['construction-sites'] })
      setOpen(false)
    },
    onError: (err) => {
      // Nếu mã trùng bản ghi đã xóa mềm, hỏi người dùng có khôi phục không.
      const handled = handleInactiveCodeConflict({
        err,
        entityLabel: 'công trường',
        summarize: (existing) => {
          const s = existing as ConstructionSite
          return {
            primary: [s.site_code, s.site_name].filter(Boolean).join(' — '),
          }
        },
        restore: (id) => sitesApi.restore(id),
        onRestored: () => {
          queryClient.invalidateQueries({ queryKey: ['construction-sites'] })
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

  const openEdit = (record: ConstructionSite) => {
    setEditing(record)
    setOpen(true)
  }

  return (
    <>
      <CrudPage<ConstructionSite>
        title="Quản lý công trường"
        queryKey="construction-sites"
        entityLabel="công trường"
        list={(params) => sitesApi.list(params)}
        remove={(id) => sitesApi.remove(id)}
        rowKey={(r) => r.id}
        columns={columns}
        onCreate={openCreate}
        onEdit={openEdit}
      />

      <Modal
        title={editing ? 'Sửa công trường' : 'Thêm công trường'}
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={saveMutation.isPending}
        okText={editing ? 'Cập nhật' : 'Thêm mới'}
        cancelText="Hủy"
        destroyOnClose
      >
        <Form<ConstructionSiteInput>
          form={form}
          layout="vertical"
          onFinish={(values) => saveMutation.mutate(values)}
          preserve={false}
          initialValues={
            editing
              ? {
                site_code: editing.site_code,
                site_name: editing.site_name,
              }
              : {}
          }
        >
          <Form.Item
            name="site_code"
            label="Mã công trường"
            rules={[
              { required: true, message: 'Vui lòng nhập mã công trường' },
              { max: 50, message: 'Tối đa 50 ký tự' },
            ]}
          >
            <Input placeholder="CT001" />
          </Form.Item>
          <Form.Item
            name="site_name"
            label="Tên công trường"
            rules={[{ required: true, message: 'Vui lòng nhập tên công trường' }]}
          >
            <Input placeholder="Dự án ABC" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
