package get_authorization_code

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"

	"authorization-server/internal/pkg/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ClientsRepository interface {
	GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error)
}

type AuthorizationCodesRepository interface {
	InsertCode(ctx context.Context, code domain.AuthorizationCode) error
}

type UseCase struct {
	clientsRepo ClientsRepository
	codesRepo   AuthorizationCodesRepository
}

func New(clientRepo ClientsRepository, codesRepo AuthorizationCodesRepository) *UseCase {
	return &UseCase{
		clientsRepo: clientRepo,
		codesRepo:   codesRepo,
	}
}

func (u *UseCase) GetCode(ctx context.Context, client *domain.Client, scope *string) (string, error) {
	// check if client exists in db
	dbClient, err := u.clientsRepo.GetClient(ctx, client.ID)
	if err != nil {
		return "", errors.Wrap(err, "get client")
	}

	// if no redirect URI in request get it from db
	if client.RedirectURI == "" {
		client.RedirectURI = dbClient.RedirectURI
	}

	//todo: use scope

	// generate code
	authCode := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, authCode)
	if err != nil {
		return "", errors.Wrap(err, "generate authorization code")
	}
	encodedCode := base64.StdEncoding.EncodeToString(authCode)

	// save code to db
	code := domain.AuthorizationCode{
		Code:        encodedCode,
		ClientID:    client.ID,
		RedirectURI: client.RedirectURI,
		// todo: move expiration time to config
		ExpirationTime: time.Now().In(time.UTC).Add(time.Minute * 5),
		Used:           false,
	}

	err = u.codesRepo.InsertCode(ctx, code)
	if err != nil {
		return "", errors.Wrap(err, "insert code")
	}

	return encodedCode, nil
}
