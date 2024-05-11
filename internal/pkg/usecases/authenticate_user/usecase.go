package authenticate_user

import (
	"context"
	"time"

	"authorization-server/internal/config/secret"
	"authorization-server/internal/pkg/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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

	match, err := matchPasswords(storedUser.Password, user.Password)
	if err != nil {
		return "", errors.Wrap(err, "match passwords")
	}
	if !match {
		return "", ErrWrongPassword
	}

	jwtToken, err := buildJwt(ctx, user.Login)
	if err != nil {
		return "", errors.Wrap(err, "build jwt")
	}

	return jwtToken, nil
}

func matchPasswords(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func buildJwt(ctx context.Context, login string) (string, error) {
	//todo: move expiration time to config
	//todo: return 5 minutes
	expirationTime := time.Now().In(time.UTC).Add(time.Minute * 30)

	// todo: move setting claims to jwt package
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  expirationTime,
		"name": login,
	})

	jwtSecretVal, err := secret.GetValue(ctx, secret.JWTSecretKey)
	if err != nil {
		return "", errors.Wrap(err, "get jwt secret")
	}
	jwtSecret := jwtSecretVal.(string)

	return token.SignedString([]byte(jwtSecret))
}
