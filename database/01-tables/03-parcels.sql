-- 包裹表
-- 存储包裹信息，是系统的核心数据表

CREATE TABLE parcels (
    id BIGSERIAL PRIMARY KEY,
    tracking_number VARCHAR(50) NOT NULL UNIQUE,
    pickup_code VARCHAR(20) NOT NULL UNIQUE,
    recipient_name VARCHAR(100) NOT NULL,
    recipient_phone VARCHAR(20) NOT NULL,
    recipient_id_card VARCHAR(18),
    courier_company VARCHAR(50) NOT NULL,
    size VARCHAR(20) NOT NULL,
    weight DECIMAL(10, 2),
    status VARCHAR(20) NOT NULL DEFAULT 'ready_for_pickup',
    shelf_id BIGINT REFERENCES shelves(id) ON DELETE SET NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    shelved_at TIMESTAMPTZ,
    picked_up_at TIMESTAMPTZ,
    expected_overdue_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    -- 约束条件
    CONSTRAINT chk_size_valid CHECK (
        size IN ('small', 'medium', 'large')
    ),
    CONSTRAINT chk_weight_positive CHECK (weight IS NULL OR weight > 0),
    CONSTRAINT chk_status_valid CHECK (status IN (
        'ready_for_pickup', 'picked_up', 'overdue', 'returned'
    )),
    CONSTRAINT chk_pickup_code_format CHECK (
        pickup_code ~ '^[A-Z][0-9A-Z]{5}$'
    ),
    CONSTRAINT chk_shelved_after_received CHECK (
        shelved_at IS NULL OR shelved_at >= received_at
    ),
    CONSTRAINT chk_picked_up_after_shelved CHECK (
        picked_up_at IS NULL OR
        (shelved_at IS NOT NULL AND picked_up_at >= shelved_at)
    )
);

-- 创建索引
CREATE INDEX idx_parcels_tracking ON parcels(tracking_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_parcels_pickup_code ON parcels(pickup_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_parcels_recipient_phone ON parcels(recipient_phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_parcels_status ON parcels(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_parcels_shelf ON parcels(shelf_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_parcels_received_at ON parcels(received_at DESC) WHERE deleted_at IS NULL;
