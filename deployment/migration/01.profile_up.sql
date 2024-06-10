-- Note: No FOREIGN KEY constraints due to future potential module separation

-- Profiles table
-- Note: No FOREIGN KEY constraints due to future potential module separation

-- Profiles table
CREATE TABLE profiles (
    id BIGINT UNSIGNED PRIMARY KEY COMMENT "Profile unique identifier, equal to id from iam service",

    nickname VARCHAR(255) COMMENT "User's nick name",
    email VARCHAR(255) UNIQUE default NULL,
    avatar_url VARCHAR(255) COMMENT "URL to user's avatar",
    bio TEXT,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT "Record creation time in UTC",
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "Record update time in UTC",
    deleted_at DATETIME DEFAULT NULL COMMENT "Record delete time in UTC",

    INDEX `idx_profiles_created_at` (`created_at`) -- For sorting/filtering by creation date
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Settings table - Memorization
CREATE TABLE profile_memorization_settings (
    id BIGINT UNSIGNED PRIMARY KEY COMMENT "rofile unique identifier, equal to id from iam service",

    review_interval_setting VARCHAR(255) COMMENT "Interval for review in days",

    difficulty_preference TINYINT UNSIGNED COMMENT "User's preference for difficulty",

    quiz_mode VARCHAR(255) COMMENT "Preferred quiz mode"

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Settings table - Advance
CREATE TABLE profile_advance_settings (
    id BIGINT UNSIGNED PRIMARY KEY COMMENT "rofile unique identifier, equal to id from iam service",

    theme VARCHAR(255) DEFAULT 'light',
    language VARCHAR(255) DEFAULT 'en',

    email_notifications BOOLEAN DEFAULT TRUE COMMENT "Flag for email notification settings",
    push_notifications BOOLEAN DEFAULT TRUE COMMENT "Flag for push notification settings"

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Points table - Advance
CREATE TABLE profile_points (
    id BIGINT UNSIGNED PRIMARY KEY COMMENT "rofile unique identifier, equal to id from iam service",

    cash BIGINT UNSIGNED DEFAULT 0,
    gem BIGINT UNSIGNED DEFAULT 0,
    vip_score  BIGINT UNSIGNED DEFAULT 0,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT "Record creation time in UTC",
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "Record update time in UTC",
    deleted_at DATETIME DEFAULT NULL COMMENT "Record delete time in UTC"

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
