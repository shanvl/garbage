package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	healthv1pb "github.com/shanvl/garbage/api/health/v1/pb"
	"github.com/shanvl/garbage/internal/eventssvc/aggregating"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/internal/eventssvc/schooling"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	agSvc aggregating.Service
	evSvc eventing.Service
	scSvc schooling.Service
	log   *zap.Logger
}

func NewServer(agSvc aggregating.Service, evSvc eventing.Service, scSvc schooling.Service, log *zap.Logger) *Server {
	server := &Server{
		agSvc: agSvc,
		evSvc: evSvc,
		scSvc: scSvc,
		log:   log,
	}
	return server
}

// Run configures and starts gRPC server
func (s *Server) Run(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// add interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			// logging
			grpc_zap.UnaryServerInterceptor(s.log),
			// payload logging
			grpc_zap.PayloadUnaryServerInterceptor(s.log, s.payloadDecider),
			// panic recovery
			grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandler(s.handleRecovery)),
		)),
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			// logging
			grpc_zap.StreamServerInterceptor(s.log),
			// payload logging
			grpc_zap.PayloadStreamServerInterceptor(s.log, s.payloadDecider),
			// panic recovery
			grpcRecovery.StreamServerInterceptor(grpcRecovery.WithRecoveryHandler(s.handleRecovery)),
		)),
	)

	// reflection for tools like Evans
	reflection.Register(grpcServer)

	// register the services
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
		s.log.Info("stopping gRPC server",
			zap.Int("port", port),
			zap.String("protocol", "gRPC"),
		)
	}()

	s.log.Info("starting gRPC server",
		zap.Int("port", port),
		zap.String("protocol", "gRPC"),
	)
	return grpcServer.Serve(listener)
}
