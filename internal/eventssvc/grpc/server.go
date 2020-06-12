package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

	grpcServer := grpc.NewServer()
	eventsv1pb.RegisterEventsServiceServer(grpcServer, s)

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

func (s *Server) FindClasses(ctx context.Context, req *eventsv1pb.FindClassesRequest) (*eventsv1pb.
	FindClassesResponse, error) {

	fmt.Println(req.Letter, req.DateFormed, req.Amount, req.EventFilters, req.EventSorting, req.Sorting, req.Skip)
	fmt.Printf("%+v", req.EventFilters)

	return &eventsv1pb.FindClassesResponse{
		Classes: nil,
	}, nil
}
