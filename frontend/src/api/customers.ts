import { apiClient } from './client'
import type { Customer, CustomerInput, ListParams, PageResult } from '@/types'

// CRUD khách hàng (khớp backend/internal/customers).
export const customersApi = {
  async list(params: ListParams = {}): Promise<PageResult<Customer>> {
    const { data } = await apiClient.get<PageResult<Customer>>('/customers', {
      params,
    })
    return data
  },

  async get(id: number): Promise<Customer> {
    const { data } = await apiClient.get<Customer>(`/customers/${id}`)
    return data
  },

  async create(input: CustomerInput): Promise<Customer> {
    const { data } = await apiClient.post<Customer>('/customers', input)
    return data
  },

  async update(id: number, input: CustomerInput): Promise<Customer> {
    const { data } = await apiClient.put<Customer>(`/customers/${id}`, input)
    return data
  },

  async remove(id: number): Promise<void> {
    await apiClient.delete(`/customers/${id}`)
  },

  /** Khôi phục khách hàng đã xóa mềm về trạng thái hoạt động. */
  async restore(id: number): Promise<Customer> {
    const { data } = await apiClient.post<Customer>(`/customers/${id}/restore`)
    return data
  },
}
