# API接口设计

## 1. 接口设计原则

### 1.1 RESTful规范
- 使用标准的HTTP方法：GET, POST, PUT, DELETE
- 资源使用复数名词：/api/v1/parcels, /api/v1/users
- 使用HTTP状态码表示请求结果
- 统一的响应格式

### 1.2 版本控制
- API版本通过URL路径标识：/api/v1/
- 便于后续版本迭代和兼容性管理

### 1.3 认证授权
- 使用JWT Bearer Token认证
- 公开接口：注册、登录
- 受保护接口：需要在请求头中携带Token
  ```
  Authorization: Bearer <token>
  ```

### 1.4 统一响应格式

#### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 业务数据
  }
}
```

#### 错误响应
```json
{
  "code": 40001,
  "message": "参数错误：取件码不能为空",
  "data": null
}
```

#### 分页响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

### 1.5 错误码规范

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 系统错误 |
| 20001 | 认证失败 |
| 20002 | Token过期 |
| 20003 | 无权限 |
| 40001 | 参数错误 |
| 40002 | 资源不存在 |
| 40003 | 资源已存在 |
| 40004 | 业务逻辑错误 |

## 2. 用户模块接口

### 2.1 用户注册
- **接口**: POST /api/v1/users/register
- **权限**: 公开
- **描述**: 用户注册

#### 请求参数
```json
{
  "username": "zhangsan",
  "password": "123456",
  "real_name": "张三",
  "phone": "13800138000",
  "id_card": "440101199001011234",
  "email": "zhangsan@example.com"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "id": 1,
    "username": "zhangsan",
    "real_name": "张三",
    "phone": "13800138000",
    "role": "user",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

### 2.2 用户登录
- **接口**: POST /api/v1/users/login
- **权限**: 公开
- **描述**: 用户登录，返回JWT Token

#### 请求参数
```json
{
  "username": "zhangsan",
  "password": "123456"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "zhangsan",
      "real_name": "张三",
      "role": "user"
    }
  }
}
```

### 2.3 获取当前用户信息
- **接口**: GET /api/v1/users/profile
- **权限**: 需要认证
- **描述**: 获取当前登录用户的信息

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "zhangsan",
    "real_name": "张三",
    "phone": "13800138000",
    "email": "zhangsan@example.com",
    "role": "user",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

### 2.4 更新用户信息
- **接口**: PUT /api/v1/users/profile
- **权限**: 需要认证
- **描述**: 更新当前用户的信息

#### 请求参数
```json
{
  "real_name": "张三",
  "email": "new_email@example.com"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "zhangsan",
    "real_name": "张三",
    "email": "new_email@example.com"
  }
}
```

### 2.5 修改密码
- **接口**: POST /api/v1/users/change-password
- **权限**: 需要认证
- **描述**: 修改当前用户的密码

#### 请求参数
```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "密码修改成功",
  "data": null
}
```

## 3. 包裹模块接口

### 3.1 包裹入库
- **接口**: POST /api/v1/parcels
- **权限**: staff, admin
- **描述**: 包裹入库，自动生成取件码并分配货架

#### 请求参数
```json
{
  "tracking_number": "SF1234567890",
  "recipient_name": "张三",
  "recipient_phone": "13800138000",
  "recipient_id_card": "440101199001011234",
  "courier_company": "顺丰速运",
  "size": "small",
  "weight": 1.5,
  "notes": "易碎品"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "入库成功",
  "data": {
    "id": 1001,
    "tracking_number": "SF1234567890",
    "pickup_code": "A12345",
    "recipient_name": "张三",
    "recipient_phone": "13800138000",
    "courier_company": "顺丰速运",
    "size": "small",
    "weight": 1.5,
    "status": "received",
    "shelf": {
      "id": 1,
      "shelf_code": "A-01-01",
      "area": "A",
      "floor": 1,
      "column": 1
    },
    "received_at": "2025-01-15T10:30:00Z",
    "expected_overdue_at": "2025-01-18T10:30:00Z"
  }
}
```

### 3.2 批量包裹入库（OCR识别）
- **接口**: POST /api/v1/parcels/batch
- **权限**: staff, admin
- **描述**: 通过OCR识别批量导入包裹信息

#### 请求参数（multipart/form-data）
```
file: <图片文件>
```

#### 响应示例
```json
{
  "code": 0,
  "message": "批量入库成功",
  "data": {
    "success_count": 10,
    "failed_count": 2,
    "parcels": [
      {
        "tracking_number": "SF1234567890",
        "pickup_code": "A12345",
        "status": "success"
      },
      {
        "tracking_number": "SF1234567891",
        "error": "识别失败"
      }
    ]
  }
}
```

### 3.3 包裹上架
- **接口**: PUT /api/v1/parcels/{id}/shelve
- **权限**: staff, admin
- **描述**: 将包裹状态更新为已上架

#### 请求参数
```json
{
  "shelf_id": 1
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "上架成功",
  "data": {
    "id": 1001,
    "status": "shelved",
    "shelf": {
      "shelf_code": "A-01-01",
      "area": "A",
      "floor": 1,
      "column": 1
    },
    "shelved_at": "2025-01-15T10:35:00Z"
  }
}
```

### 3.4 包裹取件
- **接口**: POST /api/v1/parcels/pickup
- **权限**: 需要认证（用户本人或员工）
- **描述**: 通过取件码取件

#### 请求参数
```json
{
  "pickup_code": "A12345",
  "recipient_phone": "13800138000"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "取件成功",
  "data": {
    "id": 1001,
    "tracking_number": "SF1234567890",
    "pickup_code": "A12345",
    "recipient_name": "张三",
    "courier_company": "顺丰速运",
    "shelf": {
      "shelf_code": "A-01-01",
      "area": "A",
      "floor": 1,
      "column": 1
    },
    "status": "picked_up",
    "picked_up_at": "2025-01-15T14:30:00Z"
  }
}
```

### 3.5 查询包裹（取件码）
- **接口**: GET /api/v1/parcels/by-pickup-code/{pickup_code}
- **权限**: 需要认证
- **描述**: 通过取件码查询包裹信息

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1001,
    "tracking_number": "SF1234567890",
    "pickup_code": "A12345",
    "recipient_name": "张三",
    "recipient_phone": "13800138000",
    "courier_company": "顺丰速运",
    "size": "small",
    "weight": 1.5,
    "status": "ready_for_pickup",
    "shelf": {
      "shelf_code": "A-01-01",
      "area": "A",
      "floor": 1,
      "column": 1
    },
    "received_at": "2025-01-15T10:30:00Z",
    "expected_overdue_at": "2025-01-18T10:30:00Z"
  }
}
```

### 3.6 查询包裹（快递单号）
- **接口**: GET /api/v1/parcels/by-tracking-number/{tracking_number}
- **权限**: 需要认证
- **描述**: 通过快递单号查询包裹信息

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1001,
    "tracking_number": "SF1234567890",
    "pickup_code": "A12345",
    "status": "ready_for_pickup",
    "shelf": {
      "shelf_code": "A-01-01"
    }
  }
}
```

