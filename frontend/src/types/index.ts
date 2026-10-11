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

// ---- Cấu hình trạm (GET/PATCH /station) ----

export interface Station {
  id: number
  station_code: string
  station_name: string
  address?: string
  phone?: string
  tp1_name: string
  tp2_name: string
  tp3_name: string
  tp4_name: string
  tp5_name: string
  tp6_name: string
  xi1_name: string
  xi2_name: string
  xi3_name: string
  xi4_name: string
  xi5_name: string
  pg1_name: string
  pg2_name: string
  pg3_name: string
  pg4_name: string
  pg5_name: string
  pg6_name: string
  /** Mask đúng 17 ký tự 0/1 theo thứ tự TP1..TP6, XI1..XI5, PG1..PG6. */
  component_visibility_mask: string
  accumulate_components: boolean
  accumulate_cement: boolean
  batch_min_volume_m3: number
  batch_max_volume_m3: number
  created_at: string
  updated_at: string
}

/** Body cập nhật cấu hình trạm (PATCH). */
export type StationInput = Omit<Station, 'id' | 'created_at' | 'updated_at'>

// ---- Mác bê tông / cấp phối (concrete-grades) ----

export interface ConcreteGrade {
  id: number
  grade_code: string
  printed_name: string
  tp1: number
  tp2: number
  tp3: number
  tp4: number
  tp5: number
  tp6: number
  xi1: number
  xi2: number
  xi3: number
  xi4: number
  xi5: number
  nuoc: number
  pg1: number
  pg2: number
  pg3: number
  pg4: number
  pg5: number
  pg6: number
  /** Độ sụt lưu dạng văn bản (ví dụ "10-12") — không phải số. */
  do_sut: string
  /** Cường độ lưu dạng văn bản (ví dụ "B25", "M300") — không phải số. */
  cuong_do: string
  notes?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

/** Body tạo/cập nhật mác bê tông. Thành phần không dùng để 0. */
export type GradeInput = Omit<
  ConcreteGrade,
  'id' | 'is_active' | 'created_at' | 'updated_at'
>

// ---- Danh mục thành phần (component-catalog) ----

export interface ComponentCatalog {
  id: number
  slot_name: string
  humidity_values: string
  unit?: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface CatalogInput {
  slot_name: string
  humidity_values: string
  unit?: string
  notes?: string
}

// ---- Cấp phối chờ (waiting-mixes) ----

export interface WaitingMix {
  id: number
  stt: number
  concrete_grade_id?: number | null
  concrete_grade_code?: string
  concrete_grade_name?: string
  customer_id?: number | null
  customer_name?: string
  volume_m3: number
  m3me: number
  status: number
  created_at: string
  updated_at: string
}

/** Body cập nhật một vị trí cấp phối chờ (PATCH). */
export interface WaitingMixInput {
  concrete_grade_id: number
  customer_id: number
  volume_m3: number
  m3me: number
  status: number
}

// ---- Phiếu cân (weighing-tickets) ----

export interface WeighingTicketDetail {
  id: number
  weighing_ticket_id: number
  batch_no: number
  batch_volume_m3: number
  mixed_at: string
  tp1: number
  tp2: number
  tp3: number
  tp4: number
  tp5: number
  tp6: number
  xi1: number
  xi2: number
  xi3: number
  xi4: number
  xi5: number
  pg1: number
  pg2: number
  pg3: number
  pg4: number
  pg5: number
  pg6: number
  manual_override_json?: string
  notes?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface WeighingTicket {
  id: number
  ticket_code: string
  volume_m3: number
  batch_count: number
  m3me: number
  order_id: number
  order_number?: string
  concrete_grade_id: number
  concrete_grade_code: string
  concrete_grade_name: string
  printed_concrete_grade_name: string
  humidity_values: string
  standard_per_cube: number
  standard_per_batch: number
  customer_name: string
  construction_site_address: string
  material_name: string
  runner_user_id?: number | null
  notes?: string
  is_active: boolean
  created_at: string
  updated_at: string
  /** Chỉ có khi gọi GET /:id (chi tiết). */
  details?: WeighingTicketDetail[]
}

/** Một mẻ chi tiết khi tạo phiếu cân (POST). */
export interface WeighingTicketDetailInput {
  batch_no: number
  batch_volume_m3: number
  mixed_at: string
  tp1?: number
  tp2?: number
  tp3?: number
  tp4?: number
  tp5?: number
  tp6?: number
  xi1?: number
  xi2?: number
  xi3?: number
  xi4?: number
  xi5?: number
  pg1?: number
  pg2?: number
  pg3?: number
  pg4?: number
  pg5?: number
  pg6?: number
  manual_override_json?: string
  notes?: string
}

export interface WeighingTicketInput {
  ticket_code: string
  volume_m3: number
  batch_count: number
  m3me?: number
  order_id: number
  concrete_grade_id: number
  concrete_grade_code: string
  concrete_grade_name: string
  printed_concrete_grade_name: string
  humidity_values?: string
  standard_per_cube?: number
  standard_per_batch?: number
  customer_name: string
  construction_site_address?: string
  material_name?: string
  runner_user_id?: number | null
  notes?: string
  details: WeighingTicketDetailInput[]
}

// ---- Thống kê phiếu cân (weighing-tickets/stats) ----

export interface GradeStat {
  concrete_grade_code: string
  concrete_grade_name: string
  ticket_count: number
  total_volume_m3: number
  total_batch_count: number
}

export interface DailyStat {
  day: string
  ticket_count: number
  total_volume_m3: number
}

export interface WeighingStats {
  by_grade: GradeStat[]
  daily: DailyStat[]
}
