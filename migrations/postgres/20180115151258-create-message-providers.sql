
-- +migrate Up
CREATE TABLE message_providers(
  id BIGSERIAL UNIQUE,
  name VARCHAR
);

CREATE TABLE user_message_providers(
  id BIGSERIAL UNIQUE,
  user_id BIGINT,
  message_provider_id BIGINT
);

-- +migrate Down
DROP TABLE message_providers;
DROP TABLE user_message_providers;
