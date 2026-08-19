ALTER DOMAIN user_name
DROP CONSTRAINT user_name_check;

ALTER DOMAIN user_name
ADD CONSTRAINT user_name_check CHECK (
    char_length(VALUE) >= 2 AND
    char_length(VALUE) <= 16 AND
    VALUE ~ '^[a-zA-Z0-9_-]+$'
);
