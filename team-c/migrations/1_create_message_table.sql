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

INSERT INTO image(path) VALUES
    ('image01.png'),
    ('image02.png'),
    ('image03.png'),
    ('image04.png'),
    ('image05.png'),
    ('image06.png'),
    ('image07.png'),
    ('image08.png'),
    ('image09.png'),
    ('image10.png');

-- +migrate Down
DROP TABLE message;

