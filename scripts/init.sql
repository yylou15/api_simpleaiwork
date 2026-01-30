CREATE TABLE users
(
    id                BIGINT UNSIGNED NOT NULL,
    email             VARCHAR(320) NOT NULL,
    email_norm        VARCHAR(320) NOT NULL,
    email_verified_at DATETIME(3)     NULL,

    status            TINYINT      NOT NULL DEFAULT 1,
    is_pro            TINYINT(1)      NOT NULL DEFAULT 0,

    created_at        DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at        DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),

    PRIMARY KEY (id),
    UNIQUE KEY `ux_users_email_norm` (email_norm)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE user_identities
(
    id            BIGINT UNSIGNED NOT NULL,
    user_id       BIGINT UNSIGNED NOT NULL,

    provider      VARCHAR(32) NOT NULL,
    provider_sub  VARCHAR(255) NULL,

    email         VARCHAR(320) NULL,
    email_norm    VARCHAR(320) NULL,

    created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    last_login_at DATETIME(3)  NULL,

    PRIMARY KEY (id),

    UNIQUE KEY ux_identity_user_provider (user_id, provider),
    UNIQUE KEY ux_identity_provider_sub (provider, provider_sub),

    KEY           ix_identity_email_norm (provider, email_norm),

    CONSTRAINT fk_identity_user
        FOREIGN KEY (user_id) REFERENCES users (id)
            ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE email_verifications
(
    id            BIGINT UNSIGNED NOT NULL,

    email         VARCHAR(320) NOT NULL,
    email_norm    VARCHAR(320) NOT NULL,

    purpose       VARCHAR(32)  NOT NULL,

    -- 不存明文验证码，只存 hash（推荐 HMAC-SHA256 的原始 32 bytes）
    code_hash     VARBINARY(32)   NOT NULL,

    expires_at    DATETIME(3)     NOT NULL,
    consumed_at   DATETIME(3)     NULL,

    attempt_count INT          NOT NULL DEFAULT 0,

    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    request_ip    VARBINARY(16)   NULL,
    user_agent    VARCHAR(512) NULL,

    PRIMARY KEY (id),
    KEY           ix_ev_lookup (email_norm, purpose, created_at),
    KEY           ix_ev_expire (expires_at),
    KEY           ix_ev_consumed (consumed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 目录表
CREATE TABLE categories
(
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name        VARCHAR(64) NOT NULL,
    description VARCHAR(255)         DEFAULT NULL,
    sort_order  INT         NOT NULL DEFAULT 0,
    icon        VARCHAR(255)         DEFAULT NULL is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_categories_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 模板表（归属目录，一对多）
CREATE TABLE templates
(
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    category_id BIGINT UNSIGNED NOT NULL,
    title       VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL,
    tags_text   VARCHAR(512) NOT NULL, -- 冗余文本，如 "Interrupt,Hard Reject"
    is_pro      TINYINT(1) NOT NULL DEFAULT 0,
    sort_order  INT          NOT NULL DEFAULT 0,
    is_active   TINYINT(1) NOT NULL DEFAULT 1,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY         idx_templates_category (category_id, sort_order),
    CONSTRAINT fk_templates_category
        FOREIGN KEY (category_id) REFERENCES categories (id)
            ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 模板详情表（1:1）
CREATE TABLE template_details
(
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    template_id     BIGINT UNSIGNED NOT NULL,
    headline        VARCHAR(128) NOT NULL, -- 详情页标题
    summary         VARCHAR(512) NOT NULL, -- 详情页描述/引言
    reply_soft      TEXT         NOT NULL, -- 软
    reply_neutral   TEXT         NOT NULL, -- 中性
    reply_firm      TEXT         NOT NULL, -- 强烈
    when_not_to_use TEXT         NOT NULL, -- 何时不用
    best_practices  TEXT         NOT NULL, -- 最佳实践
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_template_details_template (template_id),
    CONSTRAINT fk_template_details_template
        FOREIGN KEY (template_id) REFERENCES templates (id)
            ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 模板表（归属目录，一对多）
CREATE TABLE templates
(
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    category_id BIGINT UNSIGNED NOT NULL,
    title       VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL,
    tags_text   VARCHAR(512) NOT NULL, -- 冗余文本，如 "Interrupt,Hard Reject"
    is_pro      TINYINT(1) NOT NULL DEFAULT 0,
    sort_order  INT          NOT NULL DEFAULT 0,
    is_active   TINYINT(1) NOT NULL DEFAULT 1,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY         idx_templates_category (category_id, sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 模板详情表（1:1）
CREATE TABLE template_details
(
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    template_id     BIGINT UNSIGNED NOT NULL,
    headline        VARCHAR(128) NOT NULL, -- 详情页标题
    summary         VARCHAR(512) NOT NULL, -- 详情页描述/引言
    reply_soft      TEXT         NOT NULL, -- 软
    reply_neutral   TEXT         NOT NULL, -- 中性
    reply_firm      TEXT         NOT NULL, -- 强烈
    when_not_to_use TEXT         NOT NULL, -- 何时不用
    best_practices  TEXT         NOT NULL, -- 最佳实践
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_template_details_template (template_id),
    CONSTRAINT fk_template_details_template
        FOREIGN KEY (template_id) REFERENCES templates (id)
            ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;