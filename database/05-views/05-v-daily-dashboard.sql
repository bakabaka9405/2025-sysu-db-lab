CREATE OR REPLACE VIEW v_daily_dashboard AS
SELECT
    (SELECT COUNT(*) FROM parcels
     WHERE DATE(received_at) = CURRENT_DATE AND deleted_at IS NULL) as today_received,
    (SELECT COUNT(*) FROM parcels
     WHERE DATE(picked_up_at) = CURRENT_DATE AND deleted_at IS NULL) as today_picked_up,
    (SELECT COUNT(*) FROM parcels
     WHERE status = 'overdue' AND deleted_at IS NULL) as overdue_count,

    (SELECT COUNT(*) FROM parcels
     WHERE status IN ('received', 'shelved', 'ready') AND deleted_at IS NULL) as pending_pickup,

    (SELECT SUM(capacity) FROM shelves WHERE deleted_at IS NULL AND status = 'active') as total_capacity,
    (SELECT SUM(current_count) FROM shelves WHERE deleted_at IS NULL AND status = 'active') as total_used,
    (SELECT ROUND(
        (SUM(current_count)::NUMERIC / NULLIF(SUM(capacity), 0)) * 100, 2
    ) FROM shelves WHERE deleted_at IS NULL AND status = 'active') as utilization_rate,

    (SELECT COUNT(*) FROM shipments
     WHERE DATE(created_at) = CURRENT_DATE AND deleted_at IS NULL) as today_shipments,

    (SELECT COUNT(*) FROM parcels
     WHERE received_at >= DATE_TRUNC('week', CURRENT_DATE) AND deleted_at IS NULL) as week_received,
    (SELECT COUNT(*) FROM parcels
     WHERE picked_up_at >= DATE_TRUNC('week', CURRENT_DATE) AND deleted_at IS NULL) as week_picked_up,

    CURRENT_TIMESTAMP as generated_at;
