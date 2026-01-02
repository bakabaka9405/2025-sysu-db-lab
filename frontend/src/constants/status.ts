import type { TagProps } from 'naive-ui'

// 状态配置接口
export interface StatusConfig {
  text: string
  type: TagProps['type']
}

// 包裹状态映射
export const PARCEL_STATUS_MAP: Record<string, StatusConfig> = {
  received: { text: '待上架', type: 'info' },
  ready: { text: '待取件', type: 'success' },
  picked_up: { text: '已取件', type: 'default' },
  overdue: { text: '滞留', type: 'warning' },
  returned: { text: '已退回', type: 'error' }
}

// 寄件状态映射
export const SHIPMENT_STATUS_MAP: Record<string, StatusConfig> = {
  pending: { text: '待发货', type: 'warning' },
  shipped: { text: '已发货', type: 'info' },
  delivered: { text: '已送达', type: 'success' },
  cancelled: { text: '已取消', type: 'error' }
}

// 滞留件状态映射
export const OVERDUE_STATUS_MAP: Record<string, StatusConfig> = {
  pending: { text: '待处理', type: 'warning' },
  reminded: { text: '已提醒', type: 'info' },
  returned: { text: '已退回', type: 'default' }
}

// 货架状态映射
export const SHELF_STATUS_MAP: Record<string, StatusConfig> = {
  active: { text: '启用', type: 'success' },
  inactive: { text: '停用', type: 'warning' }
}

// 包裹状态筛选选项
export const PARCEL_STATUS_OPTIONS = [
  { label: '全部', value: '' },
  { label: '待上架', value: 'received' },
  { label: '待取件', value: 'ready' },
  { label: '已取件', value: 'picked_up' },
  { label: '滞留', value: 'overdue' },
  { label: '已退回', value: 'returned' }
]

// 滞留件状态筛选选项
export const OVERDUE_STATUS_OPTIONS = [
  { label: '全部', value: '' },
  { label: '待处理', value: 'pending' },
  { label: '已提醒', value: 'reminded' },
  { label: '已退回', value: 'returned' }
]
