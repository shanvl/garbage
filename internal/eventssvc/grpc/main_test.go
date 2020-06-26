package grpc

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/shanvl/garbage/internal/eventssvc/aggregating"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/internal/eventssvc/postgres"
	"github.com/shanvl/garbage/internal/eventssvc/schooling"
	"github.com/shanvl/garbage/pkg/env"
)

var (
	server          *Server
	aggregatingRepo *postgres.AggregatingRepo
	eventingRepo    *postgres.EventingRepo
	schoolingRepo   *postgres.SchoolingRepo
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// testMain initializes gRPC server with an instance of test db
func testMain(m *testing.M) int {
	// connect to the test db. Config values are hardcoded in order not to corrupt production db in case the wrong
	// compose file is used
	db, err := postgres.Connect(postgres.Config{
		Database:             env.String("POSTGRES_DB", "garbage1"),
		Host:                 env.String("POSTGRES_HOST", "localhost"),
		User:                 env.String("POSTGRES_USER", "jynweythek223"),
		Password:             env.String("POSTGRES_PASSWORD", "postgres"),
		Port:                 env.Int("POSTGRES_PORT", 5432),
		MaxConns:             env.Int("POSTGRES_MAX_CONN", 25),
		MaxConnLifetime:      env.Duration("POSTGRES_CON_LIFE", 5*time.Minute),
		PreferSimpleProtocol: env.Bool("POSTGRES_SIMPLE_PROTOCOL", false)})
	if err != nil {
		log.Printf("couldn't connect to testdb: %s\n", err)
		return 1
	}
	defer db.Close()

	// repos with a test db
	aggregatingRepo = postgres.NewAggregatingRepo(db)
	eventingRepo = postgres.NewEventingRepo(db)
	schoolingRepo = postgres.NewSchoolingRepo(db)
	// create services
	aggregatingService := aggregating.NewService(aggregatingRepo)
	eventingService := eventing.NewService(eventingRepo)
	schoolingService := schooling.NewService(schoolingRepo)
	// create gRPC server
	server = NewServer(aggregatingService, eventingService, schoolingService, nil)
	return m.Run()
}
