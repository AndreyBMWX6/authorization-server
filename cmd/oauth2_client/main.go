package main

import (
	"context"
	"log"
	"net/http"

	"authorization-server/internal/config/secret"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

//todo: add tls

func main() {
	clientSecret, err := getClientSecret()
	if err != nil {
		log.Fatalf("can't create client: %s", err.Error())
	}

	// todo: move values to config
	client := New(&oauth2.Config{
		ClientID: "ff76946a-cb94-45e5-b0f3-1bda06059e10",
		// secret is generated for testing purposes
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:7000/auth",
			TokenURL: "http://localhost:7000/authorization/token",
		},
		RedirectURL: "http://localhost:8000/token",
		Scopes:      []string{"profiles", "roles"},
	})

	r := chi.NewRouter()
	r.Get("/auth", client.authHandler)
	r.Get("/token", client.tokenHandler)
	err = http.ListenAndServe("localhost:8000", r)
	if err != nil {
		log.Fatalf("can't run client: %s", err.Error())
	}
}

func getClientSecret() (string, error) {
	value, err := secret.GetValue(context.Background(), secret.ClientSecretKey)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}
