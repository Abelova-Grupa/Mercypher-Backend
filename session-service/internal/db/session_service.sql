-- SCHEMA: session_service

-- DROP SCHEMA IF EXISTS session_service ;

CREATE SCHEMA IF NOT EXISTS session_service
    AUTHORIZATION postgres;

CREATE TABLE IF NOT EXISTS session_service.last_seen_sessions(
	user_id text,
	last_seen int
);

CREATE TABLE IF NOT EXISTS session_service.sessions(
	id text,
	user_id text,
	refresh_token text,
	access_token text
);

CREATE TABLE IF NOT EXISTS session_service.user_locations(
	user_id text,
	api_ip text
);