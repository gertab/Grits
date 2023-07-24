package process

import "bytes"

// All process' bodies have to follow the Form interface
type Form interface {
	String() string
	// FreeNames()
	// Substitute()
}

type SendForm struct {
	to_c           Name
	payload_c      Name
	continuation_c Name
}

func NewSend(to_c, payload_c, continuation_c Name) *SendForm {
	return &SendForm{
		to_c:           to_c,
		payload_c:      payload_c,
		continuation_c: continuation_c}
}

func (p *SendForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("send ")
	buf.WriteString(p.to_c.String())
	buf.WriteString("<")
	buf.WriteString(p.payload_c.String())
	buf.WriteString(",")
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(">")
	return buf.String()
}

type ReceiveForm struct {
	payload_c      Name
	continuation_c Name
	from_c         Name
	continuation_e Form
}

func NewReceive(payload_c, continuation_c, from_c Name, continuation_e Form) *ReceiveForm {
	return &ReceiveForm{
		payload_c:      payload_c,
		continuation_c: continuation_c,
		from_c:         from_c,
		continuation_e: continuation_e}
}
func (p *ReceiveForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("<")
	buf.WriteString(p.payload_c.String())
	buf.WriteString(",")
	buf.WriteString(p.continuation_c.String())
	buf.WriteString("> <- recv ")
	buf.WriteString(p.from_c.String())
	buf.WriteString("; ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

type SelectForm struct {
	to_c           Name
	label          Label
	continuation_c Name
}

func NewSelect(to_c Name, label Label, continuation_c Name) *SelectForm {
	return &SelectForm{
		to_c:           to_c,
		label:          label,
		continuation_c: continuation_c}
}
func (p *SelectForm) String() string {
	var buf bytes.Buffer
	buf.WriteString(p.to_c.String())
	buf.WriteString(".")
	buf.WriteString(p.label.String())
	buf.WriteString("<")
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(">")
	return buf.String()
}

type BranchForm struct {
	label          Label
	payload_c      Name
	continuation_e Form
}

func NewBranch(label Label, payload_c Name, continuation_e Form) *BranchForm {
	return &BranchForm{
		label:          label,
		payload_c:      payload_c,
		continuation_e: continuation_e}
}
func (p *BranchForm) String() string {
	var buf bytes.Buffer
	buf.WriteString(p.label.String())
	buf.WriteString("<")
	buf.WriteString(p.payload_c.String())
	buf.WriteString("> => ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

func StringifyBranches(branches []*BranchForm) string {
	var buf bytes.Buffer
	buf.WriteString("   ")

	for i, j := range branches {
		buf.WriteString(j.String())

		if i < len(branches)-1 {
			buf.WriteString("\n | ")
		}
	}
	return buf.String()
}

type CaseForm struct {
	from_c   Name
	branches []*BranchForm
}

func NewCase(from_c Name, branches []*BranchForm) *CaseForm {
	return &CaseForm{
		from_c:   from_c,
		branches: branches}
}
func (p *CaseForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("case ")
	buf.WriteString(p.from_c.String())
	buf.WriteString(" ( \n")
	buf.WriteString(StringifyBranches(p.branches))
	buf.WriteString("\n)")
	return buf.String()
}

type NewForm struct {
	continuation_c Name
	body           Form
	continuation_e Form
}

func NewNew(continuation_c Name, body Form, continuation_e Form) *NewForm {
	return &NewForm{
		continuation_c: continuation_c,
		body:           body,
		continuation_e: continuation_e}
}
func (p *NewForm) String() string {
	var buf bytes.Buffer
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(" <- new (")
	buf.WriteString(p.body.String())
	buf.WriteString("); ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

type CloseForm struct {
	from_c Name
}

func NewClose(from_c Name) *CloseForm {
	return &CloseForm{from_c: from_c}
}

func (p *CloseForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("close ")
	buf.WriteString(p.from_c.String())
	return buf.String()
}
