package process

import (
	"bytes"
	"fmt"
	"reflect"
)

// All process' bodies have to follow the Form interface
// Form refers to AST types
type Form interface {
	String() string
	FreeNames() []Name
	Substitute(Name, Name)
	Transition(*Process, *RuntimeEnvironment)
}

// Check equality between different forms
func EqualForm(form1, form2 Form) bool {
	a := reflect.TypeOf(form1)
	b := reflect.TypeOf(form2)
	if a != b {
		return false
	}

	switch interface{}(form1).(type) {
	case *SendForm:
		f1, ok1 := form1.(*SendForm)
		f2, ok2 := form2.(*SendForm)

		if ok1 && ok2 {
			return f1.continuation_c.Equal(f2.continuation_c) && f1.payload_c.Equal(f2.payload_c) && f1.to_c.Equal(f2.to_c)
		}
	case *ReceiveForm:
		f1, ok1 := form1.(*ReceiveForm)
		f2, ok2 := form2.(*ReceiveForm)

		if ok1 && ok2 {
			return f1.payload_c.Equal(f2.payload_c) && f1.continuation_c.Equal(f2.continuation_c) && f1.from_c.Equal(f2.from_c) && EqualForm(f1.continuation_e, f2.continuation_e)
		}
	case *SelectForm:
		f1, ok1 := form1.(*SelectForm)
		f2, ok2 := form2.(*SelectForm)

		if ok1 && ok2 {
			return f1.to_c.Equal(f2.to_c) && f1.continuation_c.Equal(f2.continuation_c) && f1.label.Equal(f2.label)
		}
	case *CaseForm:
		f1, ok1 := form1.(*CaseForm)
		f2, ok2 := form2.(*CaseForm)

		if ok1 && ok2 {
			for index := range f1.branches {
				if !EqualForm(f1.branches[index], f2.branches[index]) {
					return false
				}
			}

			return f1.from_c.Equal(f2.from_c)
		}
	case *BranchForm:
		f1, ok1 := form1.(*BranchForm)
		f2, ok2 := form2.(*BranchForm)

		if ok1 && ok2 {
			return f1.label.Equal(f2.label) && f1.payload_c.Equal(f2.payload_c) && EqualForm(f1.continuation_e, f2.continuation_e)
		}
	case *CloseForm:
		f1, ok1 := form1.(*CloseForm)
		f2, ok2 := form2.(*CloseForm)

		if ok1 && ok2 {
			return f1.from_c.Equal(f2.from_c)
		}
	case *NewForm:
		f1, ok1 := form1.(*NewForm)
		f2, ok2 := form2.(*NewForm)

		if ok1 && ok2 {
			return f1.continuation_c.Equal(f2.continuation_c) && EqualForm(f1.body, f2.body) && EqualForm(f1.continuation_e, f2.continuation_e)
		}
	case *ForwardForm:
		f1, ok1 := form1.(*ForwardForm)
		f2, ok2 := form2.(*ForwardForm)

		if ok1 && ok2 {
			return f1.to_c.Equal(f2.to_c) && f1.from_c.Equal(f2.from_c)
		}
	case *SplitForm:
		f1, ok1 := form1.(*SplitForm)
		f2, ok2 := form2.(*SplitForm)

		if ok1 && ok2 {
			return f1.channel_one.Equal(f2.channel_one) && f1.channel_two.Equal(f2.channel_two) && f1.from_c.Equal(f2.from_c) && EqualForm(f1.continuation_e, f2.continuation_e)
		}
	}

	fmt.Printf("todo implement EqualForm for type %s\n", a)
	return false
}

