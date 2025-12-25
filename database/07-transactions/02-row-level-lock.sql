-- 行级锁示例
-- 演示PostgreSQL行级锁的使用场景和最佳实践

-- ============================================
-- 1. FOR UPDATE - 排他锁
-- ============================================
-- 适用场景：需要修改选中行时，防止其他事务同时修改

-- 示例：取件操作（防止重复取件）
CREATE OR REPLACE FUNCTION pickup_parcel_with_lock(
    p_pickup_code VARCHAR(6),
    p_operator_id BIGINT
) RETURNS TABLE(
    parcel_id BIGINT,
    tracking_number VARCHAR(50),
    recipient_name VARCHAR(100)
) AS $$
DECLARE
    v_parcel RECORD;
BEGIN
    -- 使用FOR UPDATE锁定行，防止并发取件
    SELECT p.* INTO v_parcel
    FROM parcels p
    WHERE p.pickup_code = p_pickup_code
      AND p.deleted_at IS NULL
    FOR UPDATE NOWAIT;  -- NOWAIT: 锁冲突时立即报错而非等待

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Parcel not found with pickup code: %', p_pickup_code;
    END IF;

    IF v_parcel.status = 'picked_up' THEN
        RAISE EXCEPTION 'Parcel already picked up';
    END IF;

    IF v_parcel.status NOT IN ('received', 'shelved', 'ready_for_pickup') THEN
        RAISE EXCEPTION 'Parcel cannot be picked up in current status: %', v_parcel.status;
    END IF;

    -- 更新包裹状态
    UPDATE parcels
    SET status = 'picked_up',
        picked_up_at = CURRENT_TIMESTAMP,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_parcel.id;

    -- 更新货架计数
    IF v_parcel.shelf_id IS NOT NULL THEN
        UPDATE shelves
        SET current_count = current_count - 1
        WHERE id = v_parcel.shelf_id;
    END IF;

    -- 记录操作日志
    INSERT INTO operation_logs (operator_id, target_type, target_id, action, detail)
    VALUES (
        p_operator_id,
        'parcel',
        v_parcel.id,
        'pickup',
        jsonb_build_object(
            'pickup_code', p_pickup_code,
            'previous_status', v_parcel.status
        )
    );

    RETURN QUERY
    SELECT v_parcel.id, v_parcel.tracking_number, v_parcel.recipient_name;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 2. FOR UPDATE SKIP LOCKED - 跳过已锁定行
-- ============================================
-- 适用场景：任务队列处理，多worker并发消费

-- 示例：批量处理滞留件提醒
CREATE OR REPLACE FUNCTION process_overdue_batch(
    p_batch_size INT DEFAULT 10
) RETURNS INT AS $$
DECLARE
    v_processed INT := 0;
    v_parcel RECORD;
BEGIN
    -- 使用SKIP LOCKED跳过其他worker正在处理的记录
    FOR v_parcel IN
        SELECT p.id, p.recipient_phone, p.tracking_number
        FROM parcels p
        WHERE p.status = 'overdue'
          AND p.deleted_at IS NULL
          AND NOT EXISTS (
              SELECT 1 FROM overdue_records o
              WHERE o.parcel_id = p.id
                AND o.last_reminder_at > CURRENT_TIMESTAMP - INTERVAL '24 hours'
          )
        ORDER BY p.expected_overdue_at
        LIMIT p_batch_size
        FOR UPDATE SKIP LOCKED
    LOOP
        -- 模拟发送提醒
        RAISE NOTICE 'Sending reminder for parcel %', v_parcel.tracking_number;

        -- 更新或插入滞留记录
        INSERT INTO overdue_records (parcel_id, reminder_count, last_reminder_at)
        VALUES (v_parcel.id, 1, CURRENT_TIMESTAMP)
        ON CONFLICT (parcel_id) WHERE deleted_at IS NULL
        DO UPDATE SET
            reminder_count = overdue_records.reminder_count + 1,
            last_reminder_at = CURRENT_TIMESTAMP;

        v_processed := v_processed + 1;
    END LOOP;

    RETURN v_processed;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 3. FOR SHARE - 共享锁
-- ============================================
-- 适用场景：读取时防止其他事务修改，但允许其他事务读取

-- 示例：验证包裹状态后进行关联操作
CREATE OR REPLACE FUNCTION create_return_request(
    p_parcel_id BIGINT,
    p_reason TEXT
) RETURNS BOOLEAN AS $$
DECLARE
    v_parcel RECORD;
BEGIN
    -- 使用FOR SHARE锁定，允许其他事务读取但不能修改
    SELECT * INTO v_parcel
    FROM parcels
    WHERE id = p_parcel_id
      AND deleted_at IS NULL
    FOR SHARE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Parcel not found';
    END IF;

    -- 在持有共享锁期间，包裹状态不会被改变
    IF v_parcel.status != 'overdue' THEN
        RAISE EXCEPTION 'Only overdue parcels can be returned';
    END IF;

    -- 记录退回请求（这里简化处理）
    UPDATE parcels
    SET status = 'returned',
        notes = COALESCE(notes, '') || ' | Return reason: ' || p_reason
    WHERE id = p_parcel_id;

    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 4. FOR KEY SHARE - 最弱的锁
-- ============================================
-- 适用场景：外键检查，只防止主键被修改或行被删除

-- 示例：创建包裹时验证货架存在
-- PostgreSQL在外键约束检查时自动使用FOR KEY SHARE

-- ============================================
-- 5. 死锁预防最佳实践
-- ============================================

-- 示例：批量更新时按固定顺序加锁
CREATE OR REPLACE FUNCTION batch_update_parcels(
    p_parcel_ids BIGINT[],
    p_new_status VARCHAR(20)
) RETURNS INT AS $$
DECLARE
    v_updated INT;
BEGIN
    -- 按ID排序后加锁，避免死锁
    WITH locked_parcels AS (
        SELECT id
        FROM parcels
        WHERE id = ANY(p_parcel_ids)
          AND deleted_at IS NULL
        ORDER BY id  -- 关键：固定顺序
        FOR UPDATE
    )
    UPDATE parcels p
    SET status = p_new_status,
        updated_at = CURRENT_TIMESTAMP
    FROM locked_parcels lp
    WHERE p.id = lp.id;

    GET DIAGNOSTICS v_updated = ROW_COUNT;
    RETURN v_updated;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 6. 锁等待超时设置
-- ============================================

-- 设置锁等待超时（毫秒）
SET lock_timeout = '5s';

-- 设置语句超时
SET statement_timeout = '30s';

-- 查看当前锁信息
-- SELECT * FROM pg_locks WHERE NOT granted;

-- 查看等待中的查询
-- SELECT pid, usename, query, wait_event_type, wait_event
-- FROM pg_stat_activity
-- WHERE wait_event_type = 'Lock';

-- ============================================
-- 锁类型对比表
-- ============================================
/*
| 锁类型          | 阻塞UPDATE | 阻塞DELETE | 阻塞SELECT FOR UPDATE | 阻塞SELECT FOR SHARE |
|----------------|-----------|-----------|---------------------|---------------------|
| FOR UPDATE      | ✓         | ✓         | ✓                   | ✓                   |
| FOR NO KEY UPDATE| ✓        | ✓         | ✓                   | ✗                   |
| FOR SHARE       | ✓         | ✓         | ✗                   | ✗                   |
| FOR KEY SHARE   | ✗*        | ✓         | ✗                   | ✗                   |

* FOR KEY SHARE只阻塞修改主键的UPDATE
*/

