package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/timestamp"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
)

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

func (s *Server) FindPupilByID(ctx context.Context, request *eventsv1pb.FindPupilByIDRequest) (*eventsv1pb.
	FindPupilByIDResponse, error) {

	return &eventsv1pb.FindPupilByIDResponse{
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
