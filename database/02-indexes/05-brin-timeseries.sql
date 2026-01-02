CREATE INDEX idx_operation_logs_operated_at_brin
ON operation_logs USING BRIN(operated_at)
WITH (pages_per_range = 128);
CREATE INDEX idx_parcels_created_at_brin
ON parcels USING BRIN(created_at)
WITH (pages_per_range = 128);
CREATE INDEX idx_shipments_created_at_brin
ON shipments USING BRIN(created_at)
WITH (pages_per_range = 128);
