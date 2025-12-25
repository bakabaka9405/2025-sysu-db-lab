<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NTag,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NInputNumber,
  NStatistic,
  NGrid,
  NGi,
  NProgress,
  useMessage,
  type DataTableColumns
} from 'naive-ui'
import {
  getShelfList,
  createShelf,
  updateShelf,
  deleteShelf,
  getShelfStatistics,
  type Shelf,
  type CreateShelfRequest,
  type UpdateShelfRequest,
  type ShelfStatistics
} from '../api/shelf'

const message = useMessage()

// 货架列表
const loading = ref(false)
const shelves = ref<Shelf[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const areaFilter = ref<string | null>(null)
const typeFilter = ref<string | null>(null)
const statusFilter = ref<string | null>(null)

// 统计数据
const statistics = ref<ShelfStatistics | null>(null)

// 弹窗状态
const showModal = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const modalLoading = ref(false)
const editingId = ref<number | null>(null)

// 表单数据
const formData = ref<CreateShelfRequest>({
  shelf_code: '',
  area: '',
  floor: 1,
  column: 1,
  type: 'small',
  capacity: 20
})
const editFormData = ref<UpdateShelfRequest>({})

// 筛选选项
const areaOptions = [
  { label: '全部', value: null },
  { label: 'A区', value: 'A' },
  { label: 'B区', value: 'B' },
  { label: 'C区', value: 'C' }
]

const typeOptions = [
  { label: '全部', value: null },
  { label: '小件', value: 'small' },
  { label: '中件', value: 'medium' },
  { label: '大件', value: 'large' }
]

const statusOptions = [
  { label: '全部', value: null },
  { label: '启用', value: 'active' },
  { label: '停用', value: 'inactive' }
]

// 类型映射
const typeMap: Record<string, string> = {
  small: '小件',
  medium: '中件',
  large: '大件'
}

const statusMap: Record<string, { text: string; type: any }> = {
  active: { text: '启用', type: 'success' },
  inactive: { text: '停用', type: 'warning' }
}

// 负载等级颜色
const getLoadColor = (rate: number) => {
  if (rate >= 90) return '#f5222d'
  if (rate >= 70) return '#fa8c16'
  if (rate >= 30) return '#52c41a'
  return '#1890ff'
}

// 表格列定义
const columns: DataTableColumns<Shelf> = [
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
    render: (row) => `${row.current_count}/${row.capacity}`
  },
  {
    title: '使用率',
    key: 'utilization',
    width: 120,
    render: (row) => {
      const rate = Math.round((row.current_count / row.capacity) * 100)
      return h(NProgress, {
        type: 'line',
        percentage: rate,
        color: getLoadColor(rate),
        showIndicator: true
      })
    }
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
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
          h(NButton, {
            size: 'small',
            type: 'error',
            disabled: row.current_count > 0,
            onClick: () => handleDelete(row)
          }, { default: () => '删除' })
        ]
      })
    }
  }
]

