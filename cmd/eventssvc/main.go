package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shanvl/garbage/internal/eventssvc/aggregating"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/internal/eventssvc/grpc"
	"github.com/shanvl/garbage/internal/eventssvc/postgres"
	"github.com/shanvl/garbage/internal/eventssvc/rest"
	"github.com/shanvl/garbage/internal/eventssvc/schooling"
	"github.com/shanvl/garbage/pkg/env"
)

func main() {
	// create postgres connection pool
	postgresPool, err := postgres.Connect(postgres.Config{
		Database:             env.String("POSTGRES_DB", "garbage1"),
		Host:                 env.String("POSTGRES_HOST", "localhost"),
		User:                 env.String("POSTGRES_USER", "jynweythek223"),
		Password:             env.String("POSTGRES_PASSWORD", "postgres"),
		Port:                 env.Int("POSTGRES_PORT", 5432),
		MaxConns:             env.Int("POSTGRES_MAX_CONN", 25),
		MaxConnLifetime:      env.Duration("POSTGRES_CON_LIFE_SEC", 5*time.Minute),
		PreferSimpleProtocol: env.Bool("POSTGRES_SIMPLE_PROTOCOL", false),
		Logger:               nil,
	})
	if err != nil {
		log.Fatal(err)
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
	// run gRPC server
	go func() {
		if err := grpc.NewServer(aggregatingService, eventingService, schoolingService).Run(grpcPort); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()
	// run REST server
	if err := rest.RunServer(restPort, fmt.Sprintf(":%d", grpcPort)); err != nil && !errors.Is(err,
		http.ErrServerClosed) {

		log.Fatalf("REST endpoint error: %v", err)
	}
}
