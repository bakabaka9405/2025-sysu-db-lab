-- 货架容量限制检查触发器
-- 在包裹分配货架前检查容量，防止超载

CREATE OR REPLACE FUNCTION check_shelf_capacity()
RETURNS TRIGGER AS $$
DECLARE
    v_capacity INT;
    v_current_count INT;
BEGIN
    IF NEW.shelf_id IS NOT NULL THEN
        -- 获取货架容量信息
        SELECT capacity, current_count
        INTO v_capacity, v_current_count
        FROM shelves
        WHERE id = NEW.shelf_id AND deleted_at IS NULL;

        IF NOT FOUND THEN
            RAISE EXCEPTION '货架不存在或已删除';
        END IF;

        -- 检查是否超过容量
        IF v_current_count >= v_capacity THEN
            RAISE EXCEPTION '货架已满，无法分配。货架ID: %, 容量: %, 当前: %',
                NEW.shelf_id, v_capacity, v_current_count;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_shelf_capacity
    BEFORE INSERT OR UPDATE OF shelf_id ON parcels
    FOR EACH ROW
    WHEN (NEW.shelf_id IS NOT NULL)
    EXECUTE FUNCTION check_shelf_capacity();
