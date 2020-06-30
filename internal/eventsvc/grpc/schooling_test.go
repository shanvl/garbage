package grpc

import (
	"context"
	"testing"

	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/internal/eventsvc/schooling"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_AddPupils(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name string
		req  *eventsv1pb.AddPupilsRequest
		code codes.Code
	}{
		{
			name: "add 1 pupil",
			req: &eventsv1pb.AddPupilsRequest{Pupils: []*eventsv1pb.AddPupilsRequest_Pupil{
				{
					FirstName: "Aa",
					LastName:  "Bb",
					Class:     "1c",
				},
			}},
			code: codes.OK,
		},
		{
			name: "add 3 pupils",
			req: &eventsv1pb.AddPupilsRequest{Pupils: []*eventsv1pb.AddPupilsRequest_Pupil{
				{
					FirstName: "Aa",
					LastName:  "Bb",
					Class:     "1c",
				},
				{
					FirstName: "Xx",
					LastName:  "Yy",
					Class:     "10c",
				},
				{
					FirstName: "Qq",
					LastName:  "Ww",
					Class:     "5c",
				},
			}},
			code: codes.OK,
		},
		{
			name: "no pupils to add",
			req:  &eventsv1pb.AddPupilsRequest{Pupils: []*eventsv1pb.AddPupilsRequest_Pupil{}},
			code: codes.InvalidArgument,
		},
		{
			name: "1 valid, 1 with invalid class",
			req: &eventsv1pb.AddPupilsRequest{Pupils: []*eventsv1pb.AddPupilsRequest_Pupil{
				{
					FirstName: "Qq",
					LastName:  "Ww",
					Class:     "5c",
				},
				{
					FirstName: "Aa",
					LastName:  "Bb",
					Class:     "12c",
				},
			}},
			code: codes.InvalidArgument,
		},
		{
			name: "1 valid, 1 with invalid first name",
			req: &eventsv1pb.AddPupilsRequest{Pupils: []*eventsv1pb.AddPupilsRequest_Pupil{
				{
					FirstName: "Qq",
					LastName:  "Ww",
					Class:     "5c",
				},
				{
					FirstName: "",
					LastName:  "Bbb",
					Class:     "12c",
				},
			}},
			code: codes.InvalidArgument,
		},
		{
			name: "1 valid, 1 with invalid last name",
			req: &eventsv1pb.AddPupilsRequest{Pupils: []*eventsv1pb.AddPupilsRequest_Pupil{
				{
					FirstName: "Qq",
					LastName:  "Ww",
					Class:     "5c",
				},
				{
					FirstName: "Aa",
					LastName:  "",
					Class:     "12c",
				},
			}},
			code: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.AddPupils(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("AddPupils() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("AddPupils() res == nil, want != nil")
				}
				if len(res.PupilIds) == 0 {
					t.Errorf("AddPupils() length of ids == 0, want > 0")
				}
				testDeletePupils(t, res.PupilIds)
			} else {
				if err == nil {
					t.Errorf("AddPupils() error == nil, wantErr == true")
					testDeletePupils(t, res.PupilIds)
				}
				if res != nil {
					t.Errorf("AddPupils() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("AddPupils() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("AddPupils() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_ChangePupilClass(t *testing.T) {
	ctx := context.Background()
	pupilID := testGetPupilID(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.ChangePupilClassRequest
		code codes.Code
	}{
		{
			name: "no pupil with that id",
			req: &eventsv1pb.ChangePupilClassRequest{
				PupilId: "nopupilwiththatid",
				Class:   "1b",
			},
			code: codes.NotFound,
		},
		{
			name: "empty id",
			req: &eventsv1pb.ChangePupilClassRequest{
				PupilId: "",
				Class:   "1b",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "invalid class",
			req: &eventsv1pb.ChangePupilClassRequest{
				PupilId: pupilID,
				Class:   "12b",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "empty class",
			req: &eventsv1pb.ChangePupilClassRequest{
				PupilId: pupilID,
				Class:   "",
			},
			code: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.ChangePupilClass(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("ChangePupilClass() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("ChangePupilClass() res == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("ChangePupilClass() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("ChangePupilClass() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("ChangePupilClass() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("ChangePupilClass() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func TestServer_RemovePupils(t *testing.T) {
	ctx := context.Background()
	pupilIDs := testAddPupils(t)
	testCases := []struct {
		name string
		req  *eventsv1pb.RemovePupilsRequest
		code codes.Code
	}{
		{
			name: "valid ids",
			req: &eventsv1pb.RemovePupilsRequest{
				PupilIds: pupilIDs,
			},
			code: codes.OK,
		},
		{
			name: "no pupils with such ids",
			req: &eventsv1pb.RemovePupilsRequest{
				PupilIds: []string{"12", "32", "sss"},
			},
			code: codes.OK,
		},
		{
			name: "empty array of ids",
			req: &eventsv1pb.RemovePupilsRequest{
				PupilIds: []string{},
			},
			code: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.RemovePupils(ctx, tc.req)
			if tc.code == codes.OK {
				if err != nil {
					t.Errorf("RemovePupils() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("RemovePupils() res == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("RemovePupils() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("RemovePupils() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("RemovePupils() couldn't get status from err %v", err)
				}
				if st.Code() != tc.code {
					t.Errorf("RemovePupils() err codes mismatch: code == %v, want == %v", st.Code(), tc.code)
				}
			}
		})
	}
}

func testDeletePupils(t *testing.T, ids []string) {
	t.Helper()
	err := schoolingRepo.RemovePupils(context.Background(), ids)
	if err != nil {
		t.Fatalf("could't delete pupils: %v", err)
	}
}

func testAddPupils(t *testing.T) []string {
	t.Helper()
	ids := []string{"xxxxxx", "yyyyyyy", "zzzzzzzz"}
	err := schoolingRepo.StorePupils(context.Background(), []*schooling.Pupil{
		{
			Pupil: eventsvc.Pupil{
				ID:        ids[0],
				FirstName: "aa",
				LastName:  "bb",
			},
			Class: eventsvc.Class{
				Letter:     "a",
				DateFormed: date(2015, 9, 1),
			},
		},
		{
			Pupil: eventsvc.Pupil{
				ID:        ids[1],
				FirstName: "aa",
				LastName:  "bb",
			},
			Class: eventsvc.Class{
				Letter:     "b",
				DateFormed: date(2016, 9, 1),
			},
		},
		{
			Pupil: eventsvc.Pupil{
				ID:        ids[2],
				FirstName: "dd",
				LastName:  "zz",
			},
			Class: eventsvc.Class{
				Letter:     "c",
				DateFormed: date(2012, 9, 1),
			},
		},
	})
	if err != nil {
		t.Fatalf("couldn't add pupils %v", err)
	}
	return ids
}
