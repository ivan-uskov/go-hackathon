CREATE TABLE `session` (
                         `session_id` BINARY(16),
                         `code` VARCHAR(255) NOT NULL,
                         `code_hash` BINARY(16) UNIQUE NOT NULL,
                         `name` VARCHAR(255) NOT NULL,
                         `type` INTEGER NOT NULL,
                         `created_at` DATETIME NOT NULL,
                         `updated_at` DATETIME NOT NULL,
                         `closed_at` DATETIME,
                         PRIMARY KEY (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `session_participant` (
                              `participant_id` BINARY(16),
                              `session_id` BINARY(16),
                              `name` VARCHAR(255) NOT NULL,
                              `name_hash` BINARY(16) UNIQUE NOT NULL,
                              `endpoint` VARCHAR(255) NOT NULL,
                              `score` INTEGER NOT NULL DEFAULT 0,
                              `created_at` DATETIME NOT NULL,
                              `scored_at` DATETIME,
                              PRIMARY KEY (participant_id),
                              FOREIGN KEY (`session_id`) REFERENCES `session` (`session_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci