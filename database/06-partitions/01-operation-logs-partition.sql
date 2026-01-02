-- 操作日志分区表
-- 按月分区，提高大数据量下的查询效率

-- 创建分区主表
CREATE TABLE IF NOT EXISTS operation_logs_partitioned (
    id BIGSERIAL,
    operator_id BIGINT NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    target_id BIGINT NOT NULL,
    action VARCHAR(50) NOT NULL,
    detail JSONB,
    ip_address INET,
    user_agent VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_operation_logs_partitioned PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- 创建2025年各月分区
CREATE TABLE IF NOT EXISTS operation_logs_2025_01
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_02
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_03
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_04
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-04-01') TO ('2025-05-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_05
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-05-01') TO ('2025-06-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_06
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-06-01') TO ('2025-07-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_07
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_08
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_09
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_10
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_11
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');

CREATE TABLE IF NOT EXISTS operation_logs_2025_12
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

CREATE TABLE IF NOT EXISTS operation_logs_2026_01
    PARTITION OF operation_logs_partitioned
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

-- 默认分区
CREATE TABLE IF NOT EXISTS operation_logs_default
    PARTITION OF operation_logs_partitioned
    DEFAULT;

-- 分区索引
CREATE INDEX IF NOT EXISTS idx_op_logs_part_operator
    ON operation_logs_partitioned (operator_id, created_at);
CREATE INDEX IF NOT EXISTS idx_op_logs_part_target
    ON operation_logs_partitioned (target_type, target_id, created_at);
CREATE INDEX IF NOT EXISTS idx_op_logs_part_action
    ON operation_logs_partitioned (action, created_at);

-- 自动创建新月份分区的函数
CREATE OR REPLACE FUNCTION create_monthly_partition(
    p_year INT,
    p_month INT
) RETURNS VOID AS $$
DECLARE
    v_partition_name TEXT;
    v_start_date DATE;
    v_end_date DATE;
BEGIN
    v_partition_name := format('operation_logs_%s_%s',
        p_year,
        LPAD(p_month::TEXT, 2, '0'));
    v_start_date := make_date(p_year, p_month, 1);
    v_end_date := v_start_date + INTERVAL '1 month';

    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF operation_logs_partitioned
         FOR VALUES FROM (%L) TO (%L)',
        v_partition_name, v_start_date, v_end_date
    );

    RAISE NOTICE 'Created partition: %', v_partition_name;
END;
$$ LANGUAGE plpgsql;

-- 删除旧分区的函数（保留最近N个月）
CREATE OR REPLACE FUNCTION drop_old_partitions(
    p_months_to_keep INT DEFAULT 12
) RETURNS INT AS $$
DECLARE
    v_partition RECORD;
    v_cutoff_date DATE;
    v_count INT := 0;
BEGIN
    v_cutoff_date := DATE_TRUNC('month', CURRENT_DATE) - (p_months_to_keep || ' months')::INTERVAL;

    FOR v_partition IN
        SELECT schemaname, tablename
        FROM pg_tables
        WHERE tablename LIKE 'operation_logs_20%'
          AND tablename ~ 'operation_logs_\d{4}_\d{2}'
    LOOP
        -- 从表名解析日期
        DECLARE
            v_year INT;
            v_month INT;
            v_partition_date DATE;
        BEGIN
            v_year := SUBSTRING(v_partition.tablename FROM 'operation_logs_(\d{4})_')::INT;
            v_month := SUBSTRING(v_partition.tablename FROM 'operation_logs_\d{4}_(\d{2})')::INT;
            v_partition_date := make_date(v_year, v_month, 1);

            IF v_partition_date < v_cutoff_date THEN
                EXECUTE format('DROP TABLE IF EXISTS %I.%I',
                    v_partition.schemaname, v_partition.tablename);
                v_count := v_count + 1;
                RAISE NOTICE 'Dropped partition: %', v_partition.tablename;
            END IF;
        EXCEPTION WHEN OTHERS THEN
            CONTINUE;
        END;
    END LOOP;

    RETURN v_count;
END;
$$ LANGUAGE plpgsql;
