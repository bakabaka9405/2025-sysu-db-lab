-- 状态变更日志触发器
-- 在包裹状态改变时自动记录操作日志

CREATE OR REPLACE FUNCTION log_parcel_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'UPDATE' AND OLD.status IS DISTINCT FROM NEW.status) THEN
        INSERT INTO operation_logs (
            operator_id,
            operation_type,
            object_type,
            object_id,
            description,
            operated_at
        ) VALUES (
            NULLIF(current_setting('app.current_user_id', true), '')::BIGINT,
            'parcel_status_changed',
            'parcel',
            NEW.id,
            format('包裹状态从 %s 变更为 %s', OLD.status, NEW.status),
            CURRENT_TIMESTAMP
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_log_parcel_status_change
    AFTER UPDATE ON parcels
    FOR EACH ROW
    EXECUTE FUNCTION log_parcel_status_change();

-- 记录包裹入库日志
CREATE OR REPLACE FUNCTION log_parcel_received()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO operation_logs (
            operator_id,
            operation_type,
            object_type,
            object_id,
            description,
            operated_at
        ) VALUES (
            NULLIF(current_setting('app.current_user_id', true), '')::BIGINT,
            'parcel_received',
            'parcel',
            NEW.id,
            format('包裹入库: %s, 取件码: %s', NEW.tracking_number, NEW.pickup_code),
            CURRENT_TIMESTAMP
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_log_parcel_received
    AFTER INSERT ON parcels
    FOR EACH ROW
    EXECUTE FUNCTION log_parcel_received();
