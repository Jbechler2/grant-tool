-- +goose Up
ALTER TABLE users
  RENAME COLUMN password_has TO password_hash;

ALTER TABLE users
  ALTER COLUMN email TYPE VARCHAR(255);

-- +goose Down
ALTER TABLE users
  RENAME COLUMN password_hash TO password_has;

ALTER TABLE users
  ALTER COLUMN email TYPE VARCHAR(225);
