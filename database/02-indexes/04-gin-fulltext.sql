ALTER TABLE parcels ADD COLUMN IF NOT EXISTS notes_tsv tsvector;
CREATE INDEX idx_parcels_notes_gin
ON parcels USING GIN(notes_tsv);
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
