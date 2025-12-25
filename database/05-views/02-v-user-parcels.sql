-- 用户包裹视图（数据脱敏）
-- 为用户查询包裹提供脱敏后的安全视图

CREATE OR REPLACE VIEW v_user_parcels AS
SELECT
    p.id,
    p.tracking_number,
    p.pickup_code,
    p.recipient_name,
    -- 手机号脱敏：只显示前3位和后4位
    SUBSTRING(p.recipient_phone, 1, 3) || '****' ||
        SUBSTRING(p.recipient_phone, 8, 4) as recipient_phone_masked,
    -- 身份证脱敏：只显示前6位和后4位
    CASE
        WHEN p.recipient_id_card IS NOT NULL THEN
            SUBSTRING(p.recipient_id_card, 1, 6) || '********' ||
            SUBSTRING(p.recipient_id_card, 15, 4)
        ELSE NULL
    END as recipient_id_card_masked,
    p.courier_company,
    p.size,
    p.status,
    s.shelf_code,
    s.area as shelf_area,
    s.floor as shelf_floor,
    s."column" as shelf_column,
    p.received_at,
    p.shelved_at,
    p.picked_up_at,
    p.expected_overdue_at,
    -- 计算剩余时间(小时)
    CASE
        WHEN p.expected_overdue_at IS NOT NULL AND p.status NOT IN ('picked_up', 'returned') THEN
            ROUND(EXTRACT(EPOCH FROM (p.expected_overdue_at - CURRENT_TIMESTAMP)) / 3600, 1)
        ELSE NULL
    END as remaining_hours,
    -- 是否即将超时（24小时内）
    CASE
        WHEN p.expected_overdue_at IS NOT NULL
             AND p.status NOT IN ('picked_up', 'returned')
             AND p.expected_overdue_at < CURRENT_TIMESTAMP + INTERVAL '24 hours'
        THEN true
        ELSE false
    END as is_expiring_soon
FROM parcels p
LEFT JOIN shelves s ON p.shelf_id = s.id
WHERE p.deleted_at IS NULL;
