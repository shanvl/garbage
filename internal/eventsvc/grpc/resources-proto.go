package grpc

import (
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventsvc"
)

var resourceProtoMap = map[eventsvc.Resource]eventsv1pb.Resource{
	eventsvc.Gadgets: eventsv1pb.Resource_RESOURCE_GADGETS,
	eventsvc.Paper:   eventsv1pb.Resource_RESOURCE_PAPER,
	eventsvc.Plastic: eventsv1pb.Resource_RESOURCE_PLASTIC,
}

var protoResourceMap = map[eventsv1pb.Resource]eventsvc.Resource{
	eventsv1pb.Resource_RESOURCE_GADGETS: eventsvc.Gadgets,
	eventsv1pb.Resource_RESOURCE_PAPER:   eventsvc.Paper,
	eventsv1pb.Resource_RESOURCE_PLASTIC: eventsvc.Plastic,
}

// converts []eventsvc.Resource to []eventsv1pb.Resource
func resourcesToProto(resources []eventsvc.Resource) []eventsv1pb.Resource {
	proto := make([]eventsv1pb.Resource, len(resources))
	for i, res := range resources {
		proto[i] = resourceProtoMap[res]
	}
	return proto
}

// converts []eventsv1pb.Resource to []eventsvc.Resource
func protoToResources(proto []eventsv1pb.Resource) ([]eventsvc.Resource, error) {
	resources := make([]eventsvc.Resource, len(proto))
	for i, res := range proto {
		if res == eventsv1pb.Resource_RESOURCE_UNKNOWN {
			return nil, eventsvc.ErrUnknownResource
		}
		resources[i] = protoResourceMap[res]
	}
	return resources, nil
}

// converts *eventsv1pb.ResourcesBrought to eventsvc.ResourceMap
func protoToResourcesMap(proto *eventsv1pb.ResourcesBrought) eventsvc.ResourceMap {
	if proto == nil {
		return eventsvc.ResourceMap{}
	}
	return eventsvc.ResourceMap{
		eventsvc.Gadgets: proto.Gadgets,
		eventsvc.Paper:   proto.Paper,
		eventsvc.Plastic: proto.Plastic,
	}
}

// converts eventsvc.ResourceMap to *eventsv1pb.ResourcesBrought
func resourceMapToProto(m eventsvc.ResourceMap) *eventsv1pb.ResourcesBrought {
	return &eventsv1pb.ResourcesBrought{
		Gadgets: m[eventsvc.Gadgets],
		Paper:   m[eventsvc.Paper],
		Plastic: m[eventsvc.Plastic],
	}
}
