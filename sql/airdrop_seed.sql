USE airdrop;

-- 基础任务配置
INSERT INTO tasks (code, name, description, score_weight, enabled)
VALUES
    ('PROMOTE', '宣传项目', '完成指定宣传任务', 10, 1),
    ('QUANT_INVEST', '投资量化策略', '投资指定量化策略产品', 20, 1),
    ('TRADE_VOLUME_10K', '交易量达到 10000 美元', '累计交易量达到 10000 美元', 30, 1),
    ('REFERRAL_10', '推荐 10 位朋友', '成功推荐 10 位朋友注册/完成任务', 25, 1),
    ('LOGIN_7_DAYS', '连续登录 7 天', '连续登录天数达到 7 天', 15, 1)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    description = VALUES(description),
    score_weight = VALUES(score_weight),
    enabled = VALUES(enabled);

-- 示例积分档位
INSERT INTO score_tiers (min_score, max_score, airdrop_amount)
VALUES
    (0, 20, 100),
    (21, 50, 300),
    (51, 100, 800),
    (101, 1000000, 2000)
ON DUPLICATE KEY UPDATE
    airdrop_amount = VALUES(airdrop_amount);

-- 示例用户与积分（仅用于测试）
INSERT INTO users (id, wallet_addr, nickname)
VALUES
    (1, '0x0000000000000000000000000000000000000001', 'alice'),
    (2, '0x0000000000000000000000000000000000000002', 'bob')
ON DUPLICATE KEY UPDATE
    wallet_addr = VALUES(wallet_addr),
    nickname = VALUES(nickname);

INSERT INTO user_scores (user_id, total_score)
VALUES
    (1, 60),
    (2, 25)
ON DUPLICATE KEY UPDATE
    total_score = VALUES(total_score);


