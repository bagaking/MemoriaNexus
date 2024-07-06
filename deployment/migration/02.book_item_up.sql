CREATE TABLE tags (
                      `user_id` BIGINT NOT NULL,
                      `tag` VARCHAR(2048) NOT NULL,  -- 注意如果加密的话要支持前缀可索引
                      `entity_type` TINYINT UNSIGNED COMMENT "entity type, item=1, book=2, dungeon=3",
                      `entity_id` BIGINT NOT NULL,
                      `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
                      `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      `deleted_at` DATETIME NULL DEFAULT NULL COMMENT "Record delete time in UTC",

                      PRIMARY KEY (user_id, tag(32), entity_id),
                      INDEX `idx_entity_tag_type` (entity_id, tag(32), entity_type)
);

CREATE TABLE `books` (
                         `id` BIGINT UNSIGNED NOT NULL,
                         `user_id` BIGINT UNSIGNED NOT NULL,

                         `title` VARCHAR(255) NOT NULL,
                         `description` TEXT,

                         `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
                         `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                         `deleted_at` DATETIME DEFAULT NULL COMMENT "Record delete time in UTC",

                         PRIMARY KEY (`id`),
                         INDEX `idx_user` (`user_id`)
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

CREATE TABLE `book_items` (
                              `book_id` BIGINT UNSIGNED NOT NULL,
                              `item_id` BIGINT UNSIGNED NOT NULL,
--     FOREIGN KEY (`book_id`) REFERENCES `books`(`id`) ON DELETE CASCADE,
--     FOREIGN KEY (`item_id`) REFERENCES `items`(`id`) ON DELETE CASCADE,
                              PRIMARY KEY (`book_id`, `item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;