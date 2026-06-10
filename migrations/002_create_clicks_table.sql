CREATE TABLE IF NOT EXISTS clicks (
    id BIGSERIAL PRIMARY KEY,
    short_code TEXT NOT NULL REFERENCES urls(short_code),
    ip_address TEXT,
    user_agent TEXT,
    referer TEXT,
    clicked_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_clicks_short_code ON clicks(short_code);
