CREATE OR REPLACE VIEW v_shelf_utilization AS
SELECT
    s.id,
    s.shelf_code,
    s.area,
    s.floor,
    s."column",
    s.type,
    s.capacity,
    s.current_count,
    s.status,
    ROUND((s.current_count::NUMERIC / s.capacity) * 100, 2) as utilization_rate,
    s.capacity - s.current_count as remaining_capacity,
    CASE
        WHEN s.current_count = 0 THEN 'empty'
        WHEN s.current_count < s.capacity * 0.3 THEN 'low'
        WHEN s.current_count < s.capacity * 0.7 THEN 'medium'
        WHEN s.current_count < s.capacity THEN 'high'
        ELSE 'full'
    END as load_level,
    CASE
        WHEN s.status = 'active' AND s.current_count < s.capacity THEN true
        ELSE false
    END as is_available
FROM shelves s
WHERE s.deleted_at IS NULL;
CREATE OR REPLACE VIEW v_shelf_summary AS
SELECT
    type,
    COUNT(*) as shelf_count,
    SUM(capacity) as total_capacity,
    SUM(current_count) as total_used,
    SUM(capacity - current_count) as total_available,
    ROUND(AVG((current_count::NUMERIC / capacity) * 100), 2) as avg_utilization_rate
FROM shelves
WHERE deleted_at IS NULL AND status = 'active'
GROUP BY type;
