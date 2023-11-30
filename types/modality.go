package types

import "strings"

type Modality interface {
	String() string
	// Check which shift are allowed
	CanBeUpshiftedTo(Modality) bool
	CanBeDownshiftedTo(Modality) bool

	Copy() Modality
	Equals(Modality) bool

	// StructuralProperties() []string
	AllowsContraction() bool
	AllowsWeakening() bool
}

//           Unrestricted {W, C}             |
//          <            >                   |
//   Affine {W}            Replicable {W, C} |
//          >            <                   |
//               Linear ø                   \/  Downshifts allowed in this direction (and vice versa for upshifts)
//
// E.g. Since Unrestricted > Linear, then you can downshift from Unrestricted to Linear (but not upshift)
// You can upshift from Affine to Linear (since Affine > Linear)

func DefaultMode() *UnrestrictedMode {
	return NewUnrestrictedMode()
}

// Unrestricted => {W, C}
type UnrestrictedMode struct{}

func NewUnrestrictedMode() *UnrestrictedMode {
	return &UnrestrictedMode{}
}

func (q *UnrestrictedMode) String() string {
	return "unr"
}

func (q *UnrestrictedMode) Copy() Modality {
	return NewUnrestrictedMode()
}

func (q *UnrestrictedMode) AllowsContraction() bool {
	return true
}

func (q *UnrestrictedMode) AllowsWeakening() bool {
	return true
}

func (q *UnrestrictedMode) CanBeUpshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *UnrestrictedMode:
		return true
	case *ReplicableMode:
		return false
	case *AffineMode:
		return false
	case *LinearMode:
		return false
	default:
		panic("todo")
	}
}

func (q *UnrestrictedMode) CanBeDownshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *UnrestrictedMode:
		return true
	case *ReplicableMode:
		return true
	case *AffineMode:
		return true
	case *LinearMode:
		return true
	default:
		panic("todo")
	}
}

func (q *UnrestrictedMode) Equals(other Modality) bool {
	_, same := other.(*UnrestrictedMode)
	return same
}

// Replicable => {W, C}
type ReplicableMode struct{}

func NewReplicableMode() *ReplicableMode {
	return &ReplicableMode{}
}

func (q *ReplicableMode) String() string {
	return "rep"
}

func (q *ReplicableMode) Copy() Modality {
	return NewReplicableMode()
}

func (q *ReplicableMode) AllowsContraction() bool {
	return true
}

func (q *ReplicableMode) AllowsWeakening() bool {
	return true
}

func (q *ReplicableMode) CanBeUpshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *UnrestrictedMode:
		return true
	case *ReplicableMode:
		return true
	case *AffineMode:
		return false
	case *LinearMode:
		return false
	default:
		panic("todo")
	}
}

func (q *ReplicableMode) CanBeDownshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *UnrestrictedMode:
		return false
	case *ReplicableMode:
		return true
	case *AffineMode:
		return false // todo check with Adrian this relationship
	case *LinearMode:
		return true
	default:
		panic("todo")
	}
}

func (q *ReplicableMode) Equals(other Modality) bool {
	_, same := other.(*ReplicableMode)
	return same
}

// Affine => {W}
type AffineMode struct{}

func NewAffineMode() *AffineMode {
	return &AffineMode{}
}

func (q *AffineMode) String() string {
	return "aff"
}

func (q *AffineMode) Copy() Modality {
	return NewAffineMode()
}

func (q *AffineMode) AllowsContraction() bool {
	return true
}

func (q *AffineMode) AllowsWeakening() bool {
	return true
}

func (q *AffineMode) CanBeUpshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *UnrestrictedMode:
		return true
	case *ReplicableMode:
		return false // todo check with Adrian this relationship
	case *AffineMode:
		return true
	case *LinearMode:
		return false
	default:
		panic("todo")
	}
}

func (q *AffineMode) CanBeDownshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *UnrestrictedMode:
		return false
	case *ReplicableMode:
		return false
	case *AffineMode:
		return true
	case *LinearMode:
		return true
	default:
		panic("todo")
	}
}

func (q *AffineMode) Equals(other Modality) bool {
	_, same := other.(*AffineMode)
	return same
}

// Linear
type LinearMode struct{}

func NewLinearMode() *LinearMode {
	return &LinearMode{}
}

