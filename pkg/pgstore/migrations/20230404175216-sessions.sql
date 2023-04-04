-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up

CREATE TABLE sessions
(
    id         serial PRIMARY KEY,
    user_id    int         NOT NULL REFERENCES users (id) UNIQUE,
    state_name varchar     NOT NULL DEFAULT 'main_screen',
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE sessions_history
(
    id         serial PRIMARY KEY,
    session_id int         NOT NULL,
    user_id    int         NOT NULL,
    state_name varchar     NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION sessions_history()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO sessions_history (session_id, user_id, state_name, updated_at)
    VALUES (NEW.id, NEW.user_id, NEW.state_name, NOW());
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

DROP TRIGGER IF EXISTS sessions_history_insert ON sessions;
CREATE TRIGGER sessions_history_insert
    AFTER INSERT
    ON sessions
    FOR EACH ROW
    EXECUTE PROCEDURE sessions_history();

-- +migrate Down

DROP TRIGGER sessions_history_insert ON sessions;
DROP FUNCTION sessions_history();
DROP TABLE sessions_history;

DROP TABLE sessions;