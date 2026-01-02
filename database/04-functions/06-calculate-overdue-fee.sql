CREATE OR REPLACE FUNCTION calculate_overdue_fee(
    p_overdue_start_at TIMESTAMPTZ,
    p_fee_per_day DECIMAL DEFAULT 1.0,
    p_free_days INT DEFAULT 0
)
RETURNS DECIMAL AS $$
DECLARE
    v_overdue_days INT;
    v_chargeable_days INT;
    v_fee DECIMAL;
BEGIN
    v_overdue_days := EXTRACT(DAY FROM (CURRENT_TIMESTAMP - p_overdue_start_at))::INT;

    v_chargeable_days := GREATEST(v_overdue_days - p_free_days, 0);

    v_fee := v_chargeable_days * p_fee_per_day;
    RETURN v_fee;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE FUNCTION update_overdue_fees()
RETURNS INT AS $$
DECLARE
    v_updated_count INT := 0;
BEGIN
    UPDATE overdue_records
    SET overdue_fee = calculate_overdue_fee(overdue_start_at),
        updated_at = CURRENT_TIMESTAMP
    WHERE status != 'returned'
      AND deleted_at IS NULL;
    GET DIAGNOSTICS v_updated_count = ROW_COUNT;
    RETURN v_updated_count;
END;
$$ LANGUAGE plpgsql;
