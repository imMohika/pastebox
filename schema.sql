CREATE TABLE IF NOT EXISTS snippets
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    title   VARCHAR(100) NOT NULL,
    content TEXT         NOT NULL,
    expires DATETIME,
    created DATETIME     NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions
(
    token  TEXT PRIMARY KEY,
    data   BLOB NOT NULL,
    expiry REAL NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
