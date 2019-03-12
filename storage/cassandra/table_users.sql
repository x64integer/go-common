CREATE TABLE users (
    id uuid PRIMARY KEY,
    name text,
    email text,
    posted_time timestamp
) WITH comment='Users';