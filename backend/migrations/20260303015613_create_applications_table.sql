-- +goose Up

CREATE TYPE application_status as ENUM ('not_started', 'draft', 'submitted', 'approved', 'denied', 'withdrawn');

CREATE TABLE applications (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  grant_writer_id UUID NOT NULL,
  grant_id UUID NOT NULL,
  client_id UUID NOT NULL,
  title VARCHAR(255) NOT NULL,
  status application_status NOT NULL DEFAULT 'not_started',
  is_exclusive BOOLEAN NOT NULL DEFAULT false,
  published_at TIMESTAMPTZ,
  notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT fk_applications_users
    FOREIGN KEY(grant_writer_id)
      REFERENCES users(id),
  CONSTRAINT fk_applications_grants
    FOREIGN KEY(grant_id)
      REFERENCES grants(id),
  CONSTRAINT fk_applications_clients
    FOREIGN KEY(client_id)
      REFERENCES clients(id)
);

  CREATE INDEX idx_applications_status ON applications(status);
  CREATE INDEX idx_applications_grant_writer_id on applications(grant_writer_id);
  CREATE INDEX idx_applications_deleted_at ON applications(deleted_at);
  CREATE INDEX idx_applications_published_at ON applications(published_at);

-- +goose Down
DROP TABLE IF EXISTS applications;
DROP TYPE IF EXISTS application_status;
