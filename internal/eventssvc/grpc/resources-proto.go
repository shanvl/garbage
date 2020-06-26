package grpc

import (
	eventsv1pb "github.com/shanvl/garbage/api/events/v1/pb"
	"github.com/shanvl/garbage/internal/eventssvc"
)

var resourceProtoMap = map[eventssvc.Resource]eventsv1pb.Resource{
	eventssvc.Gadgets: eventsv1pb.Resource_RESOURCE_GADGETS,
	eventssvc.Paper:   eventsv1pb.Resource_RESOURCE_PAPER,
	eventssvc.Plastic: eventsv1pb.Resource_RESOURCE_PLASTIC,
}

var protoResourceMap = map[eventsv1pb.Resource]eventssvc.Resource{
	eventsv1pb.Resource_RESOURCE_GADGETS: eventssvc.Gadgets,
	eventsv1pb.Resource_RESOURCE_PAPER:   eventssvc.Paper,
	eventsv1pb.Resource_RESOURCE_PLASTIC: eventssvc.Plastic,
}

// converts []eventssvc.Resource to []eventsv1pb.Resource
func resourcesToProto(resources []eventssvc.Resource) []eventsv1pb.Resource {
	proto := make([]eventsv1pb.Resource, len(resources))
	for i, res := range resources {
		proto[i] = resourceProtoMap[res]
	}
	return proto
}

// converts []eventsv1pb.Resource to []eventssvc.Resource
func protoToResources(proto []eventsv1pb.Resource) ([]eventssvc.Resource, error) {
	resources := make([]eventssvc.Resource, len(proto))
	for i, res := range proto {
		if res == eventsv1pb.Resource_RESOURCE_UNKNOWN {
			return nil, eventssvc.ErrUnknownResource
		}
		resources[i] = protoResourceMap[res]
	}
	return resources, nil
}

// converts *eventsv1pb.ResourcesBrought to eventssvc.ResourceMap
func protoToResourcesMap(proto *eventsv1pb.ResourcesBrought) eventssvc.ResourceMap {
	if proto == nil {
		return eventssvc.ResourceMap{}
	}
	return eventssvc.ResourceMap{
		eventssvc.Gadgets: proto.Gadgets,
		eventssvc.Paper:   proto.Paper,
		eventssvc.Plastic: proto.Plastic,
	}
}

// converts eventssvc.ResourceMap to *eventsv1pb.ResourcesBrought
func resourceMapToProto(m eventssvc.ResourceMap) *eventsv1pb.ResourcesBrought {
	return &eventsv1pb.ResourcesBrought{
		Gadgets: m[eventssvc.Gadgets],
		Paper:   m[eventssvc.Paper],
		Plastic: m[eventssvc.Plastic],
	}
}
