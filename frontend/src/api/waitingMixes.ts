import { apiClient } from './client'
import type { WaitingMix, WaitingMixInput } from '@/types'

// Cấp phối chờ — đúng 20 vị trí cố định, chỉ đọc/cập nhật.
//
// Lưu ý: GET trả về { items: [...] } (KHÔNG phải PageResult); không có POST/DELETE.
export const waitingMixesApi = {
  async list(): Promise<WaitingMix[]> {
    const { data } = await apiClient.get<{ items: WaitingMix[] }>(
      '/waiting-mixes',
    )
    return data.items
  },

  async get(id: number): Promise<WaitingMix> {
    const { data } = await apiClient.get<WaitingMix>(`/waiting-mixes/${id}`)
    return data
  },

  async update(id: number, input: WaitingMixInput): Promise<WaitingMix> {
    const { data } = await apiClient.patch<WaitingMix>(
      `/waiting-mixes/${id}`,
      input,
    )
    return data
  },
}
