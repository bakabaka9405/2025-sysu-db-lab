<script setup lang="ts">
import { computed, h } from 'vue'
import { NDataTable, NTag, NButton, NSpace, NProgress } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import type { Shelf } from '../../api/shelf'
import { SHELF_STATUS_MAP } from '../../constants/status'
import { getUsageRateColor } from '../../constants/styles'
import { usePaginationConfig, type PaginationProps } from '../../utils/pagination'

export interface ShelfTableProps {
  shelves: Shelf[]
  loading?: boolean
  pagination?: PaginationProps | false
  scrollX?: number
}

const props = withDefaults(defineProps<ShelfTableProps>(), {
  loading: false,
  scrollX: 1000
})

const emit = defineEmits<{
  edit: [shelf: Shelf]
  delete: [shelf: Shelf]
}>()

// 货架类型映射
const typeMap: Record<string, string> = {
  small: '小件',
  medium: '中件',
  large: '大件'
}

// 表格列定义
const columns = computed<DataTableColumns<Shelf>>(() => [
  { title: 'ID', key: 'id', width: 60 },
  { title: '货架编码', key: 'shelf_code', width: 100 },
  { title: '区域', key: 'area', width: 60 },
  { title: '层', key: 'floor', width: 60 },
  { title: '列', key: 'column', width: 60 },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: (row) => typeMap[row.type] || row.type
  },
  {
    title: '容量',
    key: 'capacity',
    width: 80,
    render: (row) => `${row.current_count ?? 0}/${row.capacity}`
  },
  {
    title: '使用率',
    key: 'utilization',
    width: 120,
    render: (row) => {
      const rate = Math.round(((row.current_count ?? 0) / row.capacity) * 100)
      return h(NProgress, {
        type: 'line',
        percentage: rate,
        color: getUsageRateColor(rate),
        showIndicator: true
      })
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) => {
      const status = SHELF_STATUS_MAP[row.status ?? 'active'] || { text: row.status ?? '未知', type: 'default' }
      return h(NTag, { type: status.type }, { default: () => status.text })
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) => h(NSpace, {}, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => emit('edit', row) }, { default: () => '编辑' }),
        h(NButton, {
          size: 'small',
          type: 'error',
          disabled: (row.current_count ?? 0) > 0,
          onClick: () => emit('delete', row)
        }, { default: () => '删除' })
      ]
    })
  }
])

// 分页配置
const paginationConfig = usePaginationConfig(() => props.pagination)
</script>

<template>
  <NDataTable
    :remote="!!pagination"
    :columns="columns"
    :data="shelves"
    :loading="loading"
    :pagination="paginationConfig"
    :scroll-x="scrollX"
  />
</template>
