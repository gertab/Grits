package process

import (
	"bytes"
	"phi/types"
)

const printTypes = false

// A 'Process' contains the body of the process and the channel it is providing on.
type Process struct {
	Body      Form
	Providers []Name
	Shape     Shape
	Type      types.SessionType
}

func NewProcess(body Form, providers []Name, session_type types.SessionType, shape Shape) *Process {
	return &Process{
		Body:      body,
		Providers: providers,
		Shape:     shape,
		Type:      session_type,
	}
}

// Returns the stringified process structure, e.g. prc[pid1]
func (p *Process) OutlineString() string {
	var buf bytes.Buffer
	buf.WriteString(shapeMap[p.Shape])
	buf.WriteString("[")
	buf.WriteString(NamesToString(p.Providers))
	buf.WriteString("]")
	return buf.String()
}

// Returns the full stringified process, e.g. prc[pid1]: send a<b,c,>
func (p *Process) String() string {
	var buf bytes.Buffer
	buf.WriteString(p.OutlineString())
	buf.WriteString(": ")
	buf.WriteString(p.Body.String())
	return buf.String()
}

type FunctionDefinition struct {
	Body                 Form
	FunctionName         string
	Parameters           []Name
	Type                 types.SessionType // Session type for 'self'
	ExplicitProvider     Name              // Optional name to be used instead of 'self'
	UsesExplicitProvider bool              // ExplicitProvider set or not
}

func (function *FunctionDefinition) Arity() int {
	return len(function.Parameters)
}

func (function *FunctionDefinition) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(function.FunctionName)
	buffer.WriteString("(")
	buffer.WriteString(NamesToString(function.Parameters))
	buffer.WriteString(")")
	return buffer.String()
}

func GetFunctionByNameArity(functions []FunctionDefinition, name string, arity int) *FunctionDefinition {
	for _, f := range functions {
		if f.FunctionName == name && f.Arity() == arity {
			return &f
		} else if f.FunctionName == name && f.Arity() == arity-1 {
			// In case self is passed as a parameter, then modify the requested function arity
			return &f
		}
	}

	return nil
}

type Shape int

const (
	LINEAR Shape = iota
	SHARED
)

var shapeMap = map[Shape]string{
	LINEAR: "prc",
	SHARED: "sprc",
}
