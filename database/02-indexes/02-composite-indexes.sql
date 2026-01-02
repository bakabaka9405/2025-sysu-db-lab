CREATE INDEX idx_parcels_status_received
ON parcels(status, received_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_parcels_phone_status
ON parcels(recipient_phone, status, received_at DESC)
WHERE deleted_at IS NULL;
CREATE INDEX idx_parcels_courier_date
ON parcels(courier_company, received_at DESC)
WHERE deleted_at IS NULL;
CREATE INDEX idx_shipments_sender_status
ON shipments(sender_id, status, created_at DESC)
WHERE deleted_at IS NULL;
