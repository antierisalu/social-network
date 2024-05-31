CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    email TEXT NOT NULL UNIQUE,
    hash TEXT NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    date_of_birth INT NOT NULL,

    avatar BLOB,
    nickname TEXT,
    about TEXT,

    privacy INT DEFAULT 0 CHECK (privacy IN (0, 1)), -- 0 means public, 1 means private
    lastonline DATE DEFAULT CURRENT_DATE,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    session TEXT UNIQUE,

    CONSTRAINT unique_session CHECK (session IS NOT NULL AND session != '')
);
