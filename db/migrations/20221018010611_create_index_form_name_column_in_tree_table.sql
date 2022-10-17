-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX `idx-tree-name` ON tree(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX `idx-tree-name` ON tree;
-- +goose StatementEnd
