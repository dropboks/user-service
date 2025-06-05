package router

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func loggingUnaryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Info().
			Str("method", info.FullMethod).
			Msg("gRPC request received")
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error().
				Str("method", info.FullMethod).
				Err(err).
				Msg("gRPC request error")
		}
		return resp, err
	}
}

func NewGRPC(logger zerolog.Logger) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingUnaryInterceptor(logger)),
	)
	return grpcServer
}
