CREATE TABLE
    IF NOT EXISTS refresh_tokens (
        id SERIAL PRIMARY KEY,
        token TEXT NOT NULL,
        user_id BIGINT NOT NULL,
        expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
        is_revoked BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
        device_id TEXT
    );