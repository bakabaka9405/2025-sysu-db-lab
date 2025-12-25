-- BRIN索引 (Block Range Index)
-- 适用于按时间顺序插入的数据，索引体积小

-- 操作日志表的BRIN索引
CREATE INDEX idx_operation_logs_operated_at_brin
ON operation_logs USING BRIN(operated_at)
WITH (pages_per_range = 128);

-- 包裹表创建时间的BRIN索引（适合历史数据查询）
CREATE INDEX idx_parcels_created_at_brin
ON parcels USING BRIN(created_at)
WITH (pages_per_range = 128);

-- 寄件单创建时间的BRIN索引
CREATE INDEX idx_shipments_created_at_brin
ON shipments USING BRIN(created_at)
WITH (pages_per_range = 128);
