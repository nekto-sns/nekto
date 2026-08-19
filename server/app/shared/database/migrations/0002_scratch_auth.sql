CREATE TABLE IF NOT EXISTS scratch_auth (
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    scratch_name user_name UNIQUE NOT NULL
);