### 3.7 查询我的包裹
- **接口**: GET /api/v1/parcels/my
- **权限**: 需要认证
- **描述**: 查询当前用户的所有包裹（通过手机号匹配）

#### 查询参数
- `status`: 状态过滤（可选）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1001,
        "tracking_number": "SF1234567890",
        "pickup_code": "A12345",
        "courier_company": "顺丰速运",
        "status": "ready_for_pickup",
        "received_at": "2025-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

### 3.8 包裹列表（管理端）
- **接口**: GET /api/v1/parcels
- **权限**: staff, admin
- **描述**: 查询包裹列表，支持多条件筛选

#### 查询参数
- `status`: 状态（可选）
- `courier_company`: 快递公司（可选）
- `recipient_phone`: 收件人手机号（可选）
- `start_date`: 开始日期（可选）
- `end_date`: 结束日期（可选）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1001,
        "tracking_number": "SF1234567890",
        "pickup_code": "A12345",
        "recipient_name": "张三",
        "recipient_phone": "13800138000",
        "courier_company": "顺丰速运",
        "status": "ready_for_pickup",
        "shelf_code": "A-01-01",
        "received_at": "2025-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 1000,
      "total_pages": 50
    }
  }
}
```

### 3.9 包裹详情
- **接口**: GET /api/v1/parcels/{id}
- **权限**: 需要认证
- **描述**: 查询包裹详细信息

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1001,
    "tracking_number": "SF1234567890",
    "pickup_code": "A12345",
    "recipient_name": "张三",
    "recipient_phone": "13800138000",
    "recipient_id_card": "440101199001011234",
    "courier_company": "顺丰速运",
    "size": "small",
    "weight": 1.5,
    "status": "ready_for_pickup",
    "shelf": {
      "id": 1,
      "shelf_code": "A-01-01",
      "area": "A",
      "floor": 1,
      "column": 1
    },
    "received_at": "2025-01-15T10:30:00Z",
    "shelved_at": "2025-01-15T10:35:00Z",
    "expected_overdue_at": "2025-01-18T10:30:00Z",
    "notes": "易碎品"
  }
}
```

