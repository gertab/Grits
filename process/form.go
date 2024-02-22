package process

import (
	"bytes"
	"fmt"
	"grits/types"
	"reflect"
)

// All process' bodies have to follow the Form interface
// Form refers to AST types
type Form interface {
	String() string
	StringShort() string
	Polarity(bool, *GlobalEnvironment) types.Polarity
	FreeNames() []Name
	Substitute(Name, Name)

	// Transition functions are used during evaluation
	Transition(*Process, *RuntimeEnvironment)
	TransitionNP(*Process, *RuntimeEnvironment)
	// Main typing judgement
	typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, a types.LabelledTypesEnv, sigma FunctionTypesEnv, globalEnv *GlobalEnvironment) *TypeError
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

func (p *SendForm) StringShort() string {
	return p.String()
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

// Polarity works by performing an ast traversal until a 'self' is reached
// Polarity of a send process can be inferred directly from itself
func (p *SendForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if p.to_c.IsSelf {
		return types.POSITIVE
	}

	// Type from continuation channel
	return p.continuation_c.Polarity(fromTypes, globalEnvironment)

	// if fromTypes {
	// 	// Get polarity from the type
	// 	return p.continuation_c.Type.Polarity()
	// }

	// return types.UNKNOWN
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

func (p *ReceiveForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString("<")
	buf.WriteString(p.payload_c.String())
	buf.WriteString(",")
	buf.WriteString(p.continuation_c.String())
	buf.WriteString("> <- recv ")
	buf.WriteString(p.from_c.String())
	buf.WriteString("; ...")
	return buf.String()
}

func (p *ReceiveForm) Substitute(old, new Name) {
	p.from_c.Substitute(old, new)

	if !p.payload_c.Equal(old) && !p.continuation_c.Equal(old) {
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

// Polarity of a receive process can be inferred directly from itself
func (p *ReceiveForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if p.from_c.IsSelf {
		return types.NEGATIVE
	}

	// Fetch polarity from the continuation
	return p.continuation_e.Polarity(fromTypes, globalEnvironment)
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

func (p *SelectForm) StringShort() string {
	return p.String()
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

func (p *SelectForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if p.to_c.IsSelf {
		return types.POSITIVE
	}

	return p.continuation_c.Polarity(fromTypes, globalEnvironment)

	// if fromTypes {
	// 	// Get polarity from the type
	// 	return p.continuation_c.Type.Polarity()
	// }

	// return types.UNKNOWN
}

// Branch: label<payload_c> => continuation_e
// The Branch Form is technically not a top-level form, but used as a sub-form by the case construct
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

func (p *BranchForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString(p.label.String())
	buf.WriteString("<")
	buf.WriteString(p.payload_c.String())
	buf.WriteString("> => ...")
	return buf.String()
}

func (p *BranchForm) Substitute(old, new Name) {
	// payload_c is a bound variable
	if !p.payload_c.Equal(old) {
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

func StringifyBranchesShort(branches []*BranchForm) string {
	var buf bytes.Buffer

	for i, j := range branches {
		buf.WriteString(j.StringShort())

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

// This refers to the polarity from the continuations of the branch
// (because if there is a case on self, then polarity analysis stops before checking each branch)
func (p *BranchForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	return p.continuation_e.Polarity(fromTypes, globalEnvironment)
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

func (p *CaseForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString("case ")
	buf.WriteString(p.from_c.String())
	buf.WriteString(" (")
	buf.WriteString(StringifyBranchesShort(p.branches))
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

func (p *CaseForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if p.from_c.IsSelf {
		return types.NEGATIVE
	}

	// Get polarity from the continuations
	polarity := types.UNKNOWN

	for i := range p.branches {
		branchPolarity := p.branches[i].Polarity(fromTypes, globalEnvironment)

		if branchPolarity != types.UNKNOWN {
			if polarity == types.UNKNOWN {
				// set common polarity to the current branch's polarity
				polarity = branchPolarity
			} else if polarity != branchPolarity {
				// mismatching branch polarities
				return types.UNKNOWN
			}
		}
	}

	return polarity
}

// New: continuation_c <- new (body); continuation_e
type NewForm struct {
	continuation_c   Name
	body             Form
	continuation_e   Form
	derivedFromMacro bool
}

func NewNew(continuation_c Name, body, continuation_e Form) *NewForm {
	return &NewForm{
		continuation_c:   continuation_c,
		body:             body,
		continuation_e:   continuation_e,
		derivedFromMacro: false,
	}
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

func (p *NewForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(" <- new ...; ...")
	return buf.String()
}

func (p *NewForm) Substitute(old, new Name) {

	// continuation_c is a bound variable
	if !p.continuation_c.Equal(old) {
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

func (p *NewForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	return p.continuation_e.Polarity(fromTypes, globalEnvironment)
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

func (p *CloseForm) StringShort() string {
	return p.String()
}

func (p *CloseForm) Substitute(old, new Name) {
	p.from_c.Substitute(old, new)
}

func (p *CloseForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.from_c, fn)
	return fn
}

func (p *CloseForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if p.from_c.IsSelf {
		return types.POSITIVE
	}
	// return p.from_c.Type.Polarity()

	return types.UNKNOWN
}

// Forward: fwd to_c from_c
type ForwardForm struct {
	to_c    Name
	from_c  Name
	to_drop bool
}

func NewForward(to_c, from_c Name) *ForwardForm {
	return &ForwardForm{to_c: to_c, from_c: from_c, to_drop: false}
}

func NewDroppableForward(to_c, from_c Name) *ForwardForm {
	return &ForwardForm{to_c: to_c, from_c: from_c, to_drop: true}
}

func (p *ForwardForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("fwd ")
	buf.WriteString(p.to_c.String())
	buf.WriteString(" ")
	buf.WriteString(p.from_c.String())
	return buf.String()
}

func (p *ForwardForm) StringShort() string {
	return p.String()
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

func (p *ForwardForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	return p.from_c.Polarity(fromTypes, globalEnvironment)

	// if fromTypes {
	// 	// Get polarity from the type
	// }

	// // if p.from_c.
	// // todo check for annotated polarity of from_c

	// return types.UNKNOWN
}

// Split: <channel_one, channel_two> <- split from_c; P
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

func (p *SplitForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString("<")
	buf.WriteString(p.channel_one.String())
	buf.WriteString(",")
	buf.WriteString(p.channel_two.String())
	buf.WriteString("> <- split ")
	buf.WriteString(p.from_c.String())
	buf.WriteString("; ...")
	return buf.String()
}

func (p *SplitForm) Substitute(old, new Name) {
	p.from_c.Substitute(old, new)

	if !p.channel_one.Equal(old) && !p.channel_two.Equal(old) {
		// channel_one: channel_one,
		// channel_two: channel_two,
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

func (p *SplitForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	// Get polarity from continuation
	return p.continuation_e.Polarity(fromTypes, globalEnvironment)
}

// Call: func(param1, ...)
type CallForm struct {
	functionName string
	parameters   []Name
	ProviderType types.SessionType
}

func NewCall(functionName string, parameters []Name) *CallForm {
	return &CallForm{
		functionName: functionName,
		parameters:   parameters,
	}
}

func (p *CallForm) String() string {
	var buf bytes.Buffer
	buf.WriteString(p.functionName)
	buf.WriteString("(")
	buf.WriteString(NamesToString(p.parameters))
	buf.WriteString(")")
	return buf.String()
}

func (p *CallForm) StringShort() string {
	return p.String()
}

func (p *CallForm) Substitute(old, new Name) {
	for i := range p.parameters {
		p.parameters[i].Substitute(old, new)
	}
}

func (p *CallForm) FreeNames() []Name {
	var fn []Name
	for i := range p.parameters {
		fn = appendIfNotSelf(p.parameters[i], fn)
	}
	return fn
}

func (p *CallForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if fromTypes {
		// Get polarity from the provider type
		return p.ProviderType.Polarity()
	}

	// Else, check the function's body
	// Fetch function by name and arity
	functionCall := GetFunctionByNameArity(*globalEnvironment.FunctionDefinitions, p.functionName, len(p.parameters))

	if functionCall != nil {
		if fromTypes {
			// If there are types, get the polarity from the type directly
			return functionCall.Type.Polarity()
		} else {
			return functionCall.Body.Polarity(fromTypes, globalEnvironment)
		}
	}

	return types.UNKNOWN
}

func (p *CallForm) FunctionName() string {
	return p.functionName
}

// Wait: wait to_c; P
type WaitForm struct {
	to_c           Name
	continuation_e Form
}

func NewWait(to_c Name, continuation_e Form) *WaitForm {
	return &WaitForm{
		to_c:           to_c,
		continuation_e: continuation_e}
}

func (p *WaitForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("wait ")
	buf.WriteString(p.to_c.String())
	buf.WriteString("; ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

func (p *WaitForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString("wait ")
	buf.WriteString(p.to_c.String())
	buf.WriteString("; ...")
	return buf.String()
}

func (p *WaitForm) Substitute(old, new Name) {
	p.to_c.Substitute(old, new)
	p.continuation_e.Substitute(old, new)
}

func (p *WaitForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.to_c, fn)
	fn = mergeTwoNamesList(fn, p.continuation_e.FreeNames())
	return fn
}

func (p *WaitForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	return p.continuation_e.Polarity(fromTypes, globalEnvironment)
}

// Cast: cast to_c<continuation_c>
type CastForm struct {
	to_c           Name
	continuation_c Name
}

func NewCast(to_c, continuation_c Name) *CastForm {
	return &CastForm{
		to_c:           to_c,
		continuation_c: continuation_c}
}

func (p *CastForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("cast ")
	buf.WriteString(p.to_c.String())
	buf.WriteString("<")
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(">")
	return buf.String()
}

func (p *CastForm) StringShort() string {
	return p.String()
}

func (p *CastForm) Substitute(old, new Name) {
	p.to_c.Substitute(old, new)
	p.continuation_c.Substitute(old, new)
}

// Free names, excluding self references
func (p *CastForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.to_c, fn)
	fn = appendIfNotSelf(p.continuation_c, fn)
	return fn
}

// Polarity of a send process can be inferred directly from itself
func (p *CastForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if p.to_c.IsSelf {
		return types.POSITIVE
	}

	// Type from continuation channel
	return p.continuation_c.Polarity(fromTypes, globalEnvironment)
}

// Shift: continuation_c <- shift from_c; P
type ShiftForm struct {
	continuation_c Name
	from_c         Name
	continuation_e Form
}

func NewShift(continuation_c, from_c Name, continuation_e Form) *ShiftForm {
	return &ShiftForm{
		continuation_c: continuation_c,
		from_c:         from_c,
		continuation_e: continuation_e}
}

func (p *ShiftForm) String() string {
	var buf bytes.Buffer
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(" <- shift ")
	buf.WriteString(p.from_c.String())
	buf.WriteString("; ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

func (p *ShiftForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString(p.continuation_c.String())
	buf.WriteString(" <- shift ")
	buf.WriteString(p.from_c.String())
	buf.WriteString("; ...")
	return buf.String()
}

func (p *ShiftForm) Substitute(old, new Name) {
	p.from_c.Substitute(old, new)

	if !p.continuation_c.Equal(old) {
		// continuation_c: continuation_c,
		p.continuation_e.Substitute(old, new)
	}
}

func (p *ShiftForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.from_c, fn)
	continuation_e_excluding_bound_names := removeBoundName(p.continuation_e.FreeNames(), p.continuation_c)
	// continuation_e_excluding_bound_names = removeBoundName(continuation_e_excluding_bound_names, p.continuation_c)
	fn = mergeTwoNamesList(fn, continuation_e_excluding_bound_names)
	return fn
}

// Polarity of a receive process can be inferred directly from itself
func (p *ShiftForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if p.from_c.IsSelf {
		return types.NEGATIVE
	}

	// Lookup polarity from the continuation
	return p.continuation_e.Polarity(fromTypes, globalEnvironment)
}

// Drop: drop client_c; P
type DropForm struct {
	client_c       Name
	continuation_e Form
}

func NewDrop(client_c Name, continuation_e Form) *DropForm {
	return &DropForm{
		client_c:       client_c,
		continuation_e: continuation_e}
}

func (p *DropForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("drop ")
	buf.WriteString(p.client_c.String())
	buf.WriteString("; ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

func (p *DropForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString("drop ")
	buf.WriteString(p.client_c.String())
	buf.WriteString("; ...")
	return buf.String()
}

func (p *DropForm) Substitute(old, new Name) {
	p.client_c.Substitute(old, new)
	p.continuation_e.Substitute(old, new)
}

func (p *DropForm) FreeNames() []Name {
	var fn []Name
	fn = appendIfNotSelf(p.client_c, fn)
	fn = mergeTwoNamesList(fn, p.continuation_e.FreeNames())
	return fn
}

func (p *DropForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	if p.client_c.IsSelf {
		// Shouldn't be self
		return types.UNKNOWN
	}

	// Lookup polarity from the continuation
	return p.continuation_e.Polarity(fromTypes, globalEnvironment)
}

// Print: print l; P
// Used to print label for debugging purposes
type PrintForm struct {
	label          Label
	continuation_e Form
}

func NewPrint(label Label, continuation_e Form) *PrintForm {
	return &PrintForm{
		label:          label,
		continuation_e: continuation_e,
	}
}

func (p *PrintForm) String() string {
	var buf bytes.Buffer
	buf.WriteString("print ")
	buf.WriteString(p.label.String())
	buf.WriteString("; ")
	buf.WriteString(p.continuation_e.String())
	return buf.String()
}

func (p *PrintForm) StringShort() string {
	var buf bytes.Buffer
	buf.WriteString("print ")
	buf.WriteString(p.label.String())
	buf.WriteString("; ...")
	return buf.String()
}

func (p *PrintForm) Substitute(old, new Name) {
	p.continuation_e.Substitute(old, new)
}

func (p *PrintForm) FreeNames() []Name {
	return p.continuation_e.FreeNames()
}

func (p *PrintForm) Polarity(fromTypes bool, globalEnvironment *GlobalEnvironment) types.Polarity {
	// Lookup polarity from the continuation
	return p.continuation_e.Polarity(fromTypes, globalEnvironment)
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
	case *CallForm:
		f1, ok1 := form1.(*CallForm)
		f2, ok2 := form2.(*CallForm)

		if ok1 && ok2 {
			if f1.functionName != f2.functionName {
				return false
			}

			if len(f1.parameters) != len(f2.parameters) {
				return false
			}

			for i := range f1.parameters {
				if !f1.parameters[i].Equal(f2.parameters[i]) {
					return false
				}
			}
			// true
			return true
		}
	case *WaitForm:
		f1, ok1 := form1.(*WaitForm)
		f2, ok2 := form2.(*WaitForm)

		if ok1 && ok2 {
			return f1.to_c.Equal(f2.to_c) && EqualForm(f1.continuation_e, f2.continuation_e)
		}
	case *CastForm:
		f1, ok1 := form1.(*CastForm)
		f2, ok2 := form2.(*CastForm)

		if ok1 && ok2 {
			return f1.continuation_c.Equal(f2.continuation_c) && f1.to_c.Equal(f2.to_c)
		}
	case *ShiftForm:
		f1, ok1 := form1.(*ShiftForm)
		f2, ok2 := form2.(*ShiftForm)

		if ok1 && ok2 {
			return f1.continuation_c.Equal(f2.continuation_c) && f1.from_c.Equal(f2.from_c) && EqualForm(f1.continuation_e, f2.continuation_e)
		}
	case *DropForm:
		f1, ok1 := form1.(*DropForm)
		f2, ok2 := form2.(*DropForm)

		if ok1 && ok2 {
			return f1.client_c.Equal(f2.client_c) && EqualForm(f1.continuation_e, f2.continuation_e)
		}
	case *PrintForm:
		f1, ok1 := form1.(*PrintForm)
		f2, ok2 := form2.(*PrintForm)

		if ok1 && ok2 {
			return f1.label.Equal(f2.label) && EqualForm(f1.continuation_e, f2.continuation_e)
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
			return NewSend(*p.to_c.Copy(), *p.payload_c.Copy(), *p.continuation_c.Copy())
		}
	case *ReceiveForm:
		p, ok := orig.(*ReceiveForm)
		if ok {
			cont := CopyForm(p.continuation_e)
			return NewReceive(*p.payload_c.Copy(), *p.continuation_c.Copy(), *p.from_c.Copy(), cont)
		}
	case *SelectForm:
		p, ok := orig.(*SelectForm)
		if ok {
			return NewSelect(*p.to_c.Copy(), p.label, *p.continuation_c.Copy())
		}
	case *CaseForm:
		p, ok := orig.(*CaseForm)
		if ok {
			branches := make([]*BranchForm, len(p.branches))

			for i := 0; i < len(p.branches); i++ {
				b := CopyForm(p.branches[i]).(*BranchForm)
				branches[i] = b
			}

			return NewCase(*p.from_c.Copy(), branches)
		}

	case *BranchForm:
		p, ok := orig.(*BranchForm)
		if ok {
			cont := CopyForm(p.continuation_e)
			return NewBranch(p.label, *p.payload_c.Copy(), cont)
		}
	case *CloseForm:
		p, ok := orig.(*CloseForm)
		if ok {
			return NewClose(*p.from_c.Copy())
		}
	case *NewForm:
		p, ok := orig.(*NewForm)
		if ok {
			body := CopyForm(p.body)
			cont := CopyForm(p.continuation_e)
			return NewNew(*p.continuation_c.Copy(), body, cont)
		}
	case *ForwardForm:
		p, ok := orig.(*ForwardForm)
		if ok {
			return NewForward(*p.to_c.Copy(), p.from_c)
		}
	case *SplitForm:
		p, ok := orig.(*SplitForm)
		if ok {
			cont := CopyForm(p.continuation_e)
			return NewSplit(*p.channel_one.Copy(), *p.channel_two.Copy(), *p.from_c.Copy(), cont)
		}
	case *CallForm:
		p, ok := orig.(*CallForm)
		if ok {
			copiedParameters := make([]Name, len(p.parameters))
			for i := 0; i < len(p.parameters); i++ {
				copiedParameters[i] = *p.parameters[i].Copy()
			}
			return NewCall(p.functionName, copiedParameters)
		}
	case *WaitForm:
		p, ok := orig.(*WaitForm)
		if ok {
			body := CopyForm(p.continuation_e)
			return NewWait(*p.to_c.Copy(), body)
		}
	case *CastForm:
		p, ok := orig.(*CastForm)
		if ok {
			return NewCast(*p.to_c.Copy(), p.continuation_c)
		}
	case *ShiftForm:
		p, ok := orig.(*ShiftForm)
		if ok {
			cont := CopyForm(p.continuation_e)
			return NewShift(*p.continuation_c.Copy(), *p.from_c.Copy(), cont)
		}
	case *DropForm:
		p, ok := orig.(*DropForm)
		if ok {
			body := CopyForm(p.continuation_e)
			return NewDrop(*p.client_c.Copy(), body)
		}
	// Debug
	case *PrintForm:
		p, ok := orig.(*PrintForm)
		if ok {
			return NewPrint(p.label, CopyForm(p.continuation_e))
		}
	}

	panic("modify CopyForm to handle new type")
}

// Return true if the given for has continuation expression, or false otherwise (i.e. follows an axiomatic rule)
func FormHasContinuation(form Form) bool {
	switch interface{}(form).(type) {
	case *SendForm:
		return false
	case *SelectForm:
		return false
	case *CloseForm:
		return false
	case *ForwardForm:
		return false
	case *CallForm:
		return false
	case *CastForm:
		return false
	default:
		// These have a continuation:
		// -> ReceiveForm:
		// -> CaseForm:
		// -> BranchForm:
		// -> NewForm:
		// -> SplitForm:
		// -> WaitForm:
		// -> ShiftForm:
		// -> DropForm:
		// -> PrintForm:
		return true
	}
}

type Label struct {
	L string
}

func (p *Label) String() string {
	return p.L
}

func (label1 *Label) Equal(label2 Label) bool {
	return label1.String() == label2.String()
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
