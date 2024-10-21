CREATE TABLE IF NOT EXISTS snippets
(
    id      INTEGER      PRIMARY KEY AUTOINCREMENT,
    title   VARCHAR(100) NOT NULL,
    content TEXT         NOT NULL,
    expires DATETIME,
    created DATETIME     NOT NULL
);
