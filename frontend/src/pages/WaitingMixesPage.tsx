import { useState } from 'react'
import {
  Alert,
  Button,
  Card,
  Form,
  InputNumber,
  Modal,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { waitingMixesApi } from '@/api/waitingMixes'
import { concreteGradesApi } from '@/api/concreteGrades'
import { customersApi } from '@/api/customers'
import { getErrorMessage } from '@/api/client'
import { formatNumber } from '@/utils/format'
import type { WaitingMix, WaitingMixInput } from '@/types'

const STATUS_LABEL: Record<number, { text: string; color: string }> = {
  0: { text: 'Chờ', color: 'default' },
  1: { text: 'Sẵn sàng', color: 'green' },
}

const columns: ColumnsType<WaitingMix> = [
  { title: 'STT', dataIndex: 'stt', key: 'stt', width: 70 },
  {
    title: 'Mác bê tông',
    dataIndex: 'concrete_grade_code',
    key: 'concrete_grade_code',
    render: (v: string, r) =>
      v ? `${v} — ${r.concrete_grade_name ?? ''}` : '— (trống)',
  },
  {
    title: 'Khách hàng',
    dataIndex: 'customer_name',
    key: 'customer_name',
    render: (v: string) => v || '— (trống)',
  },
  {
    title: 'Khối lượng (m³)',
    dataIndex: 'volume_m3',
    key: 'volume_m3',
    render: (v: number) => formatNumber(v),
  },
  {
    title: 'm3/mẻ',
    dataIndex: 'm3me',
    key: 'm3me',
    render: (v: number) => formatNumber(v),
  },
  {
    title: 'Trạng thái',
    dataIndex: 'status',
    key: 'status',
    render: (v: number) => {
      const s = STATUS_LABEL[v] ?? { text: String(v), color: 'default' }
      return <Tag color={s.color}>{s.text}</Tag>
    },
  },
]

// Trang cấp phối chờ — đúng 20 vị trí cố định, chỉ đọc/cập nhật.
export default function WaitingMixesPage() {
  const queryClient = useQueryClient()
  const [form] = Form.useForm<WaitingMixInput>()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<WaitingMix | null>(null)

  const query = useQuery({
    queryKey: ['waiting-mixes'],
    queryFn: () => waitingMixesApi.list(),
  })

  // Dữ liệu cho dropdown (chỉ bản ghi đang hoạt động).
  const grades = useQuery({
    queryKey: ['concrete-grades', { page: 1, per_page: 100 }],
    queryFn: () => concreteGradesApi.list({ page: 1, per_page: 100 }),
  })
  const customers = useQuery({
    queryKey: ['customers', { page: 1, per_page: 100 }],
    queryFn: () => customersApi.list({ page: 1, per_page: 100 }),
  })

  const gradeOptions = (grades.data?.items ?? []).map((g) => ({
    label: `${g.grade_code} — ${g.printed_name}`,
    value: g.id,
  }))
  const customerOptions = (customers.data?.items ?? []).map((c) => ({
    label: `${c.customer_code} — ${c.customer_name}`,
    value: c.id,
  }))

  const saveMutation = useMutation({
    mutationFn: (values: WaitingMixInput) =>
      waitingMixesApi.update(editing!.id, values),
    onSuccess: () => {
      message.success('Đã cập nhật vị trí cấp phối chờ')
      queryClient.invalidateQueries({ queryKey: ['waiting-mixes'] })
      setOpen(false)
    },
    onError: (err) => message.error(getErrorMessage(err)),
  })

  const openEdit = (record: WaitingMix) => {
    setEditing(record)
    form.setFieldsValue({
      concrete_grade_id: record.concrete_grade_id ?? undefined!,
      customer_id: record.customer_id ?? undefined!,
      volume_m3: record.volume_m3,
      m3me: record.m3me,
      status: record.status,
    })
    setOpen(true)
  }

  const allColumns: ColumnsType<WaitingMix> = [
    ...columns,
    {
      title: 'Hành động',
      key: 'actions',
      width: 100,
      fixed: 'right',
      render: (_, record) => (
        <Button type="link" onClick={() => openEdit(record)}>
          Sửa
        </Button>
      ),
    },
  ]

  return (
    <div>
      <Typography.Title level={3} style={{ marginTop: 0 }}>
        Cấp phối chờ
      </Typography.Title>

      <Card>
        {query.isError && (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="Không tải được cấp phối chờ"
            description={getErrorMessage(query.error)}
            action={
              <Button size="small" onClick={() => query.refetch()}>
                Thử lại
              </Button>
            }
          />
        )}

        <Table<WaitingMix>
          rowKey={(r) => r.id}
          columns={allColumns}
          dataSource={query.data ?? []}
          loading={query.isLoading}
          pagination={false}
          scroll={{ x: 'max-content' }}
        />
      </Card>

      <Modal
        title={`Cập nhật vị trí #${editing?.stt ?? ''}`}
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={saveMutation.isPending}
        okText="Cập nhật"
        cancelText="Hủy"
        destroyOnClose
      >
        <Form<WaitingMixInput>
          form={form}
          layout="vertical"
          preserve={false}
          onFinish={(values) => saveMutation.mutate(values)}
        >
          <Form.Item
            name="concrete_grade_id"
            label="Mác bê tông"
            rules={[{ required: true, message: 'Vui lòng chọn mác bê tông' }]}
          >
            <Select
              options={gradeOptions}
              loading={grades.isLoading}
              showSearch
              optionFilterProp="label"
              placeholder="Chọn mác"
            />
          </Form.Item>
          <Form.Item
            name="customer_id"
            label="Khách hàng"
            rules={[{ required: true, message: 'Vui lòng chọn khách hàng' }]}
          >
            <Select
              options={customerOptions}
              loading={customers.isLoading}
              showSearch
              optionFilterProp="label"
              placeholder="Chọn khách hàng"
            />
          </Form.Item>
          <Space size="large" style={{ width: '100%' }}>
            <Form.Item
              name="volume_m3"
              label="Khối lượng (m³)"
              rules={[
                { required: true, message: 'Nhập khối lượng' },
                {
                  validator: (_, v) =>
                    v > 0
                      ? Promise.resolve()
                      : Promise.reject(new Error('Phải lớn hơn 0')),
                },
              ]}
            >
              <InputNumber min={0} step={0.1} style={{ width: 160 }} />
            </Form.Item>
            <Form.Item
              name="m3me"
              label="m3/mẻ"
              rules={[
                {
                  validator: (_, v) =>
                    v >= 0
                      ? Promise.resolve()
                      : Promise.reject(new Error('Không được âm')),
                },
              ]}
            >
              <InputNumber min={0} step={0.1} style={{ width: 160 }} />
            </Form.Item>
            <Form.Item name="status" label="Trạng thái">
              <Select
                style={{ width: 140 }}
                options={[
                  { label: 'Chờ', value: 0 },
                  { label: 'Sẵn sàng', value: 1 },
                ]}
              />
            </Form.Item>
          </Space>
        </Form>
      </Modal>
    </div>
  )
}
