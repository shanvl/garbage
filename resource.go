package garbage

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
