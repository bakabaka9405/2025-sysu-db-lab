-- 包裹取件存储过程
-- 封装取件流程，使用行级锁防止并发冲突

CREATE OR REPLACE FUNCTION pickup_parcel(
    p_pickup_code VARCHAR,
    p_recipient_phone VARCHAR
)
RETURNS TABLE(
    parcel_id BIGINT,
    tracking_number VARCHAR,
    recipient_name VARCHAR,
    courier_company VARCHAR,
    shelf_code VARCHAR,
    picked_up_at TIMESTAMPTZ
) AS $$
DECLARE
    v_parcel RECORD;
    v_shelf_code VARCHAR;
BEGIN
    -- 使用行级锁查询包裹
    SELECT p.* INTO v_parcel
    FROM parcels p
    WHERE p.pickup_code = p_pickup_code
      AND p.deleted_at IS NULL
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION '取件码不存在: %', p_pickup_code;
    END IF;

    -- 验证手机号
    IF v_parcel.recipient_phone != p_recipient_phone THEN
        RAISE EXCEPTION '手机号不匹配';
    END IF;

    -- 检查状态
    IF v_parcel.status = 'picked_up' THEN
        RAISE EXCEPTION '该包裹已被取走';
    END IF;

    IF v_parcel.status = 'returned' THEN
        RAISE EXCEPTION '该包裹已退回';
    END IF;

    -- 获取货架编码
    SELECT s.shelf_code INTO v_shelf_code
    FROM shelves s WHERE s.id = v_parcel.shelf_id;

    -- 更新包裹状态
    UPDATE parcels
    SET status = 'picked_up',
        picked_up_at = CURRENT_TIMESTAMP
    WHERE id = v_parcel.id;

    RETURN QUERY SELECT
        v_parcel.id,
        v_parcel.tracking_number,
        v_parcel.recipient_name,
        v_parcel.courier_company,
        v_shelf_code,
        CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;
