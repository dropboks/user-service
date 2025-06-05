package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/dropboks/user-service/internal/domain/dto"
	"github.com/dropboks/user-service/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	AuthRepository interface {
		QueryUserByEmail(string) (entity.User, error)
		CreateNewUser(*entity.User) error
	}
	authRepository struct {
		pgx    *pgxpool.Pool
		logger zerolog.Logger
	}
)

func NewAuthRepository(pgx *pgxpool.Pool, logger zerolog.Logger) AuthRepository {
	return &authRepository{
		pgx:    pgx,
		logger: logger,
	}
}

func (a *authRepository) CreateNewUser(user *entity.User) error {
	query, args, err := sq.Insert("users").
		Columns("id", "full_name", "image", "email", "password").
		Values(user.ID, user.FullName, user.Image, user.Email, user.Password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to build insert query")
		return status.Error(codes.Internal, dto.Err_INTERNAL_FAILED_BUILD_QUERY.Error())
	}
	row := a.pgx.QueryRow(context.Background(), query, args...)
	if err := row.Scan(&user.ID); err != nil {
		a.logger.Error().Err(err).Msg("failed to insert user")
		return status.Error(codes.Internal, dto.Err_INTERNAL_FAILED_INSERT_USER.Error())
	}
	return nil
}

func (a *authRepository) QueryUserByEmail(email string) (entity.User, error) {
	var user entity.User
	query, args, err := sq.Select("id", "full_name", "image", "email", "password").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to build query")
		return user, status.Error(codes.Internal, dto.Err_INTERNAL_FAILED_BUILD_QUERY.Error())
	}

	row := a.pgx.QueryRow(context.Background(), query, args...)
	err = row.Scan(&user.ID, &user.FullName, &user.Image, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.logger.Warn().Str("email", email).Msg("user not found")
			return user, status.Error(codes.NotFound, dto.Err_NOTFOUND_USER_NOT_FOUND.Error())
		}
		a.logger.Error().Err(err).Msg("failed to scan user")
		return user, status.Error(codes.Internal, dto.Err_INTERNAL_FAILED_SCAN_USER.Error())
	}
	return user, nil
}
