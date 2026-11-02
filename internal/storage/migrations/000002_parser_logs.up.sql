-- Parser Service schema. Rows hold only redacted text and a hash of the user
-- ID (SEC-08) and are deleted after 30 days.
CREATE SCHEMA IF NOT EXISTS parser_svc;

CREATE TABLE parser_svc.parser_logs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id_hash  TEXT NOT NULL,
    text_redacted TEXT NOT NULL,
    result        JSONB NOT NULL,
    source        TEXT NOT NULL, -- RULE, LLM, LLM_FAILED, NO_LLM
    latency_ms    INTEGER NOT NULL,
    tokens        INTEGER NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX parser_logs_created_at ON parser_svc.parser_logs (created_at);
CREATE INDEX parser_logs_user_id_hash ON parser_svc.parser_logs (user_id_hash);
