-- +goose Up
ALTER TABLE refresh_tokens
DROP CONSTRAINT fk_refresh_tokens_users;

ALTER TABLE refresh_tokens
ADD CONSTRAINT fk_refresh_tokens_users
    FOREIGN KEY (grant_writer_id)
    REFERENCES users(id)
    ON DELETE CASCADE;

ALTER TABLE clients
DROP CONSTRAINT fk_clients_grant_writer;

ALTER TABLE clients
ADD CONSTRAINT fk_clients_grant_writer
  FOREIGN KEY(grant_writer_id)
    REFERENCES users(id)
    ON DELETE CASCADE;

ALTER TABLE grants
DROP CONSTRAINT fk_grants_grant_writer;

ALTER TABLE grants
ADD CONSTRAINT fk_grants_grant_writer
  FOREIGN KEY(grant_writer_id)
    REFERENCES users(id)
      ON DELETE CASCADE;

ALTER TABLE applications
DROP CONSTRAINT fk_applications_users;

ALTER TABLE applications
ADD CONSTRAINT fk_applications_users
  FOREIGN KEY(grant_writer_id)
    REFERENCES users(id)
    ON DELETE CASCADE;

ALTER TABLE grant_deadlines
DROP CONSTRAINT fk_grants_deadlines_grants;

ALTER TABLE grant_deadlines
ADD CONSTRAINT fk_grant_deadlines_grants
 FOREIGN KEY(grant_id)
  REFERENCES grants(id)
  ON DELETE CASCADE;


-- +goose Down
ALTER TABLE refresh_tokens
DROP CONSTRAINT fk_refresh_tokens_users;

ALTER TABLE refresh_tokens
ADD CONSTRAINT fk_refresh_tokens_users
  FOREIGN KEY (grant_writer_id)
    REFERENCES users(id);

ALTER TABLE clients
DROP CONSTRAINT fk_clients_grant_writer;

ALTER TABLE clients
ADD CONSTRAINT fk_clients_grant_writer
  FOREIGN KEY(grant_writer_id)
  REFERENCES users(id);

ALTER TABLE grants
DROP CONSTRAINT fk_grants_grant_writer;

ALTER TABLE grants
ADD CONSTRAINT fk_grants_grant_writer
  FOREIGN KEY(grant_writer_id)
  REFERENCES users(id);

ALTER TABLE applications
DROP CONSTRAINT fk_applications_users;

ALTER TABLE applications
ADD CONSTRAINT fk_applications_users
  FOREIGN KEY(grant_writer_id)
    REFERENCES users(id);

ALTER TABLE grant_deadlines
DROP CONSTRAINT fk_grant_deadlines_grants;

ALTER TABLE grant_deadlines
ADD CONSTRAINT fk_grants_deadlines_grants
 FOREIGN KEY(grant_id)
  REFERENCES grants(id);