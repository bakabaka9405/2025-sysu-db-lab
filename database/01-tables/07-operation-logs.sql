-- 操作日志表（分区表）
-- 记录系统操作日志，按月分区以优化查询性能

CREATE TABLE operation_logs (
    id BIGSERIAL,
    operator_id BIGINT REFERENCES users(id),
    operation_type VARCHAR(50) NOT NULL,
    object_type VARCHAR(50),
    object_id BIGINT,
    description TEXT,
    ip_address INET,
    operated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, operated_at)
) PARTITION BY RANGE (operated_at);

-- 创建2025年的分区表
CREATE TABLE operation_logs_2025_01 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE operation_logs_2025_02 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

CREATE TABLE operation_logs_2025_03 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');

CREATE TABLE operation_logs_2025_04 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-04-01') TO ('2025-05-01');

CREATE TABLE operation_logs_2025_05 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-05-01') TO ('2025-06-01');

CREATE TABLE operation_logs_2025_06 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-06-01') TO ('2025-07-01');

CREATE TABLE operation_logs_2025_07 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');

CREATE TABLE operation_logs_2025_08 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');

CREATE TABLE operation_logs_2025_09 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');

CREATE TABLE operation_logs_2025_10 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE operation_logs_2025_11 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');

CREATE TABLE operation_logs_2025_12 PARTITION OF operation_logs
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

CREATE TABLE operation_logs_2026_01 PARTITION OF operation_logs
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

-- 创建索引
CREATE INDEX idx_operation_logs_operator ON operation_logs(operator_id);
CREATE INDEX idx_operation_logs_type ON operation_logs(operation_type);
CREATE INDEX idx_operation_logs_object ON operation_logs(object_type, object_id);
