package types

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

// Send: *{ }
type SendType struct {
	left  SessionType
	right SessionType
}

func (q *SendType) String() string {
	return "send *{}"
}

func NewSendType(left, right SessionType) *SendType {
	return &SendType{
		left:  left,
		right: right,
	}
}

// Receive: -o { }
type ReceiveType struct {
	left  SessionType
	right SessionType
}

func (q *ReceiveType) String() string {
	return "receive -o"
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
	return "SelectCase"
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
	return "branchCase"
}

func NewBranchCaseType(branches []BranchOption) *BranchCaseType {
	return &BranchCaseType{
		branches: branches,
	}
}

// BranchCase: & { branches }
type BranchOption struct {
	label        string
	session_type SessionType
}

func (q *BranchOption) String() string {
	return "branchCase"
}

func NewBranchOption(label string, session_type SessionType) *BranchOption {
	return &BranchOption{
		label:        label,
		session_type: session_type,
	}
}
