import { apiClient } from './client'
import type { CatalogInput, ComponentCatalog } from '@/types'

// Danh mục thành phần (metadata vị trí TP/XI/PG) — không phân trang.
//
// Lưu ý: GET trả về { items: [...] } (KHÔNG phải PageResult).
export const componentCatalogApi = {
  async list(search = ''): Promise<ComponentCatalog[]> {
    const { data } = await apiClient.get<{ items: ComponentCatalog[] }>(
      '/component-catalog',
      { params: search ? { search } : {} },
    )
    return data.items
  },

  async get(id: number): Promise<ComponentCatalog> {
    const { data } = await apiClient.get<ComponentCatalog>(
      `/component-catalog/${id}`,
    )
    return data
  },

  async create(input: CatalogInput): Promise<ComponentCatalog> {
    const { data } = await apiClient.post<ComponentCatalog>(
      '/component-catalog',
      input,
    )
    return data
  },

  async update(id: number, input: CatalogInput): Promise<ComponentCatalog> {
    const { data } = await apiClient.put<ComponentCatalog>(
      `/component-catalog/${id}`,
      input,
    )
    return data
  },

  async remove(id: number): Promise<void> {
    await apiClient.delete(`/component-catalog/${id}`)
  },
}
