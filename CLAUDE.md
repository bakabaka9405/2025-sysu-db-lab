# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码仓库中工作时提供指导。

## 项目概述

这是一个数据库实验项目（2025-sysu-db-lab），由 Go 后端和 Vue.js 前端组成。后端使用 Nunu 框架构建，遵循清晰的架构原则，通过 Google Wire 进行依赖注入。

项目需求位于同目录下的 TARGET.md，每次执行任务前，请仔细阅读这篇文档。

项目设计放在 design/ 文件夹下，其中已经生成了你写的系统设计文档，实际开发时，请严格按照其中的设计，如果我在之后对项目结构提出修改，请你也修改报告文档保持两端一致。

已经完成的内容、用于讨论和评审的计划、方案等文档，写到项目的 discuss/ 目录下

项目报告放在 report/ 下，使用 typst 编写

## 开发要求

- 尽可能确保每个代码文件不要超过 300 行
- 每层文件夹中的文件，尽可能不超过 8 个。如有超过，需要规划为多层子文件夹
- 除了硬性指标以外，还需要时刻关注优雅的架构设计，避免出现以下可能侵蚀我们代码质量的「坏味道」：
 （1）僵化 (Rigidity): 系统难以变更，任何微小的改动都会引发一连串的连锁修改。
 （2）冗余 (Redundancy): 同样的代码逻辑在多处重复出现，导致维护困难且容易产生不一致。
 （3）循环依赖 (Circular Dependency): 两个或多个模块互相纠缠，形成无法解耦的“死结”，导致难以测试与复用。
 （4）脆弱性 (Fragility): 对代码一处的修改，导致了系统中其他看似无关部分功能的意外损坏。
 （5）晦涩性 (Obscurity): 代码意图不明，结构混乱，导致阅读者难以理解其功能和设计。
 （6）数据泥团 (Data Clump): 多个数据项总是一起出现在不同方法的参数中，暗示着它们应该被组合成一个独立的对象。
 （7）不必要的复杂性 (Needless Complexity): 用“杀牛刀”去解决“杀鸡”的问题，过度设计使系统变得臃肿且难以理解。
- 【非常重要！！】无论是你自己编写代码，还是阅读或审核他人代码时，都要严格遵守上述硬性指标，以及时刻关注优雅的架构设计。
- 【非常重要！！】无论何时，一旦你识别出那些可能侵蚀我们代码质量的「坏味道」，都应当立即询问用户是否需要优化，并给出合理的优化建议。


## 后端架构（Go）

后端采用分层架构模式：

- **Handler 层** (`internal/handler/`): HTTP 请求处理器，将 HTTP 请求转换为服务调用
- **Service 层** (`internal/service/`): 业务逻辑实现
- **Repository 层** (`internal/repository/`): 数据访问层，抽象数据库操作
- **Model 层** (`internal/model/`): 领域实体和数据结构

### 后端核心技术

- **框架**: Gin (HTTP 服务器)
- **ORM**: GORM，支持 MySQL、PostgreSQL 和 SQLite
- **依赖注入**: Google Wire
- **JWT**: golang-jwt 用于身份验证
- **日志**: Uber Zap
- **API 文档**: Swagger/Swaggo
- **测试**: Golang Mock、SQL Mock、httpexpect

### 后端目录结构

```
backend/
├── cmd/              # 应用程序入口
│   ├── server/       # HTTP 服务器（主应用）
│   ├── migration/    # 数据库迁移
│   └── task/         # 后台任务
├── internal/         # 私有应用代码
│   ├── handler/      # HTTP 处理器
│   ├── service/      # 业务逻辑
│   ├── repository/   # 数据访问
│   ├── model/        # 领域模型
│   ├── router/       # 路由定义
│   ├── server/       # 服务器设置
│   ├── middleware/   # HTTP 中间件
│   └── job/          # 定时任务
├── api/v1/           # API 请求/响应定义
├── pkg/              # 可复用包
├── config/           # 配置文件（local.yml, prod.yml）
├── storage/          # SQLite 数据库和日志
├── deploy/           # Docker 和部署配置
├── docs/             # Swagger 文档
└── test/             # 测试文件和 mocks
```

### 使用 Wire 进行依赖注入

项目使用 Google Wire 进行编译时依赖注入。Wire 定义位于 `cmd/*/wire/wire.go`。修改 Wire 文件后，需要重新生成：

```bash
wire gen ./cmd/server/wire
```

每一层都被组织成 Wire 集合：
- `repositorySet`: 数据库连接和仓储
- `serviceSet`: 服务实现
- `handlerSet`: HTTP 处理器
- `jobSet`: 定时任务
- `serverSet`: HTTP 服务器和任务服务器

### 事务支持

Repository 层通过 `Transaction` 接口提供事务支持。在服务层使用 `s.tm.Transaction(ctx, func(ctx context.Context) error {...})` 将操作包装在数据库事务中。

## 前端架构（Vue.js）

前端是使用以下技术构建的 Vue 3 应用：

- **构建工具**: Farm（Vite 的替代品）
- **框架**: Vue 3 with TypeScript
- **样式**: Scoped CSS

