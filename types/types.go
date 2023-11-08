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
	checkLabelledTypes(typeDefNames map[string]bool) error
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

// branches
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

// Check equality between different forms
func EqualType(type1, type2 SessionType) bool {
	a := reflect.TypeOf(type1)
	b := reflect.TypeOf(type2)
	if a != b {
		return false
	}

	switch interface{}(type1).(type) {
	case *LabelType:
		_, ok1 := type1.(*LabelType)
		_, ok2 := type2.(*LabelType)
		return ok1 && ok2

	case *UnitType:
		_, ok1 := type1.(*UnitType)
		_, ok2 := type2.(*UnitType)
		return ok1 && ok2

	case *SendType:
		f1, ok1 := type1.(*SendType)
		f2, ok2 := type2.(*SendType)

		if ok1 && ok2 {
			// todo check if send type is commutative
			return EqualType(f1.Left, f2.Left) && EqualType(f1.Right, f2.Right)
		}

	case *ReceiveType:
		f1, ok1 := type1.(*ReceiveType)
		f2, ok2 := type2.(*ReceiveType)

		if ok1 && ok2 {
			// todo check if receive type is commutative
			return EqualType(f1.Left, f2.Left) && EqualType(f1.Right, f2.Right)
		}

	case *SelectLabelType:
		f1, ok1 := type1.(*SelectLabelType)
		f2, ok2 := type2.(*SelectLabelType)

		if ok1 && ok2 {
			for index := range f1.Branches {
				if !equalTypeBranch(f1.Branches[index], f2.Branches[index]) {
					return false
				}
			}

			return true
		}

	case *BranchCaseType:
		f1, ok1 := type1.(*BranchCaseType)
		f2, ok2 := type2.(*BranchCaseType)

		if ok1 && ok2 {
			// todo check if order matters
			for index := range f1.Branches {
				if !equalTypeBranch(f1.Branches[index], f2.Branches[index]) {
					return false
				}
			}

			return true
		}
	}

	fmt.Printf("todo implement EqualType for type %s\n", a)
	return false
}

func equalTypeBranch(option1, option2 BranchOption) bool {

	if option1.Label != option2.Label {
		return false
	}

	return EqualType(option1.Session_type, option2.Session_type)
}

// Takes a type and returns a (separate) clone
func CopyType(orig SessionType) SessionType {
	// origWithType := reflect.TypeOf(orig)

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

			return NewSelectType(branches)
		}
	}

	panic("Should not happen")
	// return nil
}
