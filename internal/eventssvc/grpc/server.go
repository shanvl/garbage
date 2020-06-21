package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	healthv1pb "github.com/shanvl/garbage/api/health/v1/pb"
	"github.com/shanvl/garbage/internal/eventssvc/aggregating"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/internal/eventssvc/schooling"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	agSvc aggregating.Service
	evSvc eventing.Service
	scSvc schooling.Service
}

func NewServer(agSvc aggregating.Service, evSvc eventing.Service, scSvc schooling.Service) *Server {
	server := &Server{
		agSvc: agSvc,
		evSvc: evSvc,
		scSvc: scSvc,
	}
	return server
}

func (s *Server) Run(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			NewPanicInterceptor().Unary(),
		)),
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			NewPanicInterceptor().Stream(),
		)),
	)

	reflection.Register(grpcServer)

	eventsv1pb.RegisterAggregatingServiceServer(grpcServer, s)
	eventsv1pb.RegisterEventingServiceServer(grpcServer, s)
	eventsv1pb.RegisterSchoolingServiceServer(grpcServer, s)
	healthv1pb.RegisterHealthServer(grpcServer, s)

	// graceful shutdown on signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		grpcServer.GracefulStop()
		log.Println("stopping gRPC server")
	}()

	log.Printf("starting gRPC server on %d port\n", port)
	return grpcServer.Serve(listener)
}

func handleContextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:
		return nil
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