### 3.10 包裹退回
- **接口**: POST /api/v1/parcels/{id}/return
- **权限**: staff, admin
- **描述**: 将滞留包裹退回

#### 请求参数
```json
{
  "reason": "超期未取"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "退回成功",
  "data": {
    "id": 1001,
    "status": "returned",
    "returned_at": "2025-01-20T10:30:00Z"
  }
}
```

## 4. 货架模块接口

### 4.1 创建货架
- **接口**: POST /api/v1/shelves
- **权限**: admin
- **描述**: 创建新货架

#### 请求参数
```json
{
  "shelf_code": "A-01-01",
  "area": "A",
  "floor": 1,
  "column": 1,
  "type": "small",
  "capacity": 20
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 1,
    "shelf_code": "A-01-01",
    "area": "A",
    "floor": 1,
    "column": 1,
    "type": "small",
    "capacity": 20,
    "current_count": 0,
    "status": "available"
  }
}
```

### 4.2 货架列表
- **接口**: GET /api/v1/shelves
- **权限**: staff, admin
- **描述**: 查询货架列表

#### 查询参数
- `area`: 区域（可选）
- `type`: 类型（可选）
- `status`: 状态（可选）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "shelf_code": "A-01-01",
        "area": "A",
        "floor": 1,
        "column": 1,
        "type": "small",
        "capacity": 20,
        "current_count": 15,
        "utilization_rate": 75.0,
        "status": "available"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

### 4.3 货架详情
- **接口**: GET /api/v1/shelves/{id}
- **权限**: staff, admin
- **描述**: 查询货架详细信息，包括该货架上的包裹列表

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "shelf_code": "A-01-01",
    "area": "A",
    "floor": 1,
    "column": 1,
    "type": "small",
    "capacity": 20,
    "current_count": 15,
    "status": "available",
    "parcels": [
      {
        "id": 1001,
        "tracking_number": "SF1234567890",
        "pickup_code": "A12345",
        "recipient_name": "张三",
        "status": "ready_for_pickup"
      }
    ]
  }
}
```

### 4.4 更新货架
- **接口**: PUT /api/v1/shelves/{id}
- **权限**: admin
- **描述**: 更新货架信息

#### 请求参数
```json
{
  "capacity": 25,
  "status": "maintenance"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "shelf_code": "A-01-01",
    "capacity": 25,
    "status": "maintenance"
  }
}
```

### 4.5 删除货架
- **接口**: DELETE /api/v1/shelves/{id}
- **权限**: admin
- **描述**: 删除货架（软删除）

#### 响应示例
```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

