CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(50) UNIQUE NOT NULL,
    sub VARCHAR(50) UNIQUE NOT NULL,
    is_banned BOOLEAN DEFAULT false,
    is_deleted BOOLEAN DEFAULT false
);
