-- One database, one schema per service: in Phase 2 each service only reads
-- its own schema, so there are no cross-schema foreign keys.
CREATE SCHEMA IF NOT EXISTS user_svc;
CREATE SCHEMA IF NOT EXISTS transaction_svc;
CREATE SCHEMA IF NOT EXISTS conversation_svc;

CREATE TABLE user_svc.users (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_id  BIGINT NOT NULL,
    display_name TEXT NOT NULL DEFAULT '',
    timezone     TEXT NOT NULL DEFAULT 'Asia/Jakarta',
    onboarded    BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);
-- Soft delete, then purge after 30 days; users may /start again before the purge.
CREATE UNIQUE INDEX users_telegram_id_active ON user_svc.users (telegram_id) WHERE deleted_at IS NULL;

CREATE TABLE transaction_svc.categories (
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id  UUID, -- NULL = default category
    name     TEXT NOT NULL,
    type     SMALLINT NOT NULL, -- 1 = expense, 2 = income
    keywords TEXT[] NOT NULL DEFAULT '{}',
    UNIQUE NULLS NOT DISTINCT (user_id, name, type)
);

INSERT INTO transaction_svc.categories (name, type) VALUES
    ('Makanan & Minuman', 1), ('Transportasi', 1), ('Kos & Utilitas', 1),
    ('Pulsa & Internet', 1), ('Kuliah', 1), ('Belanja', 1), ('Hiburan', 1),
    ('Kesehatan', 1), ('Sosial & Hadiah', 1), ('Lainnya', 1),
    ('Uang Saku', 2), ('Gaji/Freelance', 2), ('Beasiswa', 2), ('Lainnya', 2);

CREATE TABLE transaction_svc.transactions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL,
    wallet_id         UUID, -- shared wallet phase
    type              SMALLINT NOT NULL,
    amount            BIGINT NOT NULL CHECK (amount > 0),
    description       TEXT NOT NULL DEFAULT '',
    category_id       UUID NOT NULL REFERENCES transaction_svc.categories (id),
    occurred_at       TIMESTAMPTZ NOT NULL,
    source_message_id TEXT NOT NULL DEFAULT '',
    parser_source     TEXT NOT NULL DEFAULT 'RULE',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ
);
CREATE INDEX transactions_user_occurred ON transaction_svc.transactions (user_id, occurred_at);
CREATE INDEX transactions_user_category_occurred ON transaction_svc.transactions (user_id, category_id, occurred_at);

CREATE TABLE transaction_svc.transaction_edits (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transaction_svc.transactions (id) ON DELETE CASCADE,
    field          TEXT NOT NULL,
    old_value      TEXT NOT NULL,
    new_value      TEXT NOT NULL,
    edited_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE conversation_svc.processed_messages (
    channel             TEXT NOT NULL,
    external_message_id TEXT NOT NULL,
    processed_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (channel, external_message_id)
);
