-- +goose Up
CREATE TABLE clients (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  grant_writer_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  contact_name VARCHAR(255),
  contact_phone varchar(15),
  contact_email varchar(255),
  notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT fk_clients_grant_writer
    FOREIGN KEY(grant_writer_id)
      REFERENCES users(id)
);

CREATE INDEX idx_clients_name ON clients(name);
CREATE INDEX idx_clients_grant_writer_id ON clients(grant_writer_id);

-- +goose Down
DROP TABLE IF EXISTS clients;
