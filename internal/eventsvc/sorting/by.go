// Package sorting specifies how various entities can be sorted
package sorting

// By is a type indicating how things can be sorted
type By int

const (
	DateAsc By = iota
	DateDes
	Gadgets
	NameAsc
	NameDes
	Paper
	Plastic
	Unspecified
)

// IsDate checks if a provided string can be used as a sorting by date
func (b By) IsDate() bool {
	switch b {
	case DateAsc, DateDes:
		return true
	}
	return false
}

// IsName checks if a provided string can be used as a sorting by name
func (b By) IsName() bool {
	switch b {
	case NameAsc, NameDes:
		return true
	}
	return false
}

// IsResources check if a provided string can be used as a sorting by resources
func (b By) IsResources() bool {
	switch b {
	case Gadgets, Paper, Plastic:
		return true
	}
	return false
}
