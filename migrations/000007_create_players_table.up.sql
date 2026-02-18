CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(50) UNIQUE,
    classic_points NUMERIC(8, 2) DEFAULT 0,
    platformer_points NUMERIC(8, 2) DEFAULT 0,
    avatar VARCHAR(255),
    region_id UUID REFERENCES regions(id) ON DELETE SET NULL,
    discord VARCHAR(50) UNIQUE,
    youtube VARCHAR(50) UNIQUE,
    twitter VARCHAR(50) UNIQUE,
    twitch VARCHAR(50) UNIQUE,
    is_flagged BOOLEAN DEFAULT false
);
