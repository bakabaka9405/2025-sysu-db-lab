#import "../lib.typ": data-table, soft-header

= 数据库设计

数据库设计是本系统的核心部分。本章将详细阐述从概念设计到物理设计的全过程 @silberschatz2019 @wangshan2014，并重点展示PostgreSQL高级特性的深度应用 @postgresql17，包括触发器、存储过程、视图、高级索引策略、事务控制等，充分满足课程设计对数据库高级功能的要求。

== 概念结构设计

概念结构设计采用实体-联系(E-R)方法 @elmasri2021，通过抽象现实世界的业务实体及其关系，建立概念模型。

=== 主要实体

系统涉及的核心实体包括:

/ 用户(User): 表示系统使用者，包括普通用户、驿站员工和管理员。主要属性有：用户ID、用户名、密码、真实姓名、手机号、身份证号、邮箱、角色、创建时间等。

/ 包裹(Parcel): 表示快递包裹，是系统管理的核心对象。主要属性有：包裹ID、快递单号、取件码、收件人姓名、收件人手机号、快递公司、包裹大小、重量、状态、货架ID、入库时间、上架时间、取件时间、预计滞留时间等。

/ 货架(Shelf): 表示存放包裹的物理货架。主要属性有：货架ID、货架编号、区域、楼层、列号、类型、容量、当前占用数、状态等。

/ 寄件单(Shipment): 表示用户寄件信息。主要属性有：寄件ID、寄件单号、寄件人ID、寄件人姓名、收件人姓名、收件人地址、快递公司、重量、体积、运费、状态等。

/ 滞留记录(OverdueRecord): 表示包裹滞留的详细信息。主要属性有：记录ID、包裹ID、滞留开始时间、提醒次数、最后提醒时间、滞留费用、处理状态等。

/ 操作日志(OperationLog): 表示系统操作记录，用于审计。主要属性有：日志ID、操作人ID、操作类型、操作对象、描述、IP地址、操作时间等。

/ 快递价格(CourierPrice): 表示快递公司的价格配置。主要属性有：配置ID、快递公司、首重、首重价格、续重价格、是否启用等。

=== 实体 E-R 图

下面展示各个核心实体的 E-R 图表示：

==== 用户实体(User)

#figure(
  include "../assets/er-diagrams/user-entity.typ",
  caption: [用户实体 E-R 图],
)

==== 包裹实体(Parcel)

#figure(
  include "../assets/er-diagrams/parcel-entity.typ",
  caption: [包裹实体 E-R 图],
)

==== 货架实体(Shelf)

#figure(
  include "../assets/er-diagrams/shelf-entity.typ",
  caption: [货架实体 E-R 图],
)

==== 寄件实体(Shipment)

#figure(
  include "../assets/er-diagrams/shipment-entity.typ",
  caption: [寄件实体 E-R 图],
)

==== 滞留记录实体(OverdueRecord)

#figure(
  include "../assets/er-diagrams/overdue-entity.typ",
  caption: [滞留记录实体 E-R 图],
)

=== 实体间联系

系统中的实体通过多种联系类型相互关联，形成完整的数据模型。

用户与包裹之间为一对多关系，一个用户(收件人)可以有多个包裹，但每个包裹只对应一个收件人，联系通过收件人手机号或身份证号实现。货架与包裹之间同样是一对多关系，一个货架可以存放多个包裹，但每个包裹只能放在一个货架上，包裹表中的shelf_id外键引用货架表的主键。用户与寄件单之间为一对多关系，一个用户可以创建多个寄件单，但每个寄件单只属于一个用户，寄件单表中的sender_id外键引用用户表的主键。

包裹与滞留记录之间为一对一关系，一个包裹最多对应一条滞留记录，滞留记录表中的parcel_id字段设置唯一约束以保证这一约束。用户与操作日志之间为一对多关系，一个用户可以产生多条操作日志，操作日志表中的operator_id外键引用用户表的主键。

=== 整体 E-R 关系图

下图展示了系统中各实体间的完整关系：

#figure(
  include "../assets/er-diagrams/overall-relationship.typ",
  caption: [系统整体 E-R 关系图],
)

== 逻辑结构设计

逻辑结构设计将E-R图转换为关系模式，并进行规范化分析，确保满足至少第三范式(3NF)的要求。

=== 关系模式

系统的主要关系模式如下:

#soft-header([用户表(users)])

