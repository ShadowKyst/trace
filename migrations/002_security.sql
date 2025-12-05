ALTER TABLE pastes ADD COLUMN password_hash VARCHAR(255);
ALTER TABLE pastes ADD COLUMN burn_after_reading BOOLEAN DEFAULT FALSE;

-- Гарантируем, что expires_at не может быть NULL (или ставим дефолт в коде)
-- Но для совместимости старых данных оставим как есть, просто будем форсировать в коде.