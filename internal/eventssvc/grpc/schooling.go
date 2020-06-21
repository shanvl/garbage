package grpc

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventssvc"
	"github.com/shanvl/garbage/internal/eventssvc/schooling"
	"github.com/shanvl/garbage/pkg/valid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddPupils adds the given pupils and returns the ids of the added
func (s *Server) AddPupils(ctx context.Context, req *eventsv1pb.AddPupilsRequest) (*eventsv1pb.AddPupilsResponse,
	error) {

	reqPupils := req.GetPupils()
	pupilsBio := make([]schooling.PupilBio, len(reqPupils))
	for i, pupil := range reqPupils {
		pupilsBio[i] = schooling.PupilBio{
			FirstName: pupil.FirstName,
			LastName:  pupil.LastName,
			ClassName: pupil.Class,
		}
	}

	pupilIDS, err := s.scSvc.AddPupils(ctx, pupilsBio)
	if err != nil {
		var validationError *valid.ErrValidation
		if errors.As(err, &validationError) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &eventsv1pb.AddPupilsResponse{PupilIds: pupilIDS}, nil
}

// ChangePupilClass changes the class of the pupil
func (s *Server) ChangePupilClass(ctx context.Context, req *eventsv1pb.ChangePupilClassRequest) (*empty.Empty, error) {

	err := s.scSvc.ChangePupilClass(ctx, req.GetId(), req.GetClass())
	if err != nil {
		var validationError *valid.ErrValidation
		if errors.As(err, &validationError) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, eventssvc.ErrUnknownPupil) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

// RemovePupils removes the pupils with the given IDs
func (s *Server) RemovePupils(ctx context.Context, req *eventsv1pb.RemovePupilsRequest) (*empty.Empty, error) {

	err := s.scSvc.RemovePupils(ctx, req.GetPupilIds())
	if err != nil {
		var validationError *valid.ErrValidation
		if errors.As(err, &validationError) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}