#data-table(
  caption: [用户表关系模式],
  columns: (auto, auto, 1fr),
  headers: ("中文名", "字段名", "描述"),
  ([用户ID], [id], [主键，自增长整型]),
  ([用户名], [username], [唯一约束，登录账号]),
  ([密码], [password], [加密存储]),
  ([真实姓名], [real_name], [用户真实姓名]),
  ([手机号], [phone], [唯一约束，用于身份验证和通知]),
  ([身份证号], [id_card], [唯一约束，实名认证]),
  ([邮箱], [email], [可选联系方式]),
  ([角色], [role], [用户角色：普通用户/员工/管理员]),
  ([创建时间], [created_at], [记录创建时间戳]),
  ([更新时间], [updated_at], [记录最后更新时间戳]),
  ([删除时间], [deleted_at], [软删除标记，NULL表示未删除]),
)

#soft-header([包裹表(parcels)])

#data-table(
  caption: [包裹表关系模式],
  columns: (auto, auto, 1fr),
  headers: ("中文名", "字段名", "描述"),
  ([包裹ID], [id], [主键，自增长整型]),
  ([快递单号], [tracking_number], [唯一约束，快递公司单号]),
  ([取件码], [pickup_code], [唯一约束，用户取件凭证]),
  ([收件人姓名], [recipient_name], [包裹收件人]),
  ([收件人手机], [recipient_phone], [收件人联系电话]),
  ([收件人身份证], [recipient_id_card], [实名认证标识]),
  ([快递公司], [courier_company], [承运快递公司名称]),
  ([包裹大小], [size], [小/中/大/特大]),
  ([重量], [weight], [包裹重量(kg)]),
  ([状态], [status], [待上架/已上架/已取件/已滞留等]),
  ([货架ID], [shelf_id], [外键，关联货架表]),
  ([入库时间], [received_at], [包裹到达驿站时间]),
  ([上架时间], [shelved_at], [包裹放入货架时间]),
  ([取件时间], [picked_up_at], [用户取件时间]),
  ([预计滞留时间], [expected_overdue_at], [超过此时间视为滞留]),
  ([备注], [notes], [包裹备注信息]),
  ([创建时间], [created_at], [记录创建时间戳]),
  ([更新时间], [updated_at], [记录最后更新时间戳]),
  ([删除时间], [deleted_at], [软删除标记]),
)

#soft-header([货架表(shelves)])

#data-table(
  caption: [货架表关系模式],
  columns: (auto, auto, 1fr),
  headers: ("中文名", "字段名", "描述"),
  ([货架ID], [id], [主键，自增长整型]),
  ([货架编号], [shelf_code], [唯一约束，人工编号]),
  ([区域], [area], [货架所在区域标识]),
  ([楼层], [floor], [货架所在楼层]),
  ([列号], [column], [货架列号]),
  ([类型], [type], [普通/大件/特殊等]),
  ([容量], [capacity], [货架最大存放数量]),
  ([当前占用], [current_count], [当前存放包裹数量]),
  ([状态], [status], [启用/停用]),
  ([创建时间], [created_at], [记录创建时间戳]),
  ([更新时间], [updated_at], [记录最后更新时间戳]),
  ([删除时间], [deleted_at], [软删除标记]),
)

#soft-header([快递价格配置表(courier_prices)])

#data-table(
  caption: [快递价格配置表关系模式],
  columns: (auto, auto, 1fr),
  headers: ("中文名", "字段名", "描述"),
  ([配置ID], [id], [主键，自增长整型]),
  ([快递公司], [courier_company], [唯一约束，快递公司名称]),
  ([首重], [first_weight], [首重重量(kg)]),
  ([首重价格], [first_price], [首重运费(元)]),
  ([续重价格], [additional_price], [每续重单价(元/kg)]),
  ([是否启用], [is_active], [价格配置启用状态]),
  ([创建时间], [created_at], [记录创建时间戳]),
  ([更新时间], [updated_at], [记录最后更新时间戳]),
)

#soft-header([寄件表(shipments)])

