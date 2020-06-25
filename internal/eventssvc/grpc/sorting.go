package grpc

import (
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventssvc/sorting"
)

var protoClassSortingMap = map[eventsv1pb.ClassSorting]sorting.By{
	eventsv1pb.ClassSorting_CLASS_SORTING_GADGETS:   sorting.Gadgets,
	eventsv1pb.ClassSorting_CLASS_SORTING_NAME_ASC:  sorting.NameAsc,
	eventsv1pb.ClassSorting_CLASS_SORTING_NAME_DESC: sorting.NameDes,
	eventsv1pb.ClassSorting_CLASS_SORTING_PAPER:     sorting.Paper,
	eventsv1pb.ClassSorting_CLASS_SORTING_PLASTIC:   sorting.Plastic,
	eventsv1pb.ClassSorting_CLASS_SORTING_UNKNOWN:   sorting.Unspecified,
}

var protoPupilSortingMap = map[eventsv1pb.PupilSorting]sorting.By{
	eventsv1pb.PupilSorting_PUPIL_SORTING_GADGETS:   sorting.Gadgets,
	eventsv1pb.PupilSorting_PUPIL_SORTING_NAME_ASC:  sorting.NameAsc,
	eventsv1pb.PupilSorting_PUPIL_SORTING_NAME_DESC: sorting.NameDes,
	eventsv1pb.PupilSorting_PUPIL_SORTING_PAPER:     sorting.Paper,
	eventsv1pb.PupilSorting_PUPIL_SORTING_PLASTIC:   sorting.Plastic,
	eventsv1pb.PupilSorting_PUPIL_SORTING_UNKNOWN:   sorting.Unspecified,
}

var protoEventSortingMap = map[eventsv1pb.EventSorting]sorting.By{
	eventsv1pb.EventSorting_EVENT_SORTING_DATE_ASC:  sorting.DateAsc,
	eventsv1pb.EventSorting_EVENT_SORTING_DATE_DESC: sorting.DateDes,
	eventsv1pb.EventSorting_EVENT_SORTING_GADGETS:   sorting.Gadgets,
	eventsv1pb.EventSorting_EVENT_SORTING_NAME_ASC:  sorting.NameAsc,
	eventsv1pb.EventSorting_EVENT_SORTING_NAME_DESC: sorting.NameDes,
	eventsv1pb.EventSorting_EVENT_SORTING_PAPER:     sorting.Paper,
	eventsv1pb.EventSorting_EVENT_SORTING_PLASTIC:   sorting.Plastic,
	eventsv1pb.EventSorting_EVENT_SORTING_UNKNOWN:   sorting.Unspecified,
}
