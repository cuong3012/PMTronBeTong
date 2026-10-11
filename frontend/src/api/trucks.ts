import { apiClient } from './client'
import type { ListParams, MixerTruck, MixerTruckInput, PageResult } from '@/types'

// CRUD xe bồn (khớp backend/internal/trucks).
export const trucksApi = {
  async list(params: ListParams = {}): Promise<PageResult<MixerTruck>> {
    const { data } = await apiClient.get<PageResult<MixerTruck>>(
      '/mixer-trucks',
      { params },
    )
    return data
  },

  async get(id: number): Promise<MixerTruck> {
    const { data } = await apiClient.get<MixerTruck>(`/mixer-trucks/${id}`)
    return data
  },

  async create(input: MixerTruckInput): Promise<MixerTruck> {
    const { data } = await apiClient.post<MixerTruck>('/mixer-trucks', input)
    return data
  },

  async update(id: number, input: MixerTruckInput): Promise<MixerTruck> {
    const { data } = await apiClient.put<MixerTruck>(
      `/mixer-trucks/${id}`,
      input,
    )
    return data
  },

  async remove(id: number): Promise<void> {
    await apiClient.delete(`/mixer-trucks/${id}`)
  },

  /** Khôi phục xe bồn đã xóa mềm về trạng thái hoạt động. */
  async restore(id: number): Promise<MixerTruck> {
    const { data } = await apiClient.post<MixerTruck>(
      `/mixer-trucks/${id}/restore`,
    )
    return data
  },
}
