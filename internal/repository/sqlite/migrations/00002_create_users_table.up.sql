CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT,
        password TEXT,
        email TEXT,
        token_version INTEGER,
        created_at DATETIME,
        updated_at DATETIME
    );