#data-table(
  caption: [寄件表关系模式],
  columns: (auto, auto, 1fr),
  headers: ("中文名", "字段名", "描述"),
  ([寄件ID], [id], [主键，自增长整型]),
  ([寄件单号], [shipment_number], [唯一约束，系统生成单号]),
  ([寄件人ID], [sender_id], [外键，关联用户表]),
  ([寄件人姓名], [sender_name], [寄件人姓名]),
  ([寄件人手机], [sender_phone], [寄件人联系电话]),
  ([收件人姓名], [recipient_name], [收件人姓名]),
  ([收件人手机], [recipient_phone], [收件人联系电话]),
  ([收件地址], [recipient_address], [收件人详细地址]),
  ([快递公司], [courier_company], [选择的承运公司]),
  ([重量], [weight], [包裹重量(kg)]),
  ([体积], [volume], [包裹体积(cm³)]),
  ([运费], [freight], [计算后的运费(元)]),
  ([状态], [status], [待揽收/已揽收/运输中/已签收等]),
  ([发出时间], [shipped_at], [快递揽收时间]),
  ([创建时间], [created_at], [记录创建时间戳]),
  ([更新时间], [updated_at], [记录最后更新时间戳]),
  ([删除时间], [deleted_at], [软删除标记]),
)

#soft-header([滞留记录表(overdue_records)])

#data-table(
  caption: [滞留记录表关系模式],
  columns: (auto, auto, 1fr),
  headers: ("中文名", "字段名", "描述"),
  ([记录ID], [id], [主键，自增长整型]),
  ([包裹ID], [parcel_id], [外键，关联包裹表，唯一约束]),
  ([滞留开始时间], [overdue_start_at], [开始计算滞留的时间]),
  ([提醒次数], [reminder_count], [已发送提醒次数]),
  ([最后提醒时间], [last_reminder_at], [最近一次发送提醒时间]),
  ([滞留费用], [overdue_fee], [累计滞留费用(元)]),
  ([状态], [status], [滞留中/已处理/已转存等]),
  ([创建时间], [created_at], [记录创建时间戳]),
  ([更新时间], [updated_at], [记录最后更新时间戳]),
  ([删除时间], [deleted_at], [软删除标记]),
)

#soft-header([操作日志表(operation_logs) - 分区表])

#data-table(
  caption: [操作日志表关系模式（按operated_at范围分区）],
  columns: (auto, auto, 1fr),
  headers: ("中文名", "字段名", "描述"),
  ([日志ID], [id], [主键(与operated_at构成复合主键)]),
  ([操作人ID], [operator_id], [外键，关联用户表]),
  ([操作类型], [operation_type], [创建/更新/删除/查询等]),
  ([对象类型], [object_type], [操作对象类型(包裹/货架/用户等)]),
  ([对象ID], [object_id], [操作对象的ID]),
  ([描述], [description], [操作详细描述]),
  ([IP地址], [ip_address], [操作者IP地址(INET类型)]),
  ([操作时间], [operated_at], [操作发生时间(分区键)]),
  ([创建时间], [created_at], [记录创建时间戳]),
)

=== 规范化分析

对关系模式进行规范化检验:

*第一范式(1NF)*：所有关系模式的属性均为原子性，不存在复合属性或多值属性。例如，地址字段虽然在现实中可以分解为省、市、区等，但在本系统中作为一个原子属性存储，满足业务需求。

*第二范式(2NF)*：所有非主属性完全依赖于主键，不存在部分依赖。例如，在parcels表中，所有属性(recipient_name, status等)都完全依赖于主键id，不存在只依赖于主键一部分的情况。

*第三范式(3NF)*：所有非主属性不传递依赖于主键。例如，在parcels表中，虽然shelf_id引用shelves表，但货架的详细信息(area, floor等)存储在shelves表中，而不是在parcels表中冗余存储，避免了传递依赖。

*BC范式(BCNF)*：大部分表满足BCNF，即每个决定因素都是候选键。部分表由于业务需要保留了少量冗余(如parcels表中的recipient_name和recipient_phone，即使可以通过users表关联获取，但为了查询效率和数据独立性仍然保留)。

== 物理结构设计

物理结构设计将逻辑模式映射到具体的数据库管理系统(PostgreSQL)，定义表结构、数据类型、约束条件和索引策略。

=== 核心表结构设计

下面展示几个核心表的完整DDL定义，突出完整性约束的应用。

*用户表*
#raw(read("../../database/01-tables/01-users.sql"), lang: "sql", block: true)

*包裹表*
#raw(read("../../database/01-tables/03-parcels.sql"), lang: "sql", block: true)

*货架表*
#raw(read("../../database/01-tables/02-shelves.sql"), lang: "sql", block: true)

=== 数据类型选择

