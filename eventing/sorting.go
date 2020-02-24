package eventing

// SortBy specifies how requested events are sorted
type SortBy string

const (
	DateAsc  = "dateAsc"
	DateDesc = "dateDesc"
	Gadgets  = "gadgets"
	Paper    = "paper"
	Plastic  = "plastic"
)

// IsValid checks if a provided sorting type is valid
func (s SortBy) IsValid() bool {
	switch s {
	case DateAsc, DateDesc, Gadgets, Paper, Plastic:
		return true
	}
	return false
}
