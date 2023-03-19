DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

CREATE TABLE users
(
	login varchar(64) PRIMARY KEY NOT NULL,
	email varchar(256) NOT NULL,
	password_hash varchar(256) NOT NULL,
	phone_number varchar(16) NOT NULL
);

CREATE TABLE sessions
(
	session_id varchar(256) PRIMARY KEY NOT NULL,
	login varchar(64) not null UNIQUE REFERENCES users(login),
	create_date timestamp NOT NULL,
	expire_date timestamp NOT NULL,
	last_access_date timestamp NOT NULL
);
