CREATE SCHEMA IF NOT EXISTS user_service;
CREATE TABLE IF NOT EXISTS user_service.contacts (
    username text NOT NULL,
    contact_name text NOT NULL,
    created_at timestamptz DEFAULT now(),
    CONSTRAINT contact_pk PRIMARY KEY (username, contact_name),
    CONSTRAINT no_self_contact CHECK (username <> contact_name),
    CONSTRAINT first_fk FOREIGN KEY(username) REFERENCES user_service.users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT second_fk FOREIGN KEY(contact_name) REFERENCES user_service.users(username) ON DELETE CASCADE ON UPDATE CASCADE
);
