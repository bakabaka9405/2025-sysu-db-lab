-- 待处理滞留件视图
-- 展示所有超期未取的包裹，包含滞留天数和费用计算

CREATE OR REPLACE VIEW v_pending_overdue_parcels AS
SELECT
    p.id as parcel_id,
    p.tracking_number,
    p.pickup_code,
    p.recipient_name,
    p.recipient_phone,
    p.courier_company,
    p.size,
    p.status,
    s.shelf_code,
    p.received_at,
    p.expected_overdue_at,
    -- 滞留天数
    EXTRACT(DAY FROM (CURRENT_TIMESTAMP - p.expected_overdue_at))::INT as overdue_days,
    -- 滞留费用（每天1元）
    GREATEST(
        EXTRACT(DAY FROM (CURRENT_TIMESTAMP - p.expected_overdue_at))::INT,
        0
    ) * 1.0 as estimated_fee,
    -- 滞留记录信息
    o.id as overdue_record_id,
    o.reminder_count,
    o.last_reminder_at,
    o.overdue_fee as recorded_fee,
    o.status as overdue_status
FROM parcels p
LEFT JOIN shelves s ON p.shelf_id = s.id
LEFT JOIN overdue_records o ON p.id = o.parcel_id AND o.deleted_at IS NULL
WHERE p.deleted_at IS NULL
  AND p.status NOT IN ('picked_up', 'returned')
  AND p.expected_overdue_at < CURRENT_TIMESTAMP
ORDER BY p.expected_overdue_at ASC;
