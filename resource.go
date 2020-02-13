package garbage

type Resource int

const (
	Gadgets Resource = iota
	Paper
	Plastic
)

func (r Resource) String() string {
	switch r {
	case Gadgets:
		return "gadgets"
	case Paper:
		return "paper"
	case Plastic:
		return "plastic"
	}
	return ""
}
