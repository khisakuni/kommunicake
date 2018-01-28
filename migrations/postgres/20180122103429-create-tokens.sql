
-- +migrate Up
CREATE TABLE tokens(
  id BIGSERIAL UNIQUE,
  value VARCHAR,
  user_id BIGINT,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE
);

-- +migrate Down
DROP TABLE tokens;
