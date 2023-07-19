package process

import "bytes"

type Process struct {
	Body                Form
	Nn                  int
	FunctionDefinitions []FunctionDefinition
	// todo []FunctionDefinition should be []*FunctionDefinition
}

func (p Process) InsertFunctionDefinitions(all []FunctionDefinition) {
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

type Form interface {
	String() string
}

type SendForm struct {
	to_c           Name
	payload_c      Name
	continuation_c Name
}

func NewSend(to_c, payload_c, continuation_c Name) *SendForm {
	return &SendForm{
		to_c:           to_c,
		payload_c:      payload_c,
		continuation_c: continuation_c}
}
func (p *SendForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("send ")
	buf.WriteString(p.to_c.String())
	buf.WriteString("<")
	buf.WriteString(p.payload_c.String())
	buf.WriteString(",")
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(">")
	return buf.String()
}

type ReceiveForm struct {
	payload_c      Name
	continuation_c Name
	from_c         Name
	continuation_e Form
}

func NewReceive(payload_c, continuation_c, from_c Name, continuation_e Form) *ReceiveForm {
	return &ReceiveForm{
		payload_c:      payload_c,
		continuation_c: continuation_c,
		from_c:         from_c,
		continuation_e: continuation_e}
}
func (p *ReceiveForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("<")
	buf.WriteString(p.payload_c.String())
	buf.WriteString(",")
	buf.WriteString(p.continuation_c.String())
	buf.WriteString("> <- recv ")
	buf.WriteString(p.from_c.String())
	buf.WriteString("; ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

// Par is parallel composition of P and Q.
type Par struct {
	Procs []Process
}

// NewPar creates a new parallel composition.
func NewPar(P, Q Process) *Par { return &Par{Procs: []Process{P, Q}} }

// FreeNames of Par is the free names of composed processes.
func (p *Par) FreeNames() []Name {
	var fn []Name
	return fn
}

// FreeVars of Par is the free names of composed processes.
func (p *Par) FreeVars() []Name {
	var fv []Name
	return fv
}

func (p *Par) String() string {
	// var buf bytes.Buffer
	// buf.WriteString("par[ ")
	// for i, proc := range p.Procs {
	// 	if i != 0 {
	// 		buf.WriteString(" | ")
	// 	}
	// 	buf.WriteString(proc.String())
	// }
	// buf.WriteString(" ]")
	return "abc"
}
