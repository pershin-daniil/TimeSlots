-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up

CREATE TABLE users_history
(
    id           serial PRIMARY KEY,
    user_id      int         NOT NULL,
    last_name    varchar     NOT NULL,
    first_name   varchar     NOT NULL,
    status       varchar,
    notification time,
    event_time   timestamptz NOT NULL DEFAULT NOW(),
    created_at   timestamptz NOT NULL DEFAULT NOW()
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION users_history()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO users_history (user_id, last_name, first_name, notification, status, event_time, created_at)
    VALUES (NEW.id, NEW.last_name, NEW.first_name, NEW.notification, NEW.status, NOW(), NEW.created_at);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

DROP TRIGGER IF EXISTS users_history_insert ON users;
CREATE TRIGGER users_history_insert
    AFTER INSERT
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE users_history();

DROP TRIGGER IF EXISTS users_history_update ON users;
CREATE TRIGGER users_history_update
    AFTER UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE users_history();

DROP TRIGGER IF EXISTS users_history_delete ON users;
CREATE TRIGGER users_history_delete
    AFTER DELETE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE users_history();

-- +migrate Down

DROP TRIGGER users_history_update ON users;
DROP TRIGGER users_history_delete ON users;
DROP TRIGGER users_history_insert ON users;
DROP FUNCTION users_history();
DROP TABLE users_history;
