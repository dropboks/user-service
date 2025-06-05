package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dropboks/user-service/internal/domain/handler"
	"github.com/dropboks/user-service/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Container   *dig.Container
	ServerReady chan bool
	Address     string
}

func (s *GRPCServer) Run() {
	err := s.Container.Invoke(func(
		httpServer *gin.Engine,
		grpcServer *grpc.Server,
		logger zerolog.Logger,
		db *pgxpool.Pool,
		svc service.AuthService,
	) {
		defer db.Close()
		listen, err := net.Listen("tcp", s.Address)
		if err != nil {
			logger.Fatal().Msgf("failed to listen:%v", err)
		}
		handler.RegisterAuthService(grpcServer, svc)

		go func() {
			if serveErr := grpcServer.Serve(listen); serveErr != nil {
				logger.Fatal().Msgf("gRPC server error: %v", serveErr)
			}
		}()

		if s.ServerReady != nil {
			s.ServerReady <- true
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info().Msg("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		logger.Info().Msg("gRPC server stopped gracefully.")
	})
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
}
