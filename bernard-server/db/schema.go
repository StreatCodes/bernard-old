package db

const Schema = `
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT,
	hashed_password BLOB,
	roles INTEGER
);
CREATE UNIQUE INDEX users_email ON users(email);

CREATE TABLE sessions (
	user_id INTEGER REFERENCES users(id),
	expiry TEXT, --TODO
	key BLOB
);
CREATE UNIQUE INDEX sessions_user_id ON sessions(user_id);
CREATE UNIQUE INDEX sessions_key ON sessions(key);

CREATE TABLE hosts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	host_name TEXT,
	description TEXT,
	key BLOB
);

CREATE TABLE checks (
	id TEXT,
	host_id INTEGER REFERENCES hosts(id),
	name TEXT,
	command TEXT,
	args TEXT, --TODO
	env TEXT, --TODO
	dir TEXT,
	interval INTEGER,
	timeout INTEGER
);
CREATE UNIQUE INDEX checks_host_id ON checks(host_id);
CREATE UNIQUE INDEX checksid ON checks(id);

CREATE TABLE check_results (
	check_id TEXT REFERENCES checks(id),
	time INTEGER,
	status INTEGER,
	output BLOB
);
CREATE UNIQUE INDEX check_results_check_id ON check_results(check_id);
CREATE UNIQUE INDEX check_results_time ON check_results(time);
`
