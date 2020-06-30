package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	log *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{logger}
}

func (s *Server) Run(port int, grpcAddress string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// add custom error handler
	customErrorsOption := runtime.WithErrorHandler(customHTTPError)

	// create new mux
	mux := runtime.NewServeMux(customErrorsOption)

	// no tls
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}

	// transform and proxy all REST requests to gRPC server
	err := eventsv1pb.RegisterEventsServiceHandlerFromEndpoint(ctx, mux, grpcAddress, dialOptions)
	if err != nil {
		return err
	}

	// create REST gateway
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.logMiddleware(mux),
	}

	// graceful shutdown on signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			s.log.Error("REST gateway shutdown error",
				zap.Int("port", port),
				zap.String("protocol", "HTTP"),
				zap.Error(err),
			)
		}
		s.log.Info("REST gateway has been shut down",
			zap.Int("port", port),
			zap.String("protocol", "HTTP"),
		)
	}()

	s.log.Info("starting REST gateway",
		zap.Int("port", port),
		zap.String("protocol", "HTTP"),
	)
	return server.ListenAndServe()
}
