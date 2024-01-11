package types

import (
	"bytes"
	"fmt"
	"phi/position"
	"reflect"

	"golang.org/x/exp/slices"
)

type SessionTypeDefinition struct {
	SessionType SessionType
	Name        string
	Position    position.Position
	Modality    Modality
}

type SessionType interface {
	String() string
	StringWithModality() string
	StringWithOuterModality() string
	Polarity() Polarity
	Modality() Modality

	// used for type structure checks
	checkTypeLabels(LabelledTypesEnv) error
	checkTypeModalities(LabelledTypesEnv, Modality) error
	inferModality(LabelledTypesEnv, map[string]bool) Modality
	assignUnsetModalities(LabelledTypesEnv, Modality)
	isContractive(LabelledTypesEnv, map[string]bool) bool
}

// Label
type LabelType struct {
	Label string
	Mode  Modality
}

func NewLabelType(i string, mode Modality) *LabelType {
	return &LabelType{
		Label: i,
		Mode:  mode,
	}
}

func (q *LabelType) String() string {
	return q.Label
}

func (q *LabelType) StringWithModality() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]")
	buffer.WriteString(q.Label)

	return buffer.String()
}

func (q *LabelType) StringWithOuterModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.Label)
	buffer.WriteString(" [")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]")

	return buffer.String()
}

func (q *LabelType) Modality() Modality {
	return q.Mode
}

// Unit: 1
type UnitType struct {
	Mode Modality
}

func NewUnitType(mode Modality) *UnitType {
	return &UnitType{
		Mode: mode,
	}
}

func (q *UnitType) String() string {
	return "1"
}

func (q *UnitType) StringWithModality() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]1")

	return buffer.String()
}

func (q *UnitType) StringWithOuterModality() string {
	var buffer bytes.Buffer
	buffer.WriteString("1 [")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]")
	return buffer.String()
}

func (q *UnitType) Modality() Modality {
	return q.Mode
}

// Send: A * B
type SendType struct {
	Left  SessionType
	Right SessionType
	Mode  Modality
}

func NewSendType(left, right SessionType, mode Modality) *SendType {
	return &SendType{
		Left:  left,
		Right: right,
		Mode:  mode,
	}
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

func (q *SendType) StringWithModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.Left.StringWithModality())
	buffer.WriteString(" [")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]* ")
	buffer.WriteString(q.Right.StringWithModality())
	return buffer.String()
}

func (q *SendType) StringWithOuterModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.Left.String())
	buffer.WriteString(" * ")
	buffer.WriteString(q.Right.String())
	buffer.WriteString(" [")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]")
	return buffer.String()
}

func (q *SendType) Modality() Modality {
	return q.Mode
}

// Receive: A -* B
type ReceiveType struct {
	Left  SessionType
	Right SessionType
	Mode  Modality
}

func NewReceiveType(left, right SessionType, mode Modality) *ReceiveType {
	return &ReceiveType{
		Left:  left,
		Right: right,
		Mode:  mode,
	}
}

func (q *ReceiveType) String() string {
	var buffer bytes.Buffer
	// buffer.WriteString("(")
	buffer.WriteString(q.Left.String())
	buffer.WriteString(" -* ")
	buffer.WriteString(q.Right.String())
	// buffer.WriteString(")")
	return buffer.String()
}

func (q *ReceiveType) StringWithModality() string {
	var buffer bytes.Buffer
	// buffer.WriteString("(")
	buffer.WriteString(q.Left.StringWithModality())
	buffer.WriteString(" [")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]-* ")
	buffer.WriteString(q.Right.StringWithModality())
	// buffer.WriteString(")")
	return buffer.String()
}

func (q *ReceiveType) StringWithOuterModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.Left.String())
	buffer.WriteString(" -* ")
	buffer.WriteString(q.Right.String())
	buffer.WriteString(" [")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]")
	return buffer.String()
}

