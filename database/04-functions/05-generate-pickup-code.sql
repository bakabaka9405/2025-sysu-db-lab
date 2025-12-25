-- 取件码生成函数
-- 生成唯一的6位取件码，格式为：首字母+5位数字/字母

CREATE OR REPLACE FUNCTION generate_pickup_code(p_prefix CHAR DEFAULT 'A')
RETURNS VARCHAR AS $$
DECLARE
    v_code VARCHAR;
    v_attempts INT := 0;
    v_max_attempts INT := 100;
BEGIN
    LOOP
        -- 生成格式：前缀字母 + 5位随机数字
        v_code := p_prefix || LPAD(FLOOR(RANDOM() * 100000)::TEXT, 5, '0');

        -- 检查是否已存在
        IF NOT EXISTS (SELECT 1 FROM parcels WHERE pickup_code = v_code) THEN
            RETURN v_code;
        END IF;

        v_attempts := v_attempts + 1;
        IF v_attempts >= v_max_attempts THEN
            -- 如果纯数字冲突太多，尝试使用字母数字混合
            v_code := p_prefix ||
                CHR(65 + FLOOR(RANDOM() * 26)::INT) ||
                LPAD(FLOOR(RANDOM() * 10000)::TEXT, 4, '0');

            IF NOT EXISTS (SELECT 1 FROM parcels WHERE pickup_code = v_code) THEN
                RETURN v_code;
            END IF;

            RAISE EXCEPTION '无法生成唯一取件码，请稍后重试';
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;
