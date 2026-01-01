-- 包裹统计视图
-- 按日期、状态、快递公司聚合包裹统计数据

CREATE OR REPLACE VIEW v_parcel_statistics AS
SELECT
    DATE(received_at) as stat_date,
    status,
    courier_company,
    size,
    COUNT(*) as parcel_count,
    ROUND(AVG(weight)::NUMERIC, 2) as avg_weight,
    ROUND(AVG(
        CASE
            WHEN picked_up_at IS NOT NULL THEN
                EXTRACT(EPOCH FROM (picked_up_at - received_at)) / 3600
            ELSE NULL
        END
    )::NUMERIC, 2) as avg_pickup_hours
FROM parcels
WHERE deleted_at IS NULL
GROUP BY DATE(received_at), status, courier_company, size;
