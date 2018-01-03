
-- +migrate Up
CREATE TABLE users(
  id BIGSERIAL UNIQUE,
  name VARCHAR,
  email VARCHAR UNIQUE,
  encrypted_password VARCHAR
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
