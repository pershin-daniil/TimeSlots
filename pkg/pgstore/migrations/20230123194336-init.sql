-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up

CREATE TABLE users
(
    id           int         NOT NULL UNIQUE,
    last_name    varchar     NOT NULL,
    first_name   varchar     NOT NULL,
    status       varchar     NOT NULL DEFAULT 'guest',
    notification time                 DEFAULT '01:00:00',
    updated_at   timestamptz NOT NULL DEFAULT NOW(),
    created_at   timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX users_id_idx ON users (id);


-- +migrate Down

DROP INDEX users_id_idx;
DROP TABLE users;