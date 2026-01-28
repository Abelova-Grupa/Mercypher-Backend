CREATE SCHEMA IF NOT EXISTS user_service;
CREATE TABLE IF NOT EXISTS user_service.contacts (
    first_username text NOT NULL PRIMARY KEY,
    second_username text NOT NULL PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    -- CONSTRAINT contact_pk PRIMARY KEY (first_username, second_username),
    CONSTRAINT no_self_contact CHECK (first_username <> second_username),
    CONSTRAINT first_fk FOREIGN KEY(first_username) REFERENCES user_service.users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT second_fk FOREIGN KEY(second_username) REFERENCES user_service.users(username) ON DELETE CASCADE ON UPDATE CASCADE
);