// 加载货架列表
const loadShelves = async () => {
  loading.value = true
  try {
    const res = await getShelfList({
      area: areaFilter.value || undefined,
      type: typeFilter.value || undefined,
      status: statusFilter.value || undefined,
      page: page.value,
      page_size: pageSize.value
    })

    if (res.code === 0 && res.data) {
      shelves.value = res.data.list
      total.value = res.data.pagination.total
    } else {
      message.error(res.message || '加载失败')
    }
  } catch (error: any) {
    message.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

// 加载统计数据
const loadStatistics = async () => {
  try {
    const res = await getShelfStatistics()
    if (res.code === 0 && res.data) {
      statistics.value = res.data
    }
  } catch (error) {
    console.error('Failed to load statistics', error)
  }
}

// 打开创建弹窗
const handleCreate = () => {
  modalMode.value = 'create'
  formData.value = {
    shelf_code: '',
    area: '',
    floor: 1,
    column: 1,
    type: 'small',
    capacity: 20
  }
  showModal.value = true
}

// 打开编辑弹窗
const handleEdit = (shelf: Shelf) => {
  modalMode.value = 'edit'
  editingId.value = shelf.id
  editFormData.value = {
    area: shelf.area,
    floor: shelf.floor,
    column: shelf.column,
    type: shelf.type,
    capacity: shelf.capacity,
    status: shelf.status
  }
  showModal.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (modalLoading.value) return

  modalLoading.value = true
  try {
    if (modalMode.value === 'create') {
      if (!formData.value.shelf_code || !formData.value.area) {
        message.error('请填写完整信息')
        return
      }
      const res = await createShelf(formData.value)
      if (res.code === 0) {
        message.success('创建成功')
        showModal.value = false
        loadShelves()
        loadStatistics()
      } else {
        message.error(res.message || '创建失败')
      }
    } else {
      if (!editingId.value) return
      const res = await updateShelf(editingId.value, editFormData.value)
      if (res.code === 0) {
        message.success('更新成功')
        showModal.value = false
        loadShelves()
        loadStatistics()
      } else {
        message.error(res.message || '更新失败')
      }
    }
  } catch (error: any) {
    message.error(error.message || '操作失败')
  } finally {
    modalLoading.value = false
  }
}

// 删除货架
const handleDelete = async (shelf: Shelf) => {
  try {
    const res = await deleteShelf(shelf.id)
    if (res.code === 0) {
      message.success('删除成功')
      loadShelves()
      loadStatistics()
    } else {
      message.error(res.message || '删除失败')
    }
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}

// 处理分页变化
const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadShelves()
}

// 处理筛选变化
const handleFilterChange = () => {
  page.value = 1
  loadShelves()
}

onMounted(() => {
  loadShelves()
  loadStatistics()
})
</script>

<template>
  <div class="page-container">
    <!-- 统计卡片 -->
    <NGrid :cols="4" :x-gap="16" style="margin-bottom: 16px;">
      <NGi>
        <NCard>
          <NStatistic label="总货架数" :value="statistics?.total_shelves || 0" />
        </NCard>
      </NGi>
      <NGi>
        <NCard>
          <NStatistic label="启用货架" :value="statistics?.active_shelves || 0" />
        </NCard>
      </NGi>
      <NGi>
        <NCard>
          <NStatistic label="总容量" :value="statistics?.total_capacity || 0" />
        </NCard>
      </NGi>
      <NGi>
        <NCard>
          <NStatistic label="总使用率">
            <template #default>
              <span :style="{ color: getLoadColor(statistics?.utilization_rate || 0) }">
                {{ (statistics?.utilization_rate || 0).toFixed(1) }}%
              </span>
            </template>
          </NStatistic>
        </NCard>
      </NGi>
    </NGrid>

    <NCard title="货架管理">
      <template #header-extra>
        <NButton type="primary" @click="handleCreate">新增货架</NButton>
      </template>

      <!-- 筛选工具栏 -->
      <NSpace style="margin-bottom: 16px;">
        <NSelect
          v-model:value="areaFilter"
          :options="areaOptions"
          placeholder="区域"
          style="width: 100px;"
          @update:value="handleFilterChange"
        />
        <NSelect
          v-model:value="typeFilter"
          :options="typeOptions"
          placeholder="类型"
          style="width: 100px;"
          @update:value="handleFilterChange"
        />
        <NSelect
          v-model:value="statusFilter"
          :options="statusOptions"
          placeholder="状态"
          style="width: 100px;"
          @update:value="handleFilterChange"
        />
      </NSpace>

      <!-- 数据表格 -->
      <NDataTable
        :columns="columns"
        :data="shelves"
        :loading="loading"
        :pagination="{
          page: page,
          pageSize: pageSize,
          itemCount: total,
          onUpdatePage: handlePageChange
        }"
        :scroll-x="1000"
      />
    </NCard>

    <!-- 创建/编辑弹窗 -->
    <NModal
      v-model:show="showModal"
      preset="card"
      :title="modalMode === 'create' ? '新增货架' : '编辑货架'"
      style="width: 500px;"
      :mask-closable="false"
    >
      <NForm
        v-if="modalMode === 'create'"
        :model="formData"
        label-placement="left"
        label-width="80"
      >
        <NFormItem label="货架编码" required>
          <NInput v-model:value="formData.shelf_code" placeholder="如 A-1-1" />
        </NFormItem>
        <NFormItem label="区域" required>
          <NInput v-model:value="formData.area" placeholder="如 A、B、C" />
        </NFormItem>
        <NFormItem label="层" required>
          <NInputNumber v-model:value="formData.floor" :min="1" :max="10" style="width: 100%;" />
        </NFormItem>
        <NFormItem label="列" required>
          <NInputNumber v-model:value="formData.column" :min="1" :max="20" style="width: 100%;" />
        </NFormItem>
        <NFormItem label="类型" required>
          <NSelect
            v-model:value="formData.type"
            :options="[
              { label: '小件', value: 'small' },
              { label: '中件', value: 'medium' },
              { label: '大件', value: 'large' }
            ]"
          />
        </NFormItem>
        <NFormItem label="容量" required>
          <NInputNumber v-model:value="formData.capacity" :min="1" :max="100" style="width: 100%;" />
        </NFormItem>
      </NForm>

      <NForm
        v-else
        :model="editFormData"
        label-placement="left"
        label-width="80"
      >
        <NFormItem label="区域">
          <NInput v-model:value="editFormData.area" />
        </NFormItem>
        <NFormItem label="层">
          <NInputNumber v-model:value="editFormData.floor" :min="1" :max="10" style="width: 100%;" />
        </NFormItem>
        <NFormItem label="列">
          <NInputNumber v-model:value="editFormData.column" :min="1" :max="20" style="width: 100%;" />
        </NFormItem>
        <NFormItem label="类型">
          <NSelect
            v-model:value="editFormData.type"
            :options="[
              { label: '小件', value: 'small' },
              { label: '中件', value: 'medium' },
              { label: '大件', value: 'large' }
            ]"
          />
        </NFormItem>
        <NFormItem label="容量">
          <NInputNumber v-model:value="editFormData.capacity" :min="1" :max="100" style="width: 100%;" />
        </NFormItem>
        <NFormItem label="状态">
          <NSelect
            v-model:value="editFormData.status"
            :options="[
              { label: '启用', value: 'active' },
              { label: '停用', value: 'inactive' }
            ]"
          />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="showModal = false">取消</NButton>
          <NButton type="primary" :loading="modalLoading" @click="handleSubmit">
            确认
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.page-container {
  padding: 24px;
  min-height: 100vh;
  background-color: #f5f5f5;
}
</style>
