package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/shanvl/garbage/internal/authsvc/authent"
	"github.com/shanvl/garbage/internal/authsvc/authoriz"
	"github.com/shanvl/garbage/internal/authsvc/grpc"
	"github.com/shanvl/garbage/internal/authsvc/jwt"
	"github.com/shanvl/garbage/internal/authsvc/postgres"
	"github.com/shanvl/garbage/internal/authsvc/rest"
	"github.com/shanvl/garbage/internal/authsvc/users"
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
		Database:             env.String("POSTGRES_DB", "authsvc"),
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
	// apply migrations
	err = postgres.ValidateSchema(context.Background(), postgresPool)
	if err != nil {
		logger.Fatal("migrations failed", zap.Error(err), zap.String("protocol", "postgres"))
	}

	// create repos
	authentRepo := postgres.NewAuthentRepo(postgresPool)
	usersRepo := postgres.NewUsersRepo(postgresPool)

	// get private and public keys for the token manager
	privateKeyPath := env.String("RSA_PRIVATE_KEY_PATH", "./internal/authsvc/jwt/keys/test.rsa")
	publicKeyPath := env.String("RSA_PUBLIC_KEY_PATH", "./internal/authsvc/jwt/keys_test/test.rsa.pub")
	privateKey, err := jwt.PrivateKeyFromFile(privateKeyPath)
	if err != nil {
		logger.Fatal("couldn't load private key for the token manager", zap.Error(err))
	}
	publicKey, err := jwt.PublicKeyFromFile(publicKeyPath)
	if err != nil {
		logger.Fatal("couldn't load public key for the token manager", zap.Error(err))
	}

	// create services
	tokenManager := jwt.NewManagerRSA(env.Duration("ACCESS_TOKEN_DURATION", 30*time.Minute),
		env.Duration("REFRESH_TOKEN_DURATION", 720*time.Hour), privateKey, publicKey)
	authentSvc := authent.NewService(authentRepo, tokenManager)
	authorizSvc := authoriz.NewService(tokenManager, authoriz.ProtectedRPCMap())
	usersSvc := users.NewService(usersRepo)

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
	if err := grpc.NewServer(authentSvc, authorizSvc, usersSvc, logger).Run(grpcPort); err != nil {
		logger.Fatal("gRPC server error",
			zap.Error(err),
			zap.Int("port", grpcPort),
			zap.String("protocol", "gRPC"),
		)
	}
}