func CopyForm(orig Form) Form {
	// origWithType := reflect.TypeOf(orig)

	switch interface{}(orig).(type) {
	case *SendForm:
		p, ok := orig.(*SendForm)
		if ok {
			return NewSend(p.to_c, p.payload_c, p.continuation_c)
		}
	case *ReceiveForm:
		p, ok := orig.(*ReceiveForm)
		if ok {
			cont := CopyForm(p.continuation_e)
			return NewReceive(p.payload_c, p.continuation_c, p.from_c, cont)
		}
	case *SelectForm:
		p, ok := orig.(*SelectForm)
		if ok {
			return NewSelect(p.to_c, p.label, p.continuation_c)
		}
	case *CaseForm:
		p, ok := orig.(*CaseForm)
		if ok {
			branches := make([]*BranchForm, len(p.branches))

			for i := 0; i < len(p.branches); i++ {
				b := CopyForm(p.branches[i]).(*BranchForm)
				branches[i] = b
			}

			return NewCase(p.from_c, branches)
		}

	case *BranchForm:
		p, ok := orig.(*BranchForm)
		if ok {
			cont := CopyForm(p.continuation_e)
			return NewBranch(p.label, p.payload_c, cont)
		}

	case *CloseForm:
		p, ok := orig.(*CloseForm)
		if ok {
			return NewClose(p.from_c)
		}
	case *NewForm:
		p, ok := orig.(*NewForm)
		if ok {
			body := CopyForm(p.body)
			cont := CopyForm(p.continuation_e)
			return NewNew(p.continuation_c, body, cont)
		}
	case *ForwardForm:
		p, ok := orig.(*ForwardForm)
		if ok {
			return NewForward(p.to_c, p.from_c)
		}
	case *SplitForm:
		p, ok := orig.(*SplitForm)
		if ok {
			cont := CopyForm(p.continuation_e)
			return NewSplit(p.channel_one, p.channel_two, p.from_c, cont)
		}
	}

	// todo panic
	fmt.Print("todo: Should not happen")
	panic("Should not happen")

	// return nil
}

///////////////////////////////
////// Different Forms ////////
///////////////////////////////

// Send: send to_c<payload_c, continuation_c>
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

func (p *SendForm) Substitute(old, new Name) {
	p.to_c.Substitute(old, new)
	p.payload_c.Substitute(old, new)
	p.continuation_c.Substitute(old, new)
}

// Free names, excluding self references
func (p *SendForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.to_c, fn)
	fn = appendIfNotSelf(p.payload_c, fn)
	fn = appendIfNotSelf(p.continuation_c, fn)
	return fn
}

// Receive: <payload_c, continuation_c> <- recv from_c; P
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

func (p *ReceiveForm) Substitute(old, new Name) {
	if old != p.payload_c && old != p.continuation_c {
		// payload_c: payload_c,
		// continuation_c: continuation_c,
		p.from_c.Substitute(old, new)
		p.continuation_e.Substitute(old, new)
	}
}

func (p *ReceiveForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.from_c, fn)
	continuation_e_excluding_bound_names := removeBoundName(p.continuation_e.FreeNames(), p.payload_c)
	continuation_e_excluding_bound_names = removeBoundName(continuation_e_excluding_bound_names, p.continuation_c)
	fn = mergeTwoNamesList(fn, continuation_e_excluding_bound_names)
	return fn
}

// Select: to_c.label<continuation_c>
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

func (p *SelectForm) Substitute(old, new Name) {
	p.to_c.Substitute(old, new)
	p.continuation_c.Substitute(old, new)
	// p.continuation_e.Substitute(old, new)
}

func (p *SelectForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.to_c, fn)
	fn = appendIfNotSelf(p.continuation_c, fn)
	return fn
}

// Branch: label<payload_c> => continuation_e
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

func (p *BranchForm) Substitute(old, new Name) {
	// payload_c is a bound variable
	if old != p.payload_c {
		p.continuation_e.Substitute(old, new)
	}
}

func StringifyBranches(branches []*BranchForm) string {
	var buf bytes.Buffer

	for i, j := range branches {
		buf.WriteString(j.String())

		if i < len(branches)-1 {
			buf.WriteString(" | ")
		}
	}
	return buf.String()
}

func (p *BranchForm) FreeNames() []Name {
	var fn []Name
	continuation_e_excluding_bound_names := removeBoundName(p.continuation_e.FreeNames(), p.payload_c)
	fn = append(fn, continuation_e_excluding_bound_names...)
	return fn
}

// Case: case from_c ( branches )
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
	buf.WriteString(" (")
	buf.WriteString(StringifyBranches(p.branches))
	buf.WriteString(")")
	return buf.String()
}

func (p *CaseForm) Substitute(old, new Name) {
	p.from_c.Substitute(old, new)

	for i := range p.branches {
		p.branches[i].Substitute(old, new)
	}
}

