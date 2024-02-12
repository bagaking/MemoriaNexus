-- Note: No FOREIGN KEY constraints due to future potential domain separation

-- Profiles table
CREATE TABLE `profiles` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT "rofile unique identifier (UUID)",
    `user_id` BIGINT UNSIGNED UNIQUE NOT NULL,

    `nickname` VARCHAR(255) COMMENT "User's nick name",
    `avatar_url` VARCHAR(255) COMMENT "URL to user's avatar",

    `more_info_id` BIGINT UNSIGNED COMMENT "User security settings ID",
    `notification_settings_id` BIGINT UNSIGNED COMMENT "User notification settings ID",
    `memorization_settings_id` BIGINT UNSIGNED COMMENT "User memorization settings ID",
    `preference_id` BIGINT UNSIGNED COMMENT "User preference settings ID",

    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT "Record creation time",
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT "Record update time",
    `deleted_at` DATETIME(3) COMMENT '记录软删除时间',

    INDEX `idx_profiles_created_at` (`created_at`) -- For sorting/filtering by creation date
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- SecuritySettings table
CREATE TABLE `profile_more_infos` (
                                     `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                     `bio` TEXT COMMENT "User's biography"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- NotificationSettings table
CREATE TABLE `profile_notification_settings` (
                                                 `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                                 `email_notifications` BOOLEAN DEFAULT TRUE COMMENT "Flag for email notification settings",
                                                 `push_notifications` BOOLEAN DEFAULT TRUE COMMENT "Flag for push notification settings"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MemorizationSettings table
CREATE TABLE `profile_memorization_settings` (
                                                 `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                                 `review_interval` INT UNSIGNED COMMENT "Interval for review in days",
                                                 `difficulty_preference` TINYINT UNSIGNED COMMENT "User's preference for difficulty",
                                                 `quiz_mode` VARCHAR(255) COMMENT "Preferred quiz mode"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Preference table
CREATE TABLE `profile_preferences` (
                                      `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                      `theme` VARCHAR(255) COMMENT "User's preferred UI theme",
                                      `language` VARCHAR(255) COMMENT "User's preferred language"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