func (q *LinearMode) String() string {
	return "lin"
}

func (q *LinearMode) Copy() Modality {
	return NewLinearMode()
}

func (q *LinearMode) AllowsContraction() bool {
	return false
}

func (q *LinearMode) AllowsWeakening() bool {
	return false
}

func (q *LinearMode) CanBeUpshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *UnrestrictedMode:
		return true
	case *ReplicableMode:
		return true
	case *AffineMode:
		return true
	case *LinearMode:
		return true
	default:
		panic("todo")
	}
}

func (q *LinearMode) CanBeDownshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *UnrestrictedMode:
		return false
	case *ReplicableMode:
		return false
	case *AffineMode:
		return false
	case *LinearMode:
		return true
	default:
		panic("todo")
	}
}

func (q *LinearMode) Equals(other Modality) bool {
	_, same := other.(*LinearMode)
	return same
}

// Special Modes: Invalid & Unset

// Invalid
type InvalidMode struct {
	mode string
}

func NewInvalidMode(mode string) *InvalidMode {
	return &InvalidMode{mode}
}

func (q *InvalidMode) String() string {
	return "invalid: " + q.mode
}

func (q *InvalidMode) Copy() Modality {
	return NewInvalidMode(q.mode)
}

func (q *InvalidMode) AllowsContraction() bool {
	return false
}

func (q *InvalidMode) AllowsWeakening() bool {
	return false
}

func (q *InvalidMode) CanBeUpshiftedTo(toMode Modality) bool {
	panic("shouldn't shift invalid modes")
}

func (q *InvalidMode) CanBeDownshiftedTo(toMode Modality) bool {
	panic("shouldn't shift invalid modes")
}

func (q *InvalidMode) Equals(other Modality) bool {
	_, same := other.(*InvalidMode)
	return same
}

// Unset
type UnsetMode struct{}

func NewUnsetMode() *UnsetMode {
	return &UnsetMode{}
}

func (q *UnsetMode) String() string {
	return "unset"
}

func (q *UnsetMode) Copy() Modality {
	return NewUnsetMode()
}

func (q *UnsetMode) AllowsContraction() bool {
	return false
}

func (q *UnsetMode) AllowsWeakening() bool {
	return false
}

func (q *UnsetMode) CanBeUpshiftedTo(toMode Modality) bool {
	panic("mode is not set: couldn't shift")
}

func (q *UnsetMode) CanBeDownshiftedTo(toMode Modality) bool {
	panic("mode is not set: couldn't shift")
}

func (q *UnsetMode) Equals(other Modality) bool {
	_, same := other.(*UnsetMode)
	return same
}

// Shareable
// type ShareableMode struct{}

// Converts a string to a mode
func StringToMode(input string) Modality {
	input = strings.ToLower(input)

	switch input {
	case "u":
		return &UnrestrictedMode{}
	case "unr":
		return &UnrestrictedMode{}
	case "unrestricted":
		return &UnrestrictedMode{}
	case "r":
		return &ReplicableMode{}
	case "rep":
		return &ReplicableMode{}
	case "replicable":
		return &ReplicableMode{}
	case "a":
		return &AffineMode{}
	case "aff":
		return &AffineMode{}
	case "affine":
		return &AffineMode{}
	case "l":
		return &LinearMode{}
	case "lin":
		return &LinearMode{}
	case "linear":
		return &LinearMode{}
	default:
		return &InvalidMode{input}
	}
}

func AddMissingModalities(t *SessionType, labelledTypesEnv LabelledTypesEnv) {
	// Infer general modality of type
	mode := (*t).inferModality(labelledTypesEnv, make(map[string]bool))

	// Assign the modality to the inner type
	_, unset := mode.(*UnsetMode)
	if unset {
		// mode = UnsetMode, so use the default mode
		(*t).assignUnsetModalities(labelledTypesEnv, DefaultMode())

	} else {
		// mode ≠ UnsetMode
		(*t).assignUnsetModalities(labelledTypesEnv, mode)
	}
}

