import { apiClient } from './client'
import type { Station, StationInput } from '@/types'

// Cấu hình trạm (một bản ghi duy nhất) — GET /station, PATCH /station.
export const stationApi = {
  async get(): Promise<Station> {
    const { data } = await apiClient.get<Station>('/station')
    return data
  },

  async update(input: StationInput): Promise<Station> {
    const { data } = await apiClient.patch<Station>('/station', input)
    return data
  },
}
