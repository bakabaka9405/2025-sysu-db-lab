CREATE OR REPLACE FUNCTION calculate_freight(
    p_courier_company VARCHAR,
    p_weight DECIMAL,
    p_volume DECIMAL DEFAULT NULL,
    p_declared_value DECIMAL DEFAULT 0
)
RETURNS TABLE(
    total_freight DECIMAL,
    base_freight DECIMAL,
    chargeable_weight DECIMAL,
    calculation_method VARCHAR
) AS $$
DECLARE
    v_first_weight DECIMAL;
    v_first_price DECIMAL;
    v_additional_price DECIMAL;
    v_volume_weight DECIMAL;
    v_chargeable_weight DECIMAL;
    v_base_freight DECIMAL;
    v_calc_method VARCHAR;
BEGIN
    SELECT cp.first_weight, cp.first_price, cp.additional_price
    INTO v_first_weight, v_first_price, v_additional_price
    FROM courier_prices cp
    WHERE cp.courier_company = p_courier_company
      AND cp.is_active = true;
    IF NOT FOUND THEN
        RAISE EXCEPTION '快递公司价格配置不存在或未启用: %', p_courier_company;
    END IF;

    IF p_volume IS NOT NULL AND p_volume > 0 THEN
        v_volume_weight := (p_volume * 1000000) / 6000 / 1000;
    ELSE
        v_volume_weight := 0;
    END IF;

    IF v_volume_weight > p_weight THEN
        v_chargeable_weight := v_volume_weight;
        v_calc_method := 'volume_weight';
    ELSE
        v_chargeable_weight := p_weight;
        v_calc_method := 'actual_weight';
    END IF;

    IF v_chargeable_weight <= v_first_weight THEN
        v_base_freight := v_first_price;
    ELSE
        v_base_freight := v_first_price +
            CEIL(v_chargeable_weight - v_first_weight) * v_additional_price;
    END IF;
    RETURN QUERY SELECT v_base_freight, v_base_freight, v_chargeable_weight, v_calc_method;
END;
$$ LANGUAGE plpgsql;