### 4.6 货架统计
- **接口**: GET /api/v1/shelves/statistics
- **权限**: staff, admin
- **描述**: 获取货架使用统计信息

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_shelves": 100,
    "available_shelves": 85,
    "maintenance_shelves": 10,
    "disabled_shelves": 5,
    "total_capacity": 2000,
    "total_used": 1500,
    "overall_utilization_rate": 75.0,
    "by_type": [
      {
        "type": "small",
        "count": 50,
        "capacity": 1000,
        "used": 800,
        "utilization_rate": 80.0
      },
      {
        "type": "medium",
        "count": 30,
        "capacity": 600,
        "used": 450,
        "utilization_rate": 75.0
      },
      {
        "type": "large",
        "count": 20,
        "capacity": 400,
        "used": 250,
        "utilization_rate": 62.5
      }
    ]
  }
}
```

## 5. 寄件模块接口

### 5.1 创建寄件单
- **接口**: POST /api/v1/shipments
- **权限**: 需要认证
- **描述**: 创建寄件单，系统自动计算运费

#### 请求参数
```json
{
  "recipient_name": "李四",
  "recipient_phone": "13900139000",
  "recipient_address": "广东省广州市天河区XXX街道XXX号",
  "courier_company": "顺丰速运",
  "weight": 2.5,
  "volume": 0.01
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 1,
    "shipment_number": "SH202501150001",
    "sender_name": "张三",
    "sender_phone": "13800138000",
    "recipient_name": "李四",
    "recipient_phone": "13900139000",
    "recipient_address": "广东省广州市天河区XXX街道XXX号",
    "courier_company": "顺丰速运",
    "weight": 2.5,
    "volume": 0.01,
    "freight": 15.00,
    "status": "pending",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

### 5.2 寄件单列表
- **接口**: GET /api/v1/shipments/my
- **权限**: 需要认证
- **描述**: 查询我的寄件单列表

#### 查询参数
- `status`: 状态（可选）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "shipment_number": "SH202501150001",
        "recipient_name": "李四",
        "courier_company": "顺丰速运",
        "weight": 2.5,
        "freight": 15.00,
        "status": "shipped",
        "created_at": "2025-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 10,
      "total_pages": 1
    }
  }
}
```

### 5.3 寄件单详情
- **接口**: GET /api/v1/shipments/{id}
- **权限**: 需要认证（本人或管理员）
- **描述**: 查询寄件单详细信息

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "shipment_number": "SH202501150001",
    "sender_name": "张三",
    "sender_phone": "13800138000",
    "recipient_name": "李四",
    "recipient_phone": "13900139000",
    "recipient_address": "广东省广州市天河区XXX街道XXX号",
    "courier_company": "顺丰速运",
    "weight": 2.5,
    "volume": 0.01,
    "freight": 15.00,
    "status": "shipped",
    "shipped_at": "2025-01-15T11:00:00Z",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

### 5.4 运费试算
- **接口**: POST /api/v1/shipments/calculate-freight
- **权限**: 需要认证
- **描述**: 根据重量、体积和快递公司计算运费

#### 请求参数
```json
{
  "courier_company": "顺丰速运",
  "weight": 2.5,
  "volume": 0.01
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "freight": 15.00,
    "calculation_details": {
      "base_freight": 10.00,
      "weight_freight": 5.00,
      "volume_freight": 0.00
    }
  }
}
```

## 6. 滞留件模块接口

### 6.1 滞留件列表
- **接口**: GET /api/v1/overdue/parcels
- **权限**: staff, admin
- **描述**: 查询滞留包裹列表

#### 查询参数
- `status`: 处理状态（可选）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "parcel": {
          "id": 1001,
          "tracking_number": "SF1234567890",
          "pickup_code": "A12345",
          "recipient_name": "张三",
          "recipient_phone": "13800138000"
        },
        "overdue_start_at": "2025-01-18T10:30:00Z",
        "overdue_days": 3,
        "reminder_count": 2,
        "last_reminder_at": "2025-01-19T10:30:00Z",
        "overdue_fee": 3.00,
        "status": "reminded"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

### 6.2 发送滞留提醒
- **接口**: POST /api/v1/overdue/parcels/{id}/remind
- **权限**: staff, admin
- **描述**: 发送滞留提醒（短信/邮件）

#### 响应示例
```json
{
  "code": 0,
  "message": "提醒发送成功",
  "data": {
    "id": 1,
    "reminder_count": 3,
    "last_reminder_at": "2025-01-20T10:30:00Z"
  }
}
```

### 6.3 批量处理滞留件
- **接口**: POST /api/v1/overdue/parcels/batch-return
- **权限**: admin
- **描述**: 批量退回滞留包裹

#### 请求参数
```json
{
  "parcel_ids": [1001, 1002, 1003]
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "批量退回成功",
  "data": {
    "success_count": 3,
    "failed_count": 0
  }
}
```

## 7. 统计分析模块接口

### 7.1 仪表盘统计
- **接口**: GET /api/v1/statistics/dashboard
- **权限**: staff, admin
- **描述**: 获取仪表盘统计数据

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "today": {
      "received_count": 150,
      "picked_up_count": 120,
      "overdue_count": 5
    },
    "parcels": {
      "total": 1000,
      "received": 100,
      "shelved": 50,
      "ready_for_pickup": 300,
      "picked_up": 500,
      "overdue": 30,
      "returned": 20
    },
    "shelves": {
      "total_capacity": 2000,
      "total_used": 1500,
      "utilization_rate": 75.0
    },
    "recent_trends": [
      {
        "date": "2025-01-15",
        "received": 150,
        "picked_up": 120
      },
      {
        "date": "2025-01-14",
        "received": 140,
        "picked_up": 135
      }
    ]
  }
}
```

