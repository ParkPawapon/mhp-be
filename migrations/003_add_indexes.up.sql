CREATE INDEX IF NOT EXISTS idx_notification_events_status_scheduled_at
    ON notification_events(status, scheduled_at);

CREATE INDEX IF NOT EXISTS idx_notification_events_user_template_status
    ON notification_events(user_id, template_code, status);

CREATE INDEX IF NOT EXISTS idx_support_chat_requests_created_at
    ON support_chat_requests(created_at);
