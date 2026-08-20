CREATE TABLE IF NOT EXISTS print_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES photobooth_sessions(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    print_type VARCHAR(20) DEFAULT '4x6',
    status VARCHAR(20) DEFAULT 'QUEUED',
    copies INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    printed_at TIMESTAMPTZ
);

CREATE INDEX idx_print_jobs_session ON print_jobs(session_id);
CREATE INDEX idx_print_jobs_status ON print_jobs(status);
