CREATE TABLE `tags` (
    `id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `books` (
    `id` BIGINT UNSIGNED NOT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,

    `title` VARCHAR(255) NOT NULL,
    `description` TEXT,

    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    deleted_at DATETIME DEFAULT NULL COMMENT "Record delete time in UTC",
    PRIMARY KEY (`id`),
    INDEX `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `book_tags` (
    `book_id` BIGINT UNSIGNED NOT NULL,
    `tag_id` BIGINT UNSIGNED NOT NULL,
--     FOREIGN KEY (`book_id`) REFERENCES `books`(`id`) ON DELETE CASCADE,
--     FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`) ON DELETE CASCADE,
    PRIMARY KEY (`book_id`, `tag_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `items` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT "The Creator",
    `creator_id` BIGINT UNSIGNED NOT NULL,

    `type` VARCHAR(50),
    `content` TEXT,

    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    `difficulty` TINYINT UNSIGNED COMMENT "Difficulty",
    `importance` TINYINT UNSIGNED COMMENT "Importance",

    deleted_at DATETIME DEFAULT NULL COMMENT "Record delete time in UTC",
    PRIMARY KEY (`id`),
    INDEX idx_user_items (`creator_id`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `item_tags` (
    `item_id` BIGINT UNSIGNED NOT NULL,
    `tag_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`item_id`, `tag_id`),
    INDEX `idx_tag_items` (`tag_id`, `item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `book_items` (
    `book_id` BIGINT UNSIGNED NOT NULL,
    `item_id` BIGINT UNSIGNED NOT NULL,
--     FOREIGN KEY (`book_id`) REFERENCES `books`(`id`) ON DELETE CASCADE,
--     FOREIGN KEY (`item_id`) REFERENCES `items`(`id`) ON DELETE CASCADE,
    PRIMARY KEY (`book_id`, `item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;