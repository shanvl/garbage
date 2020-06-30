package eventsvc

import (
	"errors"
	"fmt"
	"strings"
)

// Resource is a concrete type of recyclables brought by pupils to events
type Resource int

const (
	Gadgets Resource = iota
	Paper
	Plastic
)

var stringValues = [...]string{"gadgets", "paper", "plastic"}

// String returns the string value of the resource
func (r Resource) String() string {
	if r < 0 || int(r) >= len(stringValues) {
		return "unknown"
	}
	return stringValues[r]
}

// ResourceMap is a map of resources in kg
type ResourceMap map[Resource]float32

var ErrUnknownResource = errors.New("unknown resource")

// ResourceSliceToStringSlice converts a slice of resources to a slice of strings
func ResourceSliceToStringSlice(rr []Resource) []string {
	ss := make([]string, len(rr))
	for i, r := range rr {
		ss[i] = r.String()
	}
	return ss
}

var stringToValue = map[string]Resource{
	"gadgets": Gadgets,
	"paper":   Paper,
	"plastic": Plastic,
}

// StringSliceToResourceSlice converts a slice of strings to a slice of resources
func StringSliceToResourceSlice(ss []string) ([]Resource, error) {
	rr := make([]Resource, len(ss))
	for i, s := range ss {
		value, ok := stringToValue[strings.ToLower(s)]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownResource, ss[i])
		}
		rr[i] = value
	}
	return rr, nil
}
