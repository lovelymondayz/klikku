CREATE TABLE IF NOT EXISTS templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    campaign_id UUID REFERENCES campaigns(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    preview_url TEXT,
    layout_config JSONB DEFAULT '{}',
    overlay_url TEXT,
    output_width INT DEFAULT 1200,
    output_height INT DEFAULT 1800,
    photo_count INT DEFAULT 4,
    price DECIMAL(10,2) DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_templates_merchant ON templates(merchant_id);
CREATE INDEX idx_templates_campaign ON templates(campaign_id);
