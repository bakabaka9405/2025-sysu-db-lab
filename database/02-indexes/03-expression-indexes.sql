-- 表达式索引 (Expression Indexes)
-- 支持基于函数或表达式的查询

CREATE OR REPLACE FUNCTION to_date_immutable(timestamptz)
RETURNS date AS $$
    SELECT ($1 AT TIME ZONE 'UTC')::date;
$$ LANGUAGE sql IMMUTABLE;

CREATE OR REPLACE FUNCTION date_trunc_immutable(text, timestamptz)
RETURNS timestamptz AS $$
    SELECT date_trunc($1, $2 AT TIME ZONE 'UTC') AT TIME ZONE 'UTC';
$$ LANGUAGE sql IMMUTABLE;

-- 支持不区分大小写的用户名查询
CREATE INDEX idx_users_username_lower
ON users(LOWER(username)) WHERE deleted_at IS NULL;

-- 支持按日期范围查询包裹
CREATE INDEX idx_parcels_received_date
ON parcels(to_date_immutable(received_at)) WHERE deleted_at IS NULL;

-- 支持按日期查询寄件单
CREATE INDEX idx_shipments_created_date
ON shipments(to_date_immutable(created_at)) WHERE deleted_at IS NULL;

-- 支持按月统计查询
CREATE INDEX idx_parcels_received_month
ON parcels(date_trunc_immutable('month', received_at)) WHERE deleted_at IS NULL;
