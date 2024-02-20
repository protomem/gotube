CREATE TABLE IF NOT EXISTS videos (
    id TEXT NOT NULL PRIMARY KEY,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    title TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',

    thumbnail_path TEXT NOT NULL,
    video_path TEXT NOT NULL,

    author_id TEXT NOT NULL,

    is_public INTEGER NOT NULL DEFAULT 0,
    views INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE
;
