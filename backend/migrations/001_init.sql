CREATE DATABASE IF NOT EXISTS task_management;
USE task_management;

CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status ENUM('todo','in_progress','done') NOT NULL DEFAULT 'todo',
    assignee_id BIGINT UNSIGNED NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_tasks_title_active (title, deleted_at),
    KEY idx_tasks_status (status),
    KEY idx_tasks_assignee (assignee_id),
    KEY idx_tasks_deleted_at (deleted_at),
    KEY idx_tasks_created_at (created_at)
);

INSERT INTO tasks (title, description, status, assignee_id)
SELECT 'Build landing page', 'Create responsive landing page', 'todo', 1
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Build landing page' AND deleted_at IS NULL);

INSERT INTO tasks (title, description, status, assignee_id)
SELECT 'Fix checkout bug', 'Resolve checkout validation issue', 'in_progress', 2
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Fix checkout bug' AND deleted_at IS NULL);

INSERT INTO tasks (title, description, status, assignee_id)
SELECT 'Write API tests', 'Add coverage for task endpoints', 'done', 1
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Write API tests' AND deleted_at IS NULL);
