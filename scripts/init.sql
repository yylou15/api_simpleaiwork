-- Users table
CREATE TABLE IF NOT EXISTS users (
    id                BIGSERIAL PRIMARY KEY,
    email             VARCHAR(320) NOT NULL,
    email_norm        VARCHAR(320) NOT NULL,
    email_verified_at TIMESTAMP(3) NULL,
    status            SMALLINT NOT NULL DEFAULT 1,
    created_at        TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (email_norm)
);

-- User Identities table
CREATE TABLE IF NOT EXISTS user_identities (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL,
    provider      VARCHAR(32) NOT NULL,
    provider_sub  VARCHAR(255) NULL,
    email         VARCHAR(320) NULL,
    email_norm    VARCHAR(320) NULL,
    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP(3) NULL,
    UNIQUE (user_id, provider),
    UNIQUE (provider, provider_sub),
    CONSTRAINT fk_identity_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS ix_identity_email_norm ON user_identities (provider, email_norm);

-- Email Verifications table
CREATE TABLE IF NOT EXISTS email_verifications (
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(320) NOT NULL,
    email_norm    VARCHAR(320) NOT NULL,
    purpose       VARCHAR(32) NOT NULL,
    code_hash     BYTEA NOT NULL,
    expires_at    TIMESTAMP(3) NOT NULL,
    consumed_at   TIMESTAMP(3) NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    request_ip    BYTEA NULL,
    user_agent    VARCHAR(512) NULL
);

CREATE INDEX IF NOT EXISTS ix_ev_lookup ON email_verifications (email_norm, purpose, created_at);
CREATE INDEX IF NOT EXISTS ix_ev_expire ON email_verifications (expires_at);
CREATE INDEX IF NOT EXISTS ix_ev_consumed ON email_verifications (consumed_at);
