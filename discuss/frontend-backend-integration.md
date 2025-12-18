# 前端API集成说明

## 概述

前端已完全从 mock 数据切换到真实的后端 API 调用。所有 API 请求现在都通过 axios 实例发送到后端服务器。

## 更改内容

### 1. 用户API (`src/api/user.ts`)

- ✅ `register()` - 用户注册 → `POST /v1/register`
- ✅ `login()` - 用户登录 → `POST /v1/login`
- ✅ `getProfile()` - 获取用户信息 → `GET /v1/user`
- ✅ `updateProfile()` - 更新用户信息 → `PUT /v1/user`

### 2. 包裹API (`src/api/parcel.ts`)

- ✅ `getParcelList()` - 获取包裹列表(管理端) → `GET /v1/parcels`
- ✅ `getMyParcels()` - 获取我的包裹 → `GET /v1/my-parcels`
- ✅ `getParcelByPickupCode()` - 通过取件码查询 → `GET /v1/parcels/by-pickup-code/{code}`
- ✅ `getParcelByTrackingNumber()` - 通过快递单号查询 → `GET /v1/parcels/by-tracking/{number}`
- ✅ `receiveParcel()` - 包裹入库 → `POST /v1/parcels/receive`
- ✅ `pickupParcel()` - 包裹取件 → `POST /v1/parcels/pickup`
- ✅ `getDashboardStatistics()` - 获取统计数据 → `GET /v1/dashboard/statistics`

### 3. 寄件API (`src/api/shipment.ts`)

- ✅ `getMyShipments()` - 获取我的寄件单列表 → `GET /v1/my-shipments`
- ✅ `createShipment()` - 创建寄件单 → `POST /v1/shipments`
- ✅ `getShipmentDetail()` - 获取寄件单详情 → `GET /v1/shipments/{id}`
- ✅ `calculateFreight()` - 运费试算 → `POST /v1/shipments/calculate-freight`
- ✅ `getCourierCompanies()` - 获取快递公司列表 → `GET /v1/courier-companies`

## 技术细节

### 字段名转换

后端使用 snake_case 命名，前端使用 camelCase 命名。API 层自动进行转换：

**后端格式 → 前端格式：**
- `recipient_name` → `recipientName`
- `recipient_phone` → `recipientPhone`
- `tracking_number` → `trackingNumber`
- `pickup_code` → `pickupCode`
- `courier_company` → `courierCompany`
- `received_at` → `receivedAt`
- `picked_up_at` → `pickedUpAt`
- 等等...

### 错误处理

所有 API 函数都包含完整的错误处理：

```typescript
try {
  const response = await request.get('/endpoint');
  return response;
} catch (error: any) {
  return {
    code: error.response?.data?.code || 50000,
    message: error.response?.data?.message || '默认错误消息',
    data: undefined
  };
}
```

### Token 管理

- 登录成功后自动保存 token 到 `localStorage.accessToken`
- 请求拦截器自动添加 `Authorization: Bearer {token}` 头
- 401 错误时自动清除 token 并跳转到登录页

## 启动说明

### 1. 启动后端

```bash
cd backend
nunu run ./cmd/server
```

后端将在 `http://localhost:8000` 运行

### 2. 启动前端

```bash
cd frontend
bun dev
```

前端将在 `http://localhost:5173` (或其他端口) 运行

### 3. 环境配置

前端 `.env` 文件已配置：

```
VITE_API_BASE_URL=http://localhost:8000/v1
```

## 测试建议

### 基本流程测试：

1. **用户注册与登录**
   - 访问注册页面，创建新用户
   - 使用新账号登录
   - 检查 token 是否正确保存

2. **包裹管理**
   - 测试包裹入库功能
   - 查询包裹（取件码/快递单号）
   - 包裹取件
   - 查看我的包裹列表

3. **寄件功能**
   - 创建新的寄件单
   - 运费试算
   - 查看我的寄件单列表

4. **统计仪表盘**
   - 查看首页统计数据
   - 确认数据趋势图正确显示

## 注意事项

1. **首次使用需要注册账号** - Mock数据中的账号已不可用
2. **需要后端服务运行** - 前端现在完全依赖后端API
3. **跨域问题** - 后端已配置CORS，允许前端访问
4. **数据持久化** - 数据现在保存在数据库中(SQLite)，重启后不会丢失

## 已删除的文件

建议删除不再使用的 mock 文件：
- `src/api/mock-data.ts` (可选保留作为参考)

## 后续工作

- [ ] 添加更多的错误提示和用户反馈
- [ ] 优化加载状态显示
- [ ] 添加数据缓存机制
- [ ] 完善表单验证
