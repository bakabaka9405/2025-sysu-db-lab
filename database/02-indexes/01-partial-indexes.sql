CREATE INDEX idx_parcels_pickup_code_active
ON parcels(pickup_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_parcels_pending_pickup
ON parcels(received_at DESC)
WHERE status IN ('ready', 'shelved')
  AND deleted_at IS NULL;
CREATE INDEX idx_parcels_overdue
ON parcels(expected_overdue_at)
WHERE status = 'overdue' AND deleted_at IS NULL;
CREATE INDEX idx_shelves_active
ON shelves(type, current_count)
WHERE status = 'active' AND deleted_at IS NULL;
