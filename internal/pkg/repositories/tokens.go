package repositories

import (
	"context"
	"time"

	"authorization-server/internal/generated/authorization-server/public/model"
	"authorization-server/internal/generated/authorization-server/public/table"
	"authorization-server/internal/pkg/domain"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
)

type TokensRepository struct {
	db *sqlx.DB
}

func NewTokensRepository(db *sqlx.DB) *TokensRepository {
	return &TokensRepository{
		db: db,
	}
}

func (r *TokensRepository) UpsertToken(ctx context.Context, token domain.Token) error {
	m := mapTokenToModel(token)

	stmt := table.Tokens.
		INSERT(table.Tokens.AllColumns).
		MODEL(m).
		ON_CONFLICT(table.Tokens.AccessToken).
		DO_UPDATE(
			postgres.SET(
				table.Tokens.RefreshToken.SET(table.Tokens.EXCLUDED.RefreshToken),
			),
		)

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}

	return nil
}

func (r *TokensRepository) GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	stmt := table.Tokens.
		SELECT(table.Tokens.AllColumns).
		WHERE(table.Tokens.RefreshToken.EQ(postgres.String(refreshToken)))

	var m model.Tokens
	err := stmt.QueryContext(ctx, r.db, &m)
	if err != nil {
		return nil, err
	}

	return mapTokenToDomain(m), nil
}

func (r *TokensRepository) DeleteTokensByCode(ctx context.Context, authorizationCode string) error {
	stmt := table.Tokens.
		DELETE().
		WHERE(table.Tokens.AuthorizationCode.EQ(postgres.String(authorizationCode)))

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}

	return nil
}

func mapTokenToModel(token domain.Token) model.Tokens {
	return model.Tokens{
		AccessToken:       token.AccessToken,
		AuthorizationCode: token.AuthorizationCode,
		Type:              string(token.Type),
		CreatedAt:         token.CreatedAt,
		ExpiresIn:         int64(token.ExpiresIn.Seconds()),
		RefreshToken:      token.RefreshToken,
		Scope:             token.Scope,
	}
}

func mapTokenToDomain(token model.Tokens) *domain.Token {
	return &domain.Token{
		AccessToken:       token.AccessToken,
		AuthorizationCode: token.AuthorizationCode,
		Type:              domain.TokenType(token.Type),
		CreatedAt:         token.CreatedAt,
		ExpiresIn:         time.Second * time.Duration(token.ExpiresIn),
		RefreshToken:      token.RefreshToken,
		Scope:             token.Scope,
	}
}
