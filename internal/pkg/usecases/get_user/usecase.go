package get_user

import (
	"context"

	"github.com/pkg/errors"
	"oauth2/internal/pkg/domain"
)

type UserRepository interface {
	GetUser(ctx context.Context, login string) (*domain.User, error)
}

type UseCase struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *UseCase {
	return &UseCase{
		userRepo: userRepo,
	}
}

func (u *UseCase) GetUser(ctx context.Context, login string) (*domain.User, error) {
	user, err := u.userRepo.GetUser(ctx, login)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}
	return user, nil
}
