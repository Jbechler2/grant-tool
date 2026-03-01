-- +goose Up
ALTER INDEX idx_grants_deadlines_grant_id RENAME TO idx_grant_deadlines_grant_id;

-- +goose Down
ALTER INDEX idx_grant_deadlines_grant_id RENAME TO idx_grants_deadlines_grant_id;
