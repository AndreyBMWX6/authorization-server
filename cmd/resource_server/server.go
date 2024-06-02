package main

import (
	"net/http"
)

type ResourceServer struct{}

func NewResourceServer() *ResourceServer {
	return &ResourceServer{}
}

func (c *ResourceServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "no auth", http.StatusUnauthorized)
	}

	//todo: go to authorization server and ask if access token exists and not expired

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
	_, _ = w.Write([]byte(auth))
}
