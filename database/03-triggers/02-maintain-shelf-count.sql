-- 货架计数自动维护触发器
-- 在包裹上架、取件或删除时自动维护货架的占用计数
-- 状态流转: received(入库,无货架) -> ready(上架) -> picked_up(已取)/overdue(滞留) -> returned(退回)

CREATE OR REPLACE FUNCTION maintain_shelf_count()
RETURNS TRIGGER AS $$
DECLARE
    v_old_shelf_id BIGINT;
    v_new_shelf_id BIGINT;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        -- 新包裹入库时，如果已分配货架且状态为占用货架状态，增加计数
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

        -- 从旧货架移除
        IF v_old_shelf_id IS NOT NULL
           AND OLD.status IN ('ready', 'overdue')
           AND (NEW.status IN ('picked_up', 'returned')
                OR v_new_shelf_id IS DISTINCT FROM v_old_shelf_id) THEN
            UPDATE shelves
            SET current_count = GREATEST(current_count - 1, 0)
            WHERE id = v_old_shelf_id;
        END IF;

        -- 添加到新货架
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
        -- 包裹删除，减少货架计数
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
