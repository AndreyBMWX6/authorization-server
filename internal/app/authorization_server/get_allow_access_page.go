package authorization_server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"authorization-server/internal/config/secret"
	"authorization-server/internal/pkg/claims"
	"authorization-server/internal/pkg/utils"
	desc "authorization-server/pkg/api/authorization_server"
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

	params, err := getQueryParams(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("get params: %s", err.Error()), http.StatusBadRequest)
		return
	}

	protoReq, err := buildGetAuthorizationCodeRequestFromHttpParams(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("parse cookie paramsCookie: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if err = validateGetAuthorizationCodeRequest(ctx, protoReq); err != nil {
		http.Error(w, fmt.Errorf("validate request: %s", err.Error()).Error(), http.StatusBadRequest)
		return
	}

	jwtCookie, err := r.Cookie(authCookie)
	if err != nil {
		//todo: move login endpoint to config
		loginURL := fmt.Sprintf("http://localhost:7000/login?%s", r.URL.RawQuery)
		http.Redirect(w, r, loginURL, http.StatusMovedPermanently)
		return
	}
	token := jwtCookie.Value
	if token == "" {
		//todo: move login endpoint to config
		loginURL := fmt.Sprintf("http://localhost:7000/login?%s", r.URL.RawQuery)
		http.Redirect(w, r, loginURL, http.StatusMovedPermanently)
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

	// setting url parameters in cookie to store request parameters, because static files have no access to ctx of request
	http.SetCookie(w, &http.Cookie{Name: "paramsCookie", Value: r.URL.RawQuery})

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

// params can be passed in URL or in params cookie
func getQueryParams(r *http.Request) (string, error) {
	if r.URL.RawQuery != "" {
		return r.URL.RawQuery, nil
	}

	paramsCookie, err := r.Cookie("params")
	if err != nil {
		return "", errors.Wrap(err, "no params stored in cookie")
	}
	if paramsCookie.Value == "" {
		return "", errors.New("empty params cookie")
	}

	return paramsCookie.Value, nil
}

func buildGetAuthorizationCodeRequestFromHttpParams(query string) (*desc.GetAuthorizationCodeRequest, error) {
	params, err := url.ParseQuery(query)
	if err != nil {
		return nil, errors.Wrap(err, "parse params")
	}

	return &desc.GetAuthorizationCodeRequest{
		ResponseType: desc.ResponseType(desc.ResponseType_value[params.Get("response_type")]),
		ClientId:     params.Get("client_id"),
		RedirectUri:  utils.ToPtrIfNotEmpty(params.Get("redirect_uri")),
		Scope:        utils.ToPtrIfNotEmpty(params.Get("scope")),
		State:        utils.ToPtrIfNotEmpty(params.Get("state")),
	}, nil
}
