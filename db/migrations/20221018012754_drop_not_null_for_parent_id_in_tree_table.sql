-- +goose Up
-- +goose StatementBegin
ALTER TABLE tree CHANGE COLUMN parent_id parent_id INT(11) UNSIGNED;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tree CHANGE COLUMN parent_id parent_id INT(11) UNSIGNED NOT NULL;
-- +goose StatementEnd