数据类型的选择遵循精确性和效率并重的原则。整数类型方面，主键使用BIGSERIAL(自增长整型)以支持大规模数据，普通整数字段使用INT。字符串类型方面，定长短字符串使用VARCHAR，长文本使用TEXT。数值类型方面，金额和重量使用DECIMAL(10,2)以保证计算精度。时间类型统一使用TIMESTAMPTZ(带时区的时间戳)，便于跨时区应用。此外，IP地址使用PostgreSQL原生的INET类型，可支持IP范围查询等高级功能。

== PostgreSQL高级特性应用

本节重点展示PostgreSQL数据库高级特性的应用，这是课程设计的核心要求。系统实现了5个触发器、6个存储过程/函数、5个视图、15+种高级索引以及完善的事务控制机制。

=== 触发器(Triggers)

触发器用于在特定数据库事件发生时自动执行预定义的操作，可实现业务逻辑自动化和数据一致性维护。

==== 触发器1：自动更新时间戳

该触发器在记录更新时自动维护updated_at字段，确保时间戳的准确性。

#raw(read("../../database/03-triggers/01-update-timestamp.sql"), lang: "sql", block: true)

==== 触发器2：货架计数自动维护

这是系统最核心的触发器之一，在包裹上架、取件或删除时自动维护货架的占用计数，保证数据一致性。

#raw(read("../../database/03-triggers/02-maintain-shelf-count.sql"), lang: "sql", block: true)

*设计说明*：该触发器通过INSERT、UPDATE事件的统一处理，确保货架计数始终准确。使用AFTER触发器是因为需要在主记录成功写入后再更新关联表，避免事务冲突。

==== 触发器3：状态变更自动记录

该触发器在包裹状态改变时自动记录操作日志，实现审计功能。

#raw(read("../../database/03-triggers/03-log-status-change.sql"), lang: "sql", block: true)

*设计说明*：通过PostgreSQL的session variable机制(current_setting)，应用层可以设置当前操作用户，触发器自动记录是谁执行了状态变更操作。

==== 触发器4：自动计算预计滞留时间

包裹入库时自动计算预计滞留时间，为后续滞留管理提供依据。

#raw(read("../../database/03-triggers/04-set-overdue-time.sql"), lang: "sql", block: true)

==== 触发器5：货架容量限制检查

该触发器在包裹分配货架前检查容量，防止超载。

#raw(read("../../database/03-triggers/05-check-shelf-capacity.sql"), lang: "sql", block: true)

=== 存储过程和函数(Stored Procedures & Functions)

存储过程和函数将复杂业务逻辑封装在数据库层，提高性能并保证逻辑一致性。

==== 函数1：运费计算函数

封装运费计算规则，支持首重续重、体积重量、附加费用等复杂逻辑。

#raw(read("../../database/04-functions/01-calculate-freight.sql"), lang: "sql", block: true)

*设计说明*：该函数体现了“业务逻辑下沉”的设计理念，将运费计算规则统一封装在数据库中，避免应用层重复实现，保证了计算逻辑的一致性。

==== 函数2：货架分配函数

实现智能货架分配算法，综合考虑空间、位置和负载等因素。

#raw(read("../../database/04-functions/02-allocate-shelf.sql"), lang: "sql", block: true)

*设计说明*：评分算法中，空间因子权重0.5(优先剩余空间大的货架)，位置因子权重0.3(优先低楼层和前列)，负载因子权重0.2(均衡各货架负载)。

==== 函数3：包裹入库存储过程

封装完整的包裹入库流程，包括货架分配、取件码生成、状态初始化等。

#raw(read("../../database/04-functions/03-receive-parcel.sql"), lang: "sql", block: true)

*设计说明*：存储过程将多个步骤封装为原子操作，减少了应用层与数据库的交互次数，提高了性能并保证了数据一致性。

=== 视图(Views)

视图提供了数据的逻辑抽象，简化复杂查询，实现数据安全和权限控制。

==== 视图1：包裹统计视图

按日期、状态、快递公司聚合包裹统计数据。

#raw(read("../../database/05-views/01-v-parcel-statistics.sql"), lang: "sql", block: true)

*应用场景*：运营仪表盘、业务报表生成。通过视图隐藏复杂的聚合逻辑，应用层只需简单查询即可获得统计结果。

==== 视图2：用户包裹视图(数据脱敏)

为用户查询包裹提供脱敏后的安全视图。

#raw(read("../../database/05-views/02-v-user-parcels.sql"), lang: "sql", block: true)

*设计说明*：视图实现了数据脱敏，保护用户隐私。同时计算剩余时间等派生字段，简化应用层逻辑。

==== 视图3：货架利用率视图

实时展示货架使用情况和负载水平。

