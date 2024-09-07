CREATE TABLE IF NOT EXISTS urls (
    short_url varchar PRIMARY KEY,
    long_url varchar NOT NULL,
    user_id int NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    expires_at timestamp NOT NULL,
);

CREATE TABLE IF NOT EXISTS users (
    id int PRIMARY KEY,
    email varchar NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    urls_left int NOT NULL

);

CREATE TABLE IF NOT EXISTS subscriptions (
    id int PRIMARY KEY,
    name varchar NOT NULL,
    total_urls int NOT NULL
);