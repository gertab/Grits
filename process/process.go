package process

import "bytes"

type Process struct {
	Body                Form
	Channel             Name
	Nn                  int // placeholder; to remove
	FunctionDefinitions *[]FunctionDefinition
}

func (p *Process) InsertFunctionDefinitions(all *[]FunctionDefinition) {
	p.FunctionDefinitions = all
}

func (p *Process) String() string {
	var buf bytes.Buffer
	buf.WriteString("prc [")
	buf.WriteString(p.Channel.String())
	buf.WriteString("]:")
	buf.WriteString(p.Body.String())

	return buf.String()
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
