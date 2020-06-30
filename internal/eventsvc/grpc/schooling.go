package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventsvc/schooling"
)

// AddPupils adds the given pupils and returns the ids of the added
func (s *Server) AddPupils(ctx context.Context, req *eventsv1pb.AddPupilsRequest) (*eventsv1pb.AddPupilsResponse,
	error) {

	// parse the request
	reqPupils := req.GetPupils()
	pupilsBio := make([]schooling.PupilBio, len(reqPupils))
	for i, pupil := range reqPupils {
		pupilsBio[i] = schooling.PupilBio{
			FirstName: pupil.FirstName,
			LastName:  pupil.LastName,
			ClassName: pupil.Class,
		}
	}

	// call the service
	pupilIDS, err := s.scSvc.AddPupils(ctx, pupilsBio)
	if err != nil {
		return nil, s.handleError(err)
	}
	return &eventsv1pb.AddPupilsResponse{PupilIds: pupilIDS}, nil
}

// ChangePupilClass changes the class of the pupil
func (s *Server) ChangePupilClass(ctx context.Context, req *eventsv1pb.ChangePupilClassRequest) (*empty.Empty, error) {

	err := s.scSvc.ChangePupilClass(ctx, req.GetPupilId(), req.GetClass())
	if err != nil {
		return nil, s.handleError(err)
	}
	return &empty.Empty{}, nil
}

// RemovePupils removes the pupils with the given IDs
func (s *Server) RemovePupils(ctx context.Context, req *eventsv1pb.RemovePupilsRequest) (*empty.Empty, error) {

	err := s.scSvc.RemovePupils(ctx, req.GetPupilIds())
	if err != nil {
		return nil, s.handleError(err)
	}

	return &empty.Empty{}, nil
}
