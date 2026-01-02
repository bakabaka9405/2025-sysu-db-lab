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
    pickup_code VARCHAR
) AS $$
DECLARE
    v_parcel_id BIGINT;
    v_pickup_code VARCHAR;
BEGIN
    IF EXISTS (SELECT 1 FROM parcels WHERE tracking_number = p_tracking_number AND deleted_at IS NULL) THEN
        RAISE EXCEPTION '快递单号已存在: %', p_tracking_number;
    END IF;

    v_pickup_code := generate_pickup_code(CHR(65 + FLOOR(RANDOM() * 26)::INT));

    INSERT INTO parcels (
        tracking_number, pickup_code, recipient_name,
        recipient_phone, recipient_id_card, courier_company,
        size, weight, status, notes
    ) VALUES (
        p_tracking_number, v_pickup_code, p_recipient_name,
        p_recipient_phone, p_recipient_id_card, p_courier_company,
        p_size, p_weight, 'received', p_notes
    ) RETURNING id INTO v_parcel_id;
    RETURN QUERY SELECT v_parcel_id, v_pickup_code;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE FUNCTION shelve_parcel(
    p_parcel_id BIGINT
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
    v_parcel RECORD;
    v_shelf_id BIGINT;
    v_shelf_code VARCHAR;
    v_shelf_area VARCHAR;
    v_shelf_floor INT;
    v_shelf_column INT;
BEGIN
    SELECT * INTO v_parcel FROM parcels WHERE id = p_parcel_id AND deleted_at IS NULL;
    
    IF NOT FOUND THEN
        RAISE EXCEPTION '包裹不存在: %', p_parcel_id;
    END IF;
    
    IF v_parcel.status != 'received' THEN
        RAISE EXCEPTION '只有待上架状态的包裹才能上架，当前状态: %', v_parcel.status;
    END IF;

    v_shelf_id := allocate_shelf(v_parcel.size);

    SELECT s.shelf_code, s.area, s.floor, s."column"
    INTO v_shelf_code, v_shelf_area, v_shelf_floor, v_shelf_column
    FROM shelves s WHERE s.id = v_shelf_id;

    UPDATE parcels
    SET status = 'ready',
        shelf_id = v_shelf_id,
        shelved_at = CURRENT_TIMESTAMP,
        expected_overdue_at = CURRENT_TIMESTAMP + INTERVAL '3 days'
    WHERE id = p_parcel_id;
    RETURN QUERY SELECT p_parcel_id, v_parcel.pickup_code, v_shelf_code, v_shelf_area, v_shelf_floor, v_shelf_column;
END;
$$ LANGUAGE plpgsql;
