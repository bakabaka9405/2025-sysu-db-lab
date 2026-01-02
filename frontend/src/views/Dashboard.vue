<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import {
  NButton,
  NGrid,
  NGridItem,
  NCard,
  NStatistic,
  useMessage
} from 'naive-ui'
import { getDashboardStatistics } from '../api/parcel'
import { getMyParcels } from '../api/parcel'
import { getMyShipments } from '../api/shipment'
import type { DashboardStatistics, Parcel, Shipment } from '../api/types'
import ParcelTable from '../components/table/ParcelTable.vue'
import ShipmentTable from '../components/table/ShipmentTable.vue'

const router = useRouter()
const userStore = useUserStore()
const message = useMessage()

// 初始化用户状态
userStore.initFromStorage()

const loading = ref(false)
const statistics = ref<DashboardStatistics | null>(null)
const recentParcels = ref<Parcel[]>([])
const recentShipments = ref<Shipment[]>([])

// 计算属性
const isLoggedIn = computed(() => userStore.isLoggedIn)
const userInfo = computed(() => userStore.userInfo)
const isStaffOrAdmin = computed(() => {
  const role = userInfo.value?.role
  return role === 'staff' || role === 'admin'
})

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    // 加载统计数据
    if (isStaffOrAdmin.value) {
      const statsRes = await getDashboardStatistics()
      if (statsRes.code === 0 && statsRes.data) {
        statistics.value = statsRes.data
      }
    }

    // 加载最近的包裹
    const parcelsRes = await getMyParcels({ page: 1, page_size: 5 })
    if (parcelsRes.code === 0 && parcelsRes.data) {
      recentParcels.value = parcelsRes.data.list
    }

    // 加载最近的寄件
    const shipmentsRes = await getMyShipments({ page: 1, page_size: 5 })
    if (shipmentsRes.code === 0 && shipmentsRes.data) {
      recentShipments.value = shipmentsRes.data.list
    }
  } catch (error) {
    message.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (isLoggedIn.value) {
    loadData()
  }
})
</script>

<template>
  <div class="dashboard-container">
    <div v-if="!isLoggedIn" class="not-logged-in">
      <h2>请先登录</h2>
      <NButton type="primary" @click="router.push('/login')" style="margin-top: 16px">
        前往登录
      </NButton>
    </div>

    <div v-else>
      <!-- 管理员统计数据 -->
      <div v-if="isStaffOrAdmin && statistics" class="statistics-section">
        <h3>今日数据</h3>
        <NGrid :x-gap="24" :y-gap="24" :cols="4">
          <NGridItem>
            <NCard>
              <NStatistic label="今日入库" :value="statistics.today.received_count" />
            </NCard>
          </NGridItem>
          <NGridItem>
            <NCard>
              <NStatistic label="今日取件" :value="statistics.today.picked_up_count" />
            </NCard>
          </NGridItem>
          <NGridItem>
            <NCard>
              <NStatistic label="滞留包裹" :value="statistics.today.overdue_count" />
            </NCard>
          </NGridItem>
          <NGridItem>
            <NCard>
              <NStatistic label="货架使用率">
                <template #default>
                  {{ (statistics.shelves.utilization_rate || 0).toFixed(1) }}%
                </template>
              </NStatistic>
            </NCard>
          </NGridItem>
        </NGrid>
      </div>

      <!-- 最近的包裹 -->
      <div class="section">
        <div class="section-header">
          <h3>我的包裹</h3>
          <NButton text type="primary" @click="router.push('/my-parcels')">
            查看全部
          </NButton>
        </div>
        <NCard>
          <ParcelTable
            :parcels="recentParcels"
            :loading="loading"
            compact
            :pagination="false"
          />
          <div v-if="recentParcels.length === 0 && !loading" class="empty-state">
            暂无包裹
          </div>
        </NCard>
      </div>

      <!-- 最近的寄件 -->
      <div class="section">
        <div class="section-header">
          <h3>我的寄件</h3>
          <NButton text type="primary" @click="router.push('/my-shipments')">
            查看全部
          </NButton>
        </div>
        <NCard>
          <ShipmentTable
            :shipments="recentShipments"
            :loading="loading"
            compact
            :pagination="false"
          />
          <div v-if="recentShipments.length === 0 && !loading" class="empty-state">
            暂无寄件记录
          </div>
        </NCard>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-container {
  width: 100%;
}

.not-logged-in {
  text-align: center;
  padding: 80px 0;
}

.statistics-section {
  margin-bottom: 24px;
}

.section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h3 {
  margin: 0;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #999;
}
</style>
