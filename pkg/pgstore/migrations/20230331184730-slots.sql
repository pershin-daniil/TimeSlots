-- noinspection SqlNoDataSourceInspectionForFile

-- +migrate Up

CREATE TABLE slots
(
    id         serial PRIMARY KEY,
    coach_id   int         NOT NULL REFERENCES users (id),
    client_id  int         NOT NULL REFERENCES users (id),
    event_id   varchar     NOT NULL,
    status     varchar     NOT NULL DEFAULT 'TENTATIVE',
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE slots_history
(
    id         serial PRIMARY KEY,
    slot_id    int         NOT NULL,
    coach_id   int         NOT NULL,
    client_id  int         NOT NULL,
    event_id   varchar     NOT NULL,
    status     varchar     NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION slots_history()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO slots_history (slot_id, coach_id, client_id, event_id, status, updated_at, created_at)
    VALUES (NEW.id, NEW.coach_id, NEW.client_id, NEW.event_id, NEW.status, NOW(), NEW.created_at);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

DROP TRIGGER IF EXISTS slots_history_insert ON slots;
CREATE TRIGGER slots_history_insert
    AFTER INSERT
    ON slots
    FOR EACH ROW
EXECUTE PROCEDURE slots_history();

-- +migrate Down

DROP TRIGGER slots_history_insert ON slots;
DROP FUNCTION slots_history();
DROP TABLE slots_history;

DROP TABLE slots;
