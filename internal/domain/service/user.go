package service

import (
	"github.com/dropboks/user-service/internal/domain/dto"
	"github.com/dropboks/user-service/internal/domain/repository"
	"github.com/rs/zerolog"
)

type (
	UserService interface {
		GetProfile(userId string) (dto.GetProfileResponse, error)
	}
	userService struct {
		userRepository repository.UserRepository
		logger         zerolog.Logger
	}
)

func NewUserService(userRepo repository.UserRepository, logger zerolog.Logger) UserService {
	return &userService{
		userRepository: userRepo,
		logger:         logger,
	}
}

func (u *userService) GetProfile(userId string) (dto.GetProfileResponse, error) {
	user, err := u.userRepository.QueryUserByUserId(userId)
	if err != nil {
		return dto.GetProfileResponse{}, err
	}
	profile := dto.GetProfileResponse{
		FullName:         user.FullName,
		Image:            user.Image,
		Email:            user.Email,
		Verified:         user.Verified,
		TwoFactorEnabled: user.TwoFactorEnabled,
	}
	return profile, nil
}
