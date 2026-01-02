CREATE OR REPLACE FUNCTION set_expected_overdue_time()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        NEW.expected_overdue_at := NEW.received_at + INTERVAL '3 days';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_set_expected_overdue_time
    BEFORE INSERT ON parcels
    FOR EACH ROW
    EXECUTE FUNCTION set_expected_overdue_time();
