package garbage

// Resource is a concrete type of recyclable brought by pupils to events
type Resource string

const (
	Gadgets Resource = "gadgets"
	Paper   Resource = "paper"
	Plastic Resource = "plastic"
)

// IsKnown indicates whether a given resource is one of the known. Used, for example, in json decoding
func (r Resource) IsKnown() bool {
	for _, resource := range []Resource{Gadgets, Paper, Plastic} {
		if r == resource {
			return true
		}
	}
	return false
}
