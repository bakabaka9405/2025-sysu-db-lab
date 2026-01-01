<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NTag,
  NInput,
  NSelect,
  useMessage,
  type DataTableColumns
} from 'naive-ui'
import { getMyParcels } from '../api/parcel'
import type { Parcel } from '../api/types'

const message = useMessage()
const loading = ref(false)
const parcels = ref<Parcel[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

// 筛选条件
const statusFilter = ref<string | null>(null)
const searchKeyword = ref('')

// 包裹状态选项
const statusOptions = [
  { label: '全部', value: '' },
  { label: '待取件', value: 'ready_for_pickup' },
  { label: '已取件', value: 'picked_up' },
  { label: '滞留', value: 'overdue' }
]

// 包裹状态映射
const parcelStatusMap: Record<string, { text: string; type: any }> = {
  ready_for_pickup: { text: '待取件', type: 'success' },
  picked_up: { text: '已取件', type: 'default' },
  overdue: { text: '滞留', type: 'warning' },
  returned: { text: '已退回', type: 'error' }
}

// 表格列定义
const columns: DataTableColumns<Parcel> = [
  {
    title: '取件码',
    key: 'pickup_code',
    width: 100,
    render: (row) => {
      return h('strong', { style: 'color: #18a058; font-size: 16px;' }, row.pickup_code)
    }
  },
  {
    title: '快递单号',
    key: 'tracking_number',
    width: 150
  },
  {
    title: '快递公司',
    key: 'courier_company',
    width: 100
  },
  {
    title: '货架位置',
    key: 'shelf',
    width: 120,
    render: (row) => row.shelf?.shelf_code || '-'
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const status = parcelStatusMap[row.status] || { text: row.status, type: 'default' }
      return h(NTag, { type: status.type }, { default: () => status.text })
    }
  },
  {
    title: '入库时间',
    key: 'received_at',
    width: 160,
    render: (row) => new Date(row.received_at).toLocaleString('zh-CN')
  },
  {
    title: '备注',
    key: 'notes',
    ellipsis: { tooltip: true }
  }
]

// 加载数据
const loadData = async () => {
  loading.value = true
  const res = await getMyParcels({
    status: statusFilter.value || undefined,
    page: page.value,
    page_size: pageSize.value
  })

  if (res.code === 0 && res.data) {
    parcels.value = res.data.list
    total.value = res.data.pagination.total
  } else {
    message.error(res.message || '加载失败')
  }
  loading.value = false
}

// 处理分页变化
const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadData()
}

// 处理状态筛选
const handleStatusChange = () => {
  page.value = 1
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="page-container">
    <NCard title="我的包裹">
      <!-- 筛选工具栏 -->
      <NSpace style="margin-bottom: 16px;">
        <NSelect v-model:value="statusFilter" :options="statusOptions" placeholder="状态筛选" style="width: 150px;"
          @update:value="handleStatusChange" />
      </NSpace>

      <!-- 数据表格 -->
      <NDataTable :columns="columns" :data="parcels" :loading="loading" :pagination="{
        page: page,
        pageSize: pageSize,
        itemCount: total,
        onUpdatePage: handlePageChange
      }" :scroll-x="1000" />
    </NCard>
  </div>
</template>

<style scoped>
.page-container {
  padding: 24px;
  min-height: 100vh;
  background-color: #f5f5f5;
}
</style>
