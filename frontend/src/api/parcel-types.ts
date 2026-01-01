// 包裹状态枚举
export enum ParcelStatus {
  READY_FOR_PICKUP = 'ready_for_pickup', // 待取件
  PICKED_UP = 'picked_up', // 已取件
  OVERDUE = 'overdue', // 滞留
  RETURNED = 'returned' // 已退回
}

// 包裹尺寸

export type ParcelSize = 'small' | 'medium' | 'large'

// 货架信息
export interface Shelf {
  id: number
  shelf_code: string
  area: string
  floor: number
  column: number
  type: string
  capacity: number
  current_count?: number
  status?: string
}

// 包裹基本信息
export interface Parcel {
  id: number
  tracking_number: string
  pickup_code: string
  recipient_name: string
  recipient_phone: string
  courier_company: string
  size: ParcelSize
  weight: number
  status: ParcelStatus
  shelf?: Shelf
  received_at: string
  shelved_at?: string
  picked_up_at?: string
  expected_overdue_at?: string
  notes?: string
}

// 包裹列表响应
export interface ParcelListResponse {
  list: Parcel[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// 包裹入库请求
export interface ReceiveParcelRequest {
  tracking_number: string
  recipient_name: string
  recipient_phone: string
  recipient_id_card?: string
  courier_company: string
  size: ParcelSize
  weight: number
  notes?: string
}

// 包裹取件请求
export interface PickupParcelRequest {
  pickup_code: string
  recipient_phone: string
}

// 寄件单状态
export enum ShipmentStatus {
  PENDING = 'pending', // 待发货
  SHIPPED = 'shipped', // 已发货
  DELIVERED = 'delivered', // 已送达
  CANCELLED = 'cancelled' // 已取消
}

// 寄件单信息
export interface Shipment {
  id: number
  shipment_number: string
  sender_name: string
  sender_phone: string
  recipient_name: string
  recipient_phone: string
  recipient_address: string
  courier_company: string
  weight: number
  volume: number
  freight: number
  status: ShipmentStatus
  created_at: string
  shipped_at?: string
}

// 寄件单列表响应
export interface ShipmentListResponse {
  list: Shipment[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// 创建寄件单请求
export interface CreateShipmentRequest {
  recipient_name: string
  recipient_phone: string
  recipient_address: string
  courier_company: string
  weight: number
  volume: number
}

// 运费计算请求
export interface CalculateFreightRequest {
  courier_company: string
  weight: number
  volume: number
}

// 运费计算响应
export interface CalculateFreightResponse {
  freight: number
  calculation_details: {
    base_freight: number
    weight_freight: number
    volume_freight: number
  }
}

// 统计数据
export interface DashboardStatistics {
  today: {
    received_count: number
    picked_up_count: number
    overdue_count: number
  }
  parcels: {
    total: number
    received: number
    shelved: number
    ready_for_pickup: number
    picked_up: number
    overdue: number
    returned: number
  }
  shelves: {
    total_capacity: number
    total_used: number
    utilization_rate: number
  }
  recent_trends: Array<{
    date: string
    received: number
    picked_up: number
  }>
}

// ========== 货架相关类型 ==========

// 货架列表响应
export interface ShelfListResponse {
  list: Shelf[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// 货架统计
export interface ShelfStatistics {
  total_shelves: number
  active_shelves: number
  total_capacity: number
  total_used: number
  utilization_rate: number
  by_type: {
    type: string
    count: number
    capacity: number
    used: number
    utilization_rate: number
  }[]
}

// 创建货架请求
export interface CreateShelfRequest {
  shelf_code: string
  area: string
  floor: number
  column: number
  type: string
  capacity: number
}

// 更新货架请求
export interface UpdateShelfRequest {
  area?: string
  floor?: number
  column?: number
  type?: string
  capacity?: number
  status?: string
}

// ========== 滞留件相关类型 ==========

// 滞留件信息
export interface OverdueInfo {
  id: number
  parcel_id: number
  tracking_number: string
  pickup_code: string
  recipient_name: string
  recipient_phone: string
  courier_company: string
  size: string
  shelf_code?: string
  overdue_start_at: string
  overdue_days: number
  reminder_count: number
  last_reminder_at?: string
  overdue_fee: number
  status: string
}

// 滞留件列表响应
export interface OverdueListResponse {
  list: OverdueInfo[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// 批量操作结果
export interface BatchResult {
  success: number
  failed: number
  errors?: string[]
}

// 批量退回请求
export interface BatchReturnRequest {
  ids: number[]
}
