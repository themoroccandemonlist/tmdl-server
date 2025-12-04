CREATE TABLE IF NOT EXISTS regions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    classic_points NUMERIC(8, 2) DEFAULT 0,
    platformer_points NUMERIC(8, 2) DEFAULT 0
);
