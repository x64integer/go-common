-- table users
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username varchar(50) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(255),
    status integer,
    activation_token varchar(255) DEFAULT '',
    created_at integer DEFAULT EXTRACT(EPOCH FROM NOW()),
    UNIQUE(username),
    UNIQUE(email)
);

-- table password_reset
DROP TABLE IF EXISTS password_reset;
CREATE TABLE password_reset (
    email varchar(255) NOT NULL,
    token varchar(255) DEFAULT NULL,
    create_at integer DEFAULT EXTRACT(EPOCH FROM NOW()),
    UNIQUE(email)
);

-- function create_user
CREATE OR REPLACE FUNCTION create_user(_username varchar(50), _email varchar(255), _password varchar(255), _status integer, _activation_token varchar(255))
RETURNS INTEGER
LANGUAGE SQL
AS $$
    INSERT INTO users VALUES (DEFAULT, _username, _email, _password, _status, _activation_token) RETURNING id;
$$;

-- function activate
CREATE OR REPLACE FUNCTION activate(_activation_token varchar(255))
RETURNS VOID
LANGUAGE SQL
AS $$
    UPDATE users SET status = 1, activation_token = '' WHERE activation_token = _activation_token;
$$;

-- function get_by_email
CREATE OR REPLACE FUNCTION get_by_email(_email varchar(255))
RETURNS TABLE (id integer, username varchar(50), email varchar(255), password varchar(255), status integer, activation_token varchar(255))
LANGUAGE SQL
AS $$
    SELECT id, username, email, password, status, activation_token FROM users WHERE email = _email LIMIT 1;
$$;

-- function get_by_activation_token
CREATE OR REPLACE FUNCTION get_by_activation_token(_activation_token varchar(255))
RETURNS TABLE (id integer, username varchar(50), email varchar(255), password varchar(255), status integer, activation_token varchar(255))
LANGUAGE SQL
AS $$
    SELECT id, username, email, password, status, activation_token FROM users WHERE activation_token = _activation_token LIMIT 1;
$$;

-- function update_password
CREATE OR REPLACE FUNCTION update_password(_email varchar(255), _password varchar(255))
RETURNS VOID
LANGUAGE SQL
AS $$
    UPDATE users SET password = _password WHERE email = _email;
$$;

-- function create_reset_token
CREATE OR REPLACE FUNCTION create_reset_token(_email varchar(255), _token varchar(255))
RETURNS VOID
LANGUAGE SQL
AS $$
    INSERT INTO password_reset VALUES (_email, _token);
$$;

-- function get_reset_token
CREATE OR REPLACE FUNCTION get_reset_token(_email varchar(255))
RETURNS varchar(255)
LANGUAGE SQL
AS $$
    SELECT token FROM password_reset WHERE email = _email LIMIT 1;
$$;

-- function get_by_token
CREATE OR REPLACE FUNCTION get_by_token(_token varchar(255))
RETURNS varchar(255)
LANGUAGE SQL
AS $$
    SELECT email FROM password_reset WHERE token = _token LIMIT 1;
$$;

-- function update_reset_token
CREATE OR REPLACE FUNCTION update_reset_token(_email varchar(255), _token varchar(255))
RETURNS VOID
LANGUAGE SQL
AS $$
    UPDATE password_reset SET token = _token WHERE email = _email;
$$;

-- function delete_reset_token
CREATE OR REPLACE FUNCTION delete_reset_token(_token varchar(255))
RETURNS VOID
LANGUAGE SQL
AS $$
    DELETE FROM password_reset WHERE token = _token;
$$;