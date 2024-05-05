package main

import (
	"context"
	"log"

	authorizationServer "authorization-server/internal/app/authorization_server"
	"authorization-server/internal/app/scratch"
	"authorization-server/internal/pkg/storage"
)

func main() {
	a, err := scratch.New()
	if err != nil {
		log.Fatalf("can't create app: %s", err.Error())
	}
	inmem := storage.New()
	authorizationServ := authorizationServer.NewAuthorizationServer(inmem)
	if err := a.Run(context.Background(), authorizationServ); err != nil {
		log.Fatalf("can't run app: %s", err.Error())
	}
}
