-- 货架表
-- 存储货架信息，用于包裹存放管理

CREATE TABLE shelves (
    id BIGSERIAL PRIMARY KEY,
    shelf_code VARCHAR(50) NOT NULL UNIQUE,
    area VARCHAR(10) NOT NULL,
    floor INT NOT NULL,
    "column" INT NOT NULL,
    type VARCHAR(20) NOT NULL,
    capacity INT NOT NULL DEFAULT 10,
    current_count INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    -- 约束条件
    CONSTRAINT chk_capacity_positive CHECK (capacity > 0),
    CONSTRAINT chk_current_count_valid CHECK (
        current_count >= 0 AND current_count <= capacity
    ),
    CONSTRAINT chk_type_valid CHECK (
        type IN ('small', 'medium', 'large')
    ),
    CONSTRAINT chk_status_valid CHECK (
        status IN ('active', 'inactive')
    )
);

-- 创建唯一索引，确保位置唯一
CREATE UNIQUE INDEX idx_shelves_position
ON shelves(area, floor, "column")
WHERE deleted_at IS NULL;

-- 创建状态索引
CREATE INDEX idx_shelves_status ON shelves(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_shelves_type ON shelves(type) WHERE deleted_at IS NULL;
