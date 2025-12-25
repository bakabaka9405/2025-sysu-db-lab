-- 货架分配函数
-- 实现智能货架分配算法，综合考虑空间、位置和负载等因素

CREATE OR REPLACE FUNCTION allocate_shelf(p_parcel_size VARCHAR)
RETURNS BIGINT AS $$
DECLARE
    v_shelf_id BIGINT;
    v_candidate_types VARCHAR[];
BEGIN
    -- 根据包裹尺寸确定候选货架类型
    -- 小件可以放小/中/大货架，中件可以放中/大货架，大件只能放大货架
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

    -- 多因子评分选择最优货架
    -- 空间因子权重0.5：优先选择剩余空间大的货架
    -- 位置因子权重0.3：优先选择低楼层和前列位置，便于取件
    -- 负载因子权重0.2：用于平衡各货架负载
    SELECT id INTO v_shelf_id
    FROM shelves
    WHERE type = ANY(v_candidate_types)
      AND status = 'active'
      AND current_count < capacity
      AND deleted_at IS NULL
    ORDER BY
        -- 空间因子：剩余空间比例
        ((capacity - current_count)::NUMERIC / capacity) * 0.5 +
        -- 位置因子：低楼层和前列优先
        (1.0 / (floor + "column")) * 0.3 +
        -- 负载因子：当前负载低优先
        (1.0 - (current_count::NUMERIC / capacity)) * 0.2
    DESC
    LIMIT 1;

    IF v_shelf_id IS NULL THEN
        RAISE EXCEPTION '没有可用货架，包裹尺寸: %', p_parcel_size;
    END IF;

    RETURN v_shelf_id;
END;
$$ LANGUAGE plpgsql;
