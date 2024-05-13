package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type States struct {
	// simple unique state used for prevention of CSRF-attacks
	state string
	// PKCE code, used as more robust state
	pkce string
}

type OauthClient struct {
	cfg *oauth2.Config
	// verifiers stores verifying states for different sessions
	verifiers map[string]States
}

func New(cfg *oauth2.Config) *OauthClient {
	return &OauthClient{
		cfg:       cfg,
		verifiers: make(map[string]States),
	}
}

func (c *OauthClient) authHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, sessionID)
	if err != nil {
		http.Error(w, "failed to generate session id", http.StatusInternalServerError)
		return
	}
	encodedSessionID := base64.StdEncoding.EncodeToString(sessionID)
	http.SetCookie(w, &http.Cookie{Name: "sessionID", Value: encodedSessionID})

	// using PKCE or state to protect against CSRF attacks
	pkce := oauth2.GenerateVerifier()
	state := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, state)
	if err != nil {
		http.Error(w, "failed to generate session id", http.StatusInternalServerError)
		return
	}
	encodedState := base64.StdEncoding.EncodeToString(state)

	c.verifiers[string(sessionID)] = States{
		state: encodedState,
		pkce:  pkce,
	}

	authUrl := c.cfg.AuthCodeURL(encodedState, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(pkce))
	resp, err := http.Get(authUrl)
	if err != nil {
		log.Fatalf("get auth code: %s", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	http.SetCookie(w, &http.Cookie{Name: "params", Value: resp.Request.URL.RawQuery})

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	_, _ = w.Write(body)
}

func (c *OauthClient) tokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	encodedSessionCookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "failed to get session cookie", http.StatusBadRequest)
		return
	}
	if encodedSessionCookie.Value == "" {
		http.Error(w, "empty session cookie", http.StatusBadRequest)
		return
	}

	encodedSessionID := encodedSessionCookie.Value
	sessionID, err := base64.StdEncoding.DecodeString(encodedSessionID)
	if err != nil {
		http.Error(w, "decode sessionID cookie", http.StatusInternalServerError)
		return
	}
	verifier, ok := c.verifiers[string(sessionID)]
	if !ok {
		http.Error(w, fmt.Sprintf("no verifier found for session: %s", sessionID), http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state != verifier.state {
		http.Error(w, "wrong state", http.StatusBadRequest)
	}

	authorizationCode := r.URL.Query().Get("code")
	if authorizationCode == "" {
		http.Error(w, "no authorization code", http.StatusInternalServerError)
		return
	}

	token, err := c.cfg.Exchange(ctx, authorizationCode, oauth2.VerifierOption(verifier.pkce))
	if err != nil {
		log.Fatalf("get access token: %s", err.Error())
	}

	client := c.cfg.Client(ctx, token)
	_ = client
	// todo: make requests
}
