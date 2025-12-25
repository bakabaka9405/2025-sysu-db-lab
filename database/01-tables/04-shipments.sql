-- 寄件表
-- 存储用户寄件信息

CREATE TABLE shipments (
    id BIGSERIAL PRIMARY KEY,
    shipment_number VARCHAR(100) NOT NULL UNIQUE,
    sender_id BIGINT NOT NULL REFERENCES users(id),
    sender_name VARCHAR(100) NOT NULL,
    sender_phone VARCHAR(20) NOT NULL,
    recipient_name VARCHAR(100) NOT NULL,
    recipient_phone VARCHAR(20) NOT NULL,
    recipient_address VARCHAR(500) NOT NULL,
    courier_company VARCHAR(50) NOT NULL,
    weight DECIMAL(10, 2) NOT NULL,
    volume DECIMAL(10, 4) NOT NULL DEFAULT 0,
    freight DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    shipped_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    -- 约束条件
    CONSTRAINT chk_weight_positive CHECK (weight > 0),
    CONSTRAINT chk_volume_non_negative CHECK (volume >= 0),
    CONSTRAINT chk_freight_non_negative CHECK (freight >= 0),
    CONSTRAINT chk_shipment_status_valid CHECK (
        status IN ('pending', 'shipped', 'delivered', 'cancelled')
    )
);

-- 创建索引
CREATE INDEX idx_shipments_sender ON shipments(sender_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_shipments_status ON shipments(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_shipments_created_at ON shipments(created_at DESC) WHERE deleted_at IS NULL;
