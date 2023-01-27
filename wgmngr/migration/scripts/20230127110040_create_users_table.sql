-- +goose Up
-- +goose StatementBegin
CREATE TABLE `users` (
    `id` CHAR(64) NOT NULL UNIQUE,
    `name` VARCHAR(256) NOT NULL,
    `username` VARCHAR(64) NOT NULL UNIQUE,
    `password` VARBINARY(10000) NOT NULL,
    `creator_id` CHAR(64) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `role` ENUM('admin', 'reseller') NOT NULL,
    CONSTRAINT pk_id PRIMARY KEY (`id`),
    FULLTEXT `name_full_text` (`name`),
    CONSTRAINT fk_creator_id FOREIGN KEY (`creator_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `users`;
-- +goose StatementEnd
