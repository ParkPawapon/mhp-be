DROP TABLE IF EXISTS support_chat_requests;
DROP TABLE IF EXISTS notification_events;
DROP TABLE IF EXISTS notification_templates;
DROP TABLE IF EXISTS device_tokens;
DROP TABLE IF EXISTS user_preferences;

ALTER TABLE patient_medicines
    DROP COLUMN IF EXISTS category_item_id;

DROP TABLE IF EXISTS medicine_category_items;
DROP TABLE IF EXISTS medicine_categories;

DROP TYPE IF EXISTS notification_status;
