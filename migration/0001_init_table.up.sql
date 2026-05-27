CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,

    filename TEXT NOT NULL,
    filepath TEXT NOT NULL UNIQUE,

    extension TEXT,
    mime_type TEXT,

    size_bytes INTEGER NOT NULL,

    width INTEGER,
    height INTEGER,

    from TEXT, 

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);