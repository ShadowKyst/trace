CREATE TABLE IF NOT EXISTS pastes (
    id VARCHAR(20) PRIMARY KEY,
    content TEXT NOT NULL,
    language VARCHAR(50) NOT NULL,
    hash VARCHAR(64),
    views BIGINT DEFAULT 0,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pastes_expires_at ON pastes(expires_at);