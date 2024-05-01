package main

import (
	"context"
	"log"

	"authorization-server/internal/app"
	authorizationServer "authorization-server/internal/app/authorization_server"
)

func main() {
	a := app.NewApp()
	authorizationServ := authorizationServer.NewAuthorizationServer()
	if err := a.Run(context.Background(), authorizationServ); err != nil {
		log.Fatalf("can't run app: %s", err.Error())
	}
}
