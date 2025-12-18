# 左侧栏实现文档

## 概述

参考 hyacinth-frontend 项目的实现，为校园物流系统添加了一个始终存在的左侧导航栏。

## 实现方案

### 1. 创建 MainLayout 组件

新建了 `frontend/src/views/MainLayout.vue` 组件，作为所有需要侧边栏的页面的统一布局容器。

#### 主要特性：
- **侧边栏导航**：使用 NLayoutSider 和 NMenu 组件实现
- **可折叠**：支持侧边栏折叠/展开，宽度在 64px 和 240px 之间切换
- **用户信息显示**：侧边栏底部显示用户头像和用户名
- **下拉菜单**：点击用户信息可访问个人中心或退出登录
- **路由导航**：根据用户角色动态显示菜单项

#### 菜单配置：
- 普通用户菜单：
  - 仪表盘 (`/dashboard`)
  - 包裹查询 (`/search`)
  - 我的包裹 (`/my-parcels`)
  - 我的寄件 (`/my-shipments`)

- 员工/管理员额外菜单：
  - 包裹管理 (`/parcel-management`)

### 2. 更新路由配置

修改了 `frontend/src/router/index.ts`，使用嵌套路由结构：

```typescript
{
  path: '/',
  component: MainLayout,
  meta: { requiresAuth: true },
  children: [
    { path: 'dashboard', name: 'Dashboard', component: Dashboard },
    { path: 'search', name: 'Search', component: Search },
    { path: 'my-parcels', name: 'MyParcels', component: MyParcels },
    { path: 'my-shipments', name: 'MyShipments', component: MyShipments },
    { path: 'parcel-management', name: 'ParcelManagement', component: ParcelManagement },
    { path: 'profile', name: 'Profile', component: Profile }
  ]
}
```

这样所有需要侧边栏的页面都会自动包含在 MainLayout 中。

### 3. 简化 Dashboard 组件

原 `Dashboard.vue` 组件包含了完整的布局代码（侧边栏 + 顶栏 + 内容区）。现在已移除布局部分，只保留了核心的业务逻辑和内容展示：

- 移除了 NLayout、NLayoutSider、NLayoutHeader 等布局组件
- 移除了菜单选项和用户菜单的定义
- 只保留了仪表盘的数据加载和展示逻辑

### 4. 技术细节

#### 状态管理：
- 使用 `useUserStore()` 管理用户状态
- 从 localStorage 初始化用户信息
- 根据用户角色动态生成菜单

#### 响应式交互：
- 监听路由变化，自动更新激活的菜单项
- 折叠动画：侧边栏折叠时用户名平滑隐藏/显示

#### 样式设计：
- 使用 NaiveUI 的主题变量，支持亮/暗主题切换
- 品牌色：`#667eea`
- Flexbox 布局：菜单在上，用户信息固定在底部

## 优势

1. **代码复用**：所有页面共享同一个布局，避免重复代码
2. **一致性**：确保所有页面的侧边栏行为和样式一致
3. **易维护**：修改侧边栏只需编辑 MainLayout 组件
4. **扩展性**：添加新页面只需在路由中添加子路由即可

## 测试

开发服务器已启动在 `http://localhost:9000/`

测试要点：
- ✓ 登录后自动重定向到 Dashboard
- ✓ 侧边栏在所有页面保持可见
- ✓ 菜单项点击正确跳转
- ✓ 折叠/展开功能正常
- ✓ 用户下拉菜单功能正常
- ✓ 根据用户角色显示不同的菜单项

## 参考

本实现参考了 `C:\Workspace\Project\hyacinth-frontend` 项目中的 `UserPanel.vue` 组件设计。
