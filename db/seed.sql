-- Create tables
CREATE TABLE IF NOT EXISTS users (
    discord_id BIGINT PRIMARY KEY,
    discord_name TEXT,
    discord_global_name TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS doujins (
    melonbooks_id BIGINT PRIMARY KEY,
    title TEXT NOT NULL,
    price_in_yen INTEGER NOT NULL,
    price_in_usd REAL NOT NULL,
    is_r18 BOOLEAN NOT NULL,
    image_preview_url TEXT NOT NULL,
    url TEXT NOT NULL,
    circle TEXT,
    authors TEXT[10],
    genres TEXT[10],
    events TEXT[10],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reservations (
    reservation_id SMALLSERIAL PRIMARY KEY,
    discord_id BIGINT REFERENCES users(discord_id),
    melonbooks_id BIGINT REFERENCES doujins(melonbooks_id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

