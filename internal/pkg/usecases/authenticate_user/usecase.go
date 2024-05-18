package authenticate_user

import (
	"context"

	"authorization-server/internal/pkg/domain"
	"authorization-server/internal/pkg/jwt"
	"authorization-server/internal/pkg/passwords"
	"github.com/pkg/errors"
)

type UsersRepository interface {
	GetUser(ctx context.Context, login string) (*domain.User, error)
}

type UseCase struct {
	usersRepo UsersRepository
}

func New(usersRepo UsersRepository) *UseCase {
	return &UseCase{
		usersRepo: usersRepo,
	}
}

func (u *UseCase) Authenticate(ctx context.Context, user *domain.User) (string, error) {
	storedUser, err := u.usersRepo.GetUser(ctx, user.Login)
	if err != nil {
		return "", errors.Wrap(err, "get user")
	}

	match, err := passwords.MatchPasswords(storedUser.Password, user.Password)
	if err != nil {
		return "", errors.Wrap(err, "match passwords")
	}
	if !match {
		return "", ErrWrongPassword
	}

	jwtToken, err := jwt.NewWithClaims(ctx, map[string]interface{}{
		"name": user.Login,
	})
	if err != nil {
		return "", errors.Wrap(err, "build jwt")
	}

	return jwtToken, nil
}
