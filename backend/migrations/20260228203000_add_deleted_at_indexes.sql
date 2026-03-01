-- +goose Up
CREATE INDEX idx_clients_deleted_at ON clients(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_grants_deleted_at ON grants(deleted_at) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_clients_deleted_at;
DROP INDEX IF EXISTS idx_grants_deleted_at;
