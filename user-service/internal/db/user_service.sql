-- SCHEMA: user_service

-- DROP SCHEMA IF EXISTS user_service ;

CREATE SCHEMA IF NOT EXISTS user_service
    AUTHORIZATION postgres;

CREATE TABLE IF NOT EXISTS user_service.users
(
    id text COLLATE pg_catalog."default",
	username text COLLATE pg_catalog."default",
	email text COLLATE pg_catalog."default",
	password_hash text COLLATE pg_catalog."default",
	created_at int
);
