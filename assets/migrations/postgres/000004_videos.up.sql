BEGIN;

CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    title       TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',

    thumbnail_path TEXT NOT NULL DEFAULT '',
    video_path     TEXT NOT NULL DEFAULT '',

    author_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,

    is_public BOOLEAN NOT NULL DEFAULT false,
    views     INTEGER NOT NULL DEFAULT 0
);

COMMIT;
