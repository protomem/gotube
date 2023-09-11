BEGIN;

CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    title TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,

    thumbnail_path TEXT NOT NULL,
    video_path TEXT NOT NULL,

    is_public BOOLEAN NOT NULL DEFAULT true,

    views INTEGER NOT NULL DEFAULT 0,

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

COMMIT;
