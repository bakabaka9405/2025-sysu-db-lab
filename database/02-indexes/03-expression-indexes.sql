-- 表达式索引 (Expression Indexes)
-- 支持基于函数或表达式的查询

-- 支持不区分大小写的用户名查询
CREATE INDEX idx_users_username_lower
ON users(LOWER(username)) WHERE deleted_at IS NULL;

-- 支持按日期范围查询包裹
CREATE INDEX idx_parcels_received_date
ON parcels(DATE(received_at)) WHERE deleted_at IS NULL;

-- 支持按日期查询寄件单
CREATE INDEX idx_shipments_created_date
ON shipments(DATE(created_at)) WHERE deleted_at IS NULL;

-- 支持按月统计查询
CREATE INDEX idx_parcels_received_month
ON parcels(DATE_TRUNC('month', received_at)) WHERE deleted_at IS NULL;
