-- +goose Up
CREATE TYPE grant_visibility AS ENUM ('private', 'public');

CREATE TABLE grants (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  grant_writer_id UUID NOT NULL,
  title VARCHAR(255) NOT NULL,
  funder_name VARCHAR(255) NOT NULL,
  funder_website VARCHAR(255),
  description TEXT,
  award_amount_min NUMERIC(12, 2),
  award_amount_max NUMERIC(12, 2),
  eligibility_notes TEXT,
  estimated_application_hours NUMERIC(5, 2),
  visibility grant_visibility NOT NULL DEFAULT 'private',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT fk_grants_grant_writer
    FOREIGN KEY(grant_writer_id)
      REFERENCES users(id)
);

CREATE INDEX idx_grants_title on grants(title);
CREATE INDEX idx_grants_grant_writer_id on grants(grant_writer_id);

CREATE TYPE grant_deadline_type AS ENUM ('LOI', 'application', 'other');

CREATE TABLE grants_deadlines (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  grant_id UUID NOT NULL,
  label grant_deadline_type NOT NULL,
  date DATE NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_grants_deadlines_grants
    FOREIGN KEY(grant_id)
      REFERENCES grants(id)
);

CREATE INDEX idx_grants_deadlines_grant_id on grants_deadlines(grant_id);

-- +goose Down
DROP TABLE IF EXISTS grants_deadlines;
DROP TYPE IF EXISTS grant_deadline_type;
DROP TABLE IF EXISTS grants;
DROP TYPE IF EXISTS grant_visibility;
