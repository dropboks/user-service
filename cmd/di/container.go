package di

import (
	"github.com/dropboks/user-service/config/database"
	"github.com/dropboks/user-service/config/logger"
	"github.com/dropboks/user-service/config/router"
	"github.com/dropboks/user-service/internal/domain/handler"
	"github.com/dropboks/user-service/internal/domain/repository"
	"github.com/dropboks/user-service/internal/domain/service"
	"github.com/dropboks/user-service/internal/infrastructure/grpc"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	if err := container.Provide(logger.New); err != nil {
		panic("Failed to provide logger: " + err.Error())
	}
	if err := container.Provide(database.New); err != nil {
		panic("Failed to provide database: " + err.Error())
	}
	if err := container.Provide(grpc.NewGRPCClientManager); err != nil {
		panic("Failed to provide GRPC Client Manager: " + err.Error())
	}
	if err := container.Provide(grpc.NewFileServiceConnection); err != nil {
		panic("Failed to provide user service grpc connection: " + err.Error())
	}
	if err := container.Provide(repository.NewUserRepository); err != nil {
		panic("Failed to provide authRepository: " + err.Error())
	}
	if err := container.Provide(service.NewAuthService); err != nil {
		panic("Failed to provide auth service: " + err.Error())
	}
	if err := container.Provide(service.NewUserService); err != nil {
		panic("Failed to provide user service: " + err.Error())
	}
	if err := container.Provide(handler.NewUserHandler); err != nil {
		panic("Failed to provide user handler: " + err.Error())
	}
	if err := container.Provide(router.NewHTTP); err != nil {
		panic("Failed to provide HTTP Server: " + err.Error())
	}
	if err := container.Provide(router.NewGRPC); err != nil {
		panic("Failed to provide gRPC Server: " + err.Error())
	}
	return container
}
