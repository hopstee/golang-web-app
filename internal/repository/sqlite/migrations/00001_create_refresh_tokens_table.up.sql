CREATE TABLE
    IF NOT EXISTS refresh_tokens (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        token TEXT,
        user_id INTEGER,
        expires_at DATETIME,
        is_revoked BOOLEAN,
        created_at DATETIME,
        device_id TEXT
    );