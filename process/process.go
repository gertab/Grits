package process

import "bytes"

// Stores the states of each running process.
// May also reference a controller & a monitor to observe the running state.
type ProcessConfiguration struct {
	Processes []Process
	// ref to controller/monitor
}

// A 'Process' contains the body of the process and the channel it is providing on.
type Process struct {
	Body                Form
	Channel             Name
	FunctionDefinitions *[]FunctionDefinition
}

func (p *Process) InsertFunctionDefinitions(all *[]FunctionDefinition) {
	p.FunctionDefinitions = all
}

func (p *Process) String() string {
	var buf bytes.Buffer
	buf.WriteString("prc [")
	buf.WriteString(p.Channel.String())
	buf.WriteString("]: ")
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

func (name1 *Name) Equal(name2 Name) bool {
	return name1.String() == name2.String()
}

// Name is channel or value.
type Label struct {
	L string
}

func (p *Label) String() string {
	return p.L
}

func (label1 *Label) Equal(label2 Label) bool {
	return label1.String() == label2.String()
}
