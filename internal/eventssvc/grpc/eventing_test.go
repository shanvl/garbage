package grpc

import (
	"context"
	"testing"
	"time"

	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventssvc"
	"github.com/shanvl/garbage/internal/eventssvc/aggregating"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/internal/eventssvc/sorting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_ChangePupilResources(t *testing.T) {
	ctx := context.Background()
	// we need an event with all resources allowed
	eventID := getEventIDWithResourcesAllowed(t, []eventssvc.Resource{eventssvc.Gadgets, eventssvc.Paper,
		eventssvc.Plastic})
	pupilID := getPupilID(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.ChangePupilResourcesRequest
		code codes.Code
	}{
		{
			name: "empty event id",
			req: &eventsv1pb.ChangePupilResourcesRequest{
				EventId: "",
				PupilId: pupilID,
				ResourcesBrought: &eventsv1pb.ResourcesBrought{
					Gadgets: 50,
					Paper:   50,
					Plastic: 50,
				},
			},
			code: codes.InvalidArgument,
		},
		{
			name: "empty pupil id",
			req: &eventsv1pb.ChangePupilResourcesRequest{
				EventId: eventID,
				PupilId: "",
				ResourcesBrought: &eventsv1pb.ResourcesBrought{
					Gadgets: 50,
					Paper:   50,
					Plastic: 50,
				},
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no such event",
			req: &eventsv1pb.ChangePupilResourcesRequest{
				EventId: "somerandomid",
				PupilId: pupilID,
				ResourcesBrought: &eventsv1pb.ResourcesBrought{
					Gadgets: 50,
					Paper:   50,
					Plastic: 50,
				},
			},
			code: codes.NotFound,
		},
		{
			name: "no such pupil",
			req: &eventsv1pb.ChangePupilResourcesRequest{
				EventId: eventID,
				PupilId: "somerandomid",
				ResourcesBrought: &eventsv1pb.ResourcesBrought{
					Gadgets: 50,
					Paper:   50,
					Plastic: 50,
				},
			},
			code: codes.NotFound,
		},
		{
			name: "no ResourcesBrought",
			req: &eventsv1pb.ChangePupilResourcesRequest{
				EventId:          eventID,
				PupilId:          pupilID,
				ResourcesBrought: nil,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "one of the resources is less than 0",
			req: &eventsv1pb.ChangePupilResourcesRequest{
				EventId: eventID,
				PupilId: pupilID,
				ResourcesBrought: &eventsv1pb.ResourcesBrought{
					Gadgets: 50,
					Paper:   0,
					Plastic: -50,
				},
			},
			code: codes.InvalidArgument,
		},
		{
			name: "50.09 gadgets, 5.15 plastic",
			req: &eventsv1pb.ChangePupilResourcesRequest{
				EventId: eventID,
				PupilId: pupilID,
				ResourcesBrought: &eventsv1pb.ResourcesBrought{
					Gadgets: 50.09,
					Paper:   0,
					Plastic: 5.15,
				},
			},
			code: codes.OK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.ChangePupilResources(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("ChangePupilResources() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("ChangePupilResources() res == nil, want != nil")
				}
				// compare request's resources with repo's resources
				reqResourcesBrought := tc.req.ResourcesBrought
				pupilInRepo := getEventPupilByID(t, tc.req.PupilId, tc.req.EventId)
				repoResBrough := pupilInRepo.ResourcesBrought
				if reqResourcesBrought.Gadgets != repoResBrough[eventssvc.Gadgets] || reqResourcesBrought.
					Paper != repoResBrough[eventssvc.Paper] || reqResourcesBrought.
					Plastic != repoResBrough[eventssvc.Plastic] {
					t.Errorf("ChangePupilResources() resources don't match")
				}
				// change repo's pupil's resources to some other value in order not to break the next test run
				changePupilResources(t, tc.req.EventId, tc.req.PupilId,
					eventssvc.ResourceMap{eventssvc.Gadgets: 1000, eventssvc.Paper: 222, eventssvc.Plastic: 5})
			} else {
				if err == nil {
					t.Errorf("ChangePupilResources() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("ChangePupilResources() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("ChangePupilResources() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("ChangePupilResources() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_CreateEvent(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name string
		req  *eventsv1pb.CreateEventRequest
		code codes.Code
	}{
		{
			name: "no resourcesAllowed",
			req: &eventsv1pb.CreateEventRequest{
				Date:             testTimeToProto(t, time.Now().AddDate(1, 0, 0)),
				Name:             "event name",
				ResourcesAllowed: nil,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "unknown resource",
			req: &eventsv1pb.CreateEventRequest{
				Date:             testTimeToProto(t, time.Now().AddDate(1, 0, 0)),
				Name:             "event name",
				ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_UNKNOWN},
			},
			code: codes.InvalidArgument,
		},
		{
			name: "event's date is in the past",
			req: &eventsv1pb.CreateEventRequest{
				Date:             testTimeToProto(t, time.Now().AddDate(-1, 0, 0)),
				Name:             "event name",
				ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_PLASTIC},
			},
			code: codes.InvalidArgument,
		},
		{
			name: "has no date",
			req: &eventsv1pb.CreateEventRequest{
				Date: nil,
				Name: "event name",
				ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_PLASTIC,
					eventsv1pb.Resource_RESOURCE_GADGETS},
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no name",
			req: &eventsv1pb.CreateEventRequest{
				Date: testTimeToProto(t, time.Now().AddDate(1, 0, 0)),
				Name: "",
				ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_PLASTIC,
					eventsv1pb.Resource_RESOURCE_GADGETS},
			},
			code: codes.OK,
		},
		{
			name: "has name, has date, has resources",
			req: &eventsv1pb.CreateEventRequest{
				Date: testTimeToProto(t, time.Now().AddDate(1, 0, 0)),
				Name: "event name",
				ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_PLASTIC,
					eventsv1pb.Resource_RESOURCE_GADGETS},
			},
			code: codes.OK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.CreateEvent(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("CreateEvent() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("CreateEvent() res == nil, want != nil")
				}
				if len(res.Id) == 0 {
					t.Errorf("CreateEvent() length of the id == 0, want > 0")
				}
				testDeleteEvent(t, res.Id)
			} else {
				if err == nil {
					t.Errorf("CreateEvent() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("CreateEvent() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("CreateEvent() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("CreateEvent() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
				if st.Code() == codes.OK && res != nil {
					testDeleteEvent(t, res.Id)
				}
			}
		})
	}
}

func TestServer_DeleteEvent(t *testing.T) {
	ctx := context.Background()
	eventID := testCreateEvent(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.DeleteEventRequest
		code codes.Code
	}{
		{
			name: "no id",
			req: &eventsv1pb.DeleteEventRequest{
				Id: "",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no event with that id",
			req: &eventsv1pb.DeleteEventRequest{
				Id: "somerandomid",
			},
			code: codes.OK,
		},
		{
			name: "valid id",
			req: &eventsv1pb.DeleteEventRequest{
				Id: eventID,
			},
			code: codes.OK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.DeleteEvent(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("DeleteEvent() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("DeleteEvent() res == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("DeleteEvent() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("DeleteEvent() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("DeleteEvent() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("DeleteEvent() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
		testDeleteEvent(t, eventID)
	}
}

func TestServer_FindEventByID(t *testing.T) {
	ctx := context.Background()
	eventID := getEventID(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.FindEventByIDRequest
		code codes.Code
	}{
		{
			name: "no event id",
			req: &eventsv1pb.FindEventByIDRequest{
				Id: "",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no event with that id",
			req: &eventsv1pb.FindEventByIDRequest{
				Id: "somerandomeventid",
			},
			code: codes.NotFound,
		},
		{
			name: "event has been found",
			req: &eventsv1pb.FindEventByIDRequest{
				Id: eventID,
			},
			code: codes.OK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.FindEventByID(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("FindEventByID() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindEventByID() res == nil, want != nil")
				}
				if res.Event == nil {
					t.Errorf("FindEventByID() event == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("FindEventByID() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindEventByID() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindEventByID() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("FindEventByID() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_FindEventClasses(t *testing.T) {
	ctx := context.Background()
	eventID := getEventID(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.FindEventClassesRequest
		code codes.Code
	}{
		{
			name: "no event id",
			req: &eventsv1pb.FindEventClassesRequest{
				EventId:   "",
				ClassName: "1b",
				Sorting:   eventsv1pb.ClassSorting_CLASS_SORTING_PAPER,
				Amount:    10,
				Skip:      0,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no such event",
			req: &eventsv1pb.FindEventClassesRequest{
				EventId:   "somerandomeventid",
				ClassName: "1b",
				Sorting:   eventsv1pb.ClassSorting_CLASS_SORTING_PAPER,
				Amount:    10,
				Skip:      0,
			},
			code: codes.NotFound,
		},
		{
			name: "no filters",
			req: &eventsv1pb.FindEventClassesRequest{
				EventId:   eventID,
				ClassName: "",
				Sorting:   0,
				Amount:    0,
				Skip:      0,
			},
			code: codes.OK,
		},
		{
			name: "skip > amount of records",
			req: &eventsv1pb.FindEventClassesRequest{
				EventId:   eventID,
				ClassName: "1b",
				Sorting:   eventsv1pb.ClassSorting_CLASS_SORTING_PAPER,
				Amount:    10,
				Skip:      9999,
			},
			code: codes.OK,
		},
		{
			name: "invalid class name",
			req: &eventsv1pb.FindEventClassesRequest{
				EventId:   eventID,
				ClassName: "12bb",
				Sorting:   eventsv1pb.ClassSorting_CLASS_SORTING_PAPER,
				Amount:    10,
				Skip:      0,
			},
			code: codes.OK,
		},
		{
			name: "sorting unknown",
			req: &eventsv1pb.FindEventClassesRequest{
				EventId:   eventID,
				ClassName: "12bb",
				Sorting:   eventsv1pb.ClassSorting_CLASS_SORTING_UNKNOWN,
				Amount:    10,
				Skip:      0,
			},
			code: codes.OK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.FindEventClasses(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("FindEventClasses() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindEventClasses() res == nil, want != nil")
				}
				if res.Classes == nil {
					t.Errorf("FindEventClasses() classes == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("FindEventClasses() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindEventClasses() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindEventClasses() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("FindEventClasses() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_FindEventPupils(t *testing.T) {
	ctx := context.Background()
	eventID := getEventID(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.FindEventPupilsRequest
		code codes.Code
	}{
		{
			name: "no event id",
			req: &eventsv1pb.FindEventPupilsRequest{
				EventId:      "",
				NameAndClass: "a 1a",
				Sorting:      eventsv1pb.PupilSorting_PUPIL_SORTING_GADGETS,
				Amount:       10,
				Skip:         0,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no such event",
			req: &eventsv1pb.FindEventPupilsRequest{
				EventId:      "somerandomeventid",
				NameAndClass: "a 1a",
				Sorting:      eventsv1pb.PupilSorting_PUPIL_SORTING_GADGETS,
				Amount:       10,
				Skip:         0,
			},
			code: codes.NotFound,
		},
		{
			name: "no filters",
			req: &eventsv1pb.FindEventPupilsRequest{
				EventId:      eventID,
				NameAndClass: "",
				Sorting:      0,
				Amount:       0,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "all filters are set",
			req: &eventsv1pb.FindEventPupilsRequest{
				EventId:      eventID,
				NameAndClass: "a 1a",
				Sorting:      eventsv1pb.PupilSorting_PUPIL_SORTING_GADGETS,
				Amount:       10,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "skip > amount of records",
			req: &eventsv1pb.FindEventPupilsRequest{
				EventId:      eventID,
				NameAndClass: "a 1a",
				Sorting:      eventsv1pb.PupilSorting_PUPIL_SORTING_GADGETS,
				Amount:       9999,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "sorting unknown",
			req: &eventsv1pb.FindEventPupilsRequest{
				EventId:      eventID,
				NameAndClass: "a 1a",
				Sorting:      eventsv1pb.PupilSorting_PUPIL_SORTING_UNKNOWN,
				Amount:       10,
				Skip:         0,
			},
			code: codes.OK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.FindEventPupils(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("FindEventPupils() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindEventPupils() res == nil, want != nil")
				}
				if res.Pupils == nil {
					t.Errorf("FindEventPupils() pupils == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("FindEventPupils() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindEventPupils() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindEventPupils() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("FindEventPupils() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_FindEventPupilByID(t *testing.T) {
	ctx := context.Background()
	eventID := getEventID(t)
	pupilID := getPupilID(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.FindEventPupilByIDRequest
		code codes.Code
	}{
		{
			name: "no event id",
			req: &eventsv1pb.FindEventPupilByIDRequest{
				EventId: "",
				PupilId: pupilID,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no pupil id",
			req: &eventsv1pb.FindEventPupilByIDRequest{
				EventId: eventID,
				PupilId: "",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no such event",
			req: &eventsv1pb.FindEventPupilByIDRequest{
				EventId: "random event id",
				PupilId: pupilID,
			},
			code: codes.NotFound,
		},
		{
			name: "no such pupil",
			req: &eventsv1pb.FindEventPupilByIDRequest{
				EventId: eventID,
				PupilId: "random pupil id",
			},
			code: codes.NotFound,
		},
		{
			name: "pupil found",
			req: &eventsv1pb.FindEventPupilByIDRequest{
				EventId: eventID,
				PupilId: pupilID,
			},
			code: codes.OK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.FindEventPupilByID(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("FindEventPupilByID() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindEventPupilByID() res == nil, want != nil")
				}
				if res.Pupil == nil {
					t.Errorf("FindEventPupilByID() pupil == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("FindEventPupilByID() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindEventPupilByID() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindEventPupilByID() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("FindEventPupilByID() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func getEventID(t *testing.T) string {
	events, _, err := aggregatingRepo.Events(context.Background(), aggregating.EventFilters{}, sorting.NameDes, 1, 0)
	if err != nil || len(events) == 0 {
		t.Fatalf("couldn't find an event: %v", err)
	}
	return events[0].ID
}

func getEventIDWithResourcesAllowed(t *testing.T, resourcesAllowed []eventssvc.Resource) string {
	events, _, err := aggregatingRepo.Events(context.Background(), aggregating.EventFilters{
		ResourcesAllowed: resourcesAllowed,
	}, sorting.NameDes, 1, 0)
	if err != nil || len(events) == 0 {
		t.Fatalf("couldn't find an event: %v", err)
	}
	return events[0].ID
}

func getEventPupilByID(t *testing.T, pupilID, eventID string) *eventing.Pupil {
	pupil, err := eventingRepo.PupilByID(context.Background(), pupilID, eventID)
	if err != nil {
		t.Fatalf("couldn't get a pupil: %v", err)
	}
	return pupil
}

func changePupilResources(t *testing.T, eventID, pupilID string, resources eventssvc.ResourceMap) {
	err := eventingRepo.ChangePupilResources(context.Background(), eventID, pupilID, resources)
	if err != nil {
		t.Fatalf("wasn't able to change pupil's resources back: %v", err)
	}
}

func testDeleteEvent(t *testing.T, eventID string) {
	err := eventingRepo.DeleteEvent(context.Background(), eventID)
	if err != nil {
		t.Fatalf("event wasn't deleted: %v", err)
	}
}

func testCreateEvent(t *testing.T) string {
	eventID := "qwertylopki"
	err := eventingRepo.StoreEvent(context.Background(), &eventssvc.Event{
		ID:               eventID,
		Date:             time.Now().AddDate(1, 0, 0),
		Name:             "",
		ResourcesAllowed: []eventssvc.Resource{eventssvc.Plastic},
	})
	if err != nil {
		t.Fatalf("couldn't create an event: %v", err)
	}
	return eventID
}
