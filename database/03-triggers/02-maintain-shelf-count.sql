CREATE OR REPLACE FUNCTION maintain_shelf_count()
RETURNS TRIGGER AS $$
DECLARE
    v_old_shelf_id BIGINT;
    v_new_shelf_id BIGINT;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        IF NEW.shelf_id IS NOT NULL
           AND NEW.status IN ('ready', 'overdue') THEN
            UPDATE shelves
            SET current_count = current_count + 1
            WHERE id = NEW.shelf_id;
        END IF;
        RETURN NEW;
    END IF;
    IF (TG_OP = 'UPDATE') THEN
        v_old_shelf_id := OLD.shelf_id;
        v_new_shelf_id := NEW.shelf_id;

        IF v_old_shelf_id IS NOT NULL
           AND OLD.status IN ('ready', 'overdue')
           AND (NEW.status IN ('picked_up', 'returned')
                OR v_new_shelf_id IS DISTINCT FROM v_old_shelf_id) THEN
            UPDATE shelves
            SET current_count = GREATEST(current_count - 1, 0)
            WHERE id = v_old_shelf_id;
        END IF;

        IF v_new_shelf_id IS NOT NULL
           AND NEW.status IN ('ready', 'overdue')
           AND (OLD.status NOT IN ('ready', 'overdue')
                OR v_new_shelf_id IS DISTINCT FROM v_old_shelf_id) THEN
            UPDATE shelves
            SET current_count = current_count + 1
            WHERE id = v_new_shelf_id;
        END IF;
        RETURN NEW;
    END IF;
    IF (TG_OP = 'DELETE') THEN
        IF OLD.shelf_id IS NOT NULL
           AND OLD.status IN ('ready', 'overdue') THEN
            UPDATE shelves
            SET current_count = GREATEST(current_count - 1, 0)
            WHERE id = OLD.shelf_id;
        END IF;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_maintain_shelf_count
    AFTER INSERT OR UPDATE OR DELETE ON parcels
    FOR EACH ROW
    EXECUTE FUNCTION maintain_shelf_count();
