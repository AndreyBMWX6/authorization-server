-- +goose Up
-- +goose StatementBegin
CREATE TABLE tokens (
    access_token       text PRIMARY KEY,
    authorization_code text        NOT NULL,
    type               text        NOT NULL,
    created_at         timestamptz NOT NULL,
    expires_in         bigint      NOT NULL,
    refresh_token      text,
    scope              text
);

--TODO: CREATE INDEX on authorization_code and refresh_token

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens;
-- +goose StatementEnd
