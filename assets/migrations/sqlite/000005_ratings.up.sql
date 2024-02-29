CREATE TABLE IF NOT EXISTS ratings (
    id TEXT NOT NULL PRIMARY KEY,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    user_id TEXT NOT NULL,
    video_id TEXT NOT NULL,

    is_like INTEGER NOT NULL DEFAULT 0,

    UNIQUE(user_id, video_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE
);
