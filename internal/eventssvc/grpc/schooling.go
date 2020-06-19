package grpc

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
)

func (s *Server) AddPupils(ctx context.Context, request *eventsv1pb.AddPupilsRequest) (*eventsv1pb.AddPupilsResponse, error) {
	panic("implement me")
}

func (s *Server) ChangePupilClass(ctx context.Context, request *eventsv1pb.ChangePupilClassRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Server) RemovePupils(ctx context.Context, request *eventsv1pb.RemovePupilsRequest) (*empty.Empty, error) {
	fmt.Println(request.GetPupilIds())

	return &empty.Empty{}, nil
}
