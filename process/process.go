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

// A 'Process' contains the body of the process and the channel it is providing on.
type Process struct {
	Body                Form
	Provider            Name
	FunctionDefinitions *[]FunctionDefinition
}

func (p *Process) InsertFunctionDefinitions(all *[]FunctionDefinition) {
	p.FunctionDefinitions = all
}

func (p *Process) String() string {
	var buf bytes.Buffer
	buf.WriteString("prc [")
	buf.WriteString(p.Provider.String())
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
	// Ident refers to the original name of the channel (used for pretty printing)
	Ident string
	// One a channel is initialized (i.e. Channel != nil), the Channel becomes more important than Ident
	Channel chan Message
	// Channel ID is a unique id for each channel
	// Used only for debugging, since setting the ChannelID is a slow (& synchronous) operation
	ChannelID uint64
}

func (n *Name) Initialized() bool {
	return n.Channel != nil
}

func (n *Name) String() string {
	if n.Initialized() {
		return n.Ident + "[" + strconv.FormatUint(n.ChannelID, 10) + "]"
	} else {
		return n.Ident
	}
}

func (n *Name) IsSelf() bool {
	return n.Ident == "self"
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
		}
	}
}

// Returns channel directly or the provider channel in case of channel called self
func (n *Name) GetChannel(p *Process) chan Message {
	if n.Ident == "self" {
		return p.Provider.Channel
	} else {
		return n.Channel
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
