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