### 7.2 包裹统计
- **接口**: GET /api/v1/statistics/parcels
- **权限**: staff, admin
- **描述**: 获取包裹统计数据

#### 查询参数
- `start_date`: 开始日期
- `end_date`: 结束日期
- `group_by`: 分组方式（day/week/month）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_received": 1500,
    "total_picked_up": 1350,
    "pickup_rate": 90.0,
    "avg_pickup_time": "1.5天",
    "by_courier_company": [
      {
        "courier_company": "顺丰速运",
        "count": 500
      },
      {
        "courier_company": "中通快递",
        "count": 400
      }
    ],
    "by_size": [
      {
        "size": "small",
        "count": 800
      },
      {
        "size": "medium",
        "count": 500
      },
      {
        "size": "large",
        "count": 200
      }
    ]
  }
}
```

### 7.3 导出报表
- **接口**: GET /api/v1/statistics/export
- **权限**: admin
- **描述**: 导出统计报表（Excel/CSV）

#### 查询参数
- `report_type`: 报表类型（parcels/shipments/overdue）
- `start_date`: 开始日期
- `end_date`: 结束日期
- `format`: 导出格式（excel/csv）

#### 响应
返回文件流，浏览器自动下载

## 8. 系统管理模块接口

### 8.1 操作日志列表
- **接口**: GET /api/v1/logs/operations
- **权限**: admin
- **描述**: 查询操作日志

#### 查询参数
- `operator_id`: 操作人ID（可选）
- `operation_type`: 操作类型（可选）
- `start_date`: 开始日期（可选）
- `end_date`: 结束日期（可选）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "operator": {
          "id": 1,
          "username": "admin",
          "real_name": "管理员"
        },
        "operation_type": "parcel_received",
        "object_type": "parcel",
        "object_id": 1001,
        "description": "包裹入库：SF1234567890",
        "ip_address": "192.168.1.100",
        "operated_at": "2025-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 10000,
      "total_pages": 500
    }
  }
}
```

### 8.2 用户管理列表
- **接口**: GET /api/v1/admin/users
- **权限**: admin
- **描述**: 查询用户列表

#### 查询参数
- `role`: 角色（可选）
- `keyword`: 搜索关键词（用户名/手机号/姓名）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "zhangsan",
        "real_name": "张三",
        "phone": "13800138000",
        "role": "user",
        "created_at": "2025-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 1000,
      "total_pages": 50
    }
  }
}
```

### 8.3 更新用户角色
- **接口**: PUT /api/v1/admin/users/{id}/role
- **权限**: admin
- **描述**: 修改用户角色

#### 请求参数
```json
{
  "role": "staff"
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "角色更新成功",
  "data": {
    "id": 1,
    "username": "zhangsan",
    "role": "staff"
  }
}
```

### 8.4 系统配置
- **接口**: GET /api/v1/admin/config
- **权限**: admin
- **描述**: 获取系统配置

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "overdue_days": 3,
    "reminder_interval_hours": 24,
    "max_overdue_days": 7,
    "overdue_fee_per_day": 1.00
  }
}
```

### 8.5 更新系统配置
- **接口**: PUT /api/v1/admin/config
- **权限**: admin
- **描述**: 更新系统配置

#### 请求参数
```json
{
  "overdue_days": 3,
  "reminder_interval_hours": 24,
  "max_overdue_days": 7,
  "overdue_fee_per_day": 1.00
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "配置更新成功",
  "data": null
}
```

## 9. 接口汇总表

