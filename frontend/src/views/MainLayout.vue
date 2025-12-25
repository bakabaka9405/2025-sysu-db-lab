<script setup lang="ts">
import { ref, computed, h, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import {
  NLayout,
  NLayoutSider,
  NLayoutContent,
  NMenu,
  NAvatar,
  NDropdown,
  NSpace,
  type MenuOption,
  type DropdownOption
} from 'naive-ui'
import {
  HomeOutline,
  CubeOutline,
  SendOutline,
  SearchOutline,
  PersonOutline,
  LogOutOutline,
  BarChartOutline,
  GridOutline,
  AlertCircleOutline
} from '@vicons/ionicons5'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 初始化用户状态
userStore.initFromStorage()

// 侧边栏折叠状态
const collapsed = ref(false)

// 控制用户名显示，只在动画完成后显示
const showUsername = ref(true)

// 当前激活的菜单项
const activeKey = ref('')

// 计算属性
const userInfo = computed(() => userStore.userInfo)
const isStaffOrAdmin = computed(() => {
  const role = userInfo.value?.role
  return role === 'staff' || role === 'admin'
})

// 导航菜单选项
const menuOptions = computed<MenuOption[]>(() => {
  const baseMenus: MenuOption[] = [
    {
      label: '仪表盘',
      key: '/dashboard',
      icon: () => h(HomeOutline)
    },
    {
      label: '包裹查询',
      key: '/search',
      icon: () => h(SearchOutline)
    }
  ]

  // 只有普通用户显示"我的包裹"和"我的寄件"
  if (!isStaffOrAdmin.value) {
    baseMenus.push(
      {
        label: '我的包裹',
        key: '/my-parcels',
        icon: () => h(CubeOutline)
      },
      {
        label: '我的寄件',
        key: '/my-shipments',
        icon: () => h(SendOutline)
      }
    )
  }

  // staff/admin 显示管理菜单
  if (isStaffOrAdmin.value) {
    baseMenus.push({
      label: '包裹管理',
      key: '/parcel-management',
      icon: () => h(BarChartOutline)
    })
    baseMenus.push({
      label: '货架管理',
      key: '/shelf-management',
      icon: () => h(GridOutline)
    })
    baseMenus.push({
      label: '滞留件管理',
      key: '/overdue-management',
      icon: () => h(AlertCircleOutline)
    })
  }

  return baseMenus
})

// 用户头像下拉菜单选项
const dropdownOptions: DropdownOption[] = [
  {
    label: '个人中心',
    key: 'profile'
  },
  {
    label: '退出登录',
    key: 'logout'
  }
]

// 监听侧边栏折叠状态变化，控制用户名显示时机
watch(collapsed, (newCollapsed) => {
  if (newCollapsed) {
    // 折叠时立即隐藏用户名
    showUsername.value = false
  } else {
    // 展开时延迟显示用户名，等待动画完成
    setTimeout(() => {
      showUsername.value = true
    }, 100)
  }
})

// 监听路由变化，更新激活的菜单项
watch(
  () => route.path,
  (newPath) => {
    activeKey.value = newPath
  },
  { immediate: true }
)

// 处理导航菜单项点击事件
const handleUpdateValue = (key: string) => {
  router.push(key)
}

// 处理用户头像下拉菜单选择事件
const handleSelect = (key: string) => {
  if (key === 'profile') {
    // 跳转到个人中心页面
    router.push('/profile')
  } else if (key === 'logout') {
    // 退出登录，清除token并跳转到登录页
    userStore.logout()
    router.push('/login')
  }
}

// 组件挂载时执行的操作
onMounted(() => {
  // 检查用户是否登录，如果未登录则跳转到登录页
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }

  // 根据当前路由动态设置激活的菜单项
  const currentRoute = router.currentRoute.value.path
  const menuKeys = menuOptions.value.map((option) => option.key)
  if (menuKeys.includes(currentRoute)) {
    activeKey.value = currentRoute
  } else {
    // 如果当前路由不在菜单中，设置默认值
    activeKey.value = '/dashboard'
  }
})
</script>

<template>
  <!-- 整体布局容器，包含侧边栏和内容区 -->
  <NLayout class="main-layout" has-sider>
    <!-- 侧边导航栏 -->
    <NLayoutSider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      show-trigger
      v-model:collapsed="collapsed"
      :native-scrollbar="false"
    >
      <div class="sider-content">
        <!-- 品牌标题 -->
        <div class="brand-header">
          <h2 v-if="!collapsed" class="brand-title">校园物流</h2>
          <h2 v-else class="brand-title-collapsed">物流</h2>
        </div>

        <!-- 导航菜单 -->
        <NMenu
          v-model:value="activeKey"
          :options="menuOptions"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          @update:value="handleUpdateValue"
        />

        <!-- 用户信息区域，显示在侧边栏底部 -->
        <div class="user-profile">
          <NDropdown :options="dropdownOptions" @select="handleSelect">
            <NSpace align="center">
              <NAvatar
                round
                :size="32"
                :style="{ backgroundColor: '#18a058', color: '#fff' }"
              >
                {{ userInfo?.username?.charAt(0).toUpperCase() || 'U' }}
              </NAvatar>
              <span class="username" v-if="!collapsed && showUsername">
                {{ userInfo?.realName || userInfo?.username }}
              </span>
            </NSpace>
          </NDropdown>
        </div>
      </div>
    </NLayoutSider>

    <!-- 主内容区域 -->
    <NLayout>
      <NLayoutContent content-style="padding: 24px;">
        <!-- 路由视图，根据当前路由显示对应组件 -->
        <router-view />
      </NLayoutContent>
    </NLayout>
  </NLayout>
</template>

<style scoped>
/* 主布局容器 - 占满整个视口 */
.main-layout {
  width: 100%;
  height: 100vh;
}

/* 侧边栏内容样式，使用flex布局实现顶部菜单+底部用户信息的结构 */
.sider-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 品牌标题区域 */
.brand-header {
  flex-shrink: 0;
  padding: 20px 24px;
  text-align: center;
  border-bottom: 1px solid var(--n-border-color);
  margin-bottom: 8px;
}

.brand-title {
  color: #18a058;
  margin: 0;
  font-size: 20px;
  font-weight: bold;
  transition: opacity 0.3s;
}

.brand-title-collapsed {
  color: #18a058;
  margin: 0;
  font-size: 16px;
  font-weight: bold;
}

/* 用户信息区域样式，固定在侧边栏底部 */
.user-profile {
  margin-top: auto;
  padding: 16px;
  border-top: 1px solid var(--n-border-color);
  cursor: pointer;
}

/* 用户名称文本样式 */
.username {
  font-size: 14px;
  color: var(--n-text-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 120px;
}
</style>
