-- +goose Up
CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  grant_writer_id UUID NOT NULL,
  token TEXT NOT NULL UNIQUE,
  user_agent VARCHAR(255),
  ip_address inet,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '7 days',
  CONSTRAINT fk_refresh_tokens_users
    FOREIGN KEY(grant_writer_id)
     REFERENCES users(id)
);


CREATE INDEX idx_refresh_tokens_grant_writer_id ON refresh_tokens(grant_writer_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
