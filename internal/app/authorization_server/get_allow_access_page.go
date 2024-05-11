package authorization_server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"authorization-server/internal/config/secret"
	"authorization-server/internal/pkg/claims"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	authCookie     = "jwt"
	authBearerType = "Bearer"
)

var (
	errWrongTokenFormat    = errors.New("wrong token format")
	errUnsupportedAuthType = errors.New("unsupported authorization type")
)

func (i *Implementation) GetAllowAccessPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jwtCookie, err := r.Cookie(authCookie)
	if err != nil {
		//todo: redirect to /login endpoint
		http.Error(w, "no token provided", http.StatusUnauthorized)
		return
	}
	token := jwtCookie.Value
	if token == "" {
		//todo: redirect to /login endpoint
		http.Error(w, "no token provided", http.StatusUnauthorized)
		return
	}

	jwtToken, err := parseJwtToken(ctx, token)
	if err != nil {
		http.Error(w, "parse jwt token", http.StatusBadRequest)
		return
	}

	c, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cast claims to map claims")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	expiredTime, err := claims.GetExpiredTime(c)
	if err != nil {
		err = errors.Wrap(err, "get expired time from claims")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !time.Now().In(time.UTC).Before(expiredTime) {
		http.Error(w, "jwt token expired", http.StatusBadRequest)
		return
	}

	login, err := claims.GetLogin(c)
	if err != nil {
		err = errors.Wrap(err, "get login from claims")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if user exists
	_, err = i.getUserUseCase.GetUser(ctx, login)
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	i.fileServer.ServeHTTP(w, r)
}

func parseJwtToken(ctx context.Context, token string) (*jwt.Token, error) {
	splitted := strings.Split(token, " ")
	if len(splitted) != 2 {
		return nil, errWrongTokenFormat
	}

	if splitted[0] != authBearerType {
		return nil, errUnsupportedAuthType
	}

	tokenString := splitted[1]
	jwtSecretVal, err := secret.GetValue(ctx, secret.JWTSecretKey)
	if err != nil {
		return nil, errors.Wrap(err, "get jwt secret")
	}
	jwtSecret := jwtSecretVal.(string)

	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "parse jwt token")
	}

	return jwtToken, nil
}
