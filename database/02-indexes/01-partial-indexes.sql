-- 部分索引 (Partial Indexes)
-- 只索引满足特定条件的行，节省存储空间并提高查询效率

-- 只索引未删除的包裹（取件码查询优化）
CREATE INDEX idx_parcels_pickup_code_active
ON parcels(pickup_code) WHERE deleted_at IS NULL;

-- 只索引待取件的包裹（工作人员常用查询）
CREATE INDEX idx_parcels_pending_pickup
ON parcels(received_at DESC)
WHERE status IN ('ready_for_pickup', 'shelved')
  AND deleted_at IS NULL;

-- 只索引滞留包裹
CREATE INDEX idx_parcels_overdue
ON parcels(expected_overdue_at)
WHERE status = 'overdue' AND deleted_at IS NULL;

-- 只索引活跃货架
CREATE INDEX idx_shelves_active
ON shelves(type, current_count)
WHERE status = 'active' AND deleted_at IS NULL;
