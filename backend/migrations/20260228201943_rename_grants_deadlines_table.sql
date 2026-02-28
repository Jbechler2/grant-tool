-- +goose Up
ALTER TABLE grants_deadlines RENAME TO grant_deadlines;

-- +goose Down
ALTER TABLE grant_deadlines RENAME TO grants_deadlines
