CREATE TABLE IF NOT EXISTS subscriptions (
    id         TEXT NOT NULL PRIMARY KEY,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    from_user_id TEXT NOT NULL,
    to_user_id TEXT NOT NULL,

    UNIQUE (from_user_id, to_user_id),

    FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
);