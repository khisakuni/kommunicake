
-- +migrate Up
CREATE TABLE user_message_providers(
  id BIGSERIAL UNIQUE,
  user_id BIGINT,
  message_provider_type VARCHAR,
  refresh_token VARCHAR
);

-- +migrate Down
DROP TABLE user_message_providers;
