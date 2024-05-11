package get_user

import (
	"context"

	"authorization-server/internal/pkg/domain"
	"github.com/pkg/errors"
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
