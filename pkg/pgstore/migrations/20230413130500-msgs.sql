-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up

CREATE TABLE messages
(
    id         serial PRIMARY KEY,
    user_id    int         NOT NULL REFERENCES users (id),
    msg_id     int         NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE messages_history
(
    id          serial PRIMARY KEY,
    messaged_id int         NOT NULL,
    user_id     int         NOT NULL,
    msg_id      int         NOT NULL,
    updated_at  timestamptz NOT NULL DEFAULT NOW(),
    created_at  timestamptz NOT NULL DEFAULT NOW()
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION messages_history()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO messages_history (messaged_id, user_id, msg_id, updated_at, created_at)
    VALUES (NEW.id, NEW.user_id, NEW.msg_id, NOW(), NEW.created_at);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

DROP TRIGGER IF EXISTS messages_history_insert ON messages;
CREATE TRIGGER messages_history_insert
    AFTER INSERT
    ON messages
    FOR EACH ROW
EXECUTE PROCEDURE messages_history();

DROP TRIGGER IF EXISTS messages_history_delete ON messages;
CREATE TRIGGER messages_history_delete
    AFTER DELETE
    ON messages
    FOR EACH ROW
EXECUTE PROCEDURE messages_history();

-- +migrate Down

DROP TRIGGER messages_history_insert ON messages;
DROP TRIGGER messages_history_delete ON messages
DROP FUNCTION messages_history();
DROP TABLE messages_history;

DROP TABLE messages;
