package grpc_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authent"
	"github.com/shanvl/garbage/internal/authsvc/authoriz"
	"github.com/shanvl/garbage/internal/authsvc/grpc"
	"github.com/shanvl/garbage/internal/authsvc/jwt"
	"github.com/shanvl/garbage/internal/authsvc/postgres"
	"github.com/shanvl/garbage/internal/authsvc/users"
	"go.uber.org/zap"
)

var (
	server       *grpc.Server
	usersRepo    users.Repository
	authentRepo  authent.Repository
	tokenManager authsvc.TokenManager
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// testMain initializes gRPC server with an instance of test db
func testMain(m *testing.M) int {
	// connect to the test db. Config values are hardcoded in order not to corrupt production db in case the wrong
	// compose file is used
	db, err := postgres.Connect(postgres.Config{
		Host:            "localhost",
		Database:        "authsvc",
		User:            "jynweythek223",
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
	authentRepo = postgres.NewAuthentRepo(db)
	usersRepo = postgres.NewUsersRepo(db)
	// token manager with test keys
	prKey, err := jwt.PrivateKeyFromFile("../jwt/keys_test/test.rsa")
	if err != nil {
		log.Print(err)
		return 1
	}
	pubKey, err := jwt.PublicKeyFromFile("../jwt/keys_test/test.rsa.pub")
	if err != nil {
		log.Print(err)
		return 1
	}
	tokenManager = jwt.NewManagerRSA(30*time.Minute, 120*time.Hour, prKey, pubKey)
	// create services
	authentSvc := authent.NewService(authentRepo, tokenManager)
	authorizSvc := authoriz.NewService(tokenManager, authoriz.ProtectedRPCMap())
	usersSvc := users.NewService(usersRepo)
	// logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Print(err)
		return 1
	}
	// create gRPC server
	server = grpc.NewServer(authentSvc, authorizSvc, usersSvc, logger)
	return m.Run()
}
