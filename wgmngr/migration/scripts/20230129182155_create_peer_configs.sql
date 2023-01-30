-- +goose Up
-- +goose StatementBegin
CREATE TABLE `peer_configs` (
    `id` CHAR(64) NOT NULL UNIQUE,
    `name` VARCHAR(256) NOT NULL,
    `description` VARCHAR(10000) NOT NULL,
    `generated_by_id` CHAR(64) NOT NULL,
    `generated_at` DATETIME NOT NULL,
    `ipv4` VARCHAR(15) NOT NULL UNIQUE,
    `ipv6` VARCHAR(39) NOT NULL UNIQUE,
    `private_key` CHAR(44) NOT NULL UNIQUE,
    `public_key` CHAR(44) NOT NULL UNIQUE,
    `preshared_key` CHAR(44) NOT NULL UNIQUE,
    `is_active` BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_id PRIMARY KEY (`id`),
    FULLTEXT `name_full_text` (`name`),
    FULLTEXT `description_full_text` (`name`),
    CONSTRAINT fk_generated_by_id FOREIGN KEY (`generated_by_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `peer_configs`;
-- +goose StatementEnd