// Takes a list of type definitions and sets the the general modality (for the type definition)
// and also sets the modality for the inner type structure
// Any UnsetModes are replace with inferred modes
func SetModalityTypeDef(typesDef []SessionTypeDefinition) {
	labelledTypesEnv := ProduceLabelledSessionTypeEnvironment(typesDef)

	// First assign a general modality to each labelled type,
	// e.g. for type A = linear 1 -* 1, then A has
	for i := range typesDef {
		mode := typesDef[i].SessionType.inferModality(labelledTypesEnv, make(map[string]bool))
		// Set the found mode
		_, unset := mode.(*UnsetMode)
		if unset {
			// mode = UnsetMode
			typesDef[i].Modality = DefaultMode()
		} else {
			// mode ≠ UnsetMode
			typesDef[i].Modality = mode
		}
	}

	// Recreate the labelled types env (now each type will have a defined modality)
	labelledTypesEnv = ProduceLabelledSessionTypeEnvironment(typesDef)
	for i := range typesDef {
		typesDef[i].SessionType.assignUnsetModalities(labelledTypesEnv, typesDef[i].Modality)
	}
}

// Looks within a type to find the modality.
// Modalities can be defined explicitly (i.e. mode ≠ UnsetMode), or taken from an Up/Down shift type.
// If a label is reached, then the modality of that labelled type is checked. If a (mode-less) cycle is reached, then inference stops.
func (q *LabelType) inferModality(labelledTypesEnv LabelledTypesEnv, usedLabels map[string]bool) Modality {
	_, unset := q.Mode.(*UnsetMode)
	if !unset {
		// If the type already has a modality, then return it
		return q.Mode
	}

	// Fetch labelled type
	typeFromLabel, exists := labelledTypesEnv[q.Label]
	if exists {
		// type found
		if !usedLabels[q.Label] {
			// no cycle reached yet
			usedLabels[q.Label] = true
			return typeFromLabel.Type.inferModality(labelledTypesEnv, usedLabels)
		}
	}
	return q.Mode
}

func (q *UnitType) inferModality(labelledTypesEnv LabelledTypesEnv, usedLabels map[string]bool) Modality {
	return q.Mode
}

func (q *SendType) inferModality(labelledTypesEnv LabelledTypesEnv, usedLabels map[string]bool) Modality {
	_, unset := q.Mode.(*UnsetMode)
	if !unset {
		// If the type already has a modality, then return it
		return q.Mode
	}

	leftUsedLabel := copyMap(usedLabels)
	leftMode := q.Left.inferModality(labelledTypesEnv, leftUsedLabel)
	rightMode := q.Right.inferModality(labelledTypesEnv, usedLabels)

	commonMode := commonMode(leftMode, rightMode)

	// _, unset = commonMode.(*UnsetMode)
	// if !unset {
	// 	// If the common mode is defined/set, return it
	// 	return commonMode
	// }

	return commonMode
}

func (q *ReceiveType) inferModality(labelledTypesEnv LabelledTypesEnv, usedLabels map[string]bool) Modality {
	_, unset := q.Mode.(*UnsetMode)
	if !unset {
		// If the type already has a modality, then return it
		return q.Mode
	}

	leftUsedLabel := copyMap(usedLabels)
	leftMode := q.Left.inferModality(labelledTypesEnv, leftUsedLabel)
	rightMode := q.Right.inferModality(labelledTypesEnv, usedLabels)

	commonMode := commonMode(leftMode, rightMode)

	return commonMode
}

func (q *SelectLabelType) inferModality(labelledTypesEnv LabelledTypesEnv, usedLabels map[string]bool) Modality {
	_, unset := q.Mode.(*UnsetMode)
	if !unset {
		// If the type already has a modality, then return it
		return q.Mode
	}

	var commonModes []Modality
	for _, branch := range q.Branches {
		usedLabelsCopy := copyMap(usedLabels)
		branchMode := branch.SessionType.inferModality(labelledTypesEnv, usedLabelsCopy)
		commonModes = append(commonModes, branchMode)
	}

	commonMode := commonMode(commonModes...)

	return commonMode
}

func (q *BranchCaseType) inferModality(labelledTypesEnv LabelledTypesEnv, usedLabels map[string]bool) Modality {
	_, unset := q.Mode.(*UnsetMode)
	if !unset {
		// If the type already has a modality, then return it
		return q.Mode
	}

	var commonModes []Modality
	for _, branch := range q.Branches {
		usedLabelsCopy := copyMap(usedLabels)
		branchMode := branch.SessionType.inferModality(labelledTypesEnv, usedLabelsCopy)
		commonModes = append(commonModes, branchMode)
	}

	commonMode := commonMode(commonModes...)

	return commonMode
}

