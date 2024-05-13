package repositories

import (
	"context"

	"authorization-server/internal/generated/authorization-server/public/model"
	"authorization-server/internal/generated/authorization-server/public/table"
	"authorization-server/internal/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type AuthorizationCodesRepository struct {
	db *sqlx.DB
}

func NewAuthorizationCodesRepository(db *sqlx.DB) *AuthorizationCodesRepository {
	return &AuthorizationCodesRepository{
		db: db,
	}
}

func (r *AuthorizationCodesRepository) InsertCode(ctx context.Context, code domain.AuthorizationCode) error {
	m := mapAuthorizationCodeToModel(code)

	stmt := table.AuthorizationCodes.
		INSERT(table.AuthorizationCodes.AllColumns).
		MODEL(m)

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}

	return nil
}

func mapAuthorizationCodeToModel(code domain.AuthorizationCode) model.AuthorizationCodes {
	return model.AuthorizationCodes{
		Code:           code.Code,
		ClientID:       code.ClientID,
		RedirectURI:    code.RedirectURI,
		ExpirationTime: code.ExpirationTime,
		Used:           code.Used,
	}
}
