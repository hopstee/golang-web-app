CREATE TABLE
    IF NOT EXISTS requests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        phone TEXT,
        email TEXT,
        contact_type TEXT,
        message TEXT,
        created_at DATETIME,
        updated_at DATETIME
    );