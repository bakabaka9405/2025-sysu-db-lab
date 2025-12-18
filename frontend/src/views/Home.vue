<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import {
  NLayout,
  NLayoutHeader,
  NLayoutContent,
  NButton,
  NSpace,
  NGrid,
  NGridItem,
  NStatistic,
  NCard,
  NIcon,
  NDropdown,
  useMessage
} from 'naive-ui'
import {
  CubeOutline,
  LocationOutline,
  TimeOutline,
  ShieldCheckmarkOutline,
  SpeedometerOutline,
  NotificationsOutline,
  ArchiveOutline,
  PersonOutline
} from '@vicons/ionicons5'

const router = useRouter()
const userStore = useUserStore()
const message = useMessage()

// 初始化用户状态
userStore.initFromStorage()

// 是否已登录
const isLoggedIn = computed(() => userStore.isLoggedIn)
const userInfo = computed(() => userStore.userInfo)

// 用户下拉菜单选项
const userOptions = [
  {
    label: '个人中心',
    key: 'profile'
  },
  {
    label: '退出登录',
    key: 'logout'
  }
]

// 处理用户菜单选择
const handleUserMenuSelect = (key: string) => {
  if (key === 'logout') {
    userStore.logout()
    message.success('已退出登录')
    router.push('/')
  } else if (key === 'profile') {
    router.push('/profile')
  }
}

// 网站统计数据
const statistics = [
  { label: '配送点', value: '15' },
  { label: '注册用户', value: '5,680' },
  { label: '今日配送', value: '328' },
  { label: '累计配送', value: '156,892' }
]

// 产品特性数据
const features = [
  {
    title: '快速配送',
    description: '校内1小时送达，高效的物流配送网络覆盖校园每个角落',
    icon: SpeedometerOutline
  },
  {
    title: '安全可靠',
    description: '专业团队配送，实名认证，全程监控，保障包裹安全',
    icon: ShieldCheckmarkOutline
  },
  {
    title: '便捷取件',
    description: '多个自提点位，支持预约取件，随时随地轻松取件',
    icon: LocationOutline
  },
  {
    title: '实时追踪',
    description: '包裹状态实时更新，配送进度一目了然',
    icon: NotificationsOutline
  }
]

// 服务类型数据
const services = [
  {
    title: '包裹配送',
    description: '快递包裹统一配送服务，支持宿舍、教学楼等多种地址配送',
    icon: CubeOutline
  },
  {
    title: '暂存服务',
    description: '临时外出？我们提供安全的包裹暂存服务，随到随取',
    icon: ArchiveOutline
  },
  {
    title: '配送追踪',
    description: '实时查看包裹位置和配送进度，掌握每一个配送环节',
    icon: LocationOutline
  },
  {
    title: '到件通知',
    description: '包裹到达及时提醒，支持短信、微信等多种通知方式',
    icon: TimeOutline
  }
]
</script>

