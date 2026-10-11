import { useEffect } from 'react'
import {
  Alert,
  Button,
  Card,
  Checkbox,
  Col,
  Divider,
  Form,
  Input,
  InputNumber,
  Row,
  Space,
  Spin,
  Switch,
  Typography,
  message,
} from 'antd'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { stationApi } from '@/api/station'
import { getErrorMessage } from '@/api/client'
import { ADMIN_MANAGER_ROLES, useAuthStore } from '@/store/auth'
import type { Station, StationInput } from '@/types'

// 17 vị trí thành phần theo đúng thứ tự mask: TP1..TP6, XI1..XI5, PG1..PG6.
const SLOTS = [
  { key: 'tp', label: 'TP (cát/đá)', count: 6 },
  { key: 'xi', label: 'XI (xi măng)', count: 5 },
  { key: 'pg', label: 'PG (phụ gia)', count: 6 },
] as const

// Tạo đối tượng form rỗng (mọi tên + mask mặc định).
function emptyForm(): StationInput {
  const f: Record<string, unknown> = {
    station_code: '',
    station_name: '',
    address: '',
    phone: '',
    component_visibility_mask: '0'.repeat(17),
    accumulate_components: false,
    accumulate_cement: false,
    batch_min_volume_m3: 0,
    batch_max_volume_m3: 0,
  }
  for (const group of SLOTS) {
    for (let i = 1; i <= group.count; i++) {
      f[`${group.key}${i}_name`] = ''
    }
  }
  return f as unknown as StationInput
}

// Chuyển Station (từ API) thành giá trị form.
function toForm(st: Station): StationInput {
  const { id: _id, created_at: _c, updated_at: _u, ...rest } = st
  void _id
  void _c
  void _u
  return {
    ...emptyForm(),
    ...rest,
    component_visibility_mask:
      rest.component_visibility_mask || '0'.repeat(17),
  }
}

