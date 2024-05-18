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
		ClientID: "53ab9b19-0d52-49a2-8e34-becf2908781a",
		// secret is generated for testing purposes
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "http://localhost:7000/auth",
			TokenURL:  "http://localhost:7000/authorization/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: "http://localhost:8000/token",
		Scopes:      []string{"profiles", "roles"},
	})

	r := chi.NewRouter()
	r.Get("/auth", client.authHandler)
	r.Get("/token", client.tokenHandler)
	r.Get("/client", client.clientHandler)
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
