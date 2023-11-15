package types

import (
	"bytes"
	"fmt"
	"reflect"
)

type SessionTypeDefinition struct {
	SessionType SessionType
	Name        string
}

type SessionType interface {
	String() string

	// used for inner checks
	checkTypeLabels(LabelledTypesEnv) error
	isContractive(LabelledTypesEnv, map[string]bool) bool
}

// Label
type LabelType struct {
	Label string
}

func (q *LabelType) String() string {
	return q.Label
}

func NewLabelType(i string) *LabelType {
	return &LabelType{
		Label: i,
	}
}

type WIPType struct{}

func (q *WIPType) String() string {
	return "wip"
}
func NewWIPType() *WIPType {
	return &WIPType{}
}

// Unit: 1
type UnitType struct{}

func (q *UnitType) String() string {
	return "1"
}

func NewUnitType() *UnitType {
	return &UnitType{}
}

// Send: A * B
type SendType struct {
	Left  SessionType
	Right SessionType
}

func (q *SendType) String() string {
	var buffer bytes.Buffer
	// buffer.WriteString("(")
	buffer.WriteString(q.Left.String())
	buffer.WriteString(" * ")
	buffer.WriteString(q.Right.String())
	// buffer.WriteString(")")
	return buffer.String()
}

func NewSendType(left, right SessionType) *SendType {
	return &SendType{
		Left:  left,
		Right: right,
	}
}

// Receive: A -o B
type ReceiveType struct {
	Left  SessionType
	Right SessionType
}

func (q *ReceiveType) String() string {
	var buffer bytes.Buffer
	// buffer.WriteString("(")
	buffer.WriteString(q.Left.String())
	buffer.WriteString(" -o ")
	buffer.WriteString(q.Right.String())
	// buffer.WriteString(")")
	return buffer.String()
}

func NewReceiveType(left, right SessionType) *ReceiveType {
	return &ReceiveType{
		Left:  left,
		Right: right,
	}
}

// SelectLabel: +{ }
type SelectLabelType struct {
	Branches []BranchOption
}

func (q *SelectLabelType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("+{")
	buffer.WriteString(stringifyBranches(q.Branches))
	buffer.WriteString("}")
	return buffer.String()
}

func NewSelectType(branches []BranchOption) *SelectLabelType {
	return &SelectLabelType{
		Branches: branches,
	}
}

// BranchCase: & { }
type BranchCaseType struct {
	Branches []BranchOption
}

func (q *BranchCaseType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("&{")
	buffer.WriteString(stringifyBranches(q.Branches))
	buffer.WriteString("}")
	return buffer.String()
}

func NewBranchCaseType(branches []BranchOption) *BranchCaseType {
	return &BranchCaseType{
		Branches: branches,
	}
}

// Branches
type BranchOption struct {
	Label        string
	Session_type SessionType
}

func (branch *BranchOption) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(branch.Label)
	buffer.WriteString(" : ")
	buffer.WriteString(branch.Session_type.String())

	return buffer.String()
}

