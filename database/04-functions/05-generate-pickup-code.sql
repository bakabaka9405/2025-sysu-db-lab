-- 取件码生成函数
-- 生成唯一的6位取件码，格式为：首字母+5位数字

CREATE OR REPLACE FUNCTION generate_pickup_code(p_prefix CHAR DEFAULT 'A')
RETURNS VARCHAR AS $$
DECLARE
    v_code VARCHAR;
BEGIN
    LOOP
        v_code := p_prefix || LPAD(FLOOR(RANDOM() * 100000)::TEXT, 5, '0');

        IF NOT EXISTS (SELECT 1 FROM parcels WHERE pickup_code = v_code) THEN
            RETURN v_code;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;
