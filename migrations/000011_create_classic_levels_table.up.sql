CREATE TABLE classic_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    publisher VARCHAR(255) NOT NULL,
    difficulty difficulty NOT NULL,
    duration duration NOT NULL,
    ranking INT NOT NULL,
    list_percentage INT NOT NULL,
    points NUMERIC(10,2) NOT NULL,
    minimum_points NUMERIC(10,2) NOT NULL,
    youtube_link TEXT NOT NULL,
    thumbnail_path TEXT NOT NULL
);
