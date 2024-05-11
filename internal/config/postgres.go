package config

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToPostgres(ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", getPostgresDsn())
	if err != nil {
		log.Fatalln(err)
	}

	//todo: ping db

	return db, nil
}

// todo: put dsn to config
func getPostgresDsn() string {
	return "user=postgres dbname=authorization-server sslmode=disable"
}