<template>
  <NLayout style="width: 100%; min-height: 100vh;">
    <!-- 页面顶部导航栏 -->
    <NLayoutHeader class="header">
      <div class="header-content">
        <div class="logo">校园物流系统</div>
        <!-- 顶部导航按钮区域 -->
        <NSpace v-if="!isLoggedIn">
          <NButton @click="router.push('/login')">登录</NButton>
          <NButton type="primary" @click="router.push('/register')">注册</NButton>
        </NSpace>
        <NDropdown v-else :options="userOptions" @select="handleUserMenuSelect">
          <NButton>
            <template #icon>
              <NIcon :component="PersonOutline" />
            </template>
            {{ userInfo?.realName || userInfo?.username }}
          </NButton>
        </NDropdown>
      </div>
    </NLayoutHeader>

    <!-- 页面内容区域 -->
    <NLayoutContent content-style="padding: 24px;">
      <!-- Hero Section -->
      <div class="hero-section">
        <h1 class="hero-title">便捷、高效、安全的校园物流服务</h1>
        <p class="hero-description">让每一个包裹都能准时送达，让每一次配送都充满信赖</p>
        <NButton v-if="!isLoggedIn" type="primary" size="large" @click="router.push('/register')">
          立即开始
        </NButton>
      </div>

      <!-- 统计数据 -->
      <div class="statistics-section">
        <NGrid :x-gap="24" :y-gap="24" :cols="4" responsive="screen">
          <NGridItem v-for="stat in statistics" :key="stat.label">
            <NStatistic :label="stat.label" :value="stat.value" />
          </NGridItem>
        </NGrid>
      </div>

      <!-- 产品特性部分 -->
      <div class="section">
        <h2 class="section-title">为什么选择我们？</h2>
        <NGrid :x-gap="24" :y-gap="24" :cols="4" responsive="screen">
          <NGridItem v-for="feature in features" :key="feature.title">
            <NCard class="feature-card">
              <div class="feature-icon">
                <NIcon :component="feature.icon" :size="32" />
              </div>
              <h3 class="feature-title">{{ feature.title }}</h3>
              <p class="feature-description">{{ feature.description }}</p>
            </NCard>
          </NGridItem>
        </NGrid>
      </div>

      <!-- 主要功能 -->
      <div class="section">
        <h2 class="section-title">主要服务</h2>
        <NGrid :x-gap="24" :y-gap="24" :cols="4" responsive="screen">
          <NGridItem v-for="service in services" :key="service.title">
            <NCard class="service-card">
              <div class="service-icon">
                <NIcon :component="service.icon" :size="32" />
              </div>
              <h3 class="service-title">{{ service.title }}</h3>
              <p class="service-description">{{ service.description }}</p>
            </NCard>
          </NGridItem>
        </NGrid>
      </div>

      <!-- 行动号召部分 -->
      <div v-if="!isLoggedIn" class="cta-section">
        <h2 class="cta-title">准备好开始了吗？</h2>
        <p class="cta-description">注册账号，即刻体验便捷的校园物流服务</p>
        <NButton type="primary" size="large" @click="router.push('/register')">立即注册</NButton>
      </div>
    </NLayoutContent>
  </NLayout>
</template>

<style scoped>
/* 页面顶部导航栏样式 */
.header {
  height: 64px;
  padding: 0 24px;
  background-color: #fff;
  border-bottom: 1px solid #eee;
  width: 100%;
  position: sticky;
  top: 0;
  z-index: 100;
}

/* 导航栏内容居中对齐 */
.header-content {
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* Logo样式 */
.logo {
  font-size: 20px;
  font-weight: bold;
  color: #667eea;
  letter-spacing: 1px;
}

/* Hero区域样式 */
.hero-section {
  max-width: 1200px;
  margin: 48px auto;
  text-align: center;
  padding: 64px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  color: white;
}

.hero-title {
  font-size: 36px;
  font-weight: bold;
  margin: 0 0 16px;
}

.hero-description {
  font-size: 18px;
  margin: 0 0 32px;
  opacity: 0.95;
}

/* 统计数据区域样式 */
.statistics-section {
  max-width: 1200px;
  margin: 48px auto;
  padding: 32px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

/* 通用内容区块样式 */
.section {
  max-width: 1200px;
  margin: 48px auto;
}

/* 区块标题样式 */
.section-title {
  text-align: center;
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 40px;
  color: #333;
}

/* 特性卡片和服务卡片通用样式 */
.feature-card,
.service-card {
  text-align: center;
  padding: 32px 24px;
  height: 100%;
  transition: all 0.3s ease;
  cursor: pointer;
}

/* 卡片悬停效果 */
.feature-card:hover,
.service-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

/* 卡片图标样式 */
.feature-icon,
.service-icon {
  margin-bottom: 16px;
  color: #667eea;
  display: flex;
  justify-content: center;
}

/* 卡片标题样式 */
.feature-title,
.service-title {
  font-size: 18px;
  font-weight: 600;
  margin: 12px 0;
  color: #333;
}

/* 卡片描述文字样式 */
.feature-description,
.service-description {
  font-size: 14px;
  color: #666;
  margin: 0;
  line-height: 1.6;
}

/* 行动号召区域样式 */
.cta-section {
  max-width: 1200px;
  margin: 64px auto;
  text-align: center;
  padding: 64px 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  color: white;
}

/* 行动号召标题样式 */
.cta-title {
  font-size: 32px;
  font-weight: bold;
  margin: 0 0 16px;
}

/* 行动号召描述样式 */
.cta-description {
  font-size: 18px;
  margin: 0 0 32px;
  opacity: 0.95;
}
</style>
