package handler

import (
	"context"

	upb "github.com/dropboks/proto-user/pkg/upb"
	"github.com/dropboks/user-service/internal/domain/dto"
	"github.com/dropboks/user-service/internal/domain/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_status "google.golang.org/grpc/status"
)

type AuthGrpcHandler struct {
	authService service.AuthService
	upb.UnimplementedUserServiceServer
}

func newAuthGrpcHandler(authService service.AuthService) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		authService: authService,
	}
}

func RegisterAuthService(grpc *grpc.Server, authService service.AuthService) {
	grpcHandler := newAuthGrpcHandler(authService)
	upb.RegisterUserServiceServer(grpc, grpcHandler)
}

func (a *AuthGrpcHandler) CreateUser(c context.Context, user *upb.User) (*upb.Status, error) {
	status, err := a.authService.CreateUser(user)
	if err != nil {
		if err == dto.Err_INTERNAL_FAILED_BUILD_QUERY || err == dto.Err_INTERNAL_FAILED_INSERT_USER {
			return nil, _status.Error(codes.Internal, err.Error())
		}
		return nil, _status.Error(codes.Internal, err.Error())
	}
	return status, nil
}
func (a *AuthGrpcHandler) GetUserByEmail(c context.Context, email *upb.Email) (*upb.User, error) {
	user, err := a.authService.GetUserByEmail(email)
	if err != nil {
		if err == dto.Err_NOTFOUND_USER_NOT_FOUND {
			return nil, _status.Error(codes.NotFound, err.Error())
		} else if err == dto.Err_INTERNAL_FAILED_BUILD_QUERY || err == dto.Err_INTERNAL_FAILED_INSERT_USER {
			return nil, _status.Error(codes.Internal, err.Error())
		}
		return nil, _status.Error(codes.Internal, err.Error())
	}
	return user, nil
}

func (a *AuthGrpcHandler) UpdateUser(c context.Context, user *upb.User) (*upb.Status, error) {
	if err := a.authService.UpdateUser(c, user); err != nil {
		if err == dto.Err_NOTFOUND_USER_NOT_FOUND {
			return nil, _status.Error(codes.NotFound, err.Error())
		} else if err == dto.Err_INTERNAL_FAILED_BUILD_QUERY || err == dto.Err_INTERNAL_FAILED_INSERT_USER {
			return nil, _status.Error(codes.Internal, err.Error())
		}
		return nil, _status.Error(codes.Internal, err.Error())
	}
	return &upb.Status{Success: true}, nil
}

func (a *AuthGrpcHandler) GetUserByUserId(c context.Context, userId *upb.UserId) (*upb.User, error) {
	user, err := a.authService.GetUserByUserId(userId)
	if err != nil {
		if err == dto.Err_NOTFOUND_USER_NOT_FOUND {
			return nil, _status.Error(codes.NotFound, err.Error())
		} else if err == dto.Err_INTERNAL_FAILED_BUILD_QUERY || err == dto.Err_INTERNAL_FAILED_INSERT_USER {
			return nil, _status.Error(codes.Internal, err.Error())
		}
		return nil, _status.Error(codes.Internal, err.Error())
	}
	return user, nil
}
