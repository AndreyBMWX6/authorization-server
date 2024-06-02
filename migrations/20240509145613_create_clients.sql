-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients (
    id              uuid PRIMARY KEY,
    name            text NOT NULL,
    url             text NOT NULL,
    redirect_uri    text NOT NULL,
    secret          text NOT NULL,
    is_confidential bool NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE clients;
-- +goose StatementEnd
