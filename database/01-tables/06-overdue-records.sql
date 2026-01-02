CREATE TABLE overdue_records (
    id BIGSERIAL PRIMARY KEY,
    parcel_id BIGINT NOT NULL UNIQUE REFERENCES parcels(id),
    overdue_start_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reminder_count INT NOT NULL DEFAULT 0,
    last_reminder_at TIMESTAMPTZ,
    overdue_fee DECIMAL(10, 2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_reminder_count_non_negative CHECK (reminder_count >= 0),
    CONSTRAINT chk_overdue_fee_non_negative CHECK (overdue_fee >= 0),
    CONSTRAINT chk_overdue_status_valid CHECK (
        status IN ('pending', 'reminded', 'returned')
    )
);
CREATE INDEX idx_overdue_records_parcel ON overdue_records(parcel_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_overdue_records_status ON overdue_records(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_overdue_records_start ON overdue_records(overdue_start_at) WHERE deleted_at IS NULL;
