package eventssvc

import (
	"testing"
)

func TestEvent_IsResourceAllowed(t *testing.T) {
	tests := []struct {
		name             string
		resource         Resource
		resourcesAllowed []Resource
		want             bool
	}{
		{
			name:             "no resources allowed",
			resource:         Plastic,
			resourcesAllowed: nil,
			want:             false,
		},
		{
			name:             "given resources is not allowed",
			resource:         Plastic,
			resourcesAllowed: []Resource{Gadgets, Paper},
			want:             false,
		},
		{
			name:             "resource is allowed",
			resource:         Plastic,
			resourcesAllowed: []Resource{Plastic, Paper},
			want:             true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{
				ResourcesAllowed: tt.resourcesAllowed,
			}
			if got := e.IsResourceAllowed(tt.resource); got != tt.want {
				t.Errorf("IsResourceAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
