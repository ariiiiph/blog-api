CREATE TABLE users_new (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

INSERT INTO users_new (id, email, password_hash, created_at)
SELECT id, email, password_hash, created_at FROM users;

DROP TABLE users;

ALTER TABLE users_new RENAME TO users;