package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dropboks/user-service/internal/domain/handler"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
)

type HTTPServer struct {
	Container   *dig.Container
	ServerReady chan bool
	Address     string
}

func (s *HTTPServer) Run(ctx context.Context) {
	err := s.Container.Invoke(
		func(
			logger zerolog.Logger,
			router *gin.Engine,
			uh handler.UserHandler,
		) {
			handler.RegisterUserRoutes(router, uh)
			srv := &http.Server{
				Addr:    s.Address,
				Handler: router,
			}
			logger.Info().Msgf("HTTP Server Starting in port %s", s.Address)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal().Err(err).Msg("Failed to listen and serve http server")
				}
			}()

			if s.ServerReady != nil {
				s.ServerReady <- true
			}

			<-ctx.Done()
			logger.Info().Msg("Shutting down server...")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := srv.Shutdown(shutdownCtx); err != nil {
				logger.Fatal().Err(err).Msg("Server forced to shutdown")
			}
			logger.Info().Msg("Server exiting...")
		})
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
}
