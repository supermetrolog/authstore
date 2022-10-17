-- +goose Up
-- +goose StatementBegin
ALTER TABLE tree ADD name VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tree DROP COLUMN name;
-- +goose StatementEnd
