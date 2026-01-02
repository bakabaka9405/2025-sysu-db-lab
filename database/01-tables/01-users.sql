CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    real_name VARCHAR(100),
    phone VARCHAR(20) NOT NULL UNIQUE,
    id_card VARCHAR(18) UNIQUE,
    email VARCHAR(100),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_username_length CHECK (length(username) >= 3),
    CONSTRAINT chk_phone_format CHECK (phone ~ '^1[3-9]\d{9}$'),
    CONSTRAINT chk_id_card_format CHECK (
        id_card IS NULL OR length(id_card) = 18
    ),
    CONSTRAINT chk_role_valid CHECK (
        role IN ('user', 'staff', 'admin')
    )
);
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