func (p *CaseForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.from_c, fn)
	for _, branch := range p.branches {
		fn = mergeTwoNamesList(fn, branch.FreeNames())
	}
	return fn
}

// New: continuation_c <- new (body); continuation_e
type NewForm struct {
	continuation_c Name
	body           Form
	continuation_e Form
}

func NewNew(continuation_c Name, body, continuation_e Form) *NewForm {
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

func (p *NewForm) Substitute(old, new Name) {
	// continuation_c is a bound variable
	if old != p.continuation_c {
		p.body.Substitute(old, new)
		p.continuation_e.Substitute(old, new)
	}
}

func (p *NewForm) FreeNames() []Name {
	var fn []Name
	body_excluding_bound_names := removeBoundName(p.body.FreeNames(), p.continuation_c)
	fn = mergeTwoNamesList(fn, body_excluding_bound_names)
	continuation_e_excluding_bound_names := removeBoundName(p.continuation_e.FreeNames(), p.continuation_c)
	fn = mergeTwoNamesList(fn, continuation_e_excluding_bound_names)
	return fn
}

// Close: close from_c
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

func (p *CloseForm) Substitute(old, new Name) {
	p.from_c.Substitute(old, new)
}

func (p *CloseForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.from_c, fn)
	return fn
}

// Forward: fwd to_c from_c
type ForwardForm struct {
	to_c   Name
	from_c Name
}

func NewForward(to_c, from_c Name) *ForwardForm {
	return &ForwardForm{to_c: to_c, from_c: from_c}
}

func (p *ForwardForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("fwd ")
	buf.WriteString(p.to_c.String())
	buf.WriteString(" ")
	buf.WriteString(p.from_c.String())
	return buf.String()
}

func (p *ForwardForm) Substitute(old, new Name) {
	p.to_c.Substitute(old, new)
	p.from_c.Substitute(old, new)
}

func (p *ForwardForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.to_c, fn)
	fn = appendIfNotSelf(p.from_c, fn)
	return fn
}

// Split: <payload_c, continuation_c> <- recv from_c; P
type SplitForm struct {
	channel_one    Name
	channel_two    Name
	from_c         Name
	continuation_e Form
}

func NewSplit(channel_one, channel_two, from_c Name, continuation_e Form) *SplitForm {
	return &SplitForm{
		channel_one:    channel_one,
		channel_two:    channel_two,
		from_c:         from_c,
		continuation_e: continuation_e}
}
func (p *SplitForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("<")
	buf.WriteString(p.channel_one.String())
	buf.WriteString(",")
	buf.WriteString(p.channel_two.String())
	buf.WriteString("> <- split ")
	buf.WriteString(p.from_c.String())
	buf.WriteString("; ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

func (p *SplitForm) Substitute(old, new Name) {
	if old != p.channel_one && old != p.channel_two {
		// channel_one: channel_one,
		// channel_two: channel_two,
		p.from_c.Substitute(old, new)
		p.continuation_e.Substitute(old, new)
	}
}

func (p *SplitForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.from_c, fn)
	continuation_e_excluding_bound_names := removeBoundName(p.continuation_e.FreeNames(), p.channel_one)
	continuation_e_excluding_bound_names = removeBoundName(continuation_e_excluding_bound_names, p.channel_two)
	fn = mergeTwoNamesList(fn, continuation_e_excluding_bound_names)
	return fn
}

// Utility functions

// Add name to fn list, excluding ones with IsSelf: true
func appendIfNotSelf(name Name, fn []Name) []Name {
	if !name.IsSelf {
		fn = append(fn, name)
	}

	return fn
}

// Remove bound name from list [used when computing the list of free names]
func removeBoundName(names []Name, boundName Name) (freeNames []Name) {
	for _, n := range names {
		if !n.Equal(boundName) {
			freeNames = append(freeNames, n)
		}
	}
	return
}

// Check whether name 'check' exists within a slice 'names'
func nameExists(names []Name, check Name) bool {
	for _, n := range names {
		if n.Equal(check) {
			return true
		}
	}

	return false
}

// Merges two lists of names keeping only unique values (avoiding duplicates)
func mergeTwoNamesList(names1, names2 []Name) []Name {
	for _, n := range names2 {
		if !nameExists(names1, n) {
			names1 = append(names1, n)
		}
	}
	return names1
}
