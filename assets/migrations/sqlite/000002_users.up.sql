CREATE TABLE IF NOT EXISTS users (
    id         TEXT NOT NULL PRIMARY KEY,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    nickname TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,

    email        TEXT NOT NULL UNIQUE,
    is_verified INTEGER NOT NULL DEFAULT 0,

    avatar_path TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
);