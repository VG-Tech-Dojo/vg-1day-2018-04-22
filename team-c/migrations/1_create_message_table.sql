-- +migrate Up
CREATE TABLE message (
    id INTEGER NOT NULL PRIMARY KEY,
    body TEXT NOT NULL DEFAULT "",
    username TEXT NOT NULL DEFAULT "",
    image_id INTEGER NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    updated TIMESTAMP NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);
CREATE TABLE image (
    id INTEGER NOT NULL PRIMARY KEY,
    path TEXT NOT NULL
);
-- +migrate Down
DROP TABLE message;

