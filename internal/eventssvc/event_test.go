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
			resource:         "plastic",
			resourcesAllowed: nil,
			want:             false,
		},
		{
			name:             "given resources is not allowed",
			resource:         "plastic",
			resourcesAllowed: []Resource{"gadgets", "paper"},
			want:             false,
		},
		{
			name:             "resource is allowed",
			resource:         "plastic",
			resourcesAllowed: []Resource{"plastic", "paper"},
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
