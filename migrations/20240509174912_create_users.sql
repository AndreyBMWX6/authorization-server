-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
     login        TEXT PRIMARY KEY,
     password     TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
