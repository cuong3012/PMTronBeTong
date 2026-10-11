import { apiClient } from './client'
import type {
  ConcreteGrade,
  GradeInput,
  ListParams,
  PageResult,
} from '@/types'

// CRUD mác bê tông (khớp backend/internal/grades).
//
// Lưu ý: module này KHÔNG có endpoint restore. Khi trùng mã (kể cả với bản ghi
// đã xóa mềm) backend trả 409 thường, không kèm code INACTIVE_CODE_EXISTS.
export const concreteGradesApi = {
  async list(params: ListParams = {}): Promise<PageResult<ConcreteGrade>> {
    const { data } = await apiClient.get<PageResult<ConcreteGrade>>(
      '/concrete-grades',
      { params },
    )
    return data
  },

  async get(id: number): Promise<ConcreteGrade> {
    const { data } = await apiClient.get<ConcreteGrade>(`/concrete-grades/${id}`)
    return data
  },

  async create(input: GradeInput): Promise<ConcreteGrade> {
    const { data } = await apiClient.post<ConcreteGrade>(
      '/concrete-grades',
      input,
    )
    return data
  },

  async update(id: number, input: GradeInput): Promise<ConcreteGrade> {
    const { data } = await apiClient.put<ConcreteGrade>(
      `/concrete-grades/${id}`,
      input,
    )
    return data
  },

  async remove(id: number): Promise<void> {
    await apiClient.delete(`/concrete-grades/${id}`)
  },
}