### 用户模块
| 接口 | 方法 | 权限 | 描述 |
|------|------|------|------|
| /api/v1/users/register | POST | 公开 | 用户注册 |
| /api/v1/users/login | POST | 公开 | 用户登录 |
| /api/v1/users/profile | GET | 认证 | 获取个人信息 |
| /api/v1/users/profile | PUT | 认证 | 更新个人信息 |
| /api/v1/users/change-password | POST | 认证 | 修改密码 |

### 包裹模块
| 接口 | 方法 | 权限 | 描述 |
|------|------|------|------|
| /api/v1/parcels | POST | staff/admin | 包裹入库 |
| /api/v1/parcels/batch | POST | staff/admin | 批量入库 |
| /api/v1/parcels/{id}/shelve | PUT | staff/admin | 包裹上架 |
| /api/v1/parcels/pickup | POST | 认证 | 包裹取件 |
| /api/v1/parcels/by-pickup-code/{code} | GET | 认证 | 查询包裹（取件码） |
| /api/v1/parcels/by-tracking-number/{number} | GET | 认证 | 查询包裹（单号） |
| /api/v1/parcels/my | GET | 认证 | 我的包裹 |
| /api/v1/parcels | GET | staff/admin | 包裹列表 |
| /api/v1/parcels/{id} | GET | 认证 | 包裹详情 |
| /api/v1/parcels/{id}/return | POST | staff/admin | 包裹退回 |

### 货架模块
| 接口 | 方法 | 权限 | 描述 |
|------|------|------|------|
| /api/v1/shelves | POST | admin | 创建货架 |
| /api/v1/shelves | GET | staff/admin | 货架列表 |
| /api/v1/shelves/{id} | GET | staff/admin | 货架详情 |
| /api/v1/shelves/{id} | PUT | admin | 更新货架 |
| /api/v1/shelves/{id} | DELETE | admin | 删除货架 |
| /api/v1/shelves/statistics | GET | staff/admin | 货架统计 |

### 寄件模块
| 接口 | 方法 | 权限 | 描述 |
|------|------|------|------|
| /api/v1/shipments | POST | 认证 | 创建寄件单 |
| /api/v1/shipments/my | GET | 认证 | 我的寄件单 |
| /api/v1/shipments/{id} | GET | 认证 | 寄件单详情 |
| /api/v1/shipments/calculate-freight | POST | 认证 | 运费试算 |

### 滞留件模块
| 接口 | 方法 | 权限 | 描述 |
|------|------|------|------|
| /api/v1/overdue/parcels | GET | staff/admin | 滞留件列表 |
| /api/v1/overdue/parcels/{id}/remind | POST | staff/admin | 发送提醒 |
| /api/v1/overdue/parcels/batch-return | POST | admin | 批量退回 |

### 统计分析模块
| 接口 | 方法 | 权限 | 描述 |
|------|------|------|------|
| /api/v1/statistics/dashboard | GET | staff/admin | 仪表盘统计 |
| /api/v1/statistics/parcels | GET | staff/admin | 包裹统计 |
| /api/v1/statistics/export | GET | admin | 导出报表 |

### 系统管理模块
| 接口 | 方法 | 权限 | 描述 |
|------|------|------|------|
| /api/v1/logs/operations | GET | admin | 操作日志 |
| /api/v1/admin/users | GET | admin | 用户列表 |
| /api/v1/admin/users/{id}/role | PUT | admin | 更新角色 |
| /api/v1/admin/config | GET | admin | 获取配置 |
| /api/v1/admin/config | PUT | admin | 更新配置 |

## 10. Swagger文档

系统使用Swaggo自动生成API文档，启动服务后访问：
```
http://localhost:8000/swagger/index.html
```

在Handler层添加Swagger注解示例：
```go
// GetParcelByPickupCode godoc
// @Summary 通过取件码查询包裹
// @Description 通过取件码查询包裹详细信息
// @Tags 包裹管理
// @Accept json
// @Produce json
// @Param pickup_code path string true "取件码"
// @Success 200 {object} Response{data=ParcelDetailResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Security BearerAuth
// @Router /api/v1/parcels/by-pickup-code/{pickup_code} [get]
func (h *ParcelHandler) GetParcelByPickupCode(ctx *gin.Context) {
    // 实现逻辑
}
```

更新Swagger文档：
```bash
cd backend
make swag
```
