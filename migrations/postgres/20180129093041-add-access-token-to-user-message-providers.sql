
-- +migrate Up
ALTER TABLE user_message_providers ADD access_token VARCHAR; 
ALTER TABLE user_message_providers ADD token_type VARCHAR; 
ALTER TABLE user_message_providers ADD expiry TIMESTAMP WITH TIME ZONE; 

-- +migrate Down
ALTER TABLE user_message_providers DROP access_token;
ALTER TABLE user_message_providers DROP token_type;
ALTER TABLE user_message_providers DROP expiry;
