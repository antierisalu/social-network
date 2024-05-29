CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    -- Password (initially for testing, edit this table later)
    password TEXT NOT NULL
); 