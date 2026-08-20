CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    device_token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'online',
    current_campaign_id UUID REFERENCES campaigns(id) ON DELETE SET NULL,
    printer_config JSONB DEFAULT '{}',
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_devices_merchant ON devices(merchant_id);
CREATE INDEX idx_devices_token ON devices(device_token);