前端结构简洁：
```
frontend/
├── src/
│   ├── components/   # Vue 组件
│   ├── assets/       # 静态资源
│   ├── index.ts      # 入口文件
│   ├── App.vue       # 根组件
│   └── style.css     # 全局样式
├── public/           # 公共资源
└── farm.config.ts    # 构建配置
```

## 常用开发命令

### 后端命令

所有后端命令都应在 `backend/` 目录下运行：

```bash
# 安装必需的工具（Wire、MockGen、Swag）
make init

# 启动开发环境（Docker Compose + 迁移 + 服务器）
make bootstrap

# 直接运行服务器（数据库准备好后）
nunu run ./cmd/server

# 构建生产环境二进制文件
make build

# 生成测试用的 mocks
make mock

# 运行带覆盖率的测试
make test

# 生成 Swagger 文档
make swag

# 运行单个测试
go test -v ./test/server/user_test.go
```

### 前端命令

所有前端命令都应在 `frontend/` 目录下运行。项目使用 `bun` 作为包管理器：

```bash
# 安装依赖
bun install

# 启动开发服务器
bun start
# 或者
bun dev

# 构建生产版本
bun build

# 预览生产构建
bun preview

# 清除缓存
bun clean
```

## 配置管理

后端配置通过 Viper 和 `config/` 目录中的 YAML 文件管理：
- `local.yml`: 本地开发配置（默认使用 SQLite）
- `prod.yml`: 生产环境配置

可以通过命令行指定配置文件：
```bash
go run ./cmd/server/main.go -conf ./config/local.yml
```

主要配置项：
- `http`: 服务器主机和端口
- `data.db`: 数据库驱动（mysql/postgres/sqlite）和 DSN
- `data.redis`: Redis 连接（当前在 wire 中被禁用）
- `data.mongo`: MongoDB 连接（当前在 wire 中被禁用）
- `security.jwt`: JWT 签名密钥
- `log`: 日志配置

## 数据库支持

Repository 层通过 GORM 支持多种数据库：
- **SQLite**: 本地开发的默认选项
- **MySQL**: 取消配置中的 MySQL 注释，并在 wire 中启用 repository.NewRedis
- **PostgreSQL**: 在配置中更改驱动

数据库初始化在 `repository.NewDB()` 中进行。连接池配置为最多 10 个空闲连接，100 个打开连接。

## 测试策略

测试位于 `test/` 目录，mocks 位于 `test/mocks/`：
- 使用 `mockgen` 从接口生成 mocks
- 使用 `sqlmock` 进行数据库测试
- 使用 `httpexpect` 进行 HTTP API 测试
- 测试覆盖 handler、service 和 repository 层

运行特定测试套件：
```bash
go test -v ./test/server/user_test.go
```

## API 文档

API 文档使用 Swaggo 注解自动生成。更新 handler 注释后：

```bash
make swag
```

服务器运行时，可以在以下地址查看文档：`http://localhost:8000/swagger/index.html`

## 身份验证

项目使用 JWT bearer token 身份验证：
- 注册/登录端点是公开的
- 受保护的端点需要 `Authorization: Bearer <token>` 请求头
- JWT 中间件从 token 中提取用户 ID 并存储在 Gin 上下文中
- 在 handler 中访问用户 ID：`GetUserIdFromCtx(ctx)`
- Token 默认有效期为 90 天

## 添加新功能

按照架构添加新功能的步骤：

1. **定义 API 契约**：在 `api/v1/` 中定义请求/响应结构体
2. **创建模型**：如果需要，在 `internal/model/` 中创建
3. **添加 repository**：在 `internal/repository/` 中添加接口和实现
4. **添加 service**：在 `internal/service/` 中添加接口和实现
5. **添加 handler**：在 `internal/handler/` 中添加
6. **注册路由**：在 `internal/router/` 中注册
7. **更新 wire 集合**：在 `cmd/server/wire/wire.go` 中更新
8. **生成 wire**：`wire gen ./cmd/server/wire`
9. **生成 mocks**：如果需要，添加到 Makefile 并运行 `make mock`
10. **编写测试**：在 `test/` 中编写
11. **更新 Swagger**：添加注释并运行 `make swag`

## 开发工作流

1. 确保 Docker 运行以支持外部依赖（Redis、数据库）
2. 首次设置使用 `make bootstrap`
3. 开发期间，运行 `nunu run ./cmd/server` 以获得热重载
4. 前端开发：在 frontend 目录运行 `bun dev`
5. 接口更改后生成 mocks
6. 提交前运行测试
7. API 更改后更新 Swagger 文档

## 注意事项

- Redis 和 MongoDB 在 repository 层中有定义，但目前在 Wire 中被注释掉了
- 本地开发默认使用 SQLite（存储在 `backend/storage/` 中）
- 项目通过 gocron 支持定时任务（参见 `internal/job/`）
- 后台任务可以在 `cmd/task/` 中实现
- 默认服务器运行在 `http://127.0.0.1:8000`

## 数据库脚本
- 放在 database/ 文件夹内，注意分类存放并加上数字前缀，表示执行顺序