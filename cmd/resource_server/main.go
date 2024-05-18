package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	resourceServer := NewResourceServer()

	r := chi.NewRouter()
	r.Get("/ping", resourceServer.pingHandler)
	err := http.ListenAndServe("localhost:9000", r)
	if err != nil {
		log.Fatalf("can't run client: %s", err.Error())
	}
}
