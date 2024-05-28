CREATE TABLE `dungeons` (
    `id` BIGINT UNSIGNED NOT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,

    `type` TINYINT UNSIGNED NOT NULL COMMENT 'Dungeon (<256) 类型分段 campaign = 001, endless = 002, instance > 003',

    `title` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `rule` TEXT COMMENT '复习规则的详细配置, JSON格式',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    `deleted_at` DATETIME DEFAULT NULL COMMENT "Record delete time in UTC",
    PRIMARY KEY (`id`),
    INDEX `idx_user_dungeon` (`user_id`, `id`),
    INDEX `idx_type_user` (`type`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `dungeon_books` (
    `dungeon_id` BIGINT UNSIGNED NOT NULL,
    `book_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`dungeon_id`, `book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `dungeon_tags` (
    `dungeon_id` BIGINT UNSIGNED NOT NULL,
    `tag_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`dungeon_id`, `tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- monster 代表的是 item 对于特定 user 的属性，ID 就用 itemID
CREATE TABLE `monsters` (
    `user_id` BIGINT UNSIGNED NOT NULL,
    `item_id` BIGINT UNSIGNED NOT NULL,

    `familiarity` TINYINT UNSIGNED COMMENT "percentage: 0-100",

    PRIMARY KEY (`user_id`, `item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- dungeon 和 user 是 n 对 1 的，因此 `dungeon_id` + `item_id` 是可以对应到 monster
CREATE TABLE `dungeon_monsters` (
    `dungeon_id` BIGINT UNSIGNED NOT NULL,
    `item_id` BIGINT UNSIGNED NOT NULL,

    `visibility` TINYINT UNSIGNED NOT NULL COMMENT "percentage: 0-100",

    `source_type` TINYINT UNSIGNED COMMENT "source type of the monster, item=1, book=2, tag=3",
    `source_id` BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (`dungeon_id`, `item_id`),
    INDEX idx_dungen_monsters (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
