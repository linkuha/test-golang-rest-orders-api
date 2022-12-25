package app

import (
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/config"
	v1 "github.com/linkuha/test-golang-rest-orders-api/internal/delivery/httpserver/v1"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository"
	"github.com/linkuha/test-golang-rest-orders-api/pkg/logger"
	"github.com/linkuha/test-golang-rest-orders-api/pkg/srv/httpserver"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	logPath := cfg.Merged.LogDir + "/main.log"
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Failure to open log file " + logPath)
	}
	defer func() {
		_ = logFile.Close()
	}()
	logger.InitLogger(cfg.Merged.LogLevel, logFile)

	log.Info().Msg("Starting application...")
	log.Debug().Msgf("Config dump: %#v\n", cfg)

	// Repository
	db, err := newDB(&cfg.EnvParams)
	if err != nil {
		log.Fatal().Msgf("Can't connect to database: %s", err.Error())
	}

	repos := repository.NewRepository(db)

	// HTTP Server
	ctrl := v1.NewController(repos)
	router := ctrl.ConfigureRoutes(cfg)
	httpSrv := httpserver.New(router, httpserver.Port(cfg.EnvParams.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		log.Info().Msg("app - Run - signal: " + s.String())
	case err := <-httpSrv.Notify():
		log.Error().Msgf("app - Run - httpServer.Notify: %s", err.Error())
	}

	log.Info().Msg("Graceful shutdown...")
	err = httpSrv.Shutdown()
	if err != nil {
		log.Error().Err(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
