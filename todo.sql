DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users_todos;

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
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
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE users_todos(
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    todo_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE,
    UNIQUE (user_id, todo_id)
)