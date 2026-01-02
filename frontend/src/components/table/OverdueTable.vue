<script setup lang="ts">
import { computed, h } from 'vue'
import { NDataTable, NTag, NButton, NSpace, NPopconfirm } from 'naive-ui'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import type { OverdueInfo } from '../../api/overdue'
import { OVERDUE_STATUS_MAP } from '../../constants/status'
import { PICKUP_CODE_STYLE, getOverdueDaysColor } from '../../constants/styles'
import { formatAmount } from '../../utils/formatters'
import { usePaginationConfig, type PaginationProps } from '../../utils/pagination'

export interface OverdueTableProps {
  overdueItems: OverdueInfo[]
  loading?: boolean
  checkedKeys?: DataTableRowKey[]
  pagination?: PaginationProps | false
  scrollX?: number
}

const props = withDefaults(defineProps<OverdueTableProps>(), {
  loading: false,
  checkedKeys: () => [],
  scrollX: 1300
})

const emit = defineEmits<{
  remind: [item: OverdueInfo]
  return: [ids: number[]]
  'update:checkedKeys': [keys: DataTableRowKey[]]
}>()

// 表格列定义
const columns = computed<DataTableColumns<OverdueInfo>>(() => [
  { type: 'selection' },
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '取件码',
    key: 'pickup_code',
    width: 100,
    render: (row) => h('strong', { style: PICKUP_CODE_STYLE }, row.pickup_code)
  },
  { title: '快递单号', key: 'tracking_number', width: 140 },
  { title: '收件人', key: 'recipient_name', width: 80 },
  { title: '手机号', key: 'recipient_phone', width: 120 },
  { title: '货架', key: 'shelf_code', width: 80, render: (row) => row.shelf_code || '-' },
  {
    title: '滞留天数',
    key: 'overdue_days',
    width: 100,
    render: (row) => h('span', {
      style: { color: getOverdueDaysColor(row.overdue_days), fontWeight: 'bold' }
    }, `${row.overdue_days} 天`)
  },
  { title: '提醒次数', key: 'reminder_count', width: 80 },
  {
    title: '滞留费',
    key: 'overdue_fee',
    width: 80,
    render: (row) => formatAmount(row.overdue_fee)
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) => {
      const status = OVERDUE_STATUS_MAP[row.status] || { text: row.status, type: 'default' }
      return h(NTag, { type: status.type }, { default: () => status.text })
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) => {
      if (row.status === 'returned') {
        return h('span', { style: 'color: #999;' }, '已退回')
      }
      return h(NSpace, {}, {
        default: () => [
          h(NButton, {
            size: 'small',
            type: 'info',
            onClick: () => emit('remind', row)
          }, { default: () => '发送提醒' }),
          h(NPopconfirm, {
            onPositiveClick: () => emit('return', [row.id])
          }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '退回' }),
            default: () => '确定退回该包裹吗？'
          })
        ]
      })
    }
  }
])

// 分页配置
const paginationConfig = usePaginationConfig(() => props.pagination)

const handleCheck = (keys: DataTableRowKey[]) => {
  emit('update:checkedKeys', keys)
}
</script>

<template>
  <NDataTable
    :remote="!!pagination"
    :columns="columns"
    :data="overdueItems"
    :loading="loading"
    :pagination="paginationConfig"
    :scroll-x="scrollX"
    :row-key="(row: OverdueInfo) => row.id"
    :checked-row-keys="checkedKeys"
    @update:checked-row-keys="handleCheck"
  />
</template>
