-- +migrate Up
CREATE TABLE message (
    id INTEGER NOT NULL PRIMARY KEY,
    body TEXT NOT NULL DEFAULT "",
    username TEXT NOT NULL DEFAULT "",
    created TIMESTAMP NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    updated TIMESTAMP NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

-- +migrate Down
DROP TABLE message;
