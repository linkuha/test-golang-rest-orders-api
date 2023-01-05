package app

import (
	"context"
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
	"time"
)

func Run(cfg *config.Config) {
	logPath := cfg.Merged.LogDir + "/main.log"
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Failure to open log file " + logPath)
	}
	defer logFile.Close()
	logger.InitLogger(cfg.Merged.LogLevel, logFile)

	log.Info().Msg("Starting application...")
	log.Debug().Msgf("Config dump: %#v\n", cfg)

	ctx := context.Background()

	// Repository
	db, err := newDB(&cfg.EnvParams)
	if err != nil {
		log.Fatal().Msgf("Can't connect to database: %s", err.Error())
	}
	defer db.Close()

	repos := repository.NewRepository(db)

	// HTTP Server
	ctrl := v1.NewController(ctx, repos)
	router := ctrl.ConfigureRoutes(cfg)
	httpSrv := httpserver.New(router, httpserver.Port(cfg.EnvParams.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		log.Info().Msg("app - termination signal: " + s.String())
	case err = <-httpSrv.Notify():
		log.Error().Msgf("app - httpServer error or stopped (ErrServerClosed): %s", err.Error())
	}

	log.Info().Msg("Graceful shutdown...")
	ctxTeardown, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = httpSrv.Shutdown(ctxTeardown)
	if err != nil {
		log.Error().Err(fmt.Errorf("app - teardown - httpServer.Shutdown: %w", err))
	}
	// unplug from message broker
	// unplug from service mesh
	// remove temporary files
	// wait for all pending queue/topic processor to finish
}
