CREATE TABLE
    IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        hero_img_url TEXT,
        content TEXT,
        likes INTEGER,
        is_public BOOLEAN,
        created_at DATETIME,
        updated_at DATETIME
    );