package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/shanvl/garbage/internal/eventssvc/aggregating"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/internal/eventssvc/grpc"
	"github.com/shanvl/garbage/internal/eventssvc/postgres"
	"github.com/shanvl/garbage/internal/eventssvc/rest"
	"github.com/shanvl/garbage/internal/eventssvc/schooling"
	"github.com/shanvl/garbage/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// create logger
	loggerConf := zap.NewProductionConfig()
	loggerConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := loggerConf.Build()
	if err != nil {
		fmt.Println("couldn't initialize the logger")
		os.Exit(1)
	}
	defer logger.Sync()

	// should the db log its actions
	var dbLogger pgx.Logger = nil
	if env.Bool("POSTGRES_LOG", false) {
		dbLogger = zapadapter.NewLogger(logger)
	}
	// create postgres connection pool
	postgresConf := postgres.Config{
		Database:             env.String("POSTGRES_DB", "garbage1"),
		Host:                 env.String("POSTGRES_HOST", "localhost"),
		User:                 env.String("POSTGRES_USER", "jynweythek223"),
		Password:             env.String("POSTGRES_PASSWORD", "postgres"),
		Port:                 env.Int("POSTGRES_PORT", 5432),
		MaxConns:             env.Int("POSTGRES_MAX_CONN", 25),
		MaxConnLifetime:      env.Duration("POSTGRES_CON_LIFE", 5*time.Minute),
		PreferSimpleProtocol: env.Bool("POSTGRES_SIMPLE_PROTOCOL", false),
		Logger:               dbLogger,
	}
	postgresPool, err := postgres.Connect(postgresConf)
	if err != nil {
		logger.Fatal("postgres connection error",
			zap.Error(err),
			zap.String("protocol", "postgres"),
			zap.String("addr", fmt.Sprintf("%s:%d", postgresConf.Host, postgresConf.Port)),
		)
	}

	// create repos
	aggregatingRepo := postgres.NewAggregatingRepo(postgresPool)
	eventingRepo := postgres.NewEventingRepo(postgresPool)
	schoolingRepo := postgres.NewSchoolingRepo(postgresPool)

	// create services
	aggregatingService := aggregating.NewService(aggregatingRepo)
	eventingService := eventing.NewService(eventingRepo)
	schoolingService := schooling.NewService(schoolingRepo)

	grpcPort, restPort := env.Int("GRPC_PORT", 3000), env.Int("REST_PORT", 4000)
	// run REST gateway
	go func() {
		if err := rest.NewServer(logger).Run(restPort, fmt.Sprintf(":%d", grpcPort)); err != nil && !errors.Is(err,
			http.ErrServerClosed) {

			logger.Fatal("REST gateway error",
				zap.Error(err),
				zap.Int("port", restPort),
				zap.String("protocol", "HTTP"),
			)
		}
	}()
	// run gRPC server
	if err := grpc.NewServer(aggregatingService, eventingService, schoolingService, logger).Run(grpcPort); err != nil {

		logger.Fatal("gRPC server error",
			zap.Error(err),
			zap.Int("port", grpcPort),
			zap.String("protocol", "gRPC"),
		)
	}
}
