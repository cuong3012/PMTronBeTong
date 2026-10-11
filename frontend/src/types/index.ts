// Định nghĩa kiểu dữ liệu dùng chung cho toàn frontend.
// Tên trường giữ snake_case để khớp trực tiếp với JSON của Go backend.

/** Kết quả đăng nhập trả về từ POST /api/auth/login */
export interface LoginResult {
  token: string
  expires_at: string
  user: UserInfo
}

/** Thông tin người dùng an toàn (không chứa password_hash). */
export interface UserInfo {
  id: number
  username: string
  full_name: string
  role_code: string
  role_name: string
  station_id?: number | null
}

/** Cấu trúc danh sách có phân trang trả về từ backend. */
export interface PageResult<T> {
  items: T[]
  page: number
  per_page: number
  total: number
}

/** Tham số phân trang/tìm kiếm gửi lên API danh sách. */
export interface ListParams {
  page?: number
  per_page?: number
  search?: string
}

// ---- Các bản ghi danh mục (khớp models của backend) ----

export interface Customer {
  id: number
  customer_code: string
  customer_name: string
  address?: string
  phone?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CustomerInput {
  customer_code: string
  customer_name: string
  address?: string
  phone?: string
}

export interface ConstructionSite {
  id: number
  site_code: string
  site_name: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface ConstructionSiteInput {
  site_code: string
  site_name: string
}

export interface MixerTruck {
  id: number
  truck_code: string
  plate_number: string
  default_driver_id?: number | null
  default_driver_name?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface MixerTruckInput {
  truck_code: string
  plate_number: string
  default_driver_id?: number | null
}

export interface Driver {
  id: number
  driver_code: string
  full_name: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface DriverInput {
  driver_code: string
  full_name: string
}
