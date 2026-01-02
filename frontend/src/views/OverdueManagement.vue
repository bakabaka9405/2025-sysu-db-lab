<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NTag,
  NStatistic,
  NGrid,
  NGi,
  NPopconfirm,
  useMessage,
  type DataTableColumns,
  type DataTableRowKey
} from 'naive-ui'
import {
  getOverdueList,
  sendReminder,
  batchReturn,
  checkOverdue,
  type OverdueInfo
} from '../api/overdue'

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

// 状态选项
const statusOptions = [
  { label: '全部', value: '' },
  { label: '待处理', value: 'pending' },
  { label: '已提醒', value: 'reminded' },
  { label: '已退回', value: 'returned' }
]

// 状态映射
const statusMap: Record<string, { text: string; type: any }> = {
  pending: { text: '待处理', type: 'warning' },
  reminded: { text: '已提醒', type: 'info' },
  returned: { text: '已退回', type: 'default' }
}

// 获取滞留天数颜色
const getOverdueColor = (days: number) => {
  if (days >= 7) return '#f5222d'
  if (days >= 3) return '#fa8c16'
  return '#52c41a'
}

// 表格列定义
const columns: DataTableColumns<OverdueInfo> = [
  { type: 'selection' },
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '取件码',
    key: 'pickup_code',
    width: 100,
    render: (row) => h('strong', { style: 'color: #18a058;' }, row.pickup_code)
  },
  { title: '快递单号', key: 'tracking_number', width: 140 },
  { title: '收件人', key: 'recipient_name', width: 80 },
  { title: '手机号', key: 'recipient_phone', width: 120 },
  { title: '货架', key: 'shelf_code', width: 80, render: (row) => row.shelf_code || '-' },
  {
    title: '滞留天数',
    key: 'overdue_days',
    width: 100,
    render: (row) => {
      return h('span', {
        style: { color: getOverdueColor(row.overdue_days), fontWeight: 'bold' }
      }, `${row.overdue_days} 天`)
    }
  },
  {
    title: '提醒次数',
    key: 'reminder_count',
    width: 80
  },
  {
    title: '滞留费',
    key: 'overdue_fee',
    width: 80,
    render: (row) => `¥${row.overdue_fee.toFixed(2)}`
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) => {
      const status = statusMap[row.status] || { text: row.status, type: 'default' }
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
            onClick: () => handleSendReminder(row)
          }, { default: () => '发送提醒' }),
          h(NPopconfirm, {
            onPositiveClick: () => handleReturn([row.id])
          }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '退回' }),
            default: () => '确定退回该包裹吗？'
          })
        ]
      })
    }
  }
]

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
          <NButton
            type="error"
            :disabled="checkedRowKeys.length === 0"
            @click="handleBatchReturn"
          >
            批量退回 ({{ checkedRowKeys.length }})
          </NButton>
        </NSpace>
      </template>

      <!-- 筛选工具栏 -->
      <NSpace style="margin-bottom: 16px;">
        <NButton
          v-for="option in statusOptions"
          :key="option.value ?? 'all'"
          :type="statusFilter === option.value ? 'primary' : 'default'"
          @click="statusFilter = option.value; page = 1; loadOverdueList()"
        >
          {{ option.label }}
        </NButton>
      </NSpace>

      <!-- 数据表格 -->
      <NDataTable
        :remote="true"
        :columns="columns"
        :data="overdueList"
        :loading="loading"
        :row-key="(row: OverdueInfo) => row.id"
        :checked-row-keys="checkedRowKeys"
        @update:checked-row-keys="handleCheck"
        :pagination="{
          page: page,
          pageSize: pageSize,
          itemCount: total,
          showSizePicker: true,
          pageSizes: [10, 20, 50],
          onUpdatePage: handlePageChange,
          onUpdatePageSize: handlePageSizeChange,
          prefix: () => `共 ${total} 条`
        }"
        :scroll-x="1300"
      />
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
