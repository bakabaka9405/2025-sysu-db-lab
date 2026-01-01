<script setup lang="ts">
import { ref } from 'vue'
import {
  NCard,
  NTabs,
  NTabPane,
  NForm,
  NFormItem,
  NInput,
  NButton,
  NResult,
  NDescriptions,
  NDescriptionsItem,
  NTag,
  NModal,
  NSpace,
  useMessage
} from 'naive-ui'
import { getParcelByPickupCode, getParcelByTrackingNumber, pickupParcel } from '../api/parcel'
import type { Parcel } from '../api/types'

const message = useMessage()
const searchType = ref<'pickup_code' | 'tracking_number'>('pickup_code')
const pickupCode = ref('')
const trackingNumber = ref('')
const loading = ref(false)
const searchResult = ref<Parcel | null>(null)
const searchError = ref('')

// 取件弹窗
const showPickupModal = ref(false)
const pickupLoading = ref(false)
const verifyPhone = ref('')

// 包裹状态映射
const parcelStatusMap: Record<string, { text: string; type: any }> = {
  received: { text: '待上架', type: 'info' },
  ready: { text: '待取件', type: 'success' },
  picked_up: { text: '已取件', type: 'default' },
  overdue: { text: '滞留', type: 'warning' },
  returned: { text: '已退回', type: 'error' }
}

// 是否可以取件
const canPickup = () => {
  return searchResult.value &&
    (searchResult.value.status === 'ready' || searchResult.value.status === 'overdue')
}

// 打开取件弹窗
const handleOpenPickup = () => {
  verifyPhone.value = ''
  showPickupModal.value = true
}

// 提交取件
const handleSubmitPickup = async () => {
  if (pickupLoading.value || !searchResult.value) return

  if (!verifyPhone.value) {
    message.error('请输入手机号进行验证')
    return
  }

  pickupLoading.value = true
  const res = await pickupParcel({
    pickup_code: searchResult.value.pickup_code,
    recipient_phone: verifyPhone.value
  })

  if (res.code === 0) {
    message.success('取件成功！')
    showPickupModal.value = false
    // 更新查询结果
    searchResult.value = res.data || null
  } else {
    message.error(res.message || '取件失败')
  }
  pickupLoading.value = false
}

// 查询包裹
const handleSearch = async () => {
  if (loading.value) return

  searchResult.value = null
  searchError.value = ''

  // 验证输入
  if (searchType.value === 'pickup_code' && !pickupCode.value) {
    message.error('请输入取件码')
    return
  }
  if (searchType.value === 'tracking_number' && !trackingNumber.value) {
    message.error('请输入快递单号')
    return
  }

  loading.value = true
  try {
    let res
    if (searchType.value === 'pickup_code') {
      res = await getParcelByPickupCode(pickupCode.value)
    } else {
      res = await getParcelByTrackingNumber(trackingNumber.value)
    }

    if (res.code === 0 && res.data) {
      searchResult.value = res.data
      message.success('查询成功')
    } else {
      searchError.value = res.message || '包裹不存在'
      message.error(searchError.value)
    }
  } catch (error: any) {
    searchError.value = error.message || '查询失败'
    message.error(searchError.value)
  } finally {
    loading.value = false
  }
}

// 重置搜索
const handleReset = () => {
  pickupCode.value = ''
  trackingNumber.value = ''
  searchResult.value = null
  searchError.value = ''
}

// 格式化日期时间
const formatDateTime = (dateStr?: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<template>
  <div class="search-container">
    <NCard class="search-card" title="包裹查询">
      <NTabs v-model:value="searchType" type="segment" @update:value="handleReset">
        <NTabPane name="pickup_code" tab="取件码查询">
          <NForm label-placement="top">
            <NFormItem label="取件码">
              <NInput v-model:value="pickupCode" placeholder="请输入取件码，例如：A12345" size="large"
                @keyup.enter="handleSearch" />
            </NFormItem>
            <NFormItem>
              <NButton type="primary" size="large" block :loading="loading" @click="handleSearch">
                查询
              </NButton>
            </NFormItem>
          </NForm>
        </NTabPane>

        <NTabPane name="tracking_number" tab="快递单号查询">
          <NForm label-placement="top">
            <NFormItem label="快递单号">
              <NInput v-model:value="trackingNumber" placeholder="请输入快递单号，例如：SF1234567890" size="large"
                @keyup.enter="handleSearch" />
            </NFormItem>
            <NFormItem>
              <NButton type="primary" size="large" block :loading="loading" @click="handleSearch">
                查询
              </NButton>
            </NFormItem>
          </NForm>
        </NTabPane>
      </NTabs>

      <!-- 查询结果 -->
      <div v-if="searchResult" class="result-section">
        <NDescriptions label-placement="left" :column="1" bordered>
          <NDescriptionsItem label="取件码">
            <strong style="font-size: 20px; color: #18a058;">{{ searchResult.pickup_code }}</strong>
          </NDescriptionsItem>
          <NDescriptionsItem label="快递单号">
            {{ searchResult.tracking_number }}
          </NDescriptionsItem>
          <NDescriptionsItem label="收件人">
            {{ searchResult.recipient_name }}
          </NDescriptionsItem>
          <NDescriptionsItem label="手机号">
            {{ searchResult.recipient_phone }}
          </NDescriptionsItem>
          <NDescriptionsItem label="快递公司">
            {{ searchResult.courier_company }}
          </NDescriptionsItem>
          <NDescriptionsItem label="包裹尺寸">
            {{ searchResult.size }}
          </NDescriptionsItem>
          <NDescriptionsItem label="重量">
            {{ searchResult.weight }} kg
          </NDescriptionsItem>
          <NDescriptionsItem label="状态">
            <NTag :type="parcelStatusMap[searchResult.status]?.type || 'default'">
              {{ parcelStatusMap[searchResult.status]?.text || searchResult.status }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem v-if="searchResult.shelf" label="货架位置">
            {{ searchResult.shelf.shelf_code }} ({{ searchResult.shelf.area }}区{{ searchResult.shelf.floor }}层{{
              searchResult.shelf.column }}列)
          </NDescriptionsItem>
          <NDescriptionsItem label="入库时间">
            {{ formatDateTime(searchResult.received_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem v-if="searchResult.shelved_at" label="上架时间">
            {{ formatDateTime(searchResult.shelved_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem v-if="searchResult.picked_up_at" label="取件时间">
            {{ formatDateTime(searchResult.picked_up_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem v-if="searchResult.expected_overdue_at" label="预计滞留时间">
            {{ formatDateTime(searchResult.expected_overdue_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem v-if="searchResult.notes" label="备注">
            {{ searchResult.notes }}
          </NDescriptionsItem>
        </NDescriptions>

        <div style="margin-top: 16px; text-align: center;">
          <NButton @click="handleReset">继续查询</NButton>
        </div>
      </div>

      <!-- 错误提示 -->
      <div v-else-if="searchError" class="result-section">
        <NResult status="404" title="未找到包裹" :description="searchError">
          <template #footer>
            <NButton @click="handleReset">重新查询</NButton>
          </template>
        </NResult>
      </div>
    </NCard>
  </div>
</template>

<style scoped>
.search-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #36ad6a 0%, #18a058 100%);
  padding: 40px 20px;
}

.search-card {
  width: 100%;
  max-width: 600px;
}

.result-section {
  margin-top: 32px;
  padding-top: 32px;
  border-top: 1px solid #f0f0f0;
}
</style>
