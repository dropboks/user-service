package service

import (
	"context"

	"github.com/dropboks/proto-user/pkg/upb"
	"github.com/dropboks/sharedlib/utils"
	"github.com/dropboks/user-service/internal/domain/entity"
	"github.com/dropboks/user-service/internal/domain/repository"
	"github.com/rs/zerolog"
)

type (
	AuthService interface {
		CreateUser(*upb.User) (*upb.Status, error)
		GetUserByEmail(*upb.Email) (*upb.User, error)
		UpdateUser(c context.Context, user *upb.User) error
		GetUserByUserId(user *upb.UserId) (*upb.User, error)
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

// UpdateUser implements AuthService.
func (a *authService) UpdateUser(c context.Context, user *upb.User) error {
	u := &entity.User{
		ID:               user.GetId(),
		FullName:         user.GetFullName(),
		Image:            utils.StringPtr(user.GetImage()),
		Email:            user.GetEmail(),
		Password:         user.GetPassword(),
		Verified:         user.GetVerified(),
		TwoFactorEnabled: user.GetTwoFactorEnabled(),
	}
	if err := a.userRepository.UpdateUser(u); err != nil {
		return err
	}
	return nil
}

func (a *authService) GetUserByEmail(email *upb.Email) (*upb.User, error) {
	user, err := a.userRepository.QueryUserByEmail(email.GetEmail())
	if err != nil {
		return nil, err
	}
	return &upb.User{
		Id:               user.ID,
		FullName:         user.FullName,
		Image:            user.Image,
		Email:            user.Email,
		Password:         user.Password,
		Verified:         user.Verified,
		TwoFactorEnabled: user.TwoFactorEnabled,
	}, nil
}

func (a *authService) GetUserByUserId(user *upb.UserId) (*upb.User, error) {
	userFetched, err := a.userRepository.QueryUserByUserId(user.GetUserId())
	if err != nil {
		return nil, err
	}
	return &upb.User{
		Id:               userFetched.ID,
		FullName:         userFetched.FullName,
		Image:            userFetched.Image,
		Email:            userFetched.Email,
		Password:         userFetched.Password,
		Verified:         userFetched.Verified,
		TwoFactorEnabled: userFetched.TwoFactorEnabled,
	}, nil
}

func (a *authService) CreateUser(user *upb.User) (*upb.Status, error) {
	u := &entity.User{
		ID:       user.GetId(),
		FullName: user.GetFullName(),
		Image:    utils.StringPtr(user.GetImage()),
		Email:    user.GetEmail(),
		Password: user.GetPassword(),
	}
	err := a.userRepository.CreateNewUser(u)
	if err != nil {
		return nil, err
	}
	return &upb.Status{Success: true}, nil
}
