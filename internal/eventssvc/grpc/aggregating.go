package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventssvc"
	"github.com/shanvl/garbage/internal/eventssvc/aggregating"
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

// FindPupilByID returns a pupil with the given ID with the list of all resources they has brought to the events
// that passed the provided filters
func (s *Server) FindPupilByID(ctx context.Context, req *eventsv1pb.FindPupilByIDRequest) (*eventsv1pb.
	FindPupilByIDResponse, error) {

	eventFilters, err := protoToEventFilters(req.GetEventFilters())
	if err != nil {
		return nil, s.handleError(err)
	}
	pupil, err := s.agSvc.PupilByID(
		ctx,
		req.GetId(),
		eventFilters,
		protoEventSortingMap[req.GetEventSorting()],
	)
	if err != nil {
		return nil, s.handleError(err)
	}

	pbPupil, err := pupilAggrToProto(pupil)
	if err != nil {
		return nil, s.handleError(err)
	}
	return &eventsv1pb.FindPupilByIDResponse{Pupil: pbPupil}, nil
}

func protoToEventFilters(proto *eventsv1pb.EventFilters) (aggregating.EventFilters, error) {
	if proto == nil {
		return aggregating.EventFilters{}, nil
	}
	resourcesAllowed, err := protoToResources(proto.GetResourcesAllowed())
	if err != nil {
		return aggregating.EventFilters{}, err
	}
	from, err := protoTimeToTimestamp(proto.GetFrom())
	if err != nil {
		return aggregating.EventFilters{}, fmt.Errorf("event from: %w", err)
	}
	to, err := protoTimeToTimestamp(proto.GetTo())
	if err != nil {
		return aggregating.EventFilters{}, fmt.Errorf("event to: %w", err)
	}
	return aggregating.EventFilters{
		From:             from,
		To:               to,
		Name:             proto.GetName(),
		ResourcesAllowed: resourcesAllowed,
	}, nil
}

// converts *eventssvc.Event to *eventsv1pb.Event
func eventAggrToProto(event *aggregating.Event) (*eventsv1pb.Event, error) {
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
		ResourcesBrought: &eventsv1pb.ResourcesBrought{
			Gadgets: event.ResourcesBrought[eventssvc.Gadgets],
			Paper:   event.ResourcesBrought[eventssvc.Paper],
			Plastic: event.ResourcesBrought[eventssvc.Plastic],
		},
	}, nil
}

func pupilAggrToProto(pupil *aggregating.Pupil) (*eventsv1pb.PupilAggr, error) {
	// events to proto
	pbEvents := make([]*eventsv1pb.Event, len(pupil.Events))
	for i, event := range pupil.Events {
		pbEvent, err := eventAggrToProto(event)
		if err != nil {
			return nil, err
		}
		pbEvents[i] = pbEvent
	}
	// class date formed to proto
	pbClassDateFormed, err := ptypes.TimestampProto(pupil.DateFormed)
	if err != nil {
		return nil, fmt.Errorf("class date formed: %w", ErrInvalidTimestamp)
	}
	return &eventsv1pb.PupilAggr{
		Id:               pupil.ID,
		FirstName:        pupil.FirstName,
		LastName:         pupil.LastName,
		ClassLetter:      pupil.Letter,
		ClassDateFormed:  pbClassDateFormed,
		ResourcesBrought: resourceMapToProto(pupil.ResourcesBrought),
		Events:           pbEvents,
	}, nil
}

func protoTimeToTimestamp(proto *timestamp.Timestamp) (time.Time, error) {
	if proto == nil {
		return time.Time{}, nil
	}
	ts, err := ptypes.Timestamp(proto)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrInvalidTimestamp, err)
	}
	return ts, nil
}
