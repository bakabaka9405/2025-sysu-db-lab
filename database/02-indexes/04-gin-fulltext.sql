-- GIN索引 (Generalized Inverted Index)
-- 用于全文搜索和数组查询

-- 添加tsvector列用于全文搜索
ALTER TABLE parcels ADD COLUMN IF NOT EXISTS notes_tsv tsvector;

-- 创建GIN索引支持全文搜索
CREATE INDEX idx_parcels_notes_gin
ON parcels USING GIN(notes_tsv);

-- 创建触发器自动更新tsvector
CREATE OR REPLACE FUNCTION update_notes_tsv()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.notes IS NOT NULL THEN
        NEW.notes_tsv := to_tsvector('simple', NEW.notes);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_parcels_notes_tsv
    BEFORE INSERT OR UPDATE OF notes ON parcels
    FOR EACH ROW
    EXECUTE FUNCTION update_notes_tsv();

-- 全文搜索查询示例:
-- SELECT * FROM parcels
-- WHERE notes_tsv @@ to_tsquery('simple', '易碎');
