package users

type Sorting int

const (
	NameAsc Sorting = iota
	NameDes
	EmailAsc
	EmailDesc
	Unspecified
)
