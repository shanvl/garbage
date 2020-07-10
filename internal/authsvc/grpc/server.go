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
	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
	healthv1pb "github.com/shanvl/garbage/api/health/v1/pb"
	"github.com/shanvl/garbage/internal/authsvc/authent"
	"github.com/shanvl/garbage/internal/authsvc/authoriz"
	"github.com/shanvl/garbage/internal/authsvc/users"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	log      *zap.Logger
	authent  authent.Service
	authoriz authoriz.Service
	users    users.Service
}

func NewServer(authent authent.Service, authoriz authoriz.Service, users users.Service, log *zap.Logger) *Server {
	server := &Server{
		log:      log,
		authent:  authent,
		authoriz: authoriz,
		users:    users,
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
			// panic recovery
			grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandler(s.handleRecovery)),
		)),
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			// logging
			grpc_zap.StreamServerInterceptor(s.log),
			// panic recovery
			grpcRecovery.StreamServerInterceptor(grpcRecovery.WithRecoveryHandler(s.handleRecovery)),
		)),
	)

	// reflection for tools like Evans
	reflection.Register(grpcServer)

	// register the services
	authv1pb.RegisterAuthServiceServer(grpcServer, s)
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