func stringifyBranches(branches []BranchOption) string {
	var buf bytes.Buffer

	for i, j := range branches {
		buf.WriteString(j.Label)
		buf.WriteString(" : ")
		buf.WriteString(j.Session_type.String())

		if i < len(branches)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func NewBranchOption(label string, session_type SessionType) *BranchOption {
	return &BranchOption{
		Label:        label,
		Session_type: session_type,
	}
}

// Check for equality
func EqualType(type1, type2 SessionType, labelledTypesEnv LabelledTypesEnv) bool {
	return innerEqualType(type1, type2, make(map[string]bool), labelledTypesEnv)
}

// The snapshots maps keeps a snapshot of both types in case the types are unfolded. This ensures that the types do not keep unfolding infinitely.
func innerEqualType(type1, type2 SessionType, snapshots map[string]bool, labelledTypesEnv LabelledTypesEnv) bool {
	a := reflect.TypeOf(type1)
	b := reflect.TypeOf(type2)

	f1, isLabel1 := type1.(*LabelType)
	f2, isLabel2 := type2.(*LabelType)

	// If neither is a label and neither type's match
	if a != b && !isLabel1 && !isLabel2 {
		return false
	}

	if isLabel1 || isLabel2 {
		// Compare with existing snapshots
		presentSnapshot := type1.String() + type2.String()

		_, exists := snapshots[presentSnapshot]
		if exists {
			return true
		}

		if isLabel1 && isLabel2 && f1.Label == f2.Label {
			return true
		}

		// Expand label/s
		// This fetch operation (from the map) should succeed since we already check that all labels used are defined
		if isLabel1 {
			labelledType, ok1 := labelledTypesEnv[f1.Label]
			if ok1 {
				type1 = labelledType.Type
			} else {
				return false
			}
		}

		if isLabel2 {
			labelledType, ok2 := labelledTypesEnv[f2.Label]
			if ok2 {
				type2 = labelledType.Type
			} else {
				return false
			}
		}

		// Add new snapshot
		newSnapshot := type1.String() + type2.String()
		snapshots[newSnapshot] = true

		return innerEqualType(type1, type2, snapshots, labelledTypesEnv)
	}

	// At this point, neither type1 nor type2 can be of LabelType
	if a != b {
		return false
	}

	switch interface{}(type1).(type) {
	// case *LabelType:
	// 	f1, ok1 := type1.(*LabelType)
	// 	f2, ok2 := type2.(*LabelType)
	// 	return ok1 && ok2 && f1.Label == f2.Label

	case *UnitType:
		_, ok1 := type1.(*UnitType)
		_, ok2 := type2.(*UnitType)
		return ok1 && ok2

	case *SendType:
		f1, ok1 := type1.(*SendType)
		f2, ok2 := type2.(*SendType)

		if ok1 && ok2 {
			// todo check if send type is commutative
			return innerEqualType(f1.Left, f2.Left, snapshots, labelledTypesEnv) && innerEqualType(f1.Right, f2.Right, snapshots, labelledTypesEnv)
		}

	case *ReceiveType:
		f1, ok1 := type1.(*ReceiveType)
		f2, ok2 := type2.(*ReceiveType)

		if ok1 && ok2 {
			// todo check if receive type is commutative
			return innerEqualType(f1.Left, f2.Left, snapshots, labelledTypesEnv) && innerEqualType(f1.Right, f2.Right, snapshots, labelledTypesEnv)
		}

	case *SelectLabelType:
		f1, ok1 := type1.(*SelectLabelType)
		f2, ok2 := type2.(*SelectLabelType)

		if ok1 && ok2 && len(f1.Branches) == len(f2.Branches) {
			return equalTypeBranch(f1.Branches, f2.Branches, snapshots, labelledTypesEnv)
		}

	case *BranchCaseType:
		f1, ok1 := type1.(*BranchCaseType)
		f2, ok2 := type2.(*BranchCaseType)

		if ok1 && ok2 && len(f1.Branches) == len(f2.Branches) {
			// todo check if order matters
			return equalTypeBranch(f1.Branches, f2.Branches, snapshots, labelledTypesEnv)
		}
	case *WIPType:
		return true
	}

	fmt.Printf("issue in EqualType for type %s\n", a)
	return false
}

// Compare branches in an unordered way. Here we are assuming that both branches contain unique labels
func equalTypeBranch(options1, options2 []BranchOption, snapshots map[string]bool, labelledTypesEnv LabelledTypesEnv) bool {
	if len(options1) != len(options2) {
		return false
	}

	// Match each label to the other set
	for _, b := range options1 {
		matchingBranch, foundMatchingBranch := matchBranchByLabel(options2, b.Label)
		if foundMatchingBranch {
			if !innerEqualType(b.Session_type, matchingBranch.Session_type, snapshots, labelledTypesEnv) {
				// If inner types do not match, then stop checking
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// Lookup branches by label
func matchBranchByLabel(branches []BranchOption, label string) (*BranchOption, bool) {
	for index := range branches {
		if branches[index].Label == label {
			return &branches[index], true
		}
	}

	return nil, false
}

// Takes a type and returns a (separate) clone
func CopyType(orig SessionType) SessionType {
	if orig == nil {
		return nil
	}

	switch interface{}(orig).(type) {
	case *LabelType:
		p, ok := orig.(*LabelType)
		if ok {
			return NewLabelType(p.Label)
		}
	case *UnitType:
		_, ok := orig.(*UnitType)
		if ok {
			return NewUnitType()
		}
	case *SendType:
		p, ok := orig.(*SendType)
		if ok {
			return NewSendType(CopyType(p.Left), CopyType(p.Right))
		}
	case *ReceiveType:
		p, ok := orig.(*ReceiveType)
		if ok {
			return NewReceiveType(CopyType(p.Left), CopyType(p.Right))
		}
	case *SelectLabelType:
		p, ok := orig.(*SelectLabelType)
		if ok {
			branches := make([]BranchOption, len(p.Branches))

			for i := 0; i < len(p.Branches); i++ {
				branches[i].Label = p.Branches[i].Label
				branches[i].Session_type = CopyType(p.Branches[i].Session_type)
			}

			return NewSelectType(branches)
		}

	case *BranchCaseType:
		p, ok := orig.(*BranchCaseType)
		if ok {
			branches := make([]BranchOption, len(p.Branches))

			for i := 0; i < len(p.Branches); i++ {
				branches[i].Label = p.Branches[i].Label
				branches[i].Session_type = CopyType(p.Branches[i].Session_type)
			}

			return NewBranchCaseType(branches)
		}

	case *WIPType:
		return NewWIPType()
	}

	panic("Should not happen (type)")
	// return nil
}

// The labelled types environment is constant and set once at the beginning. The information is obtained from the 'type A = ...' definitions.
// labelledTypesEnv: map of labels to their session type (wrapped in a LabelledType struct)

type LabelledTypesEnv map[string]LabelledType
type LabelledType struct {
	Name string
	Type SessionType
}

func ProduceLabelledSessionTypeEnvironment(typeDefs []SessionTypeDefinition) LabelledTypesEnv {
	labelledTypesEnv := make(LabelledTypesEnv)
	for _, j := range typeDefs {
		labelledTypesEnv[j.Name] = LabelledType{Type: j.SessionType, Name: j.Name}
	}

	return labelledTypesEnv
}

func LabelledTypedExists(labelledTypesEnv LabelledTypesEnv, key string) bool {
	_, ok := labelledTypesEnv[key]
	return ok
}

// Weaenable types allow for channels to be dropped
func IsWeakenable(sessionType SessionType) bool {
	// todo implement
	fmt.Println("todo: IsWeakenable")

	return true
}

// func IsContractable(sessionType SessionType) bool {
// 	// todo implement
// 	fmt.Println("todo: IsContractable")

// 	return true
// }

// Used to unroll a type only if needed (i.e. reached label)
func Unfold(orig SessionType, labelledTypesEnv LabelledTypesEnv) SessionType {
	if orig == nil {
		return nil
	}

	labelSessionType, labelType := orig.(*LabelType)

	if labelType {
		unfoldedSessionType, exists := labelledTypesEnv[labelSessionType.Label]

		if exists {
			// This could potentially cause an infinite loop if non-contractive types are used, however we already make sure that only contractive types are used
			return Unfold(unfoldedSessionType.Type, labelledTypesEnv)
		} else {
			return nil
		}
	} else {
		return orig
	}
}

// Lookup branches by label
func FetchSelectBranch(branches []BranchOption, label string) (SessionType, bool) {
	for _, branch := range branches {
		if branch.Label == label {
			return branch.Session_type, true
		}
	}

	return nil, false
}
