-- 事务隔离级别示例
-- 演示PostgreSQL四种隔离级别的使用场景

-- ============================================
-- 1. READ UNCOMMITTED（PostgreSQL中等同于READ COMMITTED）
-- ============================================
-- PostgreSQL不支持脏读，READ UNCOMMITTED实际行为与READ COMMITTED相同

-- ============================================
-- 2. READ COMMITTED（默认级别）
-- ============================================
-- 适用场景：普通的CRUD操作，对一致性要求不高的查询
-- 特点：每条SQL语句看到的是该语句开始时已提交的数据

-- 示例：查询包裹列表（允许不可重复读）
BEGIN;
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;

SELECT * FROM parcels WHERE status = 'received';
-- 此时另一个事务可能更新了包裹状态
-- 再次查询可能看到不同的结果
SELECT * FROM parcels WHERE status = 'received';

COMMIT;

-- ============================================
-- 3. REPEATABLE READ（可重复读）
-- ============================================
-- 适用场景：需要在事务中多次读取相同数据且要求一致的场景
-- 特点：事务开始时创建快照，整个事务期间看到的数据一致

-- 示例：生成统计报表（需要数据一致性）
BEGIN;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;

-- 统计今日入库数量
SELECT COUNT(*) as received_count
FROM parcels
WHERE DATE(received_at) = CURRENT_DATE;

-- 统计今日取件数量
SELECT COUNT(*) as picked_up_count
FROM parcels
WHERE DATE(picked_up_at) = CURRENT_DATE;

-- 计算货架利用率
SELECT
    SUM(current_count) as total_used,
    SUM(capacity) as total_capacity,
    ROUND(SUM(current_count)::NUMERIC / SUM(capacity) * 100, 2) as utilization
FROM shelves
WHERE status = 'active';

-- 整个事务期间数据保持一致
COMMIT;

-- ============================================
-- 4. SERIALIZABLE（可串行化）
-- ============================================
-- 适用场景：需要最高一致性保证的关键业务操作
-- 特点：完全隔离，如同串行执行，可能发生序列化失败需要重试

-- 示例：货架分配（防止超分配）
CREATE OR REPLACE FUNCTION allocate_shelf_serializable(
    p_parcel_id BIGINT,
    p_size VARCHAR(20)
) RETURNS BIGINT AS $$
DECLARE
    v_shelf_id BIGINT;
    v_retry_count INT := 0;
    v_max_retries INT := 3;
BEGIN
    LOOP
        BEGIN
            -- 使用SERIALIZABLE隔离级别
            SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;

            -- 查找合适的货架
            SELECT id INTO v_shelf_id
            FROM shelves
            WHERE status = 'active'
              AND type = p_size
              AND current_count < capacity
            ORDER BY current_count ASC
            LIMIT 1;

            IF v_shelf_id IS NULL THEN
                RAISE EXCEPTION 'No available shelf for size: %', p_size;
            END IF;

            -- 更新货架计数
            UPDATE shelves
            SET current_count = current_count + 1
            WHERE id = v_shelf_id;

            -- 更新包裹
            UPDATE parcels
            SET shelf_id = v_shelf_id,
                status = 'shelved',
                shelved_at = CURRENT_TIMESTAMP
            WHERE id = p_parcel_id;

            RETURN v_shelf_id;

        EXCEPTION
            WHEN serialization_failure OR deadlock_detected THEN
                v_retry_count := v_retry_count + 1;
                IF v_retry_count >= v_max_retries THEN
                    RAISE EXCEPTION 'Failed after % retries', v_max_retries;
                END IF;
                -- 等待随机时间后重试
                PERFORM pg_sleep(random() * 0.1);
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 隔离级别对比表
-- ============================================
/*
| 隔离级别         | 脏读 | 不可重复读 | 幻读 | 性能 | 使用场景 |
|-----------------|------|-----------|------|------|---------|
| READ UNCOMMITTED | ✗    | ✓         | ✓    | 最高 | PostgreSQL不支持 |
| READ COMMITTED   | ✗    | ✓         | ✓    | 高   | 普通CRUD |
| REPEATABLE READ  | ✗    | ✗         | ✗*   | 中   | 报表统计 |
| SERIALIZABLE     | ✗    | ✗         | ✗    | 低   | 关键业务 |

注：PostgreSQL的REPEATABLE READ通过快照隔离实现，也能防止幻读
*/

-- ============================================
-- 查看当前隔离级别
-- ============================================
SHOW transaction_isolation;

-- 设置会话默认隔离级别
SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL REPEATABLE READ;

