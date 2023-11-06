package process

import (
	"bytes"
	"strconv"
)

// Stores the states of each running process.
// May also reference a controller & a monitor to observe the running state.
type ProcessConfiguration struct {
	Processes []Process
	// ref to controller/monitor
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

// A 'Process' contains the body of the process and the channel it is providing on.
type Process struct {
	Body                Form
	Providers           []Name
	Shape               Shape
	FunctionDefinitions *[]FunctionDefinition
}

func NewProcess(body Form, providers []Name, shape Shape, functionDefinitions *[]FunctionDefinition) *Process {
	return &Process{
		Body:                body,
		Providers:           providers,
		Shape:               shape,
		FunctionDefinitions: functionDefinitions,
	}
}

// func (p *Process) InsertFunctionDefinitions(all *[]FunctionDefinition) {
// 	p.FunctionDefinitions = all
// }

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

// Name is channel or value.
type Name struct {
	// Ident refers to the original name of the channel (used for pretty printing)
	Ident string
	// If IsSelf, then this channel should reference the the channel from the provider
	IsSelf bool
	// One a channel is initialized (i.e. Channel != nil), the Channel becomes more important than Ident
	Channel chan Message
	// Polarity used in NORMAL_[A]SYNC
	// todo: currently not being set/used, but might be useful to use it
	Polarity Polarity
	// Used for control commands (i.e. fwd, split, ...) in NON_POLARIZED_SYNC
	ControlChannel chan ControlMessage
	// Channel ID is a unique id for each channel
	// Used only for debugging, since setting the ChannelID is a slow (& synchronous) operation
	ChannelID uint64
}

func (n *Name) Initialized() bool {
	return n.Channel != nil
}

func (n *Name) String() string {
	m := ""

	if n.IsSelf {
		m = "self"
	} else {
		if n.Ident == "" {
			// Set anonymous process names to *
			m = "*"
		} else {
			m = n.Ident
		}
	}

	if n.Initialized() {
		return m + "[" + strconv.FormatUint(n.ChannelID, 10) + "]"
	} else {
		return m
	}
}

func NamesToString(names []Name) string {
	var buf bytes.Buffer

	for i, n := range names {
		buf.WriteString(n.String())

		if i < len(names)-1 {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

func (name1 *Name) Equal(name2 Name) bool {
	if name1.Initialized() && name2.Initialized() {
		// If the channel is initialized, then only compare the actual channel reference
		return name1.Channel == name2.Channel
	}
	return name1.String() == name2.String() && name1.Initialized() == name2.Initialized()
}

func (n *Name) Substitute(old, new Name) {
	if n.Initialized() && n.Channel == old.Channel {
		// If a channel is initialized, then compare using the channel value
		n.Channel = new.Channel
		n.IsSelf = new.IsSelf
		n.ControlChannel = new.ControlChannel
		n.Polarity = new.Polarity
		if new.Ident != "" {
			// not sure if this works
			n.Ident = new.Ident
			n.ChannelID = new.ChannelID
		}
	} else {
		if n.Ident == old.Ident {
			// fmt.Print("\t")
			// fmt.Print(n.String())
			// fmt.Println("match")
			n.Ident = new.Ident
			n.Channel = new.Channel
			n.ChannelID = new.ChannelID
			n.IsSelf = new.IsSelf
			n.Polarity = new.Polarity
			n.ControlChannel = new.ControlChannel
		}
	}
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

type FunctionDefinition struct {
	Body         Form
	FunctionName string
	Parameters   []Name
}

func (function *FunctionDefinition) Arity() int {
	return len(function.Parameters)
}

func GetFunctionByNameArity(functions []FunctionDefinition, name string, arity int) *FunctionDefinition {
	for _, f := range functions {
		if f.FunctionName == name && f.Arity() == arity {
			return &f
		}
	}

	return nil
}
