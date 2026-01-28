-- =========================
-- 1) 用户主表 users
-- =========================
CREATE TABLE IF NOT EXISTS users (
  id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  email             VARCHAR(320) NOT NULL,
  email_norm        VARCHAR(320) NOT NULL,     -- lower(trim(email))，由应用层写入
  email_verified_at DATETIME(3) NULL,

  status            TINYINT NOT NULL DEFAULT 1, -- 1=active, 0=disabled, 2=deleted(软删)
  created_at        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3),
  updated_at        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),

  PRIMARY KEY (id),
  UNIQUE KEY ux_users_email_norm (email_norm)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- =========================
-- 2) 登录身份绑定表 user_identities
--    provider: 'email' | 'google' | 'apple' ...
--    provider_sub: OAuth 的 subject（Google 的 sub）
-- =========================
CREATE TABLE IF NOT EXISTS user_identities (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id       BIGINT UNSIGNED NOT NULL,

  provider      VARCHAR(32) NOT NULL,
  provider_sub  VARCHAR(255) NULL,

  email         VARCHAR(320) NULL,
  email_norm    VARCHAR(320) NULL,

  created_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3),
  last_login_at DATETIME(3) NULL,

  PRIMARY KEY (id),

  -- 一个用户同一种 provider 只能绑定一次（避免同一用户绑两个 google）
  UNIQUE KEY ux_identity_user_provider (user_id, provider),

  -- 同一个第三方账号不能绑定到多个用户（provider + sub 唯一）
  UNIQUE KEY ux_identity_provider_sub (provider, provider_sub),

  KEY ix_identity_email_norm (provider, email_norm),

  CONSTRAINT fk_identity_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- =========================
-- 3) 邮箱验证码表 email_verifications
--    purpose: 'signup' | 'login' | 'bind' | 'change_email' ...
-- =========================
CREATE TABLE IF NOT EXISTS email_verifications (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,

  email         VARCHAR(320) NOT NULL,
  email_norm    VARCHAR(320) NOT NULL,

  purpose       VARCHAR(32) NOT NULL,

  -- 不存明文验证码，只存 hash（推荐 HMAC-SHA256 的 hex 或 base64）
  code_hash     VARBINARY(32) NOT NULL, -- 若你存 SHA256 原始 32 bytes

  expires_at    DATETIME(3) NOT NULL,
  consumed_at   DATETIME(3) NULL,

  attempt_count INT NOT NULL DEFAULT 0,

  created_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3),

  request_ip    VARBINARY(16) NULL, -- IPv4/IPv6 都能塞（应用层转 16 bytes）
  user_agent    VARCHAR(512) NULL,

  PRIMARY KEY (id),
  KEY ix_ev_lookup (email_norm, purpose, created_at),
  KEY ix_ev_expire (expires_at),
  KEY ix_ev_consumed (consumed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
