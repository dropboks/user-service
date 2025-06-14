package service

import (
	"github.com/dropboks/proto-user/pkg/upb"
	"github.com/dropboks/user-service/internal/domain/entity"
	"github.com/dropboks/user-service/internal/domain/repository"
	"github.com/rs/zerolog"
)

type (
	AuthService interface {
		CreateUser(*upb.User) (*upb.Status, error)
		GetUserByEmail(*upb.Email) (*upb.User, error)
	}
	authService struct {
		userRepository repository.UserRepository
		logger         zerolog.Logger
	}
)

func NewAuthService(userRepository repository.UserRepository, logger zerolog.Logger) AuthService {
	return &authService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (a *authService) GetUserByEmail(email *upb.Email) (*upb.User, error) {
	user, err := a.userRepository.QueryUserByEmail(email.GetEmail())
	if err != nil {
		return nil, err
	}
	return &upb.User{
		Id:       user.ID,
		FullName: user.FullName,
		Image:    user.Image,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (a *authService) CreateUser(user *upb.User) (*upb.Status, error) {
	u := &entity.User{
		ID:       user.GetId(),
		FullName: user.GetFullName(),
		Image:    user.GetImage(),
		Email:    user.GetEmail(),
		Password: user.GetPassword(),
	}
	err := a.userRepository.CreateNewUser(u)
	if err != nil {
		return nil, err
	}
	return &upb.Status{Success: true}, nil
}
