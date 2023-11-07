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
}

// Label
type LabelType struct {
	label string
}

func (q *LabelType) String() string {
	return q.label
}

func NewLabelType(i string) *LabelType {
	return &LabelType{
		label: i,
	}
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
	left  SessionType
	right SessionType
}

func (q *SendType) String() string {
	var buffer bytes.Buffer
	// buffer.WriteString("(")
	buffer.WriteString(q.left.String())
	buffer.WriteString(" * ")
	buffer.WriteString(q.right.String())
	// buffer.WriteString(")")
	return buffer.String()
}

func NewSendType(left, right SessionType) *SendType {
	return &SendType{
		left:  left,
		right: right,
	}
}

// Receive: A -o B
type ReceiveType struct {
	left  SessionType
	right SessionType
}

func (q *ReceiveType) String() string {
	var buffer bytes.Buffer
	// buffer.WriteString("(")
	buffer.WriteString(q.left.String())
	buffer.WriteString(" -o ")
	buffer.WriteString(q.right.String())
	// buffer.WriteString(")")
	return buffer.String()
}

func NewReceiveType(left, right SessionType) *ReceiveType {
	return &ReceiveType{
		left:  left,
		right: right,
	}
}

// SelectLabel: +{ }
type SelectLabelType struct {
	branches []BranchOption
}

func (q *SelectLabelType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("+{")
	buffer.WriteString(stringifyBranches(q.branches))
	buffer.WriteString("}")
	return buffer.String()
}

func NewSelectType(branches []BranchOption) *SelectLabelType {
	return &SelectLabelType{
		branches: branches,
	}
}

// BranchCase: & { }
type BranchCaseType struct {
	branches []BranchOption
}

func (q *BranchCaseType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("&{")
	buffer.WriteString(stringifyBranches(q.branches))
	buffer.WriteString("}")
	return buffer.String()
}

func NewBranchCaseType(branches []BranchOption) *BranchCaseType {
	return &BranchCaseType{
		branches: branches,
	}
}

// branches
type BranchOption struct {
	label        string
	session_type SessionType
}

func (branch *BranchOption) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(branch.label)
	buffer.WriteString(" : ")
	buffer.WriteString(branch.session_type.String())

	return buffer.String()
}

func stringifyBranches(branches []BranchOption) string {
	var buf bytes.Buffer

	for i, j := range branches {
		buf.WriteString(j.label)
		buf.WriteString(" : ")
		buf.WriteString(j.session_type.String())

		if i < len(branches)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func NewBranchOption(label string, session_type SessionType) *BranchOption {
	return &BranchOption{
		label:        label,
		session_type: session_type,
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
			return EqualType(f1.left, f2.left) && EqualType(f1.right, f2.right)
		}

	case *ReceiveType:
		f1, ok1 := type1.(*ReceiveType)
		f2, ok2 := type2.(*ReceiveType)

		if ok1 && ok2 {
			// todo check if receive type is commutative
			return EqualType(f1.left, f2.left) && EqualType(f1.right, f2.right)
		}

	case *SelectLabelType:
		f1, ok1 := type1.(*SelectLabelType)
		f2, ok2 := type2.(*SelectLabelType)

		if ok1 && ok2 {
			for index := range f1.branches {
				if !equalTypeBranch(f1.branches[index], f2.branches[index]) {
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
			for index := range f1.branches {
				if !equalTypeBranch(f1.branches[index], f2.branches[index]) {
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

	if option1.label != option2.label {
		return false
	}

	return EqualType(option1.session_type, option2.session_type)
}