func (q *ReceiveType) Modality() Modality {
	return q.Mode
}

// SelectLabel: +{ }
type SelectLabelType struct {
	Branches []Option
	Mode     Modality
}

func NewSelectLabelType(branches []Option, mode Modality) *SelectLabelType {
	return &SelectLabelType{
		Branches: branches,
		Mode:     mode,
	}
}

func (q *SelectLabelType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("+{")
	buffer.WriteString(stringifyBranches(q.Branches))
	buffer.WriteString("}")
	return buffer.String()
}

func (q *SelectLabelType) StringWithModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("+{")
	buffer.WriteString(stringifyBranchesWithModalities(q.Branches))
	buffer.WriteString("}")
	return buffer.String()
}

func (q *SelectLabelType) StringWithOuterModality() string {
	var buffer bytes.Buffer
	buffer.WriteString("+{")
	buffer.WriteString(stringifyBranches(q.Branches))
	buffer.WriteString("}")
	buffer.WriteString(" [")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]")
	return buffer.String()
}

func (q *SelectLabelType) Modality() Modality {
	return q.Mode
}

// BranchCase: & { }
type BranchCaseType struct {
	Branches []Option
	Mode     Modality
}

func NewBranchCaseType(branches []Option, mode Modality) *BranchCaseType {
	return &BranchCaseType{
		Branches: branches,
		Mode:     mode,
	}
}

func (q *BranchCaseType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("&{")
	buffer.WriteString(stringifyBranches(q.Branches))
	buffer.WriteString("}")
	return buffer.String()
}

func (q *BranchCaseType) StringWithModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("&{")
	buffer.WriteString(stringifyBranchesWithModalities(q.Branches))
	buffer.WriteString("}")
	return buffer.String()
}

func (q *BranchCaseType) StringWithOuterModality() string {
	var buffer bytes.Buffer
	buffer.WriteString("&{")
	buffer.WriteString(stringifyBranches(q.Branches))
	buffer.WriteString("}")
	buffer.WriteString(" [")
	buffer.WriteString(q.Mode.String())
	buffer.WriteString("]")
	return buffer.String()
}

func (q *BranchCaseType) Modality() Modality {
	return q.Mode
}

// Up shift: m /\ n ...
type UpType struct {
	From         Modality
	To           Modality
	Continuation SessionType
}

func NewUpType(From, To Modality, Continuation SessionType) *UpType {
	return &UpType{
		From:         From,
		To:           To,
		Continuation: Continuation,
	}
}

func (q *UpType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.From.String())
	buffer.WriteString("/\\")
	buffer.WriteString(q.To.String())
	buffer.WriteString(" ")
	buffer.WriteString(q.Continuation.String())
	return buffer.String()
}

func (q *UpType) StringWithModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.From.String())
	buffer.WriteString("/\\")
	buffer.WriteString(q.To.String())
	buffer.WriteString(" ")
	buffer.WriteString(q.Continuation.StringWithModality())
	return buffer.String()
}

func (q *UpType) StringWithOuterModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.From.String())
	buffer.WriteString("/\\")
	buffer.WriteString(q.To.String())
	buffer.WriteString(" ")
	buffer.WriteString(q.Continuation.String())
	return buffer.String()
}

func (q *UpType) Modality() Modality {
	return q.To
}

// Down shift: m \/ n ...
type DownType struct {
	From         Modality
	To           Modality
	Continuation SessionType
}

func NewDownType(From, To Modality, Continuation SessionType) *DownType {
	return &DownType{
		From:         From,
		To:           To,
		Continuation: Continuation,
	}
}

func (q *DownType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.From.String())
	buffer.WriteString("\\/")
	buffer.WriteString(q.To.String())
	buffer.WriteString(" ")
	buffer.WriteString(q.Continuation.String())
	return buffer.String()
}

