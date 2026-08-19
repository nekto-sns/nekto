CREATE EXTENSION IF NOT EXISTS citext;

CREATE DOMAIN user_name AS citext
CHECK (
    char_length(VALUE) >= 2 AND
    char_length(VALUE) <= 16 AND
    VALUE ~ '^[a-zA-Z0-9_]+$'
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name user_name UNIQUE NOT NULL,
    display_name VARCHAR(32) NOT NULL,
    bio VARCHAR(256),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_modtime
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

