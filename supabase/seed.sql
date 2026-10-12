-- =============================================================================
-- WEYANA PLATFORM — SEED DATA FOR SUPABASE
-- =============================================================================

-- Seed Archive Categories
INSERT INTO public.archive_categories (cat_id, name, "order")
VALUES 
    (1, 'محاضرات PDF', 1),
    (2, 'ملخصات PDF', 2),
    (3, 'محاضرات فيديوية', 3),
    (4, 'محاضرات مترجمة', 4)
ON CONFLICT (cat_id) DO NOTHING;

-- Seed Platform Settings
INSERT INTO public.platform_settings (key, value)
VALUES 
    ('site_config', '{"platform_name": "WEYANA | ويانا", "version": "1.0.0", "maintenance_mode": false}'::jsonb),
    ('bot_config', '{"welcome_message_enabled": true, "broadcast_rate_limit": 30}'::jsonb)
ON CONFLICT (key) DO NOTHING;
