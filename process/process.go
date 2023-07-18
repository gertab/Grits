package process

type Process interface {
	FreeNames() []Name
	FreeVars() []Name

	Calculi() string
	String() string
}

// Name is channel or value.
type Name interface {
	Ident() string
	String() string
}
