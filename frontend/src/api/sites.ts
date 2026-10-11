import { apiClient } from './client'
import type {
  ConstructionSite,
  ConstructionSiteInput,
  ListParams,
  PageResult,
} from '@/types'

// CRUD công trường (khớp backend/internal/sites).
export const sitesApi = {
  async list(params: ListParams = {}): Promise<PageResult<ConstructionSite>> {
    const { data } = await apiClient.get<PageResult<ConstructionSite>>(
      '/construction-sites',
      { params },
    )
    return data
  },

  async get(id: number): Promise<ConstructionSite> {
    const { data } = await apiClient.get<ConstructionSite>(
      `/construction-sites/${id}`,
    )
    return data
  },

  async create(input: ConstructionSiteInput): Promise<ConstructionSite> {
    const { data } = await apiClient.post<ConstructionSite>(
      '/construction-sites',
      input,
    )
    return data
  },

  async update(
    id: number,
    input: ConstructionSiteInput,
  ): Promise<ConstructionSite> {
    const { data } = await apiClient.put<ConstructionSite>(
      `/construction-sites/${id}`,
      input,
    )
    return data
  },

  async remove(id: number): Promise<void> {
    await apiClient.delete(`/construction-sites/${id}`)
  },

  /** Khôi phục công trường đã xóa mềm về trạng thái hoạt động. */
  async restore(id: number): Promise<ConstructionSite> {
    const { data } = await apiClient.post<ConstructionSite>(
      `/construction-sites/${id}/restore`,
    )
    return data
  },
}
