<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import {
  NCard,
  NTabs,
  NTabPane,
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
  useMessage,
  type DataTableColumns
} from 'naive-ui'
import {
  getParcelList,
  receiveParcel,
  shelveParcel,
  pickupParcel
} from '../api/parcel'
import { COURIER_COMPANIES } from '../api/mock-data'
import type { Parcel, ReceiveParcelRequest, PickupParcelRequest, ParcelSize } from '../api/types'

const message = useMessage()

// 包裹列表
const loading = ref(false)
const parcels = ref<Parcel[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref<string>('')

// 入库弹窗
const showReceiveModal = ref(false)
const receiveLoading = ref(false)
const receiveForm = ref<ReceiveParcelRequest>({
  tracking_number: '',
  recipient_name: '',
  recipient_phone: '',
  recipient_id_card: '',
  courier_company: '',
  size: 'small',
  weight: 1,
  notes: ''
})

// 取件弹窗
const showPickupModal = ref(false)
const pickupLoading = ref(false)
const pickupForm = ref<PickupParcelRequest>({
  pickup_code: '',
  recipient_phone: ''
})

// 状态选项
const statusOptions = [
  { label: '全部', value: '' },
  { label: '待上架', value: 'received' },
  { label: '待取件', value: 'ready' },
  { label: '已取件', value: 'picked_up' },
  { label: '滞留', value: 'overdue' },
  { label: '已退回', value: 'returned' }
]

// 包裹尺寸选项
const sizeOptions = [
  { label: '小', value: 'small' },
  { label: '中', value: 'medium' },
  { label: '大', value: 'large' }
]

// 快递公司选项
const courierOptions = COURIER_COMPANIES.map(c => ({ label: c, value: c }))

// 包裹状态映射
const parcelStatusMap: Record<string, { text: string; type: any }> = {
  received: { text: '待上架', type: 'info' },
  ready: { text: '待取件', type: 'success' },
  picked_up: { text: '已取件', type: 'default' },
  overdue: { text: '滞留', type: 'warning' },
  returned: { text: '已退回', type: 'error' }
}

// 表格列定义
const columns: DataTableColumns<Parcel> = [
  {
    title: 'ID',
    key: 'id',
    width: 60
  },
  {
    title: '取件码',
    key: 'pickup_code',
    width: 100,
    render: (row) => {
      return h('strong', { style: 'color: #18a058; font-size: 14px;' }, row.pickup_code)
    }
  },
  {
    title: '快递单号',
    key: 'tracking_number',
    width: 140
  },
  {
    title: '收件人',
    key: 'recipient_name',
    width: 80
  },
  {
    title: '手机号',
    key: 'recipient_phone',
    width: 120
  },
  {
    title: '快递公司',
    key: 'courier_company',
    width: 100
  },
  {
    title: '货架',
    key: 'shelf',
    width: 100,
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
    title: '操作',
    key: 'actions',
    width: 120,
    render: (row) => {
      if (row.status === 'received') {
        return h(NButton, {
          size: 'small',
          type: 'primary',
          onClick: () => handleShelve(row)
        }, { default: () => '上架' })
      }
      return null
    }
  }
]

// 加载包裹列表
const loadParcels = async () => {
  loading.value = true
  const res = await getParcelList({
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

// 打开入库弹窗
const handleOpenReceive = () => {
  receiveForm.value = {
    tracking_number: '',
    recipient_name: '',
    recipient_phone: '',
    recipient_id_card: '',
    courier_company: '',
    size: 'small',
    weight: 1,
    notes: ''
  }
  showReceiveModal.value = true
}

// 包裹上架
const handleShelve = async (parcel: Parcel) => {
  const res = await shelveParcel(parcel.id)
  if (res.code === 0) {
    message.success(`上架成功！货架位置：${res.data?.shelf?.shelf_code}`)
    loadParcels()
  } else {
    message.error(res.message || '上架失败')
  }
}

// 提交入库
const handleSubmitReceive = async () => {
  if (receiveLoading.value) return

  // 验证
  if (!receiveForm.value.tracking_number) {
    message.error('请输入快递单号')
    return
  }
  if (!receiveForm.value.recipient_name) {
    message.error('请输入收件人姓名')
    return
  }
  if (!receiveForm.value.recipient_phone) {
    message.error('请输入收件人电话')
    return
  }
  if (!receiveForm.value.courier_company) {
    message.error('请选择快递公司')
    return
  }

  receiveLoading.value = true
  const res = await receiveParcel(receiveForm.value)
  if (res.code === 0) {
    message.success(`入库成功！取件码：${res.data?.pickup_code}`)
    showReceiveModal.value = false
    loadParcels()
  } else {
    message.error(res.message || '入库失败')
  }
  receiveLoading.value = false
}

// 打开取件弹窗
const handleOpenPickup = () => {
  pickupForm.value = {
    pickup_code: '',
    recipient_phone: ''
  }
  showPickupModal.value = true
}

// 提交取件
const handleSubmitPickup = async () => {
  if (pickupLoading.value) return

  // 验证
  if (!pickupForm.value.pickup_code) {
    message.error('请输入取件码')
    return
  }
  if (!pickupForm.value.recipient_phone) {
    message.error('请输入手机号')
    return
  }

  pickupLoading.value = true
  const res = await pickupParcel(pickupForm.value)
  if (res.code === 0) {
    message.success('取件成功')
    showPickupModal.value = false
    loadParcels()
  } else {
    message.error(res.message || '取件失败')
  }
  pickupLoading.value = false
}

// 处理分页变化
const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadParcels()
}

// 处理每页条数变化
const handlePageSizeChange = (newPageSize: number) => {
  pageSize.value = newPageSize
  page.value = 1
  loadParcels()
}

// 处理状态筛选
const handleStatusChange = () => {
  page.value = 1
  loadParcels()
}

onMounted(() => {
  loadParcels()
})
</script>

<template>
  <div class="page-container">
    <NCard title="包裹管理">
      <template #header-extra>
        <NSpace>
          <NButton type="primary" @click="handleOpenReceive">
            包裹入库
          </NButton>
          <NButton type="success" @click="handleOpenPickup">
            包裹取件
          </NButton>
        </NSpace>
      </template>

      <!-- 筛选工具栏 -->
      <NSpace style="margin-bottom: 16px;">
        <NSelect v-model:value="statusFilter" :options="statusOptions" placeholder="状态筛选" style="width: 150px;"
          @update:value="handleStatusChange" />
      </NSpace>

      <!-- 数据表格 -->
      <NDataTable
        :remote="true"
        :columns="columns"
        :data="parcels"
        :loading="loading"
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
        :scroll-x="1200"
      />
    </NCard>

    <!-- 入库弹窗 -->
    <NModal v-model:show="showReceiveModal" preset="card" title="包裹入库" style="width: 600px;" :mask-closable="false">
      <NForm :model="receiveForm" label-placement="left" label-width="110">
        <NFormItem label="快递单号" required>
          <NInput v-model:value="receiveForm.tracking_number" placeholder="请输入快递单号" />
        </NFormItem>

        <NFormItem label="收件人姓名" required>
          <NInput v-model:value="receiveForm.recipient_name" placeholder="请输入收件人姓名" />
        </NFormItem>

        <NFormItem label="收件人电话" required>
          <NInput v-model:value="receiveForm.recipient_phone" placeholder="请输入收件人电话" />
        </NFormItem>

        <NFormItem label="身份证号">
          <NInput v-model:value="receiveForm.recipient_id_card" placeholder="选填" />
        </NFormItem>

        <NFormItem label="快递公司" required>
          <NSelect v-model:value="receiveForm.courier_company" :options="courierOptions" placeholder="请选择快递公司" />
        </NFormItem>

        <NFormItem label="包裹尺寸" required>
          <NSelect v-model:value="receiveForm.size" :options="sizeOptions" placeholder="请选择包裹尺寸" />
        </NFormItem>

        <NFormItem label="重量 (kg)" required>
          <NInputNumber v-model:value="receiveForm.weight" :min="0.1" :max="50" :step="0.1" placeholder="请输入重量"
            style="width: 100%;" />
        </NFormItem>

        <NFormItem label="备注">
          <NInput v-model:value="receiveForm.notes" type="textarea" placeholder="选填" :rows="3" />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="showReceiveModal = false">取消</NButton>
          <NButton type="primary" :loading="receiveLoading" @click="handleSubmitReceive">
            确认入库
          </NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- 取件弹窗 -->
    <NModal v-model:show="showPickupModal" preset="card" title="包裹取件" style="width: 500px;" :mask-closable="false">
      <NForm :model="pickupForm" label-placement="left" label-width="80">
        <NFormItem label="取件码" required>
          <NInput v-model:value="pickupForm.pickup_code" placeholder="请输入取件码" size="large" />
        </NFormItem>

        <NFormItem label="手机号" required>
          <NInput v-model:value="pickupForm.recipient_phone" placeholder="请输入收件人手机号进行验证" size="large" />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="showPickupModal = false">取消</NButton>
          <NButton type="primary" :loading="pickupLoading" @click="handleSubmitPickup">
            确认取件
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
