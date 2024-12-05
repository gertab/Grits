package process

import (
	"bytes"
	"grits/types"
	"strconv"
)

const (
	showPolarities    = false
	showChannelNumber = false
	printTypes        = false
)

// Name is channel or value.
type Name struct {
	// Ident refers to the original name of the channel (used for pretty printing)
	Ident string
	// Session type which the channel must follow
	Type types.SessionType
	// If IsSelf, then this channel should reference the the channel from the provider
	IsSelf bool
	// ExplicitPolarity
	ExplicitPolarity *types.Polarity
	// One a channel is initialized (i.e. Channel != nil), the Channel becomes more important than Ident
	Channel chan Message
	// Polarity used in NORMAL_[A]SYNC
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
	var buffer bytes.Buffer

	if n.Ident != "" {
		buffer.WriteString(n.Ident)
		if n.IsSelf {
			buffer.WriteString("|self")
		}
	} else if n.IsSelf {
		buffer.WriteString("self")
	} else {
		// Set anonymous process names to *
		buffer.WriteString("*")
	}

	if printTypes && n.Type != nil {
		buffer.WriteString("|")
		buffer.WriteString(n.Type.String())
		buffer.WriteString("|")

	}

	if showChannelNumber && n.Initialized() {
		buffer.WriteString("[")
		buffer.WriteString(strconv.FormatUint(n.ChannelID, 10))
		buffer.WriteString("]")
	}

	if showPolarities {
		buffer.WriteString("{")
		if n.Type != nil {
			_, labelType := n.Type.(*types.LabelType)
			if !labelType {
				buffer.WriteString(types.PolarityMap[n.Type.Polarity()])
			}
		} else if n.ExplicitPolarity != nil {
			buffer.WriteString(types.PolarityMap[*n.ExplicitPolarity])
		} else {
			buffer.WriteString("")
		}
		buffer.WriteString("}")
	}

	return buffer.String()
}

func NewSelf(Ident string) Name {
	return Name{
		Ident:  Ident, // technically, this should now be "" (except when using an explicit provider name)
		IsSelf: true,
		// todo add polarity (?)
	}
}

func (n *Name) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	// Fetch the polarity either directly from the type, or the user inputted polarity
	if fromTypes {
		// unfold if required
		n.Type = types.UnfoldIfNeeded(n.Type, globalEnvironment.Types)
		return n.Type.Polarity()
	} else if n.ExplicitPolarity != nil {
		return *n.ExplicitPolarity
	}

	return types.UNKNOWN
}

// Checks whether the explicit polarity (set by the user), matches the (more precise) polarity inferred from the type
func (n *Name) ExplicitPolarityValid() bool {
	if n.ExplicitPolarity != nil && n.Type != nil {
		return *n.ExplicitPolarity == n.Type.Polarity()
	}

	return true
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

// Create fresh name copy
func (n *Name) Copy() *Name {
	pol := n.ExplicitPolarity
	var new_pol types.Polarity
	if pol != nil {
		new_pol = *pol
	}

	return &Name{
		Ident:            n.Ident,
		Channel:          n.Channel,
		ChannelID:        n.ChannelID,
		IsSelf:           n.IsSelf,
		ControlChannel:   n.ControlChannel,
		Type:             types.CopyType(n.Type),
		ExplicitPolarity: &new_pol,
	}
}

// Compare two names for equality: equality is check by the initialized unique channel or by the name
func (name1 *Name) Equal(name2 Name) bool {
	if name1.Initialized() && name2.Initialized() {
		// If the channel is initialized, then only compare the actual channel reference
		return name1.Channel == name2.Channel
	}
	return name1.Ident == name2.Ident && name1.Initialized() == name2.Initialized()
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

// If the checked name matches with the old name, then replace it with the new name
func (n *Name) Substitute(old, new Name) {
	if n.Initialized() && n.Channel == old.Channel {
		// If a channel is initialized, then compare using the channel value
		n.Channel = new.Channel
		n.IsSelf = new.IsSelf
		n.ControlChannel = new.ControlChannel
		// n.Type = new.Type [type should remain the same, as set by the typechecker]
		// n.ExplicitPolarity
		if new.Ident != "" {
			// not sure if this works
			n.Ident = new.Ident
			n.ChannelID = new.ChannelID
		}
	} else if !n.Initialized() && n.Ident == old.Ident {
		n.Ident = new.Ident
		n.Channel = new.Channel
		n.ChannelID = new.ChannelID
		n.IsSelf = new.IsSelf
		n.ControlChannel = new.ControlChannel
		// n.Type = new.Type

	}
}
