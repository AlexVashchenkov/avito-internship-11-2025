-- +goose Up
INSERT INTO pr_statuses (status)
VALUES ('OPEN'), ('MERGED')
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM pr_statuses WHERE status IN ('OPEN', 'MERGED');