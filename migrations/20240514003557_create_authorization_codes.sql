-- +goose Up
-- +goose StatementBegin
CREATE TABLE authorization_codes
(
    code            text PRIMARY KEY,
    client_id       uuid               NOT NULL,
    redirect_uri    text               NOT NULL,
    expiration_time timestamptz        NOT NULL,
    used            bool DEFAULT false NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authorization_codes;
-- +goose StatementEnd
