CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    owner_id INTEGER NOT NULL,
    media TEXT,
    chat_id INTEGER NOT NULL,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS group_members (
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status INTEGER CHECK (status BETWEEN 0 AND 2),
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    chat_seen INTEGER DEFAULT 0,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
