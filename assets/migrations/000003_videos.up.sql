BEGIN;

CREATE TABLE IF NOT EXISTS videos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  title TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,

  thumbnail_path TEXT NOT NULL,
  video_path TEXT NOT NULL,

  author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  is_public BOOLEAN NOT NULL DEFAULT false,
  views BIGINT NOT NULL DEFAULT 0
);

COMMIT;
