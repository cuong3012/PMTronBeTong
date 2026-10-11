import { apiClient } from './client'
import type {
  ListParams,
  PageResult,
  WeighingStats,
  WeighingTicket,
  WeighingTicketInput,
} from '@/types'

// Phiếu cân (header + chi tiết mẻ) — khớp backend/internal/weighing.
export const weighingTicketsApi = {
  /** Danh sách phiếu cân (KHÔNG kèm details). */
  async list(params: ListParams = {}): Promise<PageResult<WeighingTicket>> {
    const { data } = await apiClient.get<PageResult<WeighingTicket>>(
      '/weighing-tickets',
      { params },
    )
    return data
  },

  /** Chi tiết một phiếu (kèm details). */
  async get(id: number): Promise<WeighingTicket> {
    const { data } = await apiClient.get<WeighingTicket>(
      `/weighing-tickets/${id}`,
    )
    return data
  },

  async create(input: WeighingTicketInput): Promise<WeighingTicket> {
    const { data } = await apiClient.post<WeighingTicket>(
      '/weighing-tickets',
      input,
    )
    return data
  },

  async remove(id: number): Promise<void> {
    await apiClient.delete(`/weighing-tickets/${id}`)
  },

  /** Thống kê theo mác và theo ngày (from/to dạng YYYY-MM-DD). */
  async stats(params: { from?: string; to?: string } = {}): Promise<WeighingStats> {
    const { data } = await apiClient.get<WeighingStats>(
      '/weighing-tickets/stats',
      { params },
    )
    return data
  },
}
