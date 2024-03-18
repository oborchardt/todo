DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS todos;

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    token TEXT UNIQUE,
    expiration TEXT
);

CREATE TABLE todos(
    id INTEGER PRIMARY KEY,
    title TEXT,
    text TEXT,
    is_done BOOLEAN,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);