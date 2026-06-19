CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,

    filename TEXT NOT NULL,
    filepath TEXT NOT NULL UNIQUE,

    extension TEXT,
    mime_type TEXT,

    size_bytes INTEGER NOT NULL,

    width INTEGER,
    height INTEGER,

    `from` TEXT, 

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO images (id, filename, filepath, extension, mime_type, size_bytes, width, height, `from`, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'fluffy.jpg', '/images/fluffy.jpg', 'jpg', 'image/jpeg', 102400, 800, 600, 'camera', CURRENT_TIMESTAMP),
('22222222-2222-2222-2222-222222222222', 'whiskers.png', '/images/whiskers.png', 'png', 'image/png', 204800, 1024, 768, 'scanner', CURRENT_TIMESTAMP),
('33333333-3333-3333-3333-333333333333', 'shadow.gif', '/images/shadow.gif', 'gif', 'image/gif', 51200, 640, 480, 'download', CURRENT_TIMESTAMP);

