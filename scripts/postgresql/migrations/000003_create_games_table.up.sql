BEGIN;

CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    url TEXT,
    age_rating INT NOT NULL,
    publisher TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;