<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  NCard,
  NButton,
  NSpace,
  NStatistic,
  NGrid,
  NGi,
  useMessage,
  type DataTableRowKey
} from 'naive-ui'
import {
  getOverdueList,
  sendReminder,
  batchReturn,
  checkOverdue,
  type OverdueInfo
} from '../api/overdue'
import { OVERDUE_STATUS_OPTIONS } from '../constants/status'
import OverdueTable from '../components/table/OverdueTable.vue'

const message = useMessage()

// 滞留件列表
const loading = ref(false)
const overdueList = ref<OverdueInfo[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref('')

// 选中的行
const checkedRowKeys = ref<DataTableRowKey[]>([])

// 统计
const statistics = computed(() => {
  const pending = overdueList.value.filter(o => o.status === 'pending').length
  const reminded = overdueList.value.filter(o => o.status === 'reminded').length
  const returned = overdueList.value.filter(o => o.status === 'returned').length
  const totalFee = overdueList.value.reduce((sum, o) => sum + o.overdue_fee, 0)
  return { pending, reminded, returned, totalFee }
})

// 加载滞留件列表
const loadOverdueList = async () => {
  loading.value = true
  const res = await getOverdueList({
    status: statusFilter.value || undefined,
    page: page.value,
    page_size: pageSize.value
  })

  if (res.code === 0 && res.data) {
    overdueList.value = res.data.list
    total.value = res.data.pagination.total
  } else {
    message.error(res.message || '加载失败')
  }
  loading.value = false
}

// 发送提醒
const handleSendReminder = async (row: OverdueInfo) => {
  const res = await sendReminder(row.id)
  if (res.code === 0) {
    message.success('提醒已发送')
    loadOverdueList()
  } else {
    message.error(res.message || '发送失败')
  }
}

// 退回包裹
const handleReturn = async (ids: number[]) => {
  const res = await batchReturn(ids)
  if (res.code === 0 && res.data) {
    message.success(`成功退回 ${res.data.success} 个包裹`)
    checkedRowKeys.value = []
    loadOverdueList()
  } else {
    message.error(res.message || '退回失败')
  }
}

// 批量退回选中的包裹
const handleBatchReturn = () => {
  const ids = checkedRowKeys.value as number[]
  if (ids.length === 0) {
    message.warning('请选择要退回的包裹')
    return
  }
  handleReturn(ids)
}

// 手动检查滞留件
const handleCheckOverdue = async () => {
  const res = await checkOverdue()
  if (res.code === 0 && res.data) {
    message.success(`已标记 ${res.data.marked_count} 个新滞留件`)
    loadOverdueList()
  } else {
    message.error(res.message || '检查失败')
  }
}

// 处理分页变化
const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadOverdueList()
}

// 处理每页条数变化
const handlePageSizeChange = (newPageSize: number) => {
  pageSize.value = newPageSize
  page.value = 1
  loadOverdueList()
}

// 处理选择变化
const handleCheck = (rowKeys: DataTableRowKey[]) => {
  checkedRowKeys.value = rowKeys
}

// 处理筛选
const handleStatusFilter = (value: string) => {
  statusFilter.value = value
  page.value = 1
  loadOverdueList()
}

onMounted(() => {
  loadOverdueList()
})
</script>

<template>
  <div class="page-container">
    <!-- 统计卡片 -->
    <NGrid :cols="4" :x-gap="16" style="margin-bottom: 16px;">
      <NGi>
        <NCard>
          <NStatistic label="待处理" :value="statistics.pending">
            <template #suffix>件</template>
          </NStatistic>
        </NCard>
      </NGi>
      <NGi>
        <NCard>
          <NStatistic label="已提醒" :value="statistics.reminded">
            <template #suffix>件</template>
          </NStatistic>
        </NCard>
      </NGi>
      <NGi>
        <NCard>
          <NStatistic label="已退回" :value="statistics.returned">
            <template #suffix>件</template>
          </NStatistic>
        </NCard>
      </NGi>
      <NGi>
        <NCard>
          <NStatistic label="累计滞留费">
            <template #prefix>¥</template>
            <template #default>{{ statistics.totalFee.toFixed(2) }}</template>
          </NStatistic>
        </NCard>
      </NGi>
    </NGrid>

    <NCard title="滞留件管理">
      <template #header-extra>
        <NSpace>
          <NButton @click="handleCheckOverdue">检查滞留件</NButton>
          <NButton type="error" :disabled="checkedRowKeys.length === 0" @click="handleBatchReturn">
            批量退回 ({{ checkedRowKeys.length }})
          </NButton>
        </NSpace>
      </template>

      <!-- 筛选工具栏 -->
      <NSpace style="margin-bottom: 16px;">
        <NButton v-for="option in OVERDUE_STATUS_OPTIONS" :key="option.value ?? 'all'"
          :type="statusFilter === option.value ? 'primary' : 'default'" @click="handleStatusFilter(option.value)">
          {{ option.label }}
        </NButton>
      </NSpace>

      <!-- 数据表格 -->
      <OverdueTable :overdue-items="overdueList" :loading="loading" :checked-keys="checkedRowKeys" :pagination="{
        page,
        pageSize,
        total,
        onChange: handlePageChange,
        onSizeChange: handlePageSizeChange
      }" @remind="handleSendReminder" @return="handleReturn" @update:checked-keys="handleCheck" />
    </NCard>
  </div>
</template>

<style scoped>
.page-container {
  padding: 24px;
}
</style>
