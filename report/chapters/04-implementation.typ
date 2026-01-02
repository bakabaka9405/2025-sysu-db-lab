= 系统详细设计与实现

本章介绍系统架构的实现方式、核心功能模块的设计思路以及关键代码示例，重点展示如何在应用层调用数据库高级特性。系统架构设计参考了企业应用架构模式 @fowler2002 和微服务设计理念 @newman2021。

== 系统架构实现

系统采用经典的三层架构模式，包括Handler层负责HTTP处理、Service层负责业务逻辑、Repository层负责数据访问，配合依赖注入实现松耦合设计。

=== 分层架构

/ Handler层: 负责接收HTTP请求，参数验证和绑定，调用Service层处理业务，格式化响应。每个Handler对应一组相关的API端点。

/ Service层: 实现核心业务逻辑，协调多个Repository进行数据操作，管理事务边界。Service层是系统的业务中枢。

/ Repository层: 封装数据库操作，提供统一的数据访问接口。通过GORM调用SQL或存储过程，隐藏数据持久化细节。

/ Model层: 定义领域实体和数据传输对象即DTO，与数据库表结构映射。

=== 依赖注入设计

使用Google Wire进行编译时依赖注入，通过代码生成避免运行时反射开销。Wire将系统组件组织为多个Provider集合：repositorySet提供数据库连接和所有Repository实例，serviceSet提供所有Service实例并依赖于Repository，handlerSet提供所有Handler实例并依赖于Service，serverSet提供HTTP服务器配置。Wire在编译时生成依赖注入代码，确保依赖关系正确且类型安全。

== 核心功能模块实现

=== 包裹入库与上架模块

包裹入库是系统的核心功能之一，涉及货架分配、取件码生成等复杂逻辑。业务流程包括以下步骤：首先接收包裹基本信息如快递单号、收件人、包裹大小等；然后调用receive_parcel存储过程生成唯一取件码并插入包裹记录，状态设为"已入库"，暂不分配货架；返回包裹ID和取件码后，包裹可暂存于临时区域。随后工作人员选择已入库的包裹执行上架操作，调用shelve_parcel存储过程分配最优货架并更新包裹状态为"待取件"，触发器自动维护货架计数和时间戳；最后发送取件通知给用户。

以下是Repository层调用存储过程的关键代码：

#raw(read("../assets/go/receive-parcel-repository.go"), lang: "go", block: true)

=== 包裹取件模块

取件模块需要处理并发场景，防止重复取件。系统使用数据库行级锁即FOR UPDATE和事务控制，确保并发安全：

#raw(read("../assets/go/pickup-parcel-repository.go"), lang: "go", block: true)

=== 货架分配算法

货架分配算法封装在数据库函数allocate_shelf中，综合考虑多个因子。算法采用多因子评分机制，其中空间因子权重为0.5，优先选择剩余空间大的货架；位置因子权重为0.3，优先选择低楼层和前列位置，便于取件；负载因子权重为0.2，用于平衡各货架负载，避免某个货架过载。算法的SQL实现已在第3章展示，应用层只需调用该函数即可获得最优货架ID。

=== 运费计算模块

运费计算规则复杂，涉及首重续重、体积重量、附加费用等。计算流程如下：首先获取快递公司的价格配置包括首重、首重价格和续重价格；然后计算体积重量，取实际重量和体积重量的较大值作为计费重量；接着计算基础运费，若计费重量小于等于首重则运费等于首重价格，否则运费等于首重价格加上续重部分乘以续重价格；最后可选添加保价费、签收回单费等附加费用。

以下是调用数据库函数的代码：

#raw(read("../assets/go/calculate-freight-repository.go"), lang: "go", block: true)

运费计算逻辑统一封装在数据库中，避免应用层重复实现，且便于后续调整价格策略。

=== 滞留件处理

滞留件处理通过定时任务即cron和数据库视图、函数联动实现自动化。自动识别机制包括：定时任务每天凌晨扫描视图v_pending_overdue_parcels；视图自动计算滞留天数和费用，筛选出超期未取的包裹；更新包裹状态为"滞留"并创建滞留记录；调用第三方服务发送提醒通知。

以下是定时任务的代码：

#raw(read("../assets/go/overdue-job.go"), lang: "go", block: true)

通过视图简化查询逻辑，定时任务只需关注业务流程，无需关心复杂的SQL计算。

== API接口设计

系统采用RESTful API设计规范，提供统一的接口风格。

=== RESTful规范

系统使用标准HTTP方法表达操作语义：GET用于查询、POST用于创建、PUT用于更新、DELETE用于删除。资源采用复数名词命名，如`/api/v1/parcels`和`/api/v1/shelves`。HTTP状态码用于表示请求结果，其中200表示成功、400表示参数错误、401表示未认证、403表示无权限、404表示资源不存在等。

=== 统一响应格式

#raw(read("../assets/api-examples/response-success.json"), lang: "json", block: true)

错误响应示例:
#raw(read("../assets/api-examples/response-error.json"), lang: "json", block: true)

=== 核心接口示例

*包裹入库接口*
#raw(read("../assets/api-examples/receive-parcel-request.txt"), block: true)

*包裹取件接口*
#raw(read("../assets/api-examples/pickup-parcel-request.txt"), block: true)
