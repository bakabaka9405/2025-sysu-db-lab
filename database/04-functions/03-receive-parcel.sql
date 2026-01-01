-- 包裹入库存储过程
-- 封装完整的包裹入库流程，包括货架分配、取件码生成、状态初始化

CREATE OR REPLACE FUNCTION receive_parcel(
    p_tracking_number VARCHAR,
    p_recipient_name VARCHAR,
    p_recipient_phone VARCHAR,
    p_courier_company VARCHAR,
    p_size VARCHAR,
    p_weight DECIMAL,
    p_recipient_id_card VARCHAR DEFAULT NULL,
    p_notes TEXT DEFAULT NULL
)
RETURNS TABLE(
    parcel_id BIGINT,
    pickup_code VARCHAR,
    shelf_code VARCHAR,
    shelf_area VARCHAR,
    shelf_floor INT,
    shelf_column INT
) AS $$
DECLARE
    v_parcel_id BIGINT;
    v_pickup_code VARCHAR;
    v_shelf_id BIGINT;
    v_shelf_code VARCHAR;
    v_shelf_area VARCHAR;
    v_shelf_floor INT;
    v_shelf_column INT;
BEGIN
    -- 检查快递单号是否已存在
    IF EXISTS (SELECT 1 FROM parcels WHERE tracking_number = p_tracking_number AND deleted_at IS NULL) THEN
        RAISE EXCEPTION '快递单号已存在: %', p_tracking_number;
    END IF;

    -- 分配货架
    v_shelf_id := allocate_shelf(p_size);

    -- 获取货架信息
    SELECT s.shelf_code, s.area, s.floor, s."column"
    INTO v_shelf_code, v_shelf_area, v_shelf_floor, v_shelf_column
    FROM shelves s WHERE s.id = v_shelf_id;

    -- 生成取件码
    v_pickup_code := generate_pickup_code(SUBSTRING(v_shelf_code, 1, 1));

    -- 插入包裹记录
    INSERT INTO parcels (
        tracking_number, pickup_code, recipient_name,
        recipient_phone, recipient_id_card, courier_company,
        size, weight, status, shelf_id, notes
    ) VALUES (
        p_tracking_number, v_pickup_code, p_recipient_name,
        p_recipient_phone, p_recipient_id_card, p_courier_company,
        p_size, p_weight, 'ready_for_pickup', v_shelf_id, p_notes
    ) RETURNING id INTO v_parcel_id;

    RETURN QUERY SELECT v_parcel_id, v_pickup_code, v_shelf_code, v_shelf_area, v_shelf_floor, v_shelf_column;
END;
$$ LANGUAGE plpgsql;
