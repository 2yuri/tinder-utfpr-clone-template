CREATE TABLE IF NOT EXISTS interactions (
    id SERIAL,
    user_id INTEGER NOT NULL,
    target_user_id INTEGER NOT NULL,
    liked bool NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    FOREIGN KEY (target_user_id) REFERENCES users (id) ON DELETE RESTRICT,
    UNIQUE(user_id, target_user_id)
);