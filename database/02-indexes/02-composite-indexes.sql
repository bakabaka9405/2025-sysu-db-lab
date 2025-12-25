-- 复合索引 (Composite Indexes)
-- 支持多条件查询和排序

-- 支持按状态和时间排序查询
CREATE INDEX idx_parcels_status_received
ON parcels(status, received_at DESC) WHERE deleted_at IS NULL;

-- 支持用户查询自己的包裹（按手机号+状态+时间）
CREATE INDEX idx_parcels_phone_status
ON parcels(recipient_phone, status, received_at DESC)
WHERE deleted_at IS NULL;

-- 支持按快递公司和日期查询
CREATE INDEX idx_parcels_courier_date
ON parcels(courier_company, received_at DESC)
WHERE deleted_at IS NULL;

-- 寄件单复合索引
CREATE INDEX idx_shipments_sender_status
ON shipments(sender_id, status, created_at DESC)
WHERE deleted_at IS NULL;
