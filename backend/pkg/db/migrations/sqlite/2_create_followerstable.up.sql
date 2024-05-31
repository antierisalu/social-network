CREATE TABLE IF NOT EXISTS followers (
    user_id INTEGER NOT NULL,
    follower_id INTEGER NOT NULL,
    isFollowing BOOLEAN NOT NULL,
    UNIQUE(user_id, follower_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(follower_id) REFERENCES users(id) ON DELETE CASCADE
);