// Trang cấu hình trạm: chỉ đọc + cập nhật (không tạo/xóa).
export default function StationConfigPage() {
  const queryClient = useQueryClient()
  const hasRole = useAuthStore((s) => s.hasRole)
  const canEdit = hasRole(...ADMIN_MANAGER_ROLES)
  const [form] = Form.useForm<StationInput>()

  const query = useQuery({
    queryKey: ['station'],
    queryFn: () => stationApi.get(),
  })

  // Nạp dữ liệu vào form khi tải xong.
  useEffect(() => {
    if (query.data) form.setFieldsValue(toForm(query.data))
  }, [query.data, form])

  const saveMutation = useMutation({
    mutationFn: (values: StationInput) => stationApi.update(values),
    onSuccess: () => {
      message.success('Đã lưu cấu hình trạm')
      queryClient.invalidateQueries({ queryKey: ['station'] })
    },
    onError: (err) => message.error(getErrorMessage(err)),
  })

  const handleSave = () => {
    form
      .validateFields()
      .then((values) => saveMutation.mutate(values))
      .catch(() => {
        /* lỗi validate đã hiển thị cạnh field */
      })
  }

  if (query.isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Spin />
      </div>
    )
  }

  if (query.isError) {
    return (
      <Alert
        type="error"
        showIcon
        message="Không tải được cấu hình trạm"
        description={getErrorMessage(query.error)}
        action={
          <Button size="small" onClick={() => query.refetch()}>
            Thử lại
          </Button>
        }
      />
    )
  }

  return (
    <div>
      <Typography.Title level={3} style={{ marginTop: 0 }}>
        Cấu hình trạm
      </Typography.Title>

      <Form form={form} layout="vertical" disabled={!canEdit}>
        <Card title="Thông tin trạm" style={{ marginBottom: 16 }}>
          <Row gutter={16}>
            <Col xs={24} md={12}>
              <Form.Item
                name="station_code"
                label="Mã trạm"
                rules={[
                  { required: true, message: 'Vui lòng nhập mã trạm' },
                  { max: 50, message: 'Tối đa 50 ký tự' },
                ]}
              >
                <Input placeholder="TRAM01" />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="station_name"
                label="Tên trạm"
                rules={[
                  { required: true, message: 'Vui lòng nhập tên trạm' },
                  { max: 200, message: 'Tối đa 200 ký tự' },
                ]}
              >
                <Input placeholder="Trạm trộn số 1" />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item name="address" label="Địa chỉ">
                <Input placeholder="Khu công nghiệp" />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item name="phone" label="Điện thoại">
                <Input placeholder="0901234567" />
              </Form.Item>
            </Col>
          </Row>
        </Card>

        <Card title="Tên hiển thị thành phần" style={{ marginBottom: 16 }}>
          <Typography.Paragraph type="secondary">
            Bật/tắt hiển thị từng vị trí và đặt tên cho nó. Thứ tự mask:
            TP1..TP6, XI1..XI5, PG1..PG6.
          </Typography.Paragraph>
          {SLOTS.map((group) => (
            <div key={group.key}>
              <Divider orientation="left" plain>
                {group.label}
              </Divider>
              <Row gutter={16}>
                {Array.from({ length: group.count }, (_, i) => i + 1).map(
                  (idx) => {
                    const slotKey = `${group.key}${idx}`
                    return (
                      <Col xs={24} sm={12} md={8} key={slotKey}>
                        <Space align="start" style={{ width: '100%' }}>
                          <VisibilityToggle slot={slotKey} form={form} />
                          <Form.Item
                            name={`${slotKey}_name`}
                            label={slotKey.toUpperCase()}
                            style={{ flex: 1, marginBottom: 12 }}
                          >
                            <Input placeholder="Tên hiển thị" />
                          </Form.Item>
                        </Space>
                      </Col>
                    )
                  },
                )}
              </Row>
            </div>
          ))}
        </Card>

        <Card title="Cấu hình mẻ trộn" style={{ marginBottom: 16 }}>
          <Row gutter={16}>
            <Col xs={24} md={12}>
              <Form.Item
                name="batch_min_volume_m3"
                label="Khối lượng mẻ nhỏ nhất (m³)"
              >
                <InputNumber min={0} step={0.1} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="batch_max_volume_m3"
                label="Khối lượng mẻ lớn nhất (m³)"
                dependencies={['batch_min_volume_m3']}
                rules={[
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      const min = getFieldValue('batch_min_volume_m3') ?? 0
                      if (value === undefined || value >= min) {
                        return Promise.resolve()
                      }
                      return Promise.reject(
                        new Error('Phải lớn hơn hoặc bằng giá trị nhỏ nhất'),
                      )
                    },
                  }),
                ]}
              >
                <InputNumber min={0} step={0.1} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="accumulate_components"
                label="Cộng dồn thành phần"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="accumulate_cement"
                label="Cộng dồn xi măng"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>
            </Col>
          </Row>
        </Card>

        {canEdit && (
          <Button
            type="primary"
            onClick={handleSave}
            loading={saveMutation.isPending}
          >
            Lưu cấu hình
          </Button>
        )}
        {!canEdit && (
          <Alert
            type="info"
            showIcon
            message="Bạn chỉ có quyền xem cấu hình trạm."
          />
        )}
      </Form>
    </div>
  )
}

// Ô checkbox điều khiển 1 bit trong component_visibility_mask.
//
// Vị trí slot (ví dụ "tp1", "xi3", "pg6") -> chỉ số bit theo thứ tự
// TP1..TP6 (0..5), XI1..XI5 (6..10), PG1..PG6 (11..16).
function maskIndex(slot: string): number {
  const m = /^([a-z]+)(\d+)$/.exec(slot)
  if (!m) return -1
  const [, prefix, numStr] = m
  const num = Number(numStr)
  const base =
    prefix === 'tp' ? 0 : prefix === 'xi' ? 6 : prefix === 'pg' ? 11 : -1
  return base < 0 ? -1 : base + (num - 1)
}

interface VisToggleProps {
  slot: string
  form: ReturnType<typeof Form.useForm<StationInput>>[0]
}

// Checkbox điều khiển 1 bit của component_visibility_mask. Dùng Form.useWatch
// để đọc mask hiện tại và setFieldValue để cập nhật khi bấm.
function VisibilityToggle({ slot, form }: VisToggleProps) {
  const mask: string = Form.useWatch('component_visibility_mask', form) ?? ''
  const idx = maskIndex(slot)
  const checked = idx >= 0 && mask[idx] === '1'

  const toggle = (next: boolean) => {
    if (idx < 0) return
    const chars = (mask || '0'.repeat(17)).padEnd(17, '0').split('')
    chars[idx] = next ? '1' : '0'
    form.setFieldValue('component_visibility_mask', chars.join(''))
  }

  return (
    <Checkbox
      checked={checked}
      onChange={(e) => toggle(e.target.checked)}
      style={{ marginTop: 6 }}
    />
  )
}
