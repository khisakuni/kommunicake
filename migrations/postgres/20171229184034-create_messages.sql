
-- +migrate Up
CREATE TABLE users(
  id BIGSERIAL,
  name VARCHAR,
  encrypted_password VARCHAR,
  salt VARCHAR
);
CREATE TABLE messages(
  id BIGSERIAL,
  body TEXT,
  recipient_id BIGINT,
  sender_id BIGINT
);

-- +migrate Down
DROP TABLE users;
DROP TABLE messages;
