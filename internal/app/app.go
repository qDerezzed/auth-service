package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"auth-service/config"
	v1 "auth-service/internal/delivery/http/v1"
	"auth-service/internal/usecase"
	"auth-service/internal/usecase/repository"
	"auth-service/pkg/httpserver"
	"auth-service/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	// Repository
	pg, err := postgres.New(*cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer pg.Close()

	// Use case
	authUseCase := usecase.New(
		repository.New(pg),
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, authUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HttpPort))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
