package grpc

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventsvc/eventing"
)

// ChangePupilResources changes the amount of resources brought by the pupil to the event
func (s *Server) ChangePupilResources(ctx context.Context, req *eventsv1pb.ChangePupilResourcesRequest) (*empty.
	Empty, error) {

	err := s.evSvc.ChangePupilResources(ctx, req.GetEventId(), req.GetPupilId(),
		protoToResourcesMap(req.GetResourcesBrought()))
	if err != nil {
		return nil, s.handleError(err)
	}

	return &empty.Empty{}, nil
}

// CreateEvent creates and stores the event
func (s *Server) CreateEvent(ctx context.Context, req *eventsv1pb.CreateEventRequest) (*eventsv1pb.
	CreateEventResponse, error) {

	// proto to args
	eventDate, err := protoTimeToTimestamp(req.GetDate())
	if err != nil {
		return nil, s.handleError(fmt.Errorf("event date: %w", err))
	}
	resourcesAllowed, err := protoToResources(req.GetResourcesAllowed())
	if err != nil {
		return nil, s.handleError(err)
	}

	// call the svc
	eventID, err := s.evSvc.CreateEvent(ctx, eventDate, req.GetName(), resourcesAllowed)
	if err != nil {
		return nil, s.handleError(err)
	}

	return &eventsv1pb.CreateEventResponse{Id: eventID}, nil
}

// DeleteEvent deletes the event
func (s *Server) DeleteEvent(ctx context.Context, req *eventsv1pb.DeleteEventRequest) (*empty.Empty, error) {
	err := s.evSvc.DeleteEvent(ctx, req.GetId())
	if err != nil {
		return nil, s.handleError(err)
	}

	return &empty.Empty{}, nil
}

// FindEventByID returns an event with the given id and all resources collected at that event
func (s *Server) FindEventByID(ctx context.Context, req *eventsv1pb.FindEventByIDRequest) (*eventsv1pb.FindEventByIDResponse, error) {

	event, err := s.evSvc.EventByID(ctx, req.GetId())
	if err != nil {
		return nil, s.handleError(err)
	}

	e, err := eventToProto(event)
	if err != nil {
		return nil, s.handleError(err)
	}
	return &eventsv1pb.FindEventByIDResponse{Event: e}, nil
}

// FindEventClasses returns an array of sorted classes with the resources they brought to the specified event
func (s *Server) FindEventClasses(ctx context.Context, req *eventsv1pb.FindEventClassesRequest) (*eventsv1pb.
	FindEventClassesResponse, error) {

	// call the svc
	classes, total, err := s.evSvc.EventClasses(
		ctx,
		req.GetEventId(),
		eventing.EventClassFilters{Name: req.GetClassName()},
		protoClassSortingMap[req.GetSorting()],
		int(req.GetAmount()),
		int(req.GetSkip()),
	)
	if err != nil {
		return nil, s.handleError(err)
	}

	// to proto
	pbClasses := make([]*eventsv1pb.Class, len(classes))
	for i, class := range classes {
		pbClasses[i] = classToProto(class)
	}
	return &eventsv1pb.FindEventClassesResponse{
		Classes: pbClasses,
		Total:   uint32(total),
	}, err
}

// FindEventPupils returns an array of sorted pupils with the resources they brought to the specified event
func (s *Server) FindEventPupils(ctx context.Context, req *eventsv1pb.FindEventPupilsRequest) (*eventsv1pb.
	FindEventPupilsResponse, error) {

	// call the svc
	pupils, total, err := s.evSvc.EventPupils(
		ctx,
		req.GetEventId(),
		eventing.EventPupilFilters{NameAndClass: req.GetNameAndClass()},
		protoPupilSortingMap[req.GetSorting()],
		int(req.GetAmount()),
		int(req.GetSkip()),
	)
	if err != nil {
		return nil, s.handleError(err)
	}

	// to proto
	pbPupils := make([]*eventsv1pb.Pupil, len(pupils))
	for i, pupil := range pupils {
		pbPupils[i] = pupilToProto(pupil)
	}
	return &eventsv1pb.FindEventPupilsResponse{
		Pupils: pbPupils,
		Total:  uint32(total),
	}, nil
}

// FindEventByID returns an event with the given id and all resources collected at that event
func (s *Server) FindEventPupilByID(ctx context.Context, req *eventsv1pb.FindEventPupilByIDRequest) (*eventsv1pb.
	FindEventPupilByIDResponse, error) {

	pupil, err := s.evSvc.PupilByID(ctx, req.GetPupilId(), req.GetEventId())
	if err != nil {
		return nil, s.handleError(err)
	}

	return &eventsv1pb.FindEventPupilByIDResponse{Pupil: pupilToProto(pupil)}, nil
}

// converts *eventsvc.Class to *eventsv1pb.Class
func classToProto(class *eventing.Class) *eventsv1pb.Class {
	if class == nil {
		return nil
	}
	return &eventsv1pb.Class{
		Name:             class.Name,
		ResourcesBrought: resourceMapToProto(class.ResourcesBrought),
	}
}

// converts *eventsvc.Event to *eventsv1pb.Event
func eventToProto(event *eventing.Event) (*eventsv1pb.Event, error) {
	if event == nil {
		return nil, nil
	}
	date, err := ptypes.TimestampProto(event.Date)
	if err != nil {
		return nil, fmt.Errorf("event date: %w", ErrInvalidTimestamp)
	}
	return &eventsv1pb.Event{
		Id:               event.ID,
		Date:             date,
		Name:             event.Name,
		ResourcesAllowed: resourcesToProto(event.ResourcesAllowed),
		ResourcesBrought: resourceMapToProto(event.ResourcesBrought),
	}, nil
}

// converts *eventsvc.Pupil to *eventsv1pb.Pupil
func pupilToProto(pupil *eventing.Pupil) *eventsv1pb.Pupil {
	if pupil == nil {
		return nil
	}
	return &eventsv1pb.Pupil{
		Id:               pupil.ID,
		FirstName:        pupil.FirstName,
		LastName:         pupil.LastName,
		Class:            pupil.Class,
		ResourcesBrought: resourceMapToProto(pupil.ResourcesBrought),
	}
}
