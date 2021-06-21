CREATE TABLE `stored_event` (
    `event_id` BINARY(16),
    `type` VARCHAR(255) NOT NULL,
    `body` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `published_at` DATETIME,
    PRIMARY KEY (event_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci