CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    subtitle TEXT,
    content TEXT NOT NULL,
    tags TEXT[] NOT NULL CHECK (array_length(tags, 1) BETWEEN 1 AND 5),
    claps BIGINT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ(0),
    slug TEXT UNIQUE,
    version INTEGER NOT NULL DEFAULT 1
);