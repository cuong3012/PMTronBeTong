import { useState } from 'react'
import {
  Col,
  Divider,
  Form,
  Input,
  InputNumber,
  Modal,
  Row,
  Tag,
  message,
} from 'antd'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { ColumnsType } from 'antd/es/table'
import CrudPage from '@/components/CrudPage'
import { concreteGradesApi } from '@/api/concreteGrades'
import { getErrorMessage } from '@/api/client'
import { formatActive, formatDateTime } from '@/utils/format'
import type { ConcreteGrade, GradeInput } from '@/types'

// Nhóm 16 cột thành phần (theo thứ tự trong DB) để render form gọn.
const COMPONENT_GROUPS: { title: string; fields: { name: keyof GradeInput; label: string }[] }[] = [
  {
    title: 'Cốt liệu & nước',
    fields: [
      { name: 'tp1', label: 'TP1' },
      { name: 'tp2', label: 'TP2' },
      { name: 'tp3', label: 'TP3' },
      { name: 'tp4', label: 'TP4' },
      { name: 'tp5', label: 'TP5' },
      { name: 'tp6', label: 'TP6' },
      { name: 'nuoc', label: 'Nước' },
    ],
  },
  {
    title: 'Xi măng',
    fields: [
      { name: 'xi1', label: 'XI1' },
      { name: 'xi2', label: 'XI2' },
      { name: 'xi3', label: 'XI3' },
      { name: 'xi4', label: 'XI4' },
      { name: 'xi5', label: 'XI5' },
    ],
  },
  {
    title: 'Phụ gia (tối đa 4 chữ số thập phân)',
    fields: [
      { name: 'pg1', label: 'PG1' },
      { name: 'pg2', label: 'PG2' },
      { name: 'pg3', label: 'PG3' },
      { name: 'pg4', label: 'PG4' },
      { name: 'pg5', label: 'PG5' },
      { name: 'pg6', label: 'PG6' },
    ],
  },
]

const columns: ColumnsType<ConcreteGrade> = [
  { title: 'Mã mác', dataIndex: 'grade_code', key: 'grade_code' },
  { title: 'Tên in', dataIndex: 'printed_name', key: 'printed_name' },
  {
    title: 'Độ sụt',
    dataIndex: 'do_sut',
    key: 'do_sut',
    render: (v: string) => v || '—',
  },
  {
    title: 'Cường độ',
    dataIndex: 'cuong_do',
    key: 'cuong_do',
    render: (v: string) => v || '—',
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

// Trang quản lý mác bê tông (định mức cho 1 m³).
export default function ConcreteGradesPage() {
  const queryClient = useQueryClient()
  const [form] = Form.useForm<GradeInput>()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<ConcreteGrade | null>(null)

  const saveMutation = useMutation({
    mutationFn: (values: GradeInput) =>
      editing
        ? concreteGradesApi.update(editing.id, values)
        : concreteGradesApi.create(values),
    onSuccess: () => {
      message.success(editing ? 'Đã cập nhật mác bê tông' : 'Đã thêm mác bê tông')
      queryClient.invalidateQueries({ queryKey: ['concrete-grades'] })
      setOpen(false)
    },
    onError: (err) => message.error(getErrorMessage(err)),
  })

  const openCreate = () => {
    setEditing(null)
    setOpen(true)
  }

  const openEdit = (record: ConcreteGrade) => {
    setEditing(record)
    setOpen(true)
  }

  return (
    <>
      <CrudPage<ConcreteGrade>
        title="Quản lý mác bê tông"
        queryKey="concrete-grades"
        entityLabel="mác bê tông"
        list={(params) => concreteGradesApi.list(params)}
        remove={(id) => concreteGradesApi.remove(id)}
        rowKey={(r) => r.id}
        columns={columns}
        onCreate={openCreate}
        onEdit={openEdit}
      />

      <Modal
        title={editing ? 'Sửa mác bê tông' : 'Thêm mác bê tông'}
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={saveMutation.isPending}
        okText={editing ? 'Cập nhật' : 'Thêm mới'}
        cancelText="Hủy"
        width={760}
        destroyOnClose
      >
        <Form<GradeInput>
          form={form}
          layout="vertical"
          preserve={false}
          // Đảm bảo mọi trường thành phần mặc định là 0 khi tạo mới.
          initialValues={
            editing
              ? { ...defaultComponents(), ...editing }
              : defaultComponents()
          }
          onFinish={(values) => saveMutation.mutate(values)}
        >
          <Row gutter={16}>
            <Col xs={24} md={12}>
              <Form.Item
                name="grade_code"
                label="Mã mác"
                rules={[
                  { required: true, message: 'Vui lòng nhập mã mác' },
                  { max: 50, message: 'Tối đa 50 ký tự' },
                ]}
              >
                <Input placeholder="M250" />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="printed_name"
                label="Tên in"
                rules={[
                  { required: true, message: 'Vui lòng nhập tên in' },
                  { max: 200, message: 'Tối đa 200 ký tự' },
                ]}
              >
                <Input placeholder="Bê tông M250" />
              </Form.Item>
            </Col>
          </Row>

          {COMPONENT_GROUPS.map((group) => (
            <div key={group.title}>
              <Divider orientation="left" plain>
                {group.title}
              </Divider>
              <Row gutter={12}>
                {group.fields.map((f) => (
                  <Col xs={12} sm={8} md={6} key={f.name}>
                    <Form.Item name={f.name} label={f.label}>
                      <InputNumber
                        min={0}
                        step={0.1}
                        style={{ width: '100%' }}
                      />
                    </Form.Item>
                  </Col>
                ))}
              </Row>
            </div>
          ))}

          <Row gutter={16}>
            <Col xs={24} md={8}>
              <Form.Item
                name="do_sut"
                label="Độ sụt"
                rules={[{ max: 100, message: 'Tối đa 100 ký tự' }]}
              >
                <Input placeholder="Ví dụ: 10-12" />
              </Form.Item>
            </Col>
            <Col xs={24} md={8}>
              <Form.Item
                name="cuong_do"
                label="Cường độ"
                rules={[{ max: 100, message: 'Tối đa 100 ký tự' }]}
              >
                <Input placeholder="Ví dụ: B25, M300" />
              </Form.Item>
            </Col>
            <Col xs={24} md={8}>
              <Form.Item name="notes" label="Ghi chú">
                <Input placeholder="Ghi chú" maxLength={500} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </>
  )
}

// Giá trị mặc định 0 cho mọi trường thành phần + số liệu.
function defaultComponents(): GradeInput {
  return {
    grade_code: '',
    printed_name: '',
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
    nuoc: 0,
    pg1: 0,
    pg2: 0,
    pg3: 0,
    pg4: 0,
    pg5: 0,
    pg6: 0,
    do_sut: '',
    cuong_do: '',
    notes: '',
  }
}
