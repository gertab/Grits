package process

type Process struct {
	Body                Form
	Nn                  int // placeholder; to remove
	FunctionDefinitions *[]FunctionDefinition
}

func (p *Process) InsertFunctionDefinitions(all *[]FunctionDefinition) {
	p.FunctionDefinitions = all
}

type FunctionDefinition struct {
	Body Form
	Name string
}

// Name is channel or value.
type Name struct {
	Ident string
	// String() string
	// channel
}

func (p *Name) String() string {
	return p.Ident
}

// Name is channel or value.
type Label struct {
	L string
}

func (p *Label) String() string {
	return p.L
}
