CREATE TABLE IF NOT EXISTS comments (
    id TEXT NOT NULL PRIMARY KEY,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

	message TEXT NOT NULL,

	video_id TEXT NOT NULL,
	author_id TEXT NOT NULL,

	FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE,
	FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

