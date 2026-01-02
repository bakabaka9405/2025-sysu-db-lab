CREATE TABLE courier_prices (
    id BIGSERIAL PRIMARY KEY,
    courier_company VARCHAR(50) NOT NULL UNIQUE,
    first_weight DECIMAL(10, 2) NOT NULL DEFAULT 1.0,
    first_price DECIMAL(10, 2) NOT NULL,
    additional_price DECIMAL(10, 2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_first_weight_positive CHECK (first_weight > 0),
    CONSTRAINT chk_first_price_positive CHECK (first_price > 0),
    CONSTRAINT chk_additional_price_non_negative CHECK (additional_price >= 0)
);
CREATE INDEX idx_courier_prices_active ON courier_prices(is_active);
INSERT INTO courier_prices (courier_company, first_weight, first_price, additional_price) VALUES
('顺丰速运', 1.0, 12.00, 2.00),
('中通快递', 1.0, 8.00, 1.50),
('圆通速递', 1.0, 8.00, 1.50),
('韵达快递', 1.0, 8.00, 1.50),
('申通快递', 1.0, 8.00, 1.50),
('京东物流', 1.0, 10.00, 1.80),
('邮政EMS', 1.0, 15.00, 3.00),
('极兔速递', 1.0, 7.00, 1.20);
