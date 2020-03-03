package eventing

// SortBy specifies how requested events are sorted
type SortBy string

const (
	DateAsc SortBy = "dateAsc"
	DateDes SortBy = "dateDes"
	Gadgets SortBy = "gadgets"
	NameAsc SortBy = "nameAsc"
	NameDes SortBy = "nameDes"
	Paper   SortBy = "paper"
	Plastic SortBy = "plastic"
)

// IsValid checks if a provided sorting type is valid
func (s SortBy) IsValid() bool {
	switch s {
	case DateAsc, DateDes, Gadgets, Paper, Plastic, NameAsc, NameDes:
		return true
	}
	return false
}
