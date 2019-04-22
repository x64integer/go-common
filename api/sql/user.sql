-- table users
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username varchar(50) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(255),
    created_at integer DEFAULT EXTRACT(EPOCH FROM NOW()),
    UNIQUE(username),
    UNIQUE(email)
);

-- function create_user
CREATE OR REPLACE FUNCTION create_user(_username varchar(50), _email varchar(255), _password varchar(255))
RETURNS VOID
LANGUAGE SQL
AS $$
    INSERT INTO users VALUES (DEFAULT, _username, _email, _password);
$$;

-- function get_by_email
CREATE OR REPLACE FUNCTION get_by_email(_email varchar(255))
RETURNS TABLE (username varchar(50), email varchar(255), password varchar(255))
LANGUAGE SQL
AS $$
    SELECT username, email, password FROM users WHERE email = _email LIMIT 1;
$$;