CREATE TABLE classic_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    classic_level_id UUID NOT NULL UNIQUE REFERENCES classic_levels(id) ON DELETE CASCADE,
    player_id UUID NOT NULL UNIQUE REFERENCES players(id) ON DELETE CASCADE,
    progress INT NOT NULL,
    device device,
    footage TEXT,
    raw_footage TEXT,
    completed_at TIMESTAMPTZ DEFAULT NOW()
);
