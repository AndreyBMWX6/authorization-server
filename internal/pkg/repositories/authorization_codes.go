package repositories

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	"oauth2/internal/generated/authorization-server/public/model"
	"oauth2/internal/generated/authorization-server/public/table"
	"oauth2/internal/pkg/domain"
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

func (r *AuthorizationCodesRepository) GetCode(ctx context.Context, code string) (*domain.AuthorizationCode, error) {
	stmt := table.AuthorizationCodes.
		SELECT(table.AuthorizationCodes.AllColumns).
		WHERE(table.AuthorizationCodes.Code.EQ(postgres.String(code)))

	var m model.AuthorizationCodes
	err := stmt.QueryContext(ctx, r.db, &m)
	if err != nil {
		return nil, err
	}

	return mapAuthorizationCodeToDomain(m), nil
}

func (r *AuthorizationCodesRepository) MatchCodeUsed(ctx context.Context, code string) error {
	stmt := table.AuthorizationCodes.
		UPDATE(table.AuthorizationCodes.Used).
		SET(true).
		WHERE(table.AuthorizationCodes.Code.EQ(postgres.String(code)))

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

func mapAuthorizationCodeToDomain(code model.AuthorizationCodes) *domain.AuthorizationCode {
	return &domain.AuthorizationCode{
		Code:           code.Code,
		ClientID:       code.ClientID,
		RedirectURI:    code.RedirectURI,
		ExpirationTime: code.ExpirationTime,
		Used:           code.Used,
	}
}
