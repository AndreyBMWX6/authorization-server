-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
     login        text PRIMARY KEY,
     password     text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
