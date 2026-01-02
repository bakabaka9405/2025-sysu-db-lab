CREATE OR REPLACE FUNCTION allocate_shelf(p_parcel_size VARCHAR)
RETURNS BIGINT AS $$
DECLARE
    v_shelf_id BIGINT;
    v_candidate_types VARCHAR[];
BEGIN

    CASE p_parcel_size
        WHEN 'small' THEN
            v_candidate_types := ARRAY['small', 'medium', 'large'];
        WHEN 'medium' THEN
            v_candidate_types := ARRAY['medium', 'large'];
        WHEN 'large' THEN
            v_candidate_types := ARRAY['large'];
        ELSE
            RAISE EXCEPTION '无效的包裹尺寸: %', p_parcel_size;
    END CASE;


    SELECT id INTO v_shelf_id
    FROM shelves
    WHERE type = ANY(v_candidate_types)
      AND status = 'active'
      AND current_count < capacity
      AND deleted_at IS NULL
    ORDER BY
        ((capacity - current_count)::NUMERIC / capacity) * 0.5 +
        (1.0 / (floor + "column")) * 0.3 +
        (1.0 - (current_count::NUMERIC / capacity)) * 0.2
    DESC
    LIMIT 1;
    IF v_shelf_id IS NULL THEN
        RAISE EXCEPTION '没有可用货架，包裹尺寸: %', p_parcel_size;
    END IF;
    RETURN v_shelf_id;
END;
$$ LANGUAGE plpgsql;
