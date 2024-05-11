package register_user

import (
	"context"

	"authorization-server/internal/pkg/domain"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository interface {
	InsertUser(ctx context.Context, user *domain.User) error
}

type UseCase struct {
	usersRepo UsersRepository
}

func New(usersRepo UsersRepository) *UseCase {
	return &UseCase{
		usersRepo: usersRepo,
	}
}

func (u *UseCase) Register(ctx context.Context, user *domain.User) error {
	// salt is user inside bcrypt lib
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "hash password")
	}
	user.Password = string(hashedPassword)

	err = u.usersRepo.InsertUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "insert user")
	}

	return nil
}
