DO $$
BEGIN
    CREATE TYPE notification_status AS ENUM ('PENDING', 'SENT', 'CANCELLED', 'FAILED');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS medicine_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS medicine_category_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES medicine_categories(id) ON DELETE CASCADE,
    display_name VARCHAR(255) NOT NULL,
    default_dosage_text VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE patient_medicines
    ADD COLUMN IF NOT EXISTS category_item_id UUID REFERENCES medicine_category_items(id);

CREATE TABLE IF NOT EXISTS device_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(20) NOT NULL,
    token TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    data JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notification_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_code VARCHAR(50) NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    sent_at TIMESTAMPTZ,
    status notification_status NOT NULL DEFAULT 'PENDING',
    payload JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    weekly_reminder_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS support_chat_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    category VARCHAR(20) NOT NULL,
    attachment_url TEXT,
    status VARCHAR(20) DEFAULT 'OPEN',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_medicine_category_items_category_id ON medicine_category_items(category_id);
CREATE INDEX IF NOT EXISTS idx_medicine_category_items_active ON medicine_category_items(category_id, is_active);
CREATE INDEX IF NOT EXISTS idx_patient_medicines_category_item_id ON patient_medicines(category_item_id);
CREATE INDEX IF NOT EXISTS idx_device_tokens_user_id_token ON device_tokens(user_id, token);
CREATE INDEX IF NOT EXISTS idx_notification_events_user_id_scheduled_status ON notification_events(user_id, scheduled_at, status);
CREATE UNIQUE INDEX IF NOT EXISTS uq_notification_events_user_template_scheduled ON notification_events(user_id, template_code, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_support_chat_requests_user_id ON support_chat_requests(user_id);
