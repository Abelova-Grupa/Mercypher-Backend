CREATE SCHEMA IF NOT EXISTS user_service;
CREATE TABLE IF NOT EXISTS user_service.users (
    username text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz DEFAULT now(),
    validated bool NOT NULL DEFAULT false,
    auth_code text NULL,
    CONSTRAINT users_pkey PRIMARY KEY (username)
);