func (q *DownType) StringWithModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.From.String())
	buffer.WriteString("\\/")
	buffer.WriteString(q.To.String())
	buffer.WriteString(" ")
	buffer.WriteString(q.Continuation.StringWithModality())
	return buffer.String()
}

func (q *DownType) StringWithOuterModality() string {
	var buffer bytes.Buffer
	buffer.WriteString(q.From.String())
	buffer.WriteString("\\/")
	buffer.WriteString(q.To.String())
	buffer.WriteString(" ")
	buffer.WriteString(q.Continuation.String())
	return buffer.String()
}

func (q *DownType) Modality() Modality {
	return q.To
}

// Branch/Case option
type Option struct {
	Label       string
	SessionType SessionType
}

func (option *Option) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(option.Label)
	buffer.WriteString(" : ")
	buffer.WriteString(option.SessionType.String())

	return buffer.String()
}

func stringifyBranches(options []Option) string {
	var buf bytes.Buffer

	for i, j := range options {
		buf.WriteString(j.Label)
		buf.WriteString(" : ")
		buf.WriteString(j.SessionType.String())

		if i < len(options)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func stringifyBranchesWithModalities(options []Option) string {
	var buf bytes.Buffer

	for i, j := range options {
		buf.WriteString(j.Label)
		buf.WriteString(" : ")
		buf.WriteString(j.SessionType.StringWithModality())

		if i < len(options)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func NewOption(label string, sessionType SessionType) *Option {
	return &Option{
		Label:       label,
		SessionType: sessionType,
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
		var presentSnapshot bytes.Buffer
		presentSnapshot.WriteString(type1.String())
		presentSnapshot.WriteString(type1.Modality().String())
		presentSnapshot.WriteString("|")
		presentSnapshot.WriteString(type2.String())
		presentSnapshot.WriteString(type2.Modality().String())

		_, exists := snapshots[presentSnapshot.String()]
		if exists {
			return true
		}

		if isLabel1 && isLabel2 && f1.Label == f2.Label {
			return f1.Modality().Equals(f2.Modality())
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
		var newSnapshot bytes.Buffer
		newSnapshot.WriteString(type1.String())
		newSnapshot.WriteString(type1.Modality().String())
		newSnapshot.WriteString("|")
		newSnapshot.WriteString(type2.String())
		newSnapshot.WriteString(type2.Modality().String())
		snapshots[newSnapshot.String()] = true

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
		f1, ok1 := type1.(*UnitType)
		f2, ok2 := type2.(*UnitType)
		return ok1 && ok2 && f1.Modality().Equals(f2.Modality())

	case *SendType:
		f1, ok1 := type1.(*SendType)
		f2, ok2 := type2.(*SendType)

		if ok1 && ok2 {
			return f1.Modality().Equals(f2.Modality()) && innerEqualType(f1.Left, f2.Left, snapshots, labelledTypesEnv) && innerEqualType(f1.Right, f2.Right, snapshots, labelledTypesEnv)
		}

	case *ReceiveType:
		f1, ok1 := type1.(*ReceiveType)
		f2, ok2 := type2.(*ReceiveType)

		if ok1 && ok2 {
			return f1.Modality().Equals(f2.Modality()) && innerEqualType(f1.Left, f2.Left, snapshots, labelledTypesEnv) && innerEqualType(f1.Right, f2.Right, snapshots, labelledTypesEnv)
		}

	case *SelectLabelType:
		f1, ok1 := type1.(*SelectLabelType)
		f2, ok2 := type2.(*SelectLabelType)

		if ok1 && ok2 && len(f1.Branches) == len(f2.Branches) {
			return f1.Modality().Equals(f2.Modality()) && equalTypeBranch(f1.Branches, f2.Branches, snapshots, labelledTypesEnv)
		}

	case *BranchCaseType:
		f1, ok1 := type1.(*BranchCaseType)
		f2, ok2 := type2.(*BranchCaseType)

		if ok1 && ok2 && len(f1.Branches) == len(f2.Branches) {
			// order doesn't matters
			return f1.Modality().Equals(f2.Modality()) && equalTypeBranch(f1.Branches, f2.Branches, snapshots, labelledTypesEnv)
		}

	case *UpType:
		f1, ok1 := type1.(*UpType)
		f2, ok2 := type2.(*UpType)

		if ok1 && ok2 {
			return f1.To.Equals(f2.To) && f1.From.Equals(f2.From) && innerEqualType(f1.Continuation, f2.Continuation, snapshots, labelledTypesEnv)
		}

	case *DownType:
		f1, ok1 := type1.(*DownType)
		f2, ok2 := type2.(*DownType)

		if ok1 && ok2 {
			return f1.To.Equals(f2.To) && f1.From.Equals(f2.From) && innerEqualType(f1.Continuation, f2.Continuation, snapshots, labelledTypesEnv)
		}
	}

	fmt.Printf("issue in EqualType for type %s\n", a)
	return false
}

// Compare branches in an unordered way. Here we are assuming that both branches contain unique labels
func equalTypeBranch(options1, options2 []Option, snapshots map[string]bool, labelledTypesEnv LabelledTypesEnv) bool {
	if len(options1) != len(options2) {
		return false
	}

	// Match each label to the other set
	for _, b := range options1 {
		matchingBranch, foundMatchingBranch := LookupBranchByLabel(options2, b.Label)
		if foundMatchingBranch {
			if !innerEqualType(b.SessionType, matchingBranch.SessionType, snapshots, labelledTypesEnv) {
				// If inner types do not match, then stop checking
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// Lookup branch by label
func LookupBranchByLabel(branches []Option, label string) (*Option, bool) {
	for index := range branches {
		if branches[index].Label == label {
			return &branches[index], true
		}
	}

	return nil, false
}

// Given a list of branch, this function returns the sub-list of branches that are not in the list of labelsChecked
func GetUncheckedBranches(branches []Option, labelsChecked []string) []Option {
	result := []Option{}

	for _, branch := range branches {
		if !slices.Contains(labelsChecked, branch.Label) {
			result = append(result, branch)
		}
	}

	return result
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
			return NewLabelType(p.Label, p.Mode.Copy())
		}
	case *UnitType:
		p, ok := orig.(*UnitType)
		if ok {
			return NewUnitType(p.Mode.Copy())
		}
	case *SendType:
		p, ok := orig.(*SendType)
		if ok {
			return NewSendType(CopyType(p.Left), CopyType(p.Right), p.Mode.Copy())
		}
	case *ReceiveType:
		p, ok := orig.(*ReceiveType)
		if ok {
			return NewReceiveType(CopyType(p.Left), CopyType(p.Right), p.Mode.Copy())
		}
	case *SelectLabelType:
		p, ok := orig.(*SelectLabelType)
		if ok {
			branches := make([]Option, len(p.Branches))

			for i := 0; i < len(p.Branches); i++ {
				branches[i].Label = p.Branches[i].Label
				branches[i].SessionType = CopyType(p.Branches[i].SessionType)
			}

			return NewSelectLabelType(branches, p.Mode.Copy())
		}
	case *BranchCaseType:
		p, ok := orig.(*BranchCaseType)
		if ok {
			branches := make([]Option, len(p.Branches))

			for i := 0; i < len(p.Branches); i++ {
				branches[i].Label = p.Branches[i].Label
				branches[i].SessionType = CopyType(p.Branches[i].SessionType)
			}

			return NewBranchCaseType(branches, p.Mode.Copy())
		}
	case *UpType:
		p, ok := orig.(*UpType)
		if ok {
			return NewUpType(p.From.Copy(), p.To.Copy(), CopyType(p.Continuation))
		}
	case *DownType:
		p, ok := orig.(*DownType)
		if ok {
			return NewDownType(p.From.Copy(), p.To.Copy(), CopyType(p.Continuation))
		}
	}

	panic("Should not happen (type)")
	// return nil
}

// The labelled types environment is constant and set once at the beginning. The information is obtained from the 'type A = ...' definitions.
// labelledTypesEnv: map of labels to their session type (wrapped in a LabelledType struct)

type LabelledTypesEnv map[string]LabelledType
type LabelledType struct {
	Name string
	Mode Modality
	Type SessionType
}

func ProduceLabelledSessionTypeEnvironment(typeDefs []SessionTypeDefinition) LabelledTypesEnv {
	labelledTypesEnv := make(LabelledTypesEnv)
	for _, j := range typeDefs {
		labelledTypesEnv[j.Name] = LabelledType{
			Type: j.SessionType,
			Name: j.Name,
			Mode: j.Modality,
		}
	}

	return labelledTypesEnv
}

func LabelledTypedExists(labelledTypesEnv LabelledTypesEnv, key string) bool {
	_, ok := labelledTypesEnv[key]
	return ok
}

// Weaenable types allow for channels to be dropped
func IsWeakenable(sessionType SessionType) bool {
	return sessionType.Modality().AllowsWeakening()
}

// Contraction types allow for channels to be copied/splits
func IsContractable(sessionType SessionType) bool {
	return sessionType.Modality().AllowsContraction()
}

func UnfoldIfNeeded(orig SessionType, typeDefs *[]SessionTypeDefinition) SessionType {
	if orig == nil {
		return nil
	}

	_, labelType := orig.(*LabelType)

	if labelType {
		return Unfold(orig, ProduceLabelledSessionTypeEnvironment(*typeDefs))
	}

	return orig
}

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
func FetchSelectBranch(branches []Option, label string) (SessionType, bool) {
	for _, branch := range branches {
		if branch.Label == label {
			return branch.SessionType, true
		}
	}

	return nil, false
}

// During session type construction (i.e. type parsing), we use construct a SessionTypeInitial.
// This allows us to define a general modality for the type.

// SessionTypeInitial can be converted into a SessionType (as used by the typechecker), where modes are then inserted into each SessionType struct (rather than just one at the beginning).

func ConvertSessionTypeInitialToSessionType(st SessionTypeInitial) SessionType {
	defaultMode := NewUnsetMode()
	return st.toSessionType(defaultMode)
}

// SessionTypeInitial defines the structure for session types with explicit modalities.
// The modes are defined as an explicit struct (usually at the beginning of the type).
type SessionTypeInitial interface {
	toSessionType(Modality) SessionType
}

// Explicit mode, e.g. unrestricted A, where the mode of A becomes unrestricted
// Sets the modality for the continuation type
type ExplicitModeTypeInitial struct {
	Modality     Modality
	Continuation SessionTypeInitial
}

func NewExplicitModeTypeInitial(modality Modality, continuation SessionTypeInitial) *ExplicitModeTypeInitial {
	return &ExplicitModeTypeInitial{
		Modality:     modality,
		Continuation: continuation,
	}
}

func (q *ExplicitModeTypeInitial) toSessionType(oldModality Modality) SessionType {
	return q.Continuation.toSessionType(q.Modality)
}

// Label
type LabelTypeInitial struct {
	Label string
}

func NewLabelTypeInitial(i string) *LabelTypeInitial {
	return &LabelTypeInitial{
		Label: i,
	}
}

func (q *LabelTypeInitial) toSessionType(mode Modality) SessionType {
	return NewLabelType(q.Label, mode)
}

// Unit: 1
type UnitTypeInitial struct{}

func NewUnitTypeInitial() *UnitTypeInitial {
	return &UnitTypeInitial{}
}

func (q *UnitTypeInitial) toSessionType(mode Modality) SessionType {
	return NewUnitType(mode)
}

// Send: A * B
type SendTypeInitial struct {
	Left  SessionTypeInitial
	Right SessionTypeInitial
}

func NewSendTypeInitial(left, right SessionTypeInitial) *SendTypeInitial {
	return &SendTypeInitial{
		Left:  left,
		Right: right,
	}
}

func (q *SendTypeInitial) toSessionType(mode Modality) SessionType {
	return NewSendType(q.Left.toSessionType(mode), q.Right.toSessionType(mode), mode)
}

// Receive: A -* B
type ReceiveTypeInitial struct {
	Left  SessionTypeInitial
	Right SessionTypeInitial
}

func NewReceiveTypeInitial(left, right SessionTypeInitial) *ReceiveTypeInitial {
	return &ReceiveTypeInitial{
		Left:  left,
		Right: right,
	}
}

func (q *ReceiveTypeInitial) toSessionType(mode Modality) SessionType {
	return NewReceiveType(q.Left.toSessionType(mode), q.Right.toSessionType(mode), mode)
}

// SelectLabel: +{ }
type SelectLabelTypeInitial struct {
	Branches []OptionInitial
}

func NewSelectLabelTypeInitial(options []OptionInitial) *SelectLabelTypeInitial {
	return &SelectLabelTypeInitial{
		Branches: options,
	}
}

func (q *SelectLabelTypeInitial) toSessionType(mode Modality) SessionType {
	branches := make([]Option, len(q.Branches))

	for i := 0; i < len(q.Branches); i++ {
		branches[i].Label = q.Branches[i].Label
		branches[i].SessionType = q.Branches[i].Session_type.toSessionType(mode)
	}

	return NewSelectLabelType(branches, mode)
}

// BranchCase: & { }
type BranchCaseTypeInitial struct {
	Branches []OptionInitial
}

func NewBranchCaseTypeInitial(options []OptionInitial) *BranchCaseTypeInitial {
	return &BranchCaseTypeInitial{
		Branches: options,
	}
}

func (q *BranchCaseTypeInitial) toSessionType(mode Modality) SessionType {
	branches := make([]Option, len(q.Branches))

	for i := 0; i < len(q.Branches); i++ {
		branches[i].Label = q.Branches[i].Label
		branches[i].SessionType = q.Branches[i].Session_type.toSessionType(mode)
	}

	return NewBranchCaseType(branches, mode)
}

// Up shift: m /\ n ...
type UpTypeInitial struct {
	From         Modality
	To           Modality
	Continuation SessionTypeInitial
}

func NewUpTypeInitial(From, To Modality, Continuation SessionTypeInitial) *UpTypeInitial {
	return &UpTypeInitial{
		From:         From,
		To:           To,
		Continuation: Continuation,
	}
}

func (q *UpTypeInitial) toSessionType(mode Modality) SessionType {
	// If 'mode' does not match the q.To, then it is an ill formed type, however a SessionTypeInitial is lenient during construct and allows this. This is checked later on during the preliminary checks
	// The mode of the continuation type has to be set to q.From
	return NewUpType(q.From, q.To, q.Continuation.toSessionType(q.From))
}

// Down shift: m \/ n ...
type DownTypeInitial struct {
	From         Modality
	To           Modality
	Continuation SessionTypeInitial
}

func NewDownTypeInitial(From, To Modality, Continuation SessionTypeInitial) *DownTypeInitial {
	return &DownTypeInitial{
		From:         From,
		To:           To,
		Continuation: Continuation,
	}
}

func (q *DownTypeInitial) toSessionType(mode Modality) SessionType {
	return NewDownType(q.From, q.To, q.Continuation.toSessionType(q.From))
}

// Branch/Case option
type OptionInitial struct {
	Label        string
	Session_type SessionTypeInitial
}

func NewOptionInitial(label string, session_type SessionTypeInitial) *OptionInitial {
	return &OptionInitial{
		Label:        label,
		Session_type: session_type,
	}
}
