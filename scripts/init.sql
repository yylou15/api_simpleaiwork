CREATE TABLE users (
                       id                BIGINT UNSIGNED NOT NULL,
                       email             VARCHAR(320)    NOT NULL,
                       email_norm        VARCHAR(320)    NOT NULL,
                       email_verified_at DATETIME(3)     NULL,

                       status            TINYINT         NOT NULL DEFAULT 1,
                       is_pro            TINYINT(1)      NOT NULL DEFAULT 0,

                       created_at        DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
                       updated_at        DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),

                       PRIMARY KEY (id),
                       UNIQUE KEY `ux_users_email_norm` (email_norm)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE user_identities (
                                 id            BIGINT UNSIGNED NOT NULL,
                                 user_id       BIGINT UNSIGNED NOT NULL,

                                 provider      VARCHAR(32)  NOT NULL,
                                 provider_sub  VARCHAR(255) NULL,

                                 email         VARCHAR(320) NULL,
                                 email_norm    VARCHAR(320) NULL,

                                 created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
                                 last_login_at DATETIME(3)  NULL,

                                 PRIMARY KEY (id),

                                 UNIQUE KEY ux_identity_user_provider (user_id, provider),
                                 UNIQUE KEY ux_identity_provider_sub (provider, provider_sub),

                                 KEY ix_identity_email_norm (provider, email_norm),

                                 CONSTRAINT fk_identity_user
                                     FOREIGN KEY (user_id) REFERENCES users (id)
                                         ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE email_verifications (
                                     id            BIGINT UNSIGNED NOT NULL ,

                                     email         VARCHAR(320)    NOT NULL,
                                     email_norm    VARCHAR(320)    NOT NULL,

                                     purpose       VARCHAR(32)     NOT NULL,

    -- 不存明文验证码，只存 hash（推荐 HMAC-SHA256 的原始 32 bytes）
                                     code_hash     VARBINARY(32)   NOT NULL,

                                     expires_at    DATETIME(3)     NOT NULL,
                                     consumed_at   DATETIME(3)     NULL,

                                     attempt_count INT             NOT NULL DEFAULT 0,

                                     created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

                                     request_ip    VARBINARY(16)   NULL,
                                     user_agent    VARCHAR(512)    NULL,

                                     PRIMARY KEY (id),
                                     KEY ix_ev_lookup   (email_norm, purpose, created_at),
                                     KEY ix_ev_expire   (expires_at),
                                     KEY ix_ev_consumed (consumed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
