-- +migrate Up
CREATE TABLE IF NOT EXISTS download_tokens (
    token VARCHAR(255) PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES photobooth_sessions(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_download_tokens_session ON download_tokens(session_id);
CREATE INDEX idx_download_tokens_expires ON download_tokens(expires_at);

-- Add missing columns to existing tables
ALTER TABLE print_jobs ADD COLUMN IF NOT EXISTS printer_name VARCHAR(255);
ALTER TABLE print_jobs ADD COLUMN IF NOT EXISTS error_message TEXT;

-- Add payment_reference column to payments
ALTER TABLE payments ADD COLUMN IF NOT EXISTS transaction_reference VARCHAR(255);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS payment_method VARCHAR(50);
