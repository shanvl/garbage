package grpc

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/shanvl/garbage/internal/eventsvc/aggregating"
	"github.com/shanvl/garbage/internal/eventsvc/eventing"
	"github.com/shanvl/garbage/internal/eventsvc/postgres"
	"github.com/shanvl/garbage/internal/eventsvc/schooling"
)

var (
	server          *Server
	aggregatingRepo aggregating.Repository
	eventingRepo    eventing.Repository
	schoolingRepo   schooling.Repository
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// testMain initializes gRPC server with an instance of test db
func testMain(m *testing.M) int {
	// connect to the test db. Config values are hardcoded in order not to corrupt production db in case the wrong
	// compose file is used
	db, err := postgres.Connect(postgres.Config{
		Host:            "db",
		Database:        "testdb",
		User:            "root",
		Password:        "root",
		Port:            5432,
		MaxConns:        20,
		MaxConnLifetime: 5 * time.Minute,
	})
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
	authService := newTestAuthService()
	aggregatingService := aggregating.NewService(aggregatingRepo)
	eventingService := eventing.NewService(eventingRepo)
	schoolingService := schooling.NewService(schoolingRepo)
	// create gRPC server
	server = NewServer(authService, aggregatingService, eventingService, schoolingService, nil)
	return m.Run()
}

func newTestAuthService() AuthorizationService {
	return &testAuthService{}
}

type testAuthService struct{}

func (t testAuthService) Authorize(_ context.Context, _, _ string) (string, error) {
	return "", nil
}