func (q *UpType) inferModality(labelledTypesEnv LabelledTypesEnv, usedLabels map[string]bool) Modality {
	return q.To
}

func (q *DownType) inferModality(labelledTypesEnv LabelledTypesEnv, usedLabels map[string]bool) Modality {
	return q.To
}

// Assigns a modality to each type (& inner types) with Unset Modes
// The modality may be given (as currentMode) or inferred from a labelled type, or an up/downshift
func (q *LabelType) assignUnsetModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) {
	_, unset := q.Mode.(*UnsetMode)
	if !unset {
		// If the type already has a modality, then do not change it
		// q.Mode ≠ Unset
		return
	}

	// Assign present mode (just in case the label doesn't exist)
	q.Mode = currentMode

	// Fetch labelled type
	foundLabelledType, exists := labelledTypesEnv[q.Label]
	if exists {
		// type found
		q.Mode = foundLabelledType.Mode
	}
}

func (q *UnitType) assignUnsetModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) {
	_, unset := q.Mode.(*UnsetMode)
	if !unset {
		// If the type already has a modality, then do not change it
		// q.Mode ≠ Unset
		return
	}

	q.Mode = currentMode
}

func (q *SendType) assignUnsetModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) {
	_, unset := q.Mode.(*UnsetMode)
	if unset {
		// If the type has no modality, so set it
		// q.Mode = Unset
		q.Mode = currentMode
	} else {
		// If the type already has a modality, then use it
		// q.Mode ≠ Unset
		currentMode = q.Mode
	}

	// Assign modes of inner type
	q.Left.assignUnsetModalities(labelledTypesEnv, currentMode)
	q.Right.assignUnsetModalities(labelledTypesEnv, currentMode)
}

func (q *ReceiveType) assignUnsetModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) {
	_, unset := q.Mode.(*UnsetMode)
	if unset {
		// If the type has no modality, so set it
		// q.Mode = Unset
		q.Mode = currentMode
	} else {
		// If the type already has a modality, then use it
		// q.Mode ≠ Unset
		currentMode = q.Mode
	}

	// Assign modes of inner type
	q.Left.assignUnsetModalities(labelledTypesEnv, currentMode)
	q.Right.assignUnsetModalities(labelledTypesEnv, currentMode)
}

func (q *SelectLabelType) assignUnsetModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) {
	_, unset := q.Mode.(*UnsetMode)
	if unset {
		// If the type has no modality, so set it
		// q.Mode = Unset
		q.Mode = currentMode
	} else {
		// If the type already has a modality, then use it
		// q.Mode ≠ Unset
		currentMode = q.Mode
	}

	// Assign modes of inner type
	for i := range q.Branches {
		q.Branches[i].SessionType.assignUnsetModalities(labelledTypesEnv, currentMode)
	}
}

func (q *BranchCaseType) assignUnsetModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) {
	_, unset := q.Mode.(*UnsetMode)
	if unset {
		// If the type has no modality, so set it
		// q.Mode = Unset
		q.Mode = currentMode
	} else {
		// If the type already has a modality, then use it
		// q.Mode ≠ Unset
		currentMode = q.Mode
	}

	// Assign modes of inner type
	for i := range q.Branches {
		q.Branches[i].SessionType.assignUnsetModalities(labelledTypesEnv, currentMode)
	}
}

func (q *UpType) assignUnsetModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) {
	q.Continuation.assignUnsetModalities(labelledTypesEnv, q.From)
}

func (q *DownType) assignUnsetModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) {
	q.Continuation.assignUnsetModalities(labelledTypesEnv, q.From)
}

// Deep copies a map
func copyMap(orig map[string]bool) map[string]bool {
	copy := make(map[string]bool)
	for k, v := range orig {
		copy[k] = v
	}

	return copy
}

// Takes a list of modalities, and returns the first non UnsetMode that there is.
// If all modes are Unset, then it returns Unset
func commonMode(modes ...Modality) Modality {
	commonMode := modes[0]

	for _, mode := range modes {
		_, unset := mode.(*UnsetMode)

		if !unset {
			commonMode = mode
			break
		}
	}

	return commonMode
}

///////
