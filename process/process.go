package process

import (
	"bytes"
	"phi/types"
	"strconv"
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
	// Session type which the channel must follow
	Type types.SessionType
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

	if printTypes && n.Type != nil {
		m = m + "|" + n.Type.String() + "|"
	}

	if n.Initialized() {
		return m + "[" + strconv.FormatUint(n.ChannelID, 10) + "]"
	} else {
		return m
	}
}

// Check whether n is in names
func (n *Name) ContainedIn(names []Name) bool {
	for _, j := range names {
		if n.Equal(j) {
			return true
		}
	}

	return false
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

// Compare two list of names and check for membership inclusion in both lists
func AreNamesEqual(first, second []Name) bool {
	if len(first) != len(second) {
		return false
	}
	exists := make(map[string]bool)
	for _, name := range first {
		exists[name.Ident] = true
	}
	for _, name := range second {
		if !exists[name.Ident] {
			return false
		}
	}
	return true
}

// Compare two list of names and check returns the lists of unique names found in each list
// E.g. [a, b, c] & [c, d] => ([a, b], [d])
func NamesNotCommon(first, second []Name) ([]Name, []Name) {
	existsInFirst := make(map[string]bool)
	for _, name := range first {
		existsInFirst[name.Ident] = true
	}

	existsInSecond := make(map[string]bool)
	for _, name := range second {
		existsInSecond[name.Ident] = true
	}

	var uniqueFirst []Name
	var uniqueSecond []Name

	for _, name := range first {
		_, common := existsInSecond[name.Ident]
		if !common {
			uniqueFirst = append(uniqueFirst, name)
		}
	}

	for _, name := range second {
		_, common := existsInFirst[name.Ident]
		if !common {
			uniqueSecond = append(uniqueSecond, name)
		}
	}

	return uniqueFirst, uniqueSecond
}

// This subtracts the names from the second list from the first list.
// E.g. [a, b, c] - [c, d] = [a, b]
func NamesInFirstListOnly(first, second []Name) []Name {
	existsInSecond := make(map[string]bool)
	for _, name := range second {
		existsInSecond[name.Ident] = true
	}

	var uniqueFirst []Name

	for _, name := range first {
		_, common := existsInSecond[name.Ident]
		if !common {
			uniqueFirst = append(uniqueFirst, name)
		}
	}

	return uniqueFirst
}

// Checks if a list of names contains unique names only
func AllNamesUnique(list []Name) bool {
	exists := make(map[string]bool)
	for _, name := range list {
		if exists[name.Ident] {
			return false
		}
		exists[name.Ident] = true
	}

	return true
}

// Returns the duplicates in a list of names
func DuplicateNames(names []Name) []Name {
	var duplicates []Name

	exists := make(map[string]bool)
	for _, name := range names {
		if exists[name.Ident] {
			duplicates = append(duplicates, name)
		}
		exists[name.Ident] = true
	}

	return duplicates
}

// Sets a common type to all names
func SetTypesToNames(names []Name, t types.SessionType) {
	for _, name := range names {
		name.Type = types.CopyType(t)
	}
}

func (name1 *Name) Equal(name2 Name) bool {
	if name1.Initialized() && name2.Initialized() {
		// If the channel is initialized, then only compare the actual channel reference
		return name1.Channel == name2.Channel
	}
	return name1.Ident == name2.Ident && name1.Initialized() == name2.Initialized()
}

func (n *Name) Substitute(old, new Name) {
	if n.Initialized() && n.Channel == old.Channel {
		// If a channel is initialized, then compare using the channel value
		n.Channel = new.Channel
		n.IsSelf = new.IsSelf
		n.ControlChannel = new.ControlChannel
		n.Polarity = new.Polarity
		// n.Type = new.Type [type should remain the same, as set by the typechecker]
		if new.Ident != "" {
			// not sure if this works
			n.Ident = new.Ident
			n.ChannelID = new.ChannelID
		}
	} else {
		if n.Ident == old.Ident {
			n.Ident = new.Ident
			n.Channel = new.Channel
			n.ChannelID = new.ChannelID
			n.IsSelf = new.IsSelf
			n.Polarity = new.Polarity
			n.ControlChannel = new.ControlChannel
			// n.Type = new.Type
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
