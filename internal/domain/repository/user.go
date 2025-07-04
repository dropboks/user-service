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
)

type (
	UserRepository interface {
		CreateNewUser(*entity.User) error
		QueryUserByEmail(string) (*entity.User, error)
		QueryUserByUserId(string) (*entity.User, error)
		UpdateUser(*entity.User) error
	}
	userRepository struct {
		pgx    *pgxpool.Pool
		logger zerolog.Logger
	}
)

func NewUserRepository(pgx *pgxpool.Pool, logger zerolog.Logger) UserRepository {
	return &userRepository{
		pgx:    pgx,
		logger: logger,
	}
}

func (a *userRepository) UpdateUser(user *entity.User) error {
	query, args, err := sq.Update("users").
		Set("full_name", user.FullName).
		Set("image", user.Image).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("verified", user.Verified).
		Set("two_factor_enabled", user.TwoFactorEnabled).
		Set("updated_at", sq.Expr("CURRENT_TIMESTAMP")).
		Where(sq.Eq{"id": user.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to build update query")
		return dto.Err_INTERNAL_FAILED_BUILD_QUERY
	}

	cmdTag, err := a.pgx.Exec(context.Background(), query, args...)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to update user")
		return dto.Err_INTERNAL_FAILED_UPDATE_USER
	}
	if cmdTag.RowsAffected() == 0 {
		a.logger.Warn().Str("id", user.ID).Msg("user not found for update")
		return dto.Err_NOTFOUND_USER_NOT_FOUND
	}
	return nil
}

func (a *userRepository) CreateNewUser(user *entity.User) error {
	query, args, err := sq.Insert("users").
		Columns("id", "full_name", "image", "email", "password", "verified", "two_factor_enabled").
		Values(user.ID, user.FullName, user.Image, user.Email, user.Password, user.Verified, user.TwoFactorEnabled).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to build insert query")
		return dto.Err_INTERNAL_FAILED_BUILD_QUERY
	}
	row := a.pgx.QueryRow(context.Background(), query, args...)
	if err := row.Scan(&user.ID); err != nil {
		a.logger.Error().Err(err).Msg("failed to insert user")
		return dto.Err_INTERNAL_FAILED_INSERT_USER
	}
	return nil
}

func (a *userRepository) QueryUserByUserId(userId string) (*entity.User, error) {
	var user entity.User
	query, args, err := sq.Select("id", "full_name", "image", "email", "password", "verified", "two_factor_enabled").
		From("users").
		Where(sq.Eq{"id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to build query")
		return nil, dto.Err_INTERNAL_FAILED_BUILD_QUERY
	}
	row := a.pgx.QueryRow(context.Background(), query, args...)
	err = row.Scan(&user.ID, &user.FullName, &user.Image, &user.Email, &user.Password, &user.Verified, &user.TwoFactorEnabled)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.logger.Warn().Str("id", userId).Msg("user not found")
			return nil, dto.Err_NOTFOUND_USER_NOT_FOUND
		}
		a.logger.Error().Err(err).Msg("failed to scan user")
		return nil, dto.Err_INTERNAL_FAILED_SCAN_USER
	}
	return &user, nil

}

func (a *userRepository) QueryUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	query, args, err := sq.Select("id", "full_name", "image", "email", "password", "verified", "two_factor_enabled").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to build query")
		return nil, dto.Err_INTERNAL_FAILED_BUILD_QUERY
	}

	row := a.pgx.QueryRow(context.Background(), query, args...)
	err = row.Scan(&user.ID, &user.FullName, &user.Image, &user.Email, &user.Password, &user.Verified, &user.TwoFactorEnabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.logger.Warn().Str("email", email).Msg("user not found")
			return nil, dto.Err_NOTFOUND_USER_NOT_FOUND
		}
		a.logger.Error().Err(err).Msg("failed to scan user")
		return nil, dto.Err_INTERNAL_FAILED_SCAN_USER
	}
	return &user, nil
}
