package eventing

// SortBy specifies how requested events are sorted
type SortBy string

const (
	DateAsc  SortBy = "dateAsc"
	DateDesc SortBy = "dateDesc"
	Gadgets  SortBy = "gadgets"
	Paper    SortBy = "paper"
	Plastic  SortBy = "plastic"
)

// IsValid checks if a provided sorting type is valid
func (s SortBy) IsValid() bool {
	switch s {
	case DateAsc, DateDesc, Gadgets, Paper, Plastic:
		return true
	}
	return false
}
