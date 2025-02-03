CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    song TEXT NOT NULL,
    group_name TEXT NOT NULL,
    release_date DATE,
    text TEXT,
    link TEXT
);