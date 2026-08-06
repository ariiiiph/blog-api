CREATE TABLE users_old (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

INSERT INTO users_old (id, email, password_hash, created_at)
SELECT id, email, password_hash, created_at FROM users;

DROP TABLE users;

ALTER TABLE users_old RENAME TO users;