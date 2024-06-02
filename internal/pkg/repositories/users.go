package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"oauth2/internal/generated/authorization-server/public/model"
	"oauth2/internal/generated/authorization-server/public/table"
	"oauth2/internal/pkg/domain"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (r *UsersRepository) InsertUser(ctx context.Context, user *domain.User) error {
	m := mapUserToModel(user)

	stmt := table.Users.
		INSERT(table.Users.AllColumns).
		MODEL(m)

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) GetUser(ctx context.Context, login string) (*domain.User, error) {
	stmt := table.Users.
		SELECT(table.Users.AllColumns).
		WHERE(table.Users.Login.EQ(postgres.String(login)))

	var user model.Users
	err := stmt.QueryContext(ctx, r.db, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return mapUserToDomain(user), nil
}

func mapUserToModel(user *domain.User) model.Users {
	return model.Users{
		Login:    user.Login,
		Password: user.Password,
	}
}

func mapUserToDomain(user model.Users) *domain.User {
	return &domain.User{
		Login:    user.Login,
		Password: user.Password,
	}
}
