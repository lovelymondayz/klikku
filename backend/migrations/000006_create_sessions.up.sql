CREATE TABLE IF NOT EXISTS photobooth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    campaign_id UUID REFERENCES campaigns(id) ON DELETE SET NULL,
    template_id UUID REFERENCES templates(id) ON DELETE SET NULL,
    status VARCHAR(30) DEFAULT 'STARTED',
    payment_status VARCHAR(20) DEFAULT 'PENDING',
    email VARCHAR(255),
    final_image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_sessions_merchant ON photobooth_sessions(merchant_id);
CREATE INDEX idx_sessions_status ON photobooth_sessions(status);
CREATE INDEX idx_sessions_created ON photobooth_sessions(created_at);
