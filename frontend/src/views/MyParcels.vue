<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NSpace, NSelect, useMessage } from 'naive-ui'
import { getMyParcels } from '../api/parcel'
import type { Parcel } from '../api/types'
import ParcelTable from '../components/table/ParcelTable.vue'

const message = useMessage()
const loading = ref(false)
const parcels = ref<Parcel[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

// 筛选条件
const statusFilter = ref<string | null>(null)

// 包裹状态选项
const statusOptions = [
  { label: '全部', value: '' },
  { label: '待取件', value: 'ready' },
  { label: '已取件', value: 'picked_up' },
  { label: '滞留', value: 'overdue' }
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

// 处理每页条数变化
const handlePageSizeChange = (newPageSize: number) => {
  pageSize.value = newPageSize
  page.value = 1
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
        <NSelect
          v-model:value="statusFilter"
          :options="statusOptions"
          placeholder="状态筛选"
          style="width: 150px;"
          @update:value="handleStatusChange"
        />
      </NSpace>

      <!-- 数据表格 -->
      <ParcelTable
        :parcels="parcels"
        :loading="loading"
        show-notes
        :pagination="{
          page,
          pageSize,
          total,
          onChange: handlePageChange,
          onSizeChange: handlePageSizeChange
        }"
      />
    </NCard>
  </div>
</template>

<style scoped>
.page-container {
  padding: 24px;
}
</style>
