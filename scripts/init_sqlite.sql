-- Users table
CREATE TABLE IF NOT EXISTS users (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    email             TEXT NOT NULL,
    email_norm        TEXT NOT NULL,
    email_verified_at DATETIME NULL,
    status            INTEGER NOT NULL DEFAULT 1,
    is_pro            INTEGER NOT NULL DEFAULT 0,
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (email_norm)
);

-- User Identities table
CREATE TABLE IF NOT EXISTS user_identities (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL,
    provider      TEXT NOT NULL,
    provider_sub  TEXT NULL,
    email         TEXT NULL,
    email_norm    TEXT NULL,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at DATETIME NULL,
    UNIQUE (user_id, provider),
    UNIQUE (provider, provider_sub),
    CONSTRAINT fk_identity_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS ix_identity_email_norm ON user_identities (provider, email_norm);

-- Email Verifications table
CREATE TABLE IF NOT EXISTS email_verifications (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         TEXT NOT NULL,
    email_norm    TEXT NOT NULL,
    purpose       TEXT NOT NULL,
    code_hash     BLOB NOT NULL,
    expires_at    DATETIME NOT NULL,
    consumed_at   DATETIME NULL,
    attempt_count INTEGER NOT NULL DEFAULT 0,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    request_ip    BLOB NULL,
    user_agent    TEXT NULL
);

CREATE INDEX IF NOT EXISTS ix_ev_lookup ON email_verifications (email_norm, purpose, created_at);
CREATE INDEX IF NOT EXISTS ix_ev_expire ON email_verifications (expires_at);
CREATE INDEX IF NOT EXISTS ix_ev_consumed ON email_verifications (consumed_at);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    icon TEXT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    tags_text TEXT NOT NULL,
    is_pro INTEGER NOT NULL DEFAULT 0,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_templates_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_templates_category ON templates (category_id, sort_order);

CREATE TABLE IF NOT EXISTS template_details (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL,
    headline TEXT NOT NULL,
    summary TEXT NOT NULL,
    reply_soft TEXT NOT NULL,
    reply_neutral TEXT NOT NULL,
    reply_firm TEXT NOT NULL,
    when_not_to_use TEXT NOT NULL,
    best_practices TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (template_id),
    CONSTRAINT fk_template_details_template FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE
);
