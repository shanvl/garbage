package eventssvc

import (
	"errors"
	"fmt"
)

// Resource is a concrete type of recyclables brought by pupils to events
type Resource string

const (
	Gadgets Resource = "gadgets"
	Paper   Resource = "paper"
	Plastic Resource = "plastic"
)

// IsKnown indicates whether a given resource is one of the known
func (r Resource) IsKnown() bool {
	switch r {
	case Gadgets, Paper, Plastic:
		return true
	}
	return false
}

var ErrUnknownResource = errors.New("unknown resource")

// ResourceMap is a map of resources in kg
type ResourceMap map[Resource]float32

// ResourceSliceToStringSlice converts a slice of resources to a slice of strings
func ResourceSliceToStringSlice(rr []Resource) []string {
	ss := make([]string, len(rr))
	for i, r := range rr {
		ss[i] = string(r)
	}
	return ss
}

// StringSliceToResourceSlice converts a slice of strings to a slice of resources
func StringSliceToResourceSlice(ss []string) ([]Resource, error) {
	rr := make([]Resource, len(ss))
	for i, s := range ss {
		r := Resource(s)
		if !r.IsKnown() {
			return nil, fmt.Errorf("%w: %s", ErrUnknownResource, rr[i])
		}
		rr[i] = r
	}
	return rr, nil
}
