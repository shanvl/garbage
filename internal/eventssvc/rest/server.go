package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"google.golang.org/grpc"
)

func RunServer(port int, grpcAddress string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	customErrorsOption := runtime.WithErrorHandler(customHTTPError)
	mux := runtime.NewServeMux(customErrorsOption)
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	err := eventsv1pb.RegisterEventsServiceHandlerFromEndpoint(ctx, mux, grpcAddress, dialOptions)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// graceful shutdown on signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("REST gateway shutdown error: %v\n", err)
		}
		log.Println("REST gateway has been shut down")
	}()

	log.Printf("Starting REST gateway on %d port\n", port)
	return server.ListenAndServe()
}
