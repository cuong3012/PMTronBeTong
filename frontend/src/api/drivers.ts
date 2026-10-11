import { apiClient } from './client'
import type { Driver, DriverInput, ListParams, PageResult } from '@/types'

// CRUD tài xế (khớp backend/internal/drivers).
export const driversApi = {
  async list(params: ListParams = {}): Promise<PageResult<Driver>> {
    const { data } = await apiClient.get<PageResult<Driver>>('/drivers', {
      params,
    })
    return data
  },

  async get(id: number): Promise<Driver> {
    const { data } = await apiClient.get<Driver>(`/drivers/${id}`)
    return data
  },

  async create(input: DriverInput): Promise<Driver> {
    const { data } = await apiClient.post<Driver>('/drivers', input)
    return data
  },

  async update(id: number, input: DriverInput): Promise<Driver> {
    const { data } = await apiClient.put<Driver>(`/drivers/${id}`, input)
    return data
  },

  async remove(id: number): Promise<void> {
    await apiClient.delete(`/drivers/${id}`)
  },

  /** Khôi phục tài xế đã xóa mềm về trạng thái hoạt động. */
  async restore(id: number): Promise<Driver> {
    const { data } = await apiClient.post<Driver>(`/drivers/${id}/restore`)
    return data
  },
}
