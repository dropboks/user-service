package handler

import (
	"context"

	"github.com/dropboks/proto-user/pkg/upb"
	"github.com/dropboks/user-service/internal/domain/service"
	"google.golang.org/grpc"
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
		return nil, err
	}
	return status, nil
}
func (a *AuthGrpcHandler) GetUserByEmail(c context.Context, email *upb.Email) (*upb.User, error) {
	user, err := a.authService.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
