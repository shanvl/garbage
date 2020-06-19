package grpc

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
)

func (s *Server) ChangePupilResources(ctx context.Context, request *eventsv1pb.ChangePupilResourcesRequest) (*empty.
	Empty, error) {

	fmt.Println(request.GetEventId(), request.GetPupilId(), request.GetResources())
	return &empty.Empty{}, nil
}

func (s *Server) CreateEvent(ctx context.Context, request *eventsv1pb.CreateEventRequest) (*eventsv1pb.
	CreateEventResponse, error) {

	fmt.Println(request.GetName(), request.GetDate(), request.GetResourcesAllowed())
	return &eventsv1pb.CreateEventResponse{Id: "some id"}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, request *eventsv1pb.DeleteEventRequest) (*empty.Empty, error) {
	fmt.Println(request.GetId())

	return &empty.Empty{}, nil
}

func (s *Server) FindEventByID(ctx context.Context, request *eventsv1pb.FindEventByIDRequest) (*eventsv1pb.FindEventByIDResponse, error) {
	panic("implement me")
}

func (s *Server) FindEventClasses(ctx context.Context, request *eventsv1pb.FindEventClassesRequest) (*eventsv1pb.FindEventClassesResponse, error) {
	panic("implement me")
}

func (s *Server) FindEventPupils(ctx context.Context, request *eventsv1pb.FindEventPupilsRequest) (*eventsv1pb.FindEventPupilsResponse, error) {
	panic("implement me")
}

func (s *Server) FindEventPupilByID(ctx context.Context, request *eventsv1pb.FindEventPupilByIDRequest) (*eventsv1pb.FindEventPupilByIDResponse, error) {
	panic("implement me")
}
