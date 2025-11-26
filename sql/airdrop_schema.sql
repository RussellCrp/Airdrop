CREATE DATABASE IF NOT EXISTS airdrop DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 以下语句需在 airdrop 数据库中执行
-- 例如: mysql> USE airdrop;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    wallet_addr   VARCHAR(64) NOT NULL COMMENT '用户钱包地址（EVM）',
    nickname      VARCHAR(64) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_wallet_addr (wallet_addr)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 任务定义表
CREATE TABLE IF NOT EXISTS tasks (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    code          VARCHAR(64) NOT NULL COMMENT '任务编码，如 PROMOTE, QUANT_INVEST',
    name          VARCHAR(128) NOT NULL,
    description   TEXT NULL,
    score_weight  INT NOT NULL DEFAULT 0 COMMENT '任务对应积分权重',
    enabled       TINYINT(1) NOT NULL DEFAULT 1,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_tasks_code (code)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 用户任务完成记录
CREATE TABLE IF NOT EXISTS user_tasks (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id         BIGINT UNSIGNED NOT NULL,
    task_id         BIGINT UNSIGNED NOT NULL,
    status          TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0 未完成 1 已完成',
    extra_data      JSON NULL COMMENT '任务相关扩展信息，如交易量、推荐人数等',
    completed_at    TIMESTAMP NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_task (user_id, task_id),
    KEY idx_user_tasks_user_id (user_id),
    CONSTRAINT fk_user_tasks_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_tasks_task FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 用户积分（可单独抽表方便统计）
CREATE TABLE IF NOT EXISTS user_scores (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id       BIGINT UNSIGNED NOT NULL,
    total_score   INT NOT NULL DEFAULT 0,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_scores_user_id (user_id),
    CONSTRAINT fk_user_scores_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 积分区间与空投数量
CREATE TABLE IF NOT EXISTS score_tiers (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    min_score      INT NOT NULL,
    max_score      INT NOT NULL,
    airdrop_amount DECIMAL(38, 0) NOT NULL COMMENT '该档位空投代币数量（整数）',
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_score_range (min_score, max_score)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 空投快照
CREATE TABLE IF NOT EXISTS airdrop_snapshots (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name           VARCHAR(128) NOT NULL COMMENT '本轮空投名称/标识',
    merkle_root    VARCHAR(66) NULL COMMENT 'Merkle Root (0x 前缀的 hex)',
    token_address  VARCHAR(64) NOT NULL COMMENT 'ERC20 代币地址',
    total_users    INT NOT NULL DEFAULT 0,
    total_amount   DECIMAL(38, 0) NOT NULL DEFAULT 0,
    status         TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0 草稿 1 已发布',
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 空投快照明细
CREATE TABLE IF NOT EXISTS airdrop_snapshot_items (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    snapshot_id    BIGINT UNSIGNED NOT NULL,
    user_id        BIGINT UNSIGNED NOT NULL,
    wallet_addr    VARCHAR(64) NOT NULL,
    idx            INT NOT NULL COMMENT 'Merkle 叶子 index',
    score          INT NOT NULL,
    amount         DECIMAL(38, 0) NOT NULL,
    leaf_hash      VARCHAR(66) NULL COMMENT '叶子哈希（0x 前缀）',
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_snapshot_idx (snapshot_id, idx),
    KEY idx_snapshot_user (snapshot_id, user_id),
    CONSTRAINT fk_snapshot_items_snapshot FOREIGN KEY (snapshot_id) REFERENCES airdrop_snapshots (id) ON DELETE CASCADE,
    CONSTRAINT fk_snapshot_items_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;


