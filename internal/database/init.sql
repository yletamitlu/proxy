CREATE TABLE requests
(
    id        serial  NOT NULL PRIMARY KEY,
    method    text NOT NULL,
    protocol  text NOT NULL,
    host      text NOT NULL,
    path      text NOT NULL,
    headers   jsonb   NOT NULL,
    body      text NOT NULL
);
