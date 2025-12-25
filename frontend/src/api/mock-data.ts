import {
  Parcel,
  ParcelStatus,
  ParcelSize,
  Shipment,
  ShipmentStatus,
  DashboardStatistics,
  Shelf
} from './parcel-types'

// Mock 快递公司列表
export const COURIER_COMPANIES = [
  '顺丰速运',
  '中通快递',
  '圆通速递',
  '申通快递',
  '韵达快递',
  '邮政EMS',
  '京东物流',
  '德邦快递'
]

// Mock 货架数据
export const MOCK_SHELVES: Shelf[] = Array.from({ length: 50 }, (_, i) => ({
  id: i + 1,
  shelf_code: `${String.fromCharCode(65 + Math.floor(i / 10))}-${String(Math.floor((i % 10) / 5) + 1).padStart(2, '0')}-${String((i % 5) + 1).padStart(2, '0')}`,
  area: String.fromCharCode(65 + Math.floor(i / 10)),
  floor: Math.floor((i % 10) / 5) + 1,
  column: (i % 5) + 1,
  type: ['small', 'medium', 'large'][i % 3],
  capacity: [20, 15, 10][i % 3],
  current_count: Math.floor(Math.random() * [20, 15, 10][i % 3]),
  status: 'available'
}))

// 生成随机包裹数据
const generateMockParcels = (): Parcel[] => {
  const statuses: ParcelStatus[] = [
    ParcelStatus.RECEIVED,
    ParcelStatus.SHELVED,
    ParcelStatus.READY_FOR_PICKUP,
    ParcelStatus.PICKED_UP,
    ParcelStatus.OVERDUE
  ]

  const sizes: ParcelSize[] = ['small', 'medium', 'large']

  const recipients = ['张三', '李四', '王五', '赵六', '钱七', '孙八', '周九', '吴十']

  return Array.from({ length: 100 }, (_, i) => {
    const status = statuses[Math.floor(Math.random() * statuses.length)]
    const receivedAt = new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000)

    return {
      id: i + 1,
      tracking_number: `SF${String(1234567890 + i).padStart(12, '0')}`,
      pickup_code: `${String.fromCharCode(65 + Math.floor(Math.random() * 26))}${String(10000 + i).slice(-5)}`,
      recipient_name: recipients[i % recipients.length],
      recipient_phone: `138${String(10000000 + i).slice(-8)}`,
      courier_company: COURIER_COMPANIES[i % COURIER_COMPANIES.length],
      size: sizes[i % 3],
      weight: parseFloat((Math.random() * 10 + 0.5).toFixed(2)),
      status,
      shelf: status !== ParcelStatus.PICKED_UP ? MOCK_SHELVES[i % 50] : undefined,
      received_at: receivedAt.toISOString(),
      shelved_at:
        status !== ParcelStatus.RECEIVED
          ? new Date(receivedAt.getTime() + 30 * 60 * 1000).toISOString()
          : undefined,
      picked_up_at:
        status === ParcelStatus.PICKED_UP
          ? new Date(receivedAt.getTime() + 2 * 24 * 60 * 60 * 1000).toISOString()
          : undefined,
      expected_overdue_at: new Date(receivedAt.getTime() + 3 * 24 * 60 * 60 * 1000).toISOString(),
      notes: i % 5 === 0 ? '易碎品，请小心搬运' : undefined
    }
  })
}

// Mock 包裹数据
export const MOCK_PARCELS = generateMockParcels()

// 生成随机寄件单数据
const generateMockShipments = (): Shipment[] => {
  const statuses: ShipmentStatus[] = [
    ShipmentStatus.PENDING,
    ShipmentStatus.SHIPPED,
    ShipmentStatus.DELIVERED
  ]

  const addresses = [
    '广东省广州市天河区五山路381号',
    '广东省深圳市南山区科技园',
    '广东省珠海市香洲区情侣路',
    '广东省佛山市顺德区大良街道'
  ]

  return Array.from({ length: 30 }, (_, i) => {
    const status = statuses[Math.floor(Math.random() * statuses.length)]
    const createdAt = new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000)

    return {
      id: i + 1,
      shipment_number: `SH${new Date().getFullYear()}${String(i + 1).padStart(8, '0')}`,
      sender_name: '张三',
      sender_phone: '13800138000',
      recipient_name: `收件人${i + 1}`,
      recipient_phone: `139${String(10000000 + i).slice(-8)}`,
      recipient_address: addresses[i % addresses.length],
      courier_company: COURIER_COMPANIES[i % COURIER_COMPANIES.length],
      weight: parseFloat((Math.random() * 5 + 0.5).toFixed(2)),
      volume: parseFloat((Math.random() * 0.05 + 0.001).toFixed(3)),
      freight: parseFloat((Math.random() * 20 + 10).toFixed(2)),
      status,
      created_at: createdAt.toISOString(),
      shipped_at:
        status !== ShipmentStatus.PENDING
          ? new Date(createdAt.getTime() + 60 * 60 * 1000).toISOString()
          : undefined
    }
  })
}

// Mock 寄件单数据
export const MOCK_SHIPMENTS = generateMockShipments()

// Mock 统计数据
export const MOCK_DASHBOARD_STATISTICS: DashboardStatistics = {
  today: {
    received_count: 150,
    picked_up_count: 120,
    overdue_count: 5
  },
  parcels: {
    total: 1000,
    received: 100,
    shelved: 50,
    ready_for_pickup: 300,
    picked_up: 500,
    overdue: 30,
    returned: 20
  },
  shelves: {
    total_capacity: 2000,
    total_used: 1500,
    utilization_rate: 75.0
  },
  recent_trends: Array.from({ length: 7 }, (_, i) => {
    const date = new Date()
    date.setDate(date.getDate() - (6 - i))
    return {
      date: date.toISOString().split('T')[0],
      received: Math.floor(Math.random() * 50 + 100),
      picked_up: Math.floor(Math.random() * 50 + 80)
    }
  })
}
