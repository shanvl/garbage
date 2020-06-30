package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/internal/eventsvc/aggregating"
	"github.com/shanvl/garbage/internal/eventsvc/sorting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_FindClasses(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name string
		req  *eventsv1pb.FindClassesRequest
		code codes.Code
	}{
		{
			name: "no filters",
			req: &eventsv1pb.FindClassesRequest{
				Letter:       "",
				DateFormed:   nil,
				EventFilters: nil,
				Sorting:      0,
				EventSorting: 0,
				Amount:       0,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "skip > amount of records",
			req: &eventsv1pb.FindClassesRequest{
				Letter:       "",
				DateFormed:   nil,
				EventFilters: nil,
				Sorting:      0,
				EventSorting: 0,
				Amount:       0,
				Skip:         9999,
			},
			code: codes.OK,
		},
		{
			name: "all filters are set",
			req: &eventsv1pb.FindClassesRequest{
				Letter:     "a",
				DateFormed: testTimeToProto(t, date(2015, 9, 1)),
				EventFilters: &eventsv1pb.EventFilters{
					From:             testTimeToProto(t, time.Now().AddDate(-5, 0, 0)),
					To:               testTimeToProto(t, time.Now().AddDate(5, 0, 0)),
					Name:             "ev",
					ResourcesAllowed: resourcesToProto([]eventsvc.Resource{eventsvc.Gadgets, eventsvc.Plastic}),
				},
				Sorting:      eventsv1pb.ClassSorting_CLASS_SORTING_PAPER,
				EventSorting: eventsv1pb.EventSorting_EVENT_SORTING_PAPER,
				Amount:       50,
				Skip:         10,
			},
			code: codes.OK,
		},
		{
			name: "invalid class letter",
			req: &eventsv1pb.FindClassesRequest{
				Letter:       "bb",
				DateFormed:   nil,
				EventFilters: nil,
				Sorting:      0,
				EventSorting: 0,
				Amount:       0,
				Skip:         0,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "unknown resource",
			req: &eventsv1pb.FindClassesRequest{
				Letter:     "bb",
				DateFormed: nil,
				EventFilters: &eventsv1pb.EventFilters{
					From:             nil,
					To:               nil,
					Name:             "",
					ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_UNKNOWN},
				},
				Sorting:      0,
				EventSorting: 0,
				Amount:       0,
				Skip:         0,
			},
			code: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.FindClasses(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("FindClasses() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindClasses() res == nil, want != nil")
				}
				if res.Classes == nil {
					t.Errorf("FindClasses() classes == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("FindClasses() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindClasses() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindClasses() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("FindClasses() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_FindEvents(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name string
		req  *eventsv1pb.FindEventsRequest
		code codes.Code
	}{
		{
			name: "no filters",
			req: &eventsv1pb.FindEventsRequest{
				Filters: nil,
				Sorting: 0,
				Amount:  0,
				Skip:    0,
			},
			code: codes.OK,
		},
		{
			name: "skip > amount of records",
			req: &eventsv1pb.FindEventsRequest{
				Filters: nil,
				Sorting: 0,
				Amount:  0,
				Skip:    9999,
			},
			code: codes.OK,
		},
		{
			name: "all filters are set",
			req: &eventsv1pb.FindEventsRequest{
				Filters: &eventsv1pb.EventFilters{
					From:             testTimeToProto(t, time.Now().AddDate(-5, 0, 0)),
					To:               testTimeToProto(t, time.Now().AddDate(5, 0, 0)),
					Name:             "ev",
					ResourcesAllowed: resourcesToProto([]eventsvc.Resource{eventsvc.Gadgets, eventsvc.Plastic}),
				},
				Sorting: eventsv1pb.EventSorting_EVENT_SORTING_GADGETS,
				Amount:  50,
				Skip:    10,
			},
			code: codes.OK,
		},
		{
			name: "no events with that name",
			req: &eventsv1pb.FindEventsRequest{
				Filters: &eventsv1pb.EventFilters{
					From:             testTimeToProto(t, time.Now().AddDate(-5, 0, 0)),
					To:               testTimeToProto(t, time.Now().AddDate(5, 0, 0)),
					Name:             "deedee megadoodoo",
					ResourcesAllowed: resourcesToProto([]eventsvc.Resource{eventsvc.Gadgets, eventsvc.Plastic}),
				},
				Sorting: eventsv1pb.EventSorting_EVENT_SORTING_GADGETS,
				Amount:  50,
				Skip:    10,
			},
			code: codes.OK,
		},
		{
			name: "invalid resources",
			req: &eventsv1pb.FindEventsRequest{
				Filters: &eventsv1pb.EventFilters{
					From:             testTimeToProto(t, time.Now().AddDate(-5, 0, 0)),
					To:               testTimeToProto(t, time.Now().AddDate(5, 0, 0)),
					Name:             "ev",
					ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_UNKNOWN},
				},
				Sorting: 0,
				Amount:  0,
				Skip:    0,
			},
			code: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.FindEvents(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("FindEvents() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindEvents() res == nil, want != nil")
				}
				if res.Events == nil {
					t.Errorf("FindEvents() events == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("FindEvents() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindEvents() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindEvents() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("FindEvents() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_FindPupils(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name string
		req  *eventsv1pb.FindPupilsRequest
		code codes.Code
	}{
		{
			name: "no filters",
			req: &eventsv1pb.FindPupilsRequest{
				NameAndClass: "",
				EventFilters: nil,
				Sorting:      0,
				EventSorting: 0,
				Amount:       0,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "no pupils by that name",
			req: &eventsv1pb.FindPupilsRequest{
				NameAndClass: "deedee megadoodoo",
				EventFilters: nil,
				Sorting:      0,
				EventSorting: 0,
				Amount:       0,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "skip > amount of records",
			req: &eventsv1pb.FindPupilsRequest{
				NameAndClass: "",
				EventFilters: nil,
				Sorting:      0,
				EventSorting: 0,
				Amount:       0,
				Skip:         9999,
			},
			code: codes.OK,
		},
		{
			name: "all filters are set",
			req: &eventsv1pb.FindPupilsRequest{
				NameAndClass: "an a",
				EventFilters: &eventsv1pb.EventFilters{
					From:             testTimeToProto(t, time.Now().AddDate(-5, 0, 0)),
					To:               testTimeToProto(t, time.Now().AddDate(5, 0, 0)),
					Name:             "ev",
					ResourcesAllowed: resourcesToProto([]eventsvc.Resource{eventsvc.Gadgets, eventsvc.Plastic}),
				},
				Sorting:      eventsv1pb.PupilSorting_PUPIL_SORTING_GADGETS,
				EventSorting: eventsv1pb.EventSorting_EVENT_SORTING_GADGETS,
				Amount:       50,
				Skip:         10,
			},
			code: codes.OK,
		},
		{
			name: "unknown resource",
			req: &eventsv1pb.FindPupilsRequest{
				NameAndClass: "an a",
				EventFilters: &eventsv1pb.EventFilters{
					From:             testTimeToProto(t, time.Now().AddDate(-5, 0, 0)),
					To:               testTimeToProto(t, time.Now().AddDate(5, 0, 0)),
					Name:             "ev",
					ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_UNKNOWN},
				},
				Sorting:      eventsv1pb.PupilSorting_PUPIL_SORTING_GADGETS,
				EventSorting: eventsv1pb.EventSorting_EVENT_SORTING_GADGETS,
				Amount:       50,
				Skip:         10,
			},
			code: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.FindPupils(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("FindPupils() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindPupils() res == nil, want != nil")
				}
				if res.Pupils == nil {
					t.Errorf("FindPupils() pupils == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("FindPupils() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindPupils() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindPupils() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("FindPupils() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_FindPupilByID(t *testing.T) {
	ctx := context.Background()
	pupilID := testGetPupilID(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.FindPupilByIDRequest
		code codes.Code
	}{
		{
			name: "no pupil with that id",
			req: &eventsv1pb.FindPupilByIDRequest{
				Id:           "somerandomid",
				EventFilters: nil,
				EventSorting: 0,
			},
			code: codes.NotFound,
		},
		{
			name: "pupil with no events",
			req: &eventsv1pb.FindPupilByIDRequest{
				Id: pupilID,
				EventFilters: &eventsv1pb.EventFilters{
					From:             testTimeToProto(t, time.Now().AddDate(15, 0, 0)),
					To:               testTimeToProto(t, time.Now().AddDate(15, 0, 0)),
					Name:             "nosuchevents",
					ResourcesAllowed: nil,
				},
				EventSorting: 0,
			},
			code: codes.OK,
		},
		{
			name: "pupil must be found",
			req: &eventsv1pb.FindPupilByIDRequest{
				Id:           pupilID,
				EventFilters: nil,
				EventSorting: 0,
			},
			code: codes.OK,
		},
		{
			name: "no pupil with that id",
			req: &eventsv1pb.FindPupilByIDRequest{
				Id:           "somerandomid",
				EventFilters: nil,
				EventSorting: 0,
			},
			code: codes.NotFound,
		},
		{
			name: "invalid resource",
			req: &eventsv1pb.FindPupilByIDRequest{
				Id: pupilID,
				EventFilters: &eventsv1pb.EventFilters{
					From:             nil,
					To:               nil,
					Name:             "",
					ResourcesAllowed: []eventsv1pb.Resource{eventsv1pb.Resource_RESOURCE_UNKNOWN},
				},
				EventSorting: 0,
			},
			code: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.FindPupilByID(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("FindPupilByID() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindPupilByID() res == nil, want != nil")
				}
				if res.Pupil == nil {
					t.Errorf("FindPupilByID() pupil == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("FindPupilByID() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindPupilByID() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindPupilByID() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("FindPupilByID() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func testGetPupilID(t *testing.T) string {
	t.Helper()
	pupils, _, err := aggregatingRepo.Pupils(context.Background(), aggregating.PupilFilters{}, sorting.NameDes,
		sorting.NameDes, 1, 0)
	if err != nil || len(pupils) == 0 {
		t.Fatalf("could't find a pupil: %v", err)
	}
	pupilID := pupils[0].ID
	return pupilID
}

func testTimeToProto(t *testing.T, ts time.Time) *timestamp.Timestamp {
	t.Helper()
	protoTS, err := ptypes.TimestampProto(ts)
	if err != nil {
		t.Fatalf("couldn't convert timestamp to protoTS: %v", err)
	}
	return protoTS
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
