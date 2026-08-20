CREATE TABLE IF NOT EXISTS merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    logo_url TEXT,
    primary_color VARCHAR(7) DEFAULT '#ff6b9d',
    secondary_color VARCHAR(7) DEFAULT '#c44dff',
    font VARCHAR(100) DEFAULT 'Inter',
    welcome_message TEXT DEFAULT 'Welcome to our Photobooth!',
    idle_background_url TEXT,
    email_design JSONB DEFAULT '{}',
    social_links JSONB DEFAULT '{}',
    subscription_status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_merchants_slug ON merchants(slug);
