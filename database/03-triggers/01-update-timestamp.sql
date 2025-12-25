-- 自动更新时间戳触发器
-- 在记录更新时自动维护updated_at字段

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为各表创建触发器
CREATE TRIGGER trg_users_update_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_shelves_update_timestamp
    BEFORE UPDATE ON shelves
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_parcels_update_timestamp
    BEFORE UPDATE ON parcels
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_shipments_update_timestamp
    BEFORE UPDATE ON shipments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_overdue_records_update_timestamp
    BEFORE UPDATE ON overdue_records
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
