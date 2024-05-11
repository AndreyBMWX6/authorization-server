package main

import (
	"context"
	"log"
	"net/http"
	"os"

	authenticationServer "authorization-server/internal/app/authentication_server"
	authorizationServer "authorization-server/internal/app/authorization_server"
	"authorization-server/internal/app/scratch"
	"authorization-server/internal/config"
)

func main() {
	ctx := context.Background()

	a, err := scratch.New()
	if err != nil {
		log.Fatalf("can't create app: %s", err.Error())
	}
	db, err := config.ConnectToPostgres(ctx)
	if err != nil {
		log.Fatalf("can't connect to db")
	}

	fileServer := createFileServer()
	authorizationServ := authorizationServer.NewAuthorizationServer(db, fileServer)
	authenticationServ := authenticationServer.NewAuthenticationServer(db, fileServer)

	// custom http without grpc
	a.PublicServer().Get("/auth/*", authorizationServ.GetAllowAccessPage)
	a.PublicServer().Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/auth/", http.StatusMovedPermanently)
	})
	a.PublicServer().Get("/login/*", authenticationServ.GetLoginPage)
	a.PublicServer().Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login/", http.StatusMovedPermanently)
	})

	if err := a.Run(context.Background(), authorizationServ, authenticationServ); err != nil {
		log.Fatalf("can't run app: %s", err.Error())
	}
}

func createFileServer() http.Handler {
	// relative paths doesn't work so path have to be absolute
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get working directory")
	}
	return http.FileServer(http.Dir(dir + "/assets"))
}