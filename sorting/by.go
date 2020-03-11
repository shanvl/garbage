// Package sorting specifies how various entities can be sorted
package sorting

// By is a type indicating how things can be sorted
type By string

const (
	DateAsc By = "dateAsc"
	DateDes By = "dateDes"
	Gadgets By = "gadgets"
	NameAsc By = "nameAsc"
	NameDes By = "nameDes"
	Paper   By = "paper"
	Plastic By = "plastic"
)

// IsValid checks if a provided string can be used for sorting
func (s By) IsValid() bool {
	switch s {
	case DateAsc, DateDes, Gadgets, Paper, Plastic, NameAsc, NameDes:
		return true
	}
	return false
}

// IsForEventPupils checks if a provided string can be used to sorting an event's pupils
func (s By) IsForEventPupils() bool {
	switch s {
	case NameAsc, NameDes, Gadgets, Paper, Plastic:
		return true
	}
	return false
}

// IsForEventPupils checks if a provided string can be used to sorting an event's classes
func (s By) IsForEventClasses() bool {
	switch s {
	case NameAsc, NameDes, Gadgets, Paper, Plastic:
		return true
	}
	return false
}
