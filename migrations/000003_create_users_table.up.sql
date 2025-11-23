CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    username TEXT UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1

)