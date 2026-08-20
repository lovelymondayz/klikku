CREATE TABLE IF NOT EXISTS photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES photobooth_sessions(id) ON DELETE CASCADE,
    original_url TEXT NOT NULL,
    processed_url TEXT,
    final_url TEXT,
    position INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_photos_session ON photos(session_id);