#raw(read("../../database/05-views/03-v-shelf-utilization.sql"), lang: "sql", block: true)

=== 高级索引策略

索引是提升查询性能的关键。系统采用了多种索引类型，针对不同查询场景进行优化。

==== 部分索引(Partial Index)

只索引满足特定条件的行，节省存储空间并提高查询效率。

#raw(read("../../database/02-indexes/01-partial-indexes.sql"), lang: "sql", block: true)

*优势*：相比全表索引，部分索引显著减少了索引大小，提高了更新性能，同时保持了查询效率。

==== 复合索引(Composite Index)

支持多条件查询和排序。

#raw(read("../../database/02-indexes/02-composite-indexes.sql"), lang: "sql", block: true)

*设计说明*：索引字段的顺序遵循"等值在前，范围在后，排序字段最后"的原则。

==== 表达式索引(Expression Index)

支持基于函数或表达式的查询。

#raw(read("../../database/02-indexes/03-expression-indexes.sql"), lang: "sql", block: true)

==== GIN索引(全文搜索)

用于包裹备注的全文搜索。

#raw(read("../../database/02-indexes/04-gin-fulltext.sql"), lang: "sql", block: true)

*设计说明*:GIN索引特别适合全文搜索、数组和JSONB类型的查询，虽然索引构建较慢，但查询性能优异。

==== BRIN索引(时间序列数据)

用于操作日志的时间范围查询。

#raw(read("../../database/02-indexes/05-brin-timeseries.sql"), lang: "sql", block: true)

*优势*:BRIN索引占用空间极小(相比B-Tree可减少99%)，非常适合按时间顺序插入的日志数据。

=== 事务控制和并发处理

事务控制保证数据的ACID特性，是数据库应用的核心机制。

==== 事务隔离级别

系统根据不同业务场景选择合适的隔离级别:

#raw(read("../../database/07-transactions/01-isolation-level.sql"), lang: "sql", block: true)

==== 行级锁(FOR UPDATE)

防止并发冲突，确保数据一致性。

#raw(read("../../database/07-transactions/02-row-level-lock.sql"), lang: "sql", block: true)

*应用场景*：在包裹取件存储过程中，使用FOR UPDATE锁定包裹记录，防止重复取件。

==== 分区表(Partitioning)

操作日志表按月分区，优化历史数据查询和维护。

#raw(read("../../database/06-partitions/01-operation-logs-partition.sql"), lang: "sql", block: true)

*优势*：分区表支持按分区删除历史数据、并行查询、分区独立备份等高级功能。

== 数据库设计总结

=== 高级特性应用统计

本系统充分应用了PostgreSQL的高级特性，统计如下:

#data-table(
  columns: (auto, auto, auto),
  align: (left, center, left),
  headers: ("特性类型", "数量", "应用场景"),
  ([触发器], [5个], [自动维护货架计数、状态日志、时间计算、容量检查]),
  ([存储过程/函数], [6个], [运费计算、货架分配、包裹入库/取件]),
  ([视图], [5个], [统计分析、数据脱敏、查询简化]),
  ([高级索引], [15+], [部分、复合、表达式、GIN、BRIN索引]),
  ([事务控制], [-], [隔离级别、行级锁、保存点]),
  ([分区表], [1个], [日志表按月分区]),
)

=== 设计亮点

本系统数据库设计的核心亮点体现在以下几个方面。

首先是业务逻辑下沉的设计理念，核心业务逻辑通过存储过程实现(如货架分配、包裹入库)，有效减少了应用层复杂度，并保证了逻辑一致性。其次是数据一致性保证机制，通过触发器自动维护冗余字段(如货架计数)，避免了手工更新可能导致的数据不一致问题。

在查询性能优化方面，系统设计了15+种索引策略针对不同查询场景进行优化，其中部分索引和表达式索引的运用显著提升了查询效率。数据安全方面，视图机制实现了敏感数据脱敏，有效保护用户隐私；操作日志记录关键操作，支持完整的审计追溯。

并发安全方面，系统使用行级锁和事务隔离级别控制，保证了高并发场景下的数据一致性。可扩展性方面，分区表设计支持海量历史数据管理，BRIN索引有效降低了存储成本。

本章充分展示了数据库系统在实际应用中的强大功能，通过合理运用触发器、存储过程、视图、索引等高级特性，不仅提升了系统性能，更保证了数据的完整性、一致性和安全性，达到了课程设计对数据库高级特性应用的要求。
