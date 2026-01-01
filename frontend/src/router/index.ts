import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import Home from '../views/Home.vue';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import MainLayout from '../views/MainLayout.vue';
import Dashboard from '../views/Dashboard.vue';
import Search from '../views/Search.vue';
import MyParcels from '../views/MyParcels.vue';
import MyShipments from '../views/MyShipments.vue';
import ParcelManagement from '../views/ParcelManagement.vue';
import ShelfManagement from '../views/ShelfManagement.vue';
import OverdueManagement from '../views/OverdueManagement.vue';
import Profile from '../views/Profile.vue';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false }
  },
  // 使用 MainLayout 作为父路由，所有需要侧边栏的页面作为子路由
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: Dashboard
      },
      {
        path: 'search',
        name: 'Search',
        component: Search,
        meta: { requiresAuth: false }
      },
      {
        path: 'my-parcels',
        name: 'MyParcels',
        component: MyParcels,
        meta: { requiresUser: true }
      },
      {
        path: 'my-shipments',
        name: 'MyShipments',
        component: MyShipments,
        meta: { requiresUser: true }
      },
      {
        path: 'parcel-management',
        name: 'ParcelManagement',
        component: ParcelManagement,
        meta: { requiresStaff: true }
      },
      {
        path: 'shelf-management',
        name: 'ShelfManagement',
        component: ShelfManagement,
        meta: { requiresStaff: true }
      },
      {
        path: 'overdue-management',
        name: 'OverdueManagement',
        component: OverdueManagement,
        meta: { requiresStaff: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: Profile
      }
    ]
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});

// 路由守卫 - 检查认证和权限
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('accessToken');
  const userStr = localStorage.getItem('userInfo');
  let userInfo = null;

  if (userStr) {
    try {
      userInfo = JSON.parse(userStr);
    } catch (e) {
      next({ name: 'Login' });
      return;
    }
  }

  // 需要认证的页面
  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login' });
    return;
  }

  // 需要员工或管理员权限的页面
  if (to.meta.requiresStaff) {
    if (!userInfo || (userInfo.role !== 'staff' && userInfo.role !== 'admin')) {
      next({ name: 'Home' });
      return;
    }
  }

  // 只允许普通用户访问的页面（我的包裹、我的寄件）
  if (to.meta.requiresUser) {
    if (!userInfo || userInfo.role === 'staff' || userInfo.role === 'admin') {
      next({ name: 'Dashboard' });
      return;
    }
  }

  // 已登录用户访问首页时重定向到仪表盘
  if (to.name === 'Home' && token) {
    next({ name: 'Dashboard' });
    return;
  }

  next();
});

export default router;
