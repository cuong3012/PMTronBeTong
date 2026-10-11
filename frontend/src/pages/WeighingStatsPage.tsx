import { useState } from 'react'
import {
  Alert,
  Button,
  Card,
  Col,
  DatePicker,
  Row,
  Statistic,
  Table,
  Typography,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import type { Dayjs } from 'dayjs'
import { useQuery } from '@tanstack/react-query'
import { weighingTicketsApi } from '@/api/weighingTickets'
import { getErrorMessage } from '@/api/client'
import { formatNumber } from '@/utils/format'
import type { DailyStat, GradeStat } from '@/types'

type RangeValue = [Dayjs | null, Dayjs | null] | null

const gradeColumns: ColumnsType<GradeStat> = [
  { title: 'Mã mác', dataIndex: 'concrete_grade_code', key: 'code' },
  { title: 'Tên mác', dataIndex: 'concrete_grade_name', key: 'name' },
  { title: 'Số phiếu', dataIndex: 'ticket_count', key: 'ticket_count' },
  {
    title: 'Tổng khối lượng (m³)',
    dataIndex: 'total_volume_m3',
    key: 'total_volume_m3',
    render: (v: number) => formatNumber(v),
  },
  {
    title: 'Tổng số mẻ',
    dataIndex: 'total_batch_count',
    key: 'total_batch_count',
  },
]

const dailyColumns: ColumnsType<DailyStat> = [
  { title: 'Ngày', dataIndex: 'day', key: 'day' },
  { title: 'Số phiếu', dataIndex: 'ticket_count', key: 'ticket_count' },
  {
    title: 'Tổng khối lượng (m³)',
    dataIndex: 'total_volume_m3',
    key: 'total_volume_m3',
    render: (v: number) => formatNumber(v),
  },
]

// Trang thống kê phiếu cân theo mác và theo ngày.
export default function WeighingStatsPage() {
  const [range, setRange] = useState<RangeValue>(null)
  const from = range?.[0]?.format('YYYY-MM-DD')
  const to = range?.[1]?.format('YYYY-MM-DD')

  const query = useQuery({
    queryKey: ['weighing-stats', { from, to }],
    queryFn: () => weighingTicketsApi.stats({ from, to }),
  })

  const totalTickets = (query.data?.by_grade ?? []).reduce(
    (sum, g) => sum + g.ticket_count,
    0,
  )
  const totalVolume = (query.data?.by_grade ?? []).reduce(
    (sum, g) => sum + g.total_volume_m3,
    0,
  )

  return (
    <div>
      <Typography.Title level={3} style={{ marginTop: 0 }}>
        Thống kê phiếu cân
      </Typography.Title>

      <Card style={{ marginBottom: 16 }}>
        <Row gutter={16} align="middle">
          <Col>
            <Typography.Text strong>Khoảng ngày:</Typography.Text>
          </Col>
          <Col>
            <DatePicker.RangePicker
              value={range}
              onChange={(value) => setRange(value as RangeValue)}
              placeholder={['Từ ngày', 'Đến ngày']}
            />
          </Col>
          <Col>
            <Button onClick={() => setRange(null)}>Xóa lọc</Button>
          </Col>
        </Row>
      </Card>

      {query.isError && (
        <Alert
          type="error"
          showIcon
          style={{ marginBottom: 16 }}
          message="Không tải được thống kê"
          description={getErrorMessage(query.error)}
          action={
            <Button size="small" onClick={() => query.refetch()}>
              Thử lại
            </Button>
          }
        />
      )}

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={12}>
          <Card>
            <Statistic
              title="Tổng số phiếu"
              value={totalTickets}
              loading={query.isLoading}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12}>
          <Card>
            <Statistic
              title="Tổng khối lượng (m³)"
              value={totalVolume}
              precision={2}
              loading={query.isLoading}
            />
          </Card>
        </Col>
      </Row>

      <Card title="Theo mác bê tông" style={{ marginBottom: 16 }}>
        <Table<GradeStat>
          rowKey={(r) => r.concrete_grade_code}
          columns={gradeColumns}
          dataSource={query.data?.by_grade ?? []}
          loading={query.isLoading}
          pagination={false}
          scroll={{ x: 'max-content' }}
        />
      </Card>

      <Card title="Theo ngày">
        <Table<DailyStat>
          rowKey={(r) => r.day}
          columns={dailyColumns}
          dataSource={query.data?.daily ?? []}
          loading={query.isLoading}
          pagination={false}
          scroll={{ x: 'max-content' }}
        />
      </Card>
    </div>
  )
}
