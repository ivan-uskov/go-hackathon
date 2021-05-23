CREATE TABLE `scoring_task` (
       `task_id` BINARY(16),
       `solution_id` BINARY(16) UNIQUE,
       `endpoint` VARCHAR(255) NOT NULL,
       `score` INTEGER NOT NULL DEFAULT 0,
       `created_at` DATETIME NOT NULL,
       `scored_at` DATETIME,
       PRIMARY KEY (task_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci