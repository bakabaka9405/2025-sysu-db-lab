<script setup lang="ts">
import { computed, h } from 'vue'
import { NDataTable, NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import type { Shipment } from '../../api/types'
import { SHIPMENT_STATUS_MAP } from '../../constants/status'
import { formatDateTime, formatAmount, formatWeight } from '../../utils/formatters'
import { usePaginationConfig, type PaginationProps } from '../../utils/pagination'

export interface ShipmentTableProps {
  shipments: Shipment[]
  loading?: boolean
  compact?: boolean
  pagination?: PaginationProps | false
  scrollX?: number
}

const props = withDefaults(defineProps<ShipmentTableProps>(), {
  loading: false,
  compact: false,
  scrollX: 1200
})

// 构建列定义
const columns = computed<DataTableColumns<Shipment>>(() => {
  const cols: DataTableColumns<Shipment> = []

  // 寄件单号
  cols.push({ title: '寄件单号', key: 'shipment_number', width: 180 })

  // 收件人
  cols.push({ title: '收件人', key: 'recipient_name', width: 100 })

  // 非紧凑模式显示更多信息
  if (!props.compact) {
    cols.push(
      { title: '收件电话', key: 'recipient_phone', width: 120 },
      { title: '收件地址', key: 'recipient_address', width: 200, ellipsis: { tooltip: true } }
    )
  }

  // 快递公司
  cols.push({ title: '快递公司', key: 'courier_company', width: 100 })

  // 非紧凑模式显示重量
  if (!props.compact) {
    cols.push({
      title: '重量',
      key: 'weight',
      width: 80,
      render: (row) => formatWeight(row.weight)
    })
  }

  // 运费
  cols.push({
    title: '运费',
    key: 'freight',
    width: 80,
    render: (row) => formatAmount(row.freight)
  })

  // 状态
  cols.push({
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const status = SHIPMENT_STATUS_MAP[row.status] || { text: row.status, type: 'default' }
      return h(NTag, { type: status.type }, { default: () => status.text })
    }
  })

  // 非紧凑模式显示创建时间
  if (!props.compact) {
    cols.push({
      title: '创建时间',
      key: 'created_at',
      width: 160,
      render: (row) => formatDateTime(row.created_at)
    })
  }

  return cols
})

// 分页配置
const paginationConfig = usePaginationConfig(() => props.pagination)
</script>

<template>
  <NDataTable
    :remote="!!pagination"
    :columns="columns"
    :data="shipments"
    :loading="loading"
    :pagination="paginationConfig"
    :scroll-x="compact ? undefined : scrollX"
  />
</template>
