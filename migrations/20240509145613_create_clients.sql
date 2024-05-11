-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients (
    id           uuid PRIMARY KEY,
    name         TEXT NOT NULL,
    url          TEXT NOT NULL,
    redirect_uri TEXT NOT NULL,
    secret       TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE clients;
-- +goose StatementEnd
