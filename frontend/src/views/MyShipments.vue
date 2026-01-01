<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
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
  useMessage,
  type DataTableColumns
} from 'naive-ui'
import { getMyShipments, createShipment, getCourierCompanies, calculateFreight } from '../api/shipment'
import type { Shipment, CreateShipmentRequest, CalculateFreightRequest } from '../api/types'

const message = useMessage()

const loading = ref(false)
const shipments = ref<Shipment[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

// 创建寄件单弹窗
const showCreateModal = ref(false)
const createLoading = ref(false)
const courierCompanies = ref<string[]>([])
const estimatedFreight = ref<number | null>(null)

// 表单数据
const formValue = ref<CreateShipmentRequest>({
  recipient_name: '',
  recipient_phone: '',
  recipient_address: '',
  courier_company: '',
  weight: 1,
  volume: 0.001
})

// 寄件状态映射
const shipmentStatusMap: Record<string, { text: string; type: any }> = {
  pending: { text: '待发货', type: 'warning' },
  shipped: { text: '已发货', type: 'info' },
  delivered: { text: '已送达', type: 'success' },
  cancelled: { text: '已取消', type: 'error' }
}

// 表格列定义
const columns: DataTableColumns<Shipment> = [
  {
    title: '寄件单号',
    key: 'shipment_number',
    width: 180
  },
  {
    title: '收件人',
    key: 'recipient_name',
    width: 100
  },
  {
    title: '收件电话',
    key: 'recipient_phone',
    width: 120
  },
  {
    title: '收件地址',
    key: 'recipient_address',
    ellipsis: { tooltip: true },
    width: 200
  },
  {
    title: '快递公司',
    key: 'courier_company',
    width: 100
  },
  {
    title: '重量',
    key: 'weight',
    width: 80,
    render: (row) => `${row.weight} kg`
  },
  {
    title: '运费',
    key: 'freight',
    width: 80,
    render: (row) => `¥${row.freight}`
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const status = shipmentStatusMap[row.status] || { text: row.status, type: 'default' }
      return h(NTag, { type: status.type }, { default: () => status.text })
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 160,
    render: (row) => new Date(row.created_at).toLocaleString('zh-CN')
  }
]

// 加载数据
const loadData = async () => {
  loading.value = true
  const res = await getMyShipments({
    page: page.value,
    page_size: pageSize.value
  })

  if (res.code === 0 && res.data) {
    shipments.value = res.data.list
    total.value = res.data.pagination.total
  } else {
    message.error(res.message || '加载失败')
  }
  loading.value = false
}

// 加载快递公司列表
const loadCourierCompanies = async () => {
  try {
    const res = await getCourierCompanies()
    if (res.code === 0 && res.data) {
      courierCompanies.value = res.data
    }
  } catch (error) {
    // 忽略错误
  }
}

// 计算运费
const handleCalculateFreight = async () => {
  if (!formValue.value.courier_company || !formValue.value.weight || !formValue.value.volume) {
    return
  }

  try {
    const params: CalculateFreightRequest = {
      courier_company: formValue.value.courier_company,
      weight: formValue.value.weight,
      volume: formValue.value.volume
    }
    const res = await calculateFreight(params)
    if (res.code === 0 && res.data) {
      estimatedFreight.value = res.data.freight
    }
  } catch (error) {
    // 忽略错误
  }
}

// 打开创建弹窗
const handleCreate = () => {
  formValue.value = {
    recipient_name: '',
    recipient_phone: '',
    recipient_address: '',
    courier_company: '',
    weight: 1,
    volume: 0.001
  }
  estimatedFreight.value = null
  showCreateModal.value = true
  loadCourierCompanies()
}

// 提交创建
const handleSubmit = async () => {
  if (createLoading.value) return

  // 验证
  if (!formValue.value.recipient_name) {
    message.error('请输入收件人姓名')
    return
  }
  if (!formValue.value.recipient_phone) {
    message.error('请输入收件人电话')
    return
  }
  if (!formValue.value.recipient_address) {
    message.error('请输入收件地址')
    return
  }
  if (!formValue.value.courier_company) {
    message.error('请选择快递公司')
    return
  }

  createLoading.value = true
  const res = await createShipment(formValue.value)
  if (res.code === 0) {
    message.success('创建成功')
    showCreateModal.value = false
    loadData()
  } else {
    message.error(res.message || '创建失败')
  }
  createLoading.value = false
}

// 处理分页变化
const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="page-container">
    <NCard title="我的寄件">
      <template #header-extra>
        <NButton type="primary" @click="handleCreate">
          创建寄件单
        </NButton>
      </template>

      <!-- 数据表格 -->
      <NDataTable
        :columns="columns"
        :data="shipments"
        :loading="loading"
        :pagination="{
          page: page,
          pageSize: pageSize,
          itemCount: total,
          onUpdatePage: handlePageChange
        }"
        :scroll-x="1200"
      />
    </NCard>

    <!-- 创建寄件单弹窗 -->
    <NModal
      v-model:show="showCreateModal"
      preset="card"
      title="创建寄件单"
      style="width: 600px;"
      :mask-closable="false"
    >
      <NForm :model="formValue" label-placement="left" label-width="100">
        <NFormItem label="收件人姓名" required>
          <NInput v-model:value="formValue.recipient_name" placeholder="请输入收件人姓名" />
        </NFormItem>

        <NFormItem label="收件人电话" required>
          <NInput v-model:value="formValue.recipient_phone" placeholder="请输入收件人电话" />
        </NFormItem>

        <NFormItem label="收件地址" required>
          <NInput
            v-model:value="formValue.recipient_address"
            type="textarea"
            placeholder="请输入详细地址"
            :rows="3"
          />
        </NFormItem>

        <NFormItem label="快递公司" required>
          <NSelect
            v-model:value="formValue.courier_company"
            :options="courierCompanies.map(c => ({ label: c, value: c }))"
            placeholder="请选择快递公司"
            @update:value="handleCalculateFreight"
          />
        </NFormItem>

        <NFormItem label="重量 (kg)" required>
          <NInputNumber
            v-model:value="formValue.weight"
            :min="0.1"
            :max="50"
            :step="0.1"
            placeholder="请输入重量"
            style="width: 100%;"
            @update:value="handleCalculateFreight"
          />
        </NFormItem>

        <NFormItem label="体积 (m³)" required>
          <NInputNumber
            v-model:value="formValue.volume"
            :min="0.001"
            :max="1"
            :step="0.001"
            placeholder="请输入体积"
            style="width: 100%;"
            @update:value="handleCalculateFreight"
          />
        </NFormItem>

        <NFormItem v-if="estimatedFreight !== null" label="预估运费">
          <strong style="font-size: 20px; color: #18a058;">¥{{ estimatedFreight }}</strong>
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="showCreateModal = false">取消</NButton>
          <NButton type="primary" :loading="createLoading" @click="handleSubmit">
            创建
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
