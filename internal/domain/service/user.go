package service

import (
	"context"

	"github.com/dropboks/proto-file/pkg/fpb"
	"github.com/dropboks/sharedlib/utils"
	"github.com/dropboks/user-service/internal/domain/dto"
	"github.com/dropboks/user-service/internal/domain/repository"
	"github.com/dropboks/user-service/pkg/constant"
	"github.com/rs/zerolog"
)

type (
	UserService interface {
		GetProfile(userId string) (dto.GetProfileResponse, error)
		UpdateUser(req *dto.UpdateUserRequest, userId string) error
	}
	userService struct {
		userRepository    repository.UserRepository
		logger            zerolog.Logger
		fileServiceClient fpb.FileServiceClient
	}
)

func NewUserService(userRepo repository.UserRepository, logger zerolog.Logger, fileServiceClient fpb.FileServiceClient) UserService {
	return &userService{
		userRepository:    userRepo,
		logger:            logger,
		fileServiceClient: fileServiceClient,
	}
}

func (u *userService) UpdateUser(req *dto.UpdateUserRequest, userId string) error {
	user, err := u.userRepository.QueryUserByUserId(userId)
	if err != nil {
		return err
	}
	us := *user
	if req.FullName != user.FullName {
		us.FullName = req.FullName
	}
	if req.TwoFactorEnabled != user.TwoFactorEnabled {
		us.TwoFactorEnabled = req.TwoFactorEnabled
	}
	ctx := context.Background()
	if req.Image != nil && req.Image.Filename != "" {
		ext := utils.GetFileNameExtension(req.Image.Filename)
		if ext != "jpg" && ext != "jpeg" && ext != "png" {
			return dto.Err_BAD_REQUEST_WRONG_EXTENTION
		}
		if req.Image.Size > constant.MAX_UPLOAD_SIZE {
			return dto.Err_BAD_REQUEST_LIMIT_SIZE_EXCEEDED
		}
		image, err := utils.FileToByte(req.Image)
		if err != nil {
			u.logger.Error().Err(err).Msg("error converting image")
			return dto.Err_INTERNAL_CONVERT_IMAGE
		}
		imageReq := &fpb.Image{
			Image: image,
			Ext:   ext,
		}
		resp, err := u.fileServiceClient.SaveProfileImage(ctx, imageReq)
		if err != nil {
			u.logger.Error().Err(err).Msg("Error uploading image to file service")
			return err
		}
		us.Image = utils.StringPtr(resp.GetName())
	}
	err = u.userRepository.UpdateUser(&us)
	if err != nil && req.Image != nil && req.Image.Filename != "" {
		_, err := u.fileServiceClient.RemoveProfileImage(ctx, &fpb.ImageName{Name: *us.Image})
		return err
	}
	return nil
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
