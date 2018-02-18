
-- +migrate Up
ALTER TABLE user_message_providers ADD history_id INTEGER;

-- +migrate Down
ALTER TABLE user_message_providers DROP history_id;
