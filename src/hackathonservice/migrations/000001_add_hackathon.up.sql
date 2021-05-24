CREATE TABLE `hackathon` (
                         `hackathon_id` BINARY(16),
                         `code` VARCHAR(255) UNIQUE NOT NULL,
                         `name` VARCHAR(255) NOT NULL,
                         `type` INTEGER NOT NULL,
                         `created_at` DATETIME NOT NULL,
                         `closed_at` DATETIME,
                         PRIMARY KEY (hackathon_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `hackathon_participant` (
                              `participant_id` BINARY(16),
                              `hackathon_id` BINARY(16),
                              `name` VARCHAR(255) UNIQUE NOT NULL,
                              `endpoint` VARCHAR(255) NOT NULL,
                              `score` INTEGER NOT NULL DEFAULT 0,
                              `created_at` DATETIME NOT NULL,
                              `scored_at` DATETIME,
                              PRIMARY KEY (participant_id),
                              FOREIGN KEY (`hackathon_id`) REFERENCES `hackathon` (`hackathon_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci