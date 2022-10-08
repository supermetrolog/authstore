-- +goose Up
-- +goose StatementBegin
CREATE TABLE `tree` (
  `id` int(11) UNSIGNED NOT NULL  AUTO_INCREMENT,
  `parent_id` int(11) UNSIGNED NOT NULL,
  `user_id` int(11) UNSIGNED NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `type` tinyint(1) COMMENT 'лист или узел', 
  `status` tinyint(1), 
  PRIMARY KEY (`id`),
  INDEX `idx-tree-parent_id` (parent_id),
  INDEX `idx-tree-user_id` (user_id),
  FOREIGN KEY (parent_id) REFERENCES tree(id) ON  DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES user(id) ON  DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `tree`
-- +goose StatementEnd
