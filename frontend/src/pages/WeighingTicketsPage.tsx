import { useState } from 'react'
import {
  Alert,
  Button,
  Card,
  Col,
  Descriptions,
  Divider,
  Form,
  Input,
  InputNumber,
  Modal,
  Popconfirm,
  Row,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from 'antd'
import {
  MinusCircleOutlined,
  PlusOutlined,
  SearchOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { weighingTicketsApi } from '@/api/weighingTickets'
import { concreteGradesApi } from '@/api/concreteGrades'
import { getErrorMessage } from '@/api/client'
import { ADMIN_MANAGER_ROLES, useAuthStore } from '@/store/auth'
import { formatActive, formatDateTime, formatNumber } from '@/utils/format'
import type {
  ConcreteGrade,
  WeighingTicket,
  WeighingTicketDetail,
  WeighingTicketInput,
} from '@/types'

// Header cần cho mỗi mẻ. Để gọn, form tạo cho nhập các thành phần chính
// (tp1, tp2, xi1, ...) tuỳ chọn; backend mặc định 0 cho thành phần bỏ trống.
const DETAIL_FIELDS = ['tp1', 'tp2', 'tp3', 'xi1', 'xi2', 'pg1'] as const

const columns: ColumnsType<WeighingTicket> = [
  { title: 'Mã phiếu', dataIndex: 'ticket_code', key: 'ticket_code' },
  {
    title: 'Khối lượng (m³)',
    dataIndex: 'volume_m3',
    key: 'volume_m3',
    render: (v: number) => formatNumber(v),
  },
  { title: 'Số mẻ', dataIndex: 'batch_count', key: 'batch_count' },
  {
    title: 'Mác bê tông',
    dataIndex: 'concrete_grade_code',
    key: 'concrete_grade_code',
    render: (v: string, r) => `${v} — ${r.concrete_grade_name}`,
  },
  { title: 'Khách hàng', dataIndex: 'customer_name', key: 'customer_name' },
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

const PAGE_SIZE = 10

// Trang quản lý phiếu cân: danh sách + tạo mới + xem chi tiết + xóa.
export default function WeighingTicketsPage() {
  const queryClient = useQueryClient()
  const canWrite = useAuthStore((s) => s.hasRole(...ADMIN_MANAGER_ROLES))

  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [searchInput, setSearchInput] = useState('')
  const [createOpen, setCreateOpen] = useState(false)
  // Id phiếu đang xem chi tiết (null = đóng).
  const [detailId, setDetailId] = useState<number | null>(null)

  const query = useQuery({
    queryKey: ['weighing-tickets', { page, search }],
    queryFn: () =>
      weighingTicketsApi.list({ page, per_page: PAGE_SIZE, search }),
  })

  const allColumns: ColumnsType<WeighingTicket> = [
    ...columns,
    {
      title: 'Hành động',
      key: 'actions',
      width: 170,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button type="link" onClick={() => setDetailId(record.id)}>
            Chi tiết
          </Button>
          {canWrite && (
            <DeleteButton
              id={record.id}
              onDeleted={() =>
                queryClient.invalidateQueries({ queryKey: ['weighing-tickets'] })
              }
            />
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Typography.Title level={3} style={{ marginTop: 0 }}>
        Quản lý phiếu cân
      </Typography.Title>

      <Card>
        {query.isError && (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="Không tải được danh sách phiếu cân"
            description={getErrorMessage(query.error)}
            action={
              <Button size="small" onClick={() => query.refetch()}>
                Thử lại
              </Button>
            }
          />
        )}

        <Space
          style={{ marginBottom: 16, width: '100%', justifyContent: 'space-between' }}
          wrap
        >
          <Input.Search
            placeholder="Tìm theo mã phiếu / khách / mác..."
            allowClear
            enterButton={<SearchOutlined />}
            style={{ maxWidth: 360 }}
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            onSearch={(value) => {
              setSearch(value.trim())
              setPage(1)
            }}
          />
          {canWrite && (
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setCreateOpen(true)}
            >
              Thêm phiếu cân
            </Button>
          )}
        </Space>

        <Table<WeighingTicket>
          rowKey={(r) => r.id}
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

      <CreateTicketModal
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreated={() =>
          queryClient.invalidateQueries({ queryKey: ['weighing-tickets'] })
        }
      />

      <TicketDetailModal id={detailId} onClose={() => setDetailId(null)} />
    </div>
  )
}

// Modal tạo phiếu cân (header + nhiều mẻ).
function CreateTicketModal({
  open,
  onClose,
  onCreated,
}: {
  open: boolean
  onClose: () => void
  onCreated: () => void
}) {
  const [form] = Form.useForm<WeighingTicketInput>()

  const gradesQuery = useQuery({
    queryKey: ['concrete-grades', { page: 1, per_page: 100 }],
    queryFn: () => concreteGradesApi.list({ page: 1, per_page: 100 }),
  })
  const gradeOptions = (gradesQuery.data?.items ?? []).map(
    (g: ConcreteGrade) => ({
      label: `${g.grade_code} — ${g.printed_name}`,
      value: g.id,
    }),
  )

  const createMutation = useMutation({
    mutationFn: (values: WeighingTicketInput) =>
      weighingTicketsApi.create(values),
    onSuccess: () => {
      message.success('Đã tạo phiếu cân')
      onCreated()
      form.resetFields()
      onClose()
    },
    onError: (err) => message.error(getErrorMessage(err)),
  })

  // Đặt giá trị mặc định + 1 mẻ trống khi mở modal.
  const handleAfterOpen = (isOpen: boolean) => {
    if (isOpen) {
      form.setFieldsValue({
        volume_m3: 0,
        batch_count: 1,
        m3me: 0,
        standard_per_cube: 0,
        standard_per_batch: 0,
        details: [
          { batch_no: 1, batch_volume_m3: 0, mixed_at: '', ...zeroComponents() },
        ],
      })
    }
  }

  // Cập nhật bản sao mác khi người dùng chọn mác.
  const onGradeChange = (gradeId: number) => {
    const g = (gradesQuery.data?.items ?? []).find((x) => x.id === gradeId)
    if (!g) return
    form.setFieldsValue({
      concrete_grade_code: g.grade_code,
      concrete_grade_name: g.printed_name,
      printed_concrete_grade_name: g.printed_name,
    })
  }

  return (
    <Modal
      title="Thêm phiếu cân"
      open={open}
      onCancel={onClose}
      onOk={() => form.submit()}
      confirmLoading={createMutation.isPending}
      okText="Tạo phiếu"
      cancelText="Hủy"
      width={900}
      destroyOnClose
      afterOpenChange={handleAfterOpen}
    >
      <Form<WeighingTicketInput>
        form={form}
        layout="vertical"
        preserve={false}
        onFinish={(values) =>
          createMutation.mutate({
            ...values,
            details: (values.details ?? []).map((d) => ({
              ...d,
              mixed_at: (d.mixed_at || '').trim(),
            })),
          })
        }
      >
        <Row gutter={16}>
          <Col xs={24} md={8}>
            <Form.Item
              name="ticket_code"
              label="Mã phiếu"
              rules={[{ required: true, message: 'Vui lòng nhập mã phiếu' }]}
            >
              <Input placeholder="PC001" />
            </Form.Item>
          </Col>
          <Col xs={24} md={8}>
            <Form.Item
              name="order_id"
              label="Mã đơn hàng (id)"
              rules={[{ required: true, message: 'Nhập id đơn hàng' }]}
            >
              <InputNumber min={1} style={{ width: '100%' }} />
            </Form.Item>
          </Col>
          <Col xs={24} md={8}>
            <Form.Item
              name="concrete_grade_id"
              label="Mác bê tông"
              rules={[{ required: true, message: 'Chọn mác bê tông' }]}
            >
              <Select
                options={gradeOptions}
                loading={gradesQuery.isLoading}
                showSearch
                optionFilterProp="label"
                onChange={onGradeChange}
              />
            </Form.Item>
          </Col>
        </Row>

        {/* Các trường bản sao mác (ẩn, gửi kèm khi chọn mác). */}
        <Form.Item name="concrete_grade_code" hidden>
          <Input />
        </Form.Item>
        <Form.Item name="concrete_grade_name" hidden>
          <Input />
        </Form.Item>
        <Form.Item name="printed_concrete_grade_name" hidden>
          <Input />
        </Form.Item>

        <Row gutter={16}>
          <Col xs={12} md={6}>
            <Form.Item
              name="volume_m3"
              label="Khối lượng (m³)"
              rules={[{ required: true, message: 'Nhập khối lượng' }]}
            >
              <InputNumber min={0} step={0.1} style={{ width: '100%' }} />
            </Form.Item>
          </Col>
          <Col xs={12} md={6}>
            <Form.Item
              name="batch_count"
              label="Số mẻ"
              rules={[{ required: true, message: 'Nhập số mẻ' }]}
            >
              <InputNumber min={1} style={{ width: '100%' }} />
            </Form.Item>
          </Col>
          <Col xs={12} md={6}>
            <Form.Item name="m3me" label="m3/mẻ">
              <InputNumber min={0} step={0.1} style={{ width: '100%' }} />
            </Form.Item>
          </Col>
          <Col xs={12} md={6}>
            <Form.Item
              name="customer_name"
              label="Khách hàng"
              rules={[{ required: true, message: 'Nhập tên khách hàng' }]}
            >
              <Input placeholder="Công ty ABC" />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} md={12}>
            <Form.Item
              name="construction_site_address"
              label="Địa chỉ công trường"
            >
              <Input placeholder="Công trường..." />
            </Form.Item>
          </Col>
          <Col xs={24} md={12}>
            <Form.Item name="material_name" label="Tên vật liệu">
              <Input placeholder="Cát vàng" />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={12} md={6}>
            <Form.Item name="standard_per_cube" label="Định mức/m³">
              <InputNumber min={0} step={1} style={{ width: '100%' }} />
            </Form.Item>
          </Col>
          <Col xs={12} md={6}>
            <Form.Item name="standard_per_batch" label="Định mức/mẻ">
              <InputNumber min={0} step={1} style={{ width: '100%' }} />
            </Form.Item>
          </Col>
          <Col xs={24} md={12}>
            <Form.Item name="humidity_values" label="Giá trị độ ẩm">
              <Input placeholder="12.5,11.8" />
            </Form.Item>
          </Col>
        </Row>

        <Divider orientation="left" plain>
          Chi tiết mẻ
        </Divider>

        <Form.List name="details">
          {(fields, { add, remove }) => (
            <>
              {fields.map((field) => (
                <Card
                  key={field.key}
                  size="small"
                  style={{ marginBottom: 12 }}
                  extra={
                    fields.length > 1 && (
                      <Button
                        type="text"
                        danger
                        icon={<MinusCircleOutlined />}
                        onClick={() => remove(field.name)}
                      />
                    )
                  }
                >
                  <Row gutter={12}>
                    <Col xs={12} md={4}>
                      <Form.Item
                        name={[field.name, 'batch_no']}
                        label="Số mẻ"
                        rules={[{ required: true, message: 'Nhập số mẻ' }]}
                      >
                        <InputNumber min={1} style={{ width: '100%' }} />
                      </Form.Item>
                    </Col>
                    <Col xs={12} md={4}>
                      <Form.Item
                        name={[field.name, 'batch_volume_m3']}
                        label="KL mẻ (m³)"
                        rules={[{ required: true, message: 'Nhập KL' }]}
                      >
                        <InputNumber min={0} step={0.1} style={{ width: '100%' }} />
                      </Form.Item>
                    </Col>
                    <Col xs={24} md={8}>
                      <Form.Item
                        name={[field.name, 'mixed_at']}
                        label="Thời điểm trộn"
                        rules={[
                          { required: true, message: 'Nhập thời điểm trộn' },
                        ]}
                        extra="Ví dụ: 2026-10-11 08:00:00"
                      >
                        <Input placeholder="YYYY-MM-DD HH:MM:SS" />
                      </Form.Item>
                    </Col>
                    {DETAIL_FIELDS.map((f) => (
                      <Col xs={8} md={2} key={f}>
                        <Form.Item name={[field.name, f]} label={f.toUpperCase()}>
                          <InputNumber
                            min={0}
                            step={0.1}
                            style={{ width: '100%' }}
                          />
                        </Form.Item>
                      </Col>
                    ))}
                  </Row>
                </Card>
              ))}
              <Button
                type="dashed"
                onClick={() =>
                  add({
                    batch_no: fields.length + 1,
                    batch_volume_m3: 0,
                    mixed_at: '',
                    ...zeroComponents(),
                  })
                }
                block
                icon={<PlusOutlined />}
              >
                Thêm mẻ
              </Button>
            </>
          )}
        </Form.List>
      </Form>
    </Modal>
  )
}

/** Trả về object 0 cho 16 thành phần chi tiết. */
function zeroComponents() {
  return {
    tp1: 0,
    tp2: 0,
    tp3: 0,
    tp4: 0,
    tp5: 0,
    tp6: 0,
    xi1: 0,
    xi2: 0,
    xi3: 0,
    xi4: 0,
    xi5: 0,
    pg1: 0,
    pg2: 0,
    pg3: 0,
    pg4: 0,
    pg5: 0,
    pg6: 0,
  }
}

// Nút xóa phiếu (chỉ hiển thị cho admin/manager).
function DeleteButton({
  id,
  onDeleted,
}: {
  id: number
  onDeleted: () => void
}) {
  const mutation = useMutation({
    mutationFn: () => weighingTicketsApi.remove(id),
    onSuccess: () => {
      message.success('Đã xóa phiếu cân')
      onDeleted()
    },
    onError: (err) => message.error(getErrorMessage(err)),
  })
  return (
    <Popconfirm
      title="Xóa phiếu cân?"
      description="Phiếu và các mẻ của nó sẽ chuyển sang trạng thái ngừng hoạt động."
      okText="Xóa"
      cancelText="Hủy"
      okButtonProps={{ danger: true }}
      onConfirm={() => mutation.mutate()}
    >
      <Button type="link" danger>
        Xóa
      </Button>
    </Popconfirm>
  )
}

// Modal hiển thị chi tiết một phiếu (gọi GET /:id để lấy details).
function TicketDetailModal({
  id,
  onClose,
}: {
  id: number | null
  onClose: () => void
}) {
  const query = useQuery({
    queryKey: ['weighing-ticket', id],
    queryFn: () => weighingTicketsApi.get(id as number),
    enabled: id !== null,
  })

  const t = query.data

  return (
    <Modal
      title={t ? `Phiếu cân ${t.ticket_code}` : 'Chi tiết phiếu cân'}
      open={id !== null}
      onCancel={onClose}
      footer={<Button onClick={onClose}>Đóng</Button>}
      width={900}
      destroyOnClose
    >
      {query.isLoading && <Typography.Text>Đang tải...</Typography.Text>}
      {query.isError && (
        <Alert
          type="error"
          showIcon
          message="Không tải được chi tiết phiếu"
          description={getErrorMessage(query.error)}
        />
      )}
      {t && (
        <>
          <Descriptions bordered size="small" column={2}>
            <Descriptions.Item label="Mã phiếu">
              {t.ticket_code}
            </Descriptions.Item>
            <Descriptions.Item label="Đơn hàng">
              {t.order_number || t.order_id}
            </Descriptions.Item>
            <Descriptions.Item label="Khối lượng">
              {formatNumber(t.volume_m3)} m³
            </Descriptions.Item>
            <Descriptions.Item label="Số mẻ">
              {t.batch_count}
            </Descriptions.Item>
            <Descriptions.Item label="Mác bê tông">
              {t.concrete_grade_code} — {t.concrete_grade_name}
            </Descriptions.Item>
            <Descriptions.Item label="Tên in">
              {t.printed_concrete_grade_name}
            </Descriptions.Item>
            <Descriptions.Item label="Khách hàng">
              {t.customer_name}
            </Descriptions.Item>
            <Descriptions.Item label="Công trường">
              {t.construction_site_address || '—'}
            </Descriptions.Item>
            <Descriptions.Item label="Vật liệu">
              {t.material_name || '—'}
            </Descriptions.Item>
            <Descriptions.Item label="Độ ẩm">
              {t.humidity_values || '—'}
            </Descriptions.Item>
            <Descriptions.Item label="Định mức/m³">
              {formatNumber(t.standard_per_cube)}
            </Descriptions.Item>
            <Descriptions.Item label="Định mức/mẻ">
              {formatNumber(t.standard_per_batch)}
            </Descriptions.Item>
          </Descriptions>

          <Divider orientation="left" plain>
            Các mẻ
          </Divider>

          <Table<WeighingTicketDetail>
            rowKey={(r) => r.id}
            size="small"
            pagination={false}
            dataSource={t.details ?? []}
            scroll={{ x: 'max-content' }}
            columns={[
              { title: 'Mẻ', dataIndex: 'batch_no', key: 'batch_no' },
              {
                title: 'KL (m³)',
                dataIndex: 'batch_volume_m3',
                key: 'batch_volume_m3',
                render: (v: number) => formatNumber(v),
              },
              {
                title: 'Thời điểm',
                dataIndex: 'mixed_at',
                key: 'mixed_at',
                render: (v: string) => formatDateTime(v),
              },
              ...DETAIL_FIELDS.map((f) => ({
                title: f.toUpperCase(),
                dataIndex: f,
                key: f,
                render: (v: number) => formatNumber(v, 3),
              })),
            ]}
          />
        </>
      )}
    </Modal>
  )
}
