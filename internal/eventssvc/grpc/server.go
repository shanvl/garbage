package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/ptypes/timestamp"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	healthv1pb "github.com/shanvl/garbage/api/health/v1/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
}

func NewServer() *Server {
	server := &Server{}
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

	eventsv1pb.RegisterEventsServiceServer(grpcServer, s)
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

// Check is used for health checks
func (s *Server) Check(ctx context.Context, request *healthv1pb.HealthCheckRequest) (*healthv1pb.HealthCheckResponse, error) {
	return &healthv1pb.HealthCheckResponse{Status: healthv1pb.HealthCheckResponse_SERVING}, nil
}

// Watch is used for stream health checks. Not implemented but required by gRPC Health Checking Protocol
func (s *Server) Watch(request *healthv1pb.HealthCheckRequest, server healthv1pb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "no stream health check")
}

func (s *Server) FindClasses(ctx context.Context, req *eventsv1pb.FindClassesRequest) (*eventsv1pb.
	FindClassesResponse, error) {

	return &eventsv1pb.FindClassesResponse{
		Classes: nil,
	}, nil
}

func (s *Server) FindEvents(ctx context.Context, request *eventsv1pb.FindEventsRequest) (*eventsv1pb.
	FindEventsResponse, error) {

	panic("implement me")
}

func (s *Server) FindPupils(ctx context.Context, request *eventsv1pb.FindPupilsRequest) (*eventsv1pb.
	FindPupilsResponse, error) {

	panic("implement me")
}

func (s *Server) FindPupil(ctx context.Context, request *eventsv1pb.FindPupilRequest) (*eventsv1pb.FindPupilResponse,
	error) {

	return &eventsv1pb.FindPupilResponse{
		Pupil: &eventsv1pb.PupilAggr{
			Id:          "",
			FirstName:   "",
			LastName:    "",
			ClassLetter: "",
			ClassDateFormed: &timestamp.Timestamp{
				Seconds: 0,
				Nanos:   0,
			},
			ResourcesBrought: &eventsv1pb.ResourcesBrought{
				Gadgets: 0,
				Paper:   0,
				Plastic: 0,
			},
			Events: nil,
		},
	}, nil
}
