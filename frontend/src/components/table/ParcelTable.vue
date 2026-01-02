<script setup lang="ts">
import { computed, h } from 'vue'
import { NDataTable, NTag, NButton } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import type { Parcel } from '../../api/types'
import { PARCEL_STATUS_MAP } from '../../constants/status'
import { PICKUP_CODE_STYLE } from '../../constants/styles'
import { formatDateTime } from '../../utils/formatters'
import { usePaginationConfig, type PaginationProps } from '../../utils/pagination'

export interface ParcelTableProps {
  parcels: Parcel[]
  loading?: boolean
  compact?: boolean // 紧凑模式（Dashboard用，少列）
  showId?: boolean // 显示ID列
  showRecipient?: boolean // 显示收件人信息
  showNotes?: boolean // 显示备注列
  showActions?: boolean // 显示操作列
  pagination?: PaginationProps | false
  scrollX?: number
}

const props = withDefaults(defineProps<ParcelTableProps>(), {
  loading: false,
  compact: false,
  showId: false,
  showRecipient: false,
  showNotes: false,
  showActions: false,
  scrollX: 1000
})

const emit = defineEmits<{
  action: [action: string, parcel: Parcel]
}>()

// 构建列定义
const columns = computed<DataTableColumns<Parcel>>(() => {
  const cols: DataTableColumns<Parcel> = []

  // ID列
  if (props.showId) {
    cols.push({ title: 'ID', key: 'id', width: 60 })
  }

  // 取件码列（始终显示）
  cols.push({
    title: '取件码',
    key: 'pickup_code',
    width: 100,
    render: (row) => h('strong', { style: PICKUP_CODE_STYLE }, row.pickup_code)
  })

  // 非紧凑模式下显示更多列
  if (!props.compact) {
    cols.push({ title: '快递单号', key: 'tracking_number', width: 140 })
  }

  // 收件人信息
  if (props.showRecipient) {
    cols.push(
      { title: '收件人', key: 'recipient_name', width: 80 },
      { title: '手机号', key: 'recipient_phone', width: 120 }
    )
  }

  // 快递公司
  cols.push({ title: '快递公司', key: 'courier_company', width: 100 })

  // 货架位置（非紧凑模式）
  if (!props.compact) {
    cols.push({
      title: '货架',
      key: 'shelf',
      width: 100,
      render: (row) => row.shelf?.shelf_code || '-'
    })
  }

  // 状态列
  cols.push({
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const status = PARCEL_STATUS_MAP[row.status] || { text: row.status, type: 'default' }
      return h(NTag, { type: status.type }, { default: () => status.text })
    }
  })

  // 入库时间
  cols.push({
    title: '入库时间',
    key: 'received_at',
    width: 160,
    render: (row) => formatDateTime(row.received_at)
  })

  // 备注列
  if (props.showNotes) {
    cols.push({ title: '备注', key: 'notes', ellipsis: { tooltip: true } })
  }

  // 操作列
  if (props.showActions) {
    cols.push({
      title: '操作',
      key: 'actions',
      width: 120,
      render: (row) => {
        if (row.status === 'received') {
          return h(NButton, {
            size: 'small',
            type: 'primary',
            onClick: () => emit('action', 'shelve', row)
          }, { default: () => '上架' })
        }
        return null
      }
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
    :data="parcels"
    :loading="loading"
    :pagination="paginationConfig"
    :scroll-x="compact ? undefined : scrollX"
  />
</template>
