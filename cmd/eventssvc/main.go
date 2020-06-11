package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/shanvl/garbage/internal/eventssvc/grpc"
	"github.com/shanvl/garbage/internal/eventssvc/rest"
)

const (
	httpPort = 5555
	grpcPort = 6666
)

func main() {
	ctx := context.Background()
	go func() {
		if err := rest.RunServer(ctx, httpPort, fmt.Sprintf(":%d", grpcPort)); err != nil && !errors.Is(err,
			http.ErrServerClosed) {

			log.Fatalf("REST endpoint error: %v", err)
		}
	}()

	if err := grpc.NewServer().Run(grpcPort); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
