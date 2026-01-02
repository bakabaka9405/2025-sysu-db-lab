-- 种子数据
-- 用于测试和演示的初始数据

-- ============================================
-- 1. 用户数据
-- ============================================
INSERT INTO users (username, password, phone, email, role) VALUES
-- 普通用户 (密码: user123)
('user', '$2a$10$nKE8HvUfx5452x8qXXgL0.nTBTjwkQDUgN68nBVExOiOhEnMHEuJq', '13212312312', 'user01@example.com', 'user'),
-- 员工用户 (密码: staff123)
('staff', '$2a$10$iCJoOySxVTNtgmoIxxl6YOax/AvecX.4cbDlgQN1t/96e7stLB1VO', '13212312313', 'staff01@example.com', 'staff'),
-- 管理员用户 (密码: admin123)
('admin', '$2a$10$Iv6/oJLRiZ22pb5VifcfEeAiXGN8J7DEnaAYguyKQfjWsb0MTkW0m', '13212312314', 'admin@example.com', 'admin')
ON CONFLICT (username) DO NOTHING;

-- ============================================
-- 2. 货架数据
-- ============================================
INSERT INTO shelves (shelf_code, area, floor, "column", type, capacity, current_count, status) VALUES
-- A区 - 小件区 (3层 x 5列 = 15个货架)
('A-1-1', 'A', 1, 1, 'small', 20, 0, 'active'),
('A-1-2', 'A', 1, 2, 'small', 20, 0, 'active'),
('A-1-3', 'A', 1, 3, 'small', 20, 0, 'active'),
('A-1-4', 'A', 1, 4, 'small', 20, 0, 'active'),
('A-1-5', 'A', 1, 5, 'small', 20, 0, 'active'),
('A-2-1', 'A', 2, 1, 'small', 20, 0, 'active'),
('A-2-2', 'A', 2, 2, 'small', 20, 0, 'active'),
('A-2-3', 'A', 2, 3, 'small', 20, 0, 'active'),
('A-2-4', 'A', 2, 4, 'small', 20, 0, 'active'),
('A-2-5', 'A', 2, 5, 'small', 20, 0, 'active'),
('A-3-1', 'A', 3, 1, 'small', 20, 0, 'active'),
('A-3-2', 'A', 3, 2, 'small', 20, 0, 'active'),
('A-3-3', 'A', 3, 3, 'small', 20, 0, 'active'),
('A-3-4', 'A', 3, 4, 'small', 20, 0, 'active'),
('A-3-5', 'A', 3, 5, 'small', 20, 0, 'active'),

-- B区 - 中件区 (2层 x 4列 = 8个货架)
('B-1-1', 'B', 1, 1, 'medium', 15, 0, 'active'),
('B-1-2', 'B', 1, 2, 'medium', 15, 0, 'active'),
('B-1-3', 'B', 1, 3, 'medium', 15, 0, 'active'),
('B-1-4', 'B', 1, 4, 'medium', 15, 0, 'active'),
('B-2-1', 'B', 2, 1, 'medium', 15, 0, 'active'),
('B-2-2', 'B', 2, 2, 'medium', 15, 0, 'active'),
('B-2-3', 'B', 2, 3, 'medium', 15, 0, 'active'),
('B-2-4', 'B', 2, 4, 'medium', 15, 0, 'active'),

-- C区 - 大件区 (1层 x 3列 = 3个货架)
('C-1-1', 'C', 1, 1, 'large', 10, 0, 'active'),
('C-1-2', 'C', 1, 2, 'large', 10, 0, 'active'),
('C-1-3', 'C', 1, 3, 'large', 10, 0, 'active')
ON CONFLICT (shelf_code) DO NOTHING;

-- ============================================
-- 3. 快递价格配置
-- ============================================
INSERT INTO courier_prices (courier_company, first_weight, first_price, additional_price, is_active) VALUES
('顺丰速运', 1.0, 13.00, 2.00, true),
('中通快递', 1.0, 6.00, 1.00, true),
('圆通速递', 1.0, 6.00, 1.00, true),
('韵达快递', 1.0, 6.00, 1.00, true),
('申通快递', 1.0, 6.00, 1.00, true),
('百世快递', 1.0, 5.00, 0.80, true),
('极兔速递', 1.0, 5.00, 0.80, true),
('邮政快递', 1.0, 8.00, 1.50, true),
('京东物流', 1.0, 10.00, 1.50, true),
('德邦快递', 1.0, 12.00, 2.00, true)
ON CONFLICT (courier_company) DO NOTHING;

-- ============================================
-- 4. 示例包裹数据
-- ============================================
-- 注意：tracking_number和pickup_code需要唯一

-- 获取货架ID的辅助函数
DO $$
DECLARE
    v_shelf_small_1 BIGINT;
    v_shelf_small_2 BIGINT;
    v_shelf_medium_1 BIGINT;
    v_shelf_large_1 BIGINT;
    v_user_1 BIGINT;
    v_user_2 BIGINT;
    v_user_3 BIGINT;
BEGIN
    SELECT id INTO v_shelf_small_1 FROM shelves WHERE shelf_code = 'A-1-1';
    SELECT id INTO v_shelf_small_2 FROM shelves WHERE shelf_code = 'A-1-2';
    SELECT id INTO v_shelf_medium_1 FROM shelves WHERE shelf_code = 'B-1-1';
    SELECT id INTO v_shelf_large_1 FROM shelves WHERE shelf_code = 'C-1-1';
    SELECT id INTO v_user_1 FROM users WHERE username = 'user01';
    SELECT id INTO v_user_2 FROM users WHERE username = 'user02';
    SELECT id INTO v_user_3 FROM users WHERE username = 'user03';

    -- 已入库待取件的包裹
    INSERT INTO parcels (tracking_number, pickup_code, recipient_name, recipient_phone, courier_company,
        size, weight, shelf_id, status, received_at, shelved_at, expected_overdue_at)
    VALUES
    ('SF1234567890001', 'A1B2C3', '张三', '13212312312', '顺丰速运', 'small', 0.5,
        v_shelf_small_1, 'ready', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', NOW() + INTERVAL '2 days'),
    ('ZT1234567890002', 'D4E5F6', '李四', '13212312312', '中通快递', 'small', 0.8,
        v_shelf_small_1, 'ready', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', NOW() + INTERVAL '1 day'),
    ('YT1234567890003', 'G7H8I9', '王五', '13212312312', '圆通速递', 'medium', 2.5,
        v_shelf_medium_1, 'ready', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', NOW() + INTERVAL '2 days'),
    ('JD1234567890004', 'J1K2L3', '赵六', '13212312312', '京东物流', 'large', 5.0,
        v_shelf_large_1, 'overdue', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 hour')
    ON CONFLICT (tracking_number) DO NOTHING;

    -- 已取件的包裹
    INSERT INTO parcels (tracking_number, pickup_code, recipient_name, recipient_phone, courier_company,
        size, weight, shelf_id, status, received_at, shelved_at, picked_up_at, expected_overdue_at)
    VALUES
    ('SF1234567890005', 'M4N5O6', '钱七', '13212312312', '顺丰速运', 'small', 0.3,
        v_shelf_small_2, 'picked_up', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '4 days', NOW() - INTERVAL '2 days')
    ON CONFLICT (tracking_number) DO NOTHING;

    -- 更新货架当前数量
    UPDATE shelves SET current_count = (
        SELECT COUNT(*) FROM parcels
        WHERE shelf_id = shelves.id
        AND status IN ('received', 'ready', 'overdue')
        AND deleted_at IS NULL
    );
END $$;
