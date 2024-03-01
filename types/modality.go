package types

import (
	"fmt"
	"strings"
)

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

//           Replicable {W, C}               |
//          <            >                   |
//   Affine {W}            Multicast {C}     |
//          >            <                   |
//               Linear ø                   \/  Downshifts allowed in this direction (and vice versa for upshifts)
//
// E.g. Since Replicable > Linear, then you can downshift from Replicable to Linear (but not upshift)
// You can upshift from Affine to Linear (since Affine > Linear)

func DefaultMode() *ReplicableMode {
	return NewReplicableMode()
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
	case *ReplicableMode:
		return true
	case *MulticastMode:
		return false
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
	case *ReplicableMode:
		return true
	case *MulticastMode:
		return true
	case *AffineMode:
		return true
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

// Multicast => {C}
type MulticastMode struct{}

func NewMulticastMode() *MulticastMode {
	return &MulticastMode{}
}

func (q *MulticastMode) String() string {
	return "rep"
}

func (q *MulticastMode) Copy() Modality {
	return NewMulticastMode()
}

func (q *MulticastMode) AllowsContraction() bool {
	return true
}

func (q *MulticastMode) AllowsWeakening() bool {
	return false
}

func (q *MulticastMode) CanBeUpshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *ReplicableMode:
		return true
	case *MulticastMode:
		return true
	case *AffineMode:
		return false
	case *LinearMode:
		return false
	default:
		panic("todo")
	}
}

func (q *MulticastMode) CanBeDownshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *ReplicableMode:
		return false
	case *MulticastMode:
		return true
	case *AffineMode:
		return false // todo check with Adrian this relationship
	case *LinearMode:
		return true
	default:
		panic("todo")
	}
}

func (q *MulticastMode) Equals(other Modality) bool {
	_, same := other.(*MulticastMode)
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
	return false
}

func (q *AffineMode) AllowsWeakening() bool {
	return true
}

func (q *AffineMode) CanBeUpshiftedTo(toMode Modality) bool {
	switch interface{}(toMode).(type) {
	case *ReplicableMode:
		return true
	case *MulticastMode:
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
	case *ReplicableMode:
		return false
	case *MulticastMode:
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
	case *ReplicableMode:
		return true
	case *MulticastMode:
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
	case *ReplicableMode:
		return false
	case *MulticastMode:
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
	case "r":
		return &ReplicableMode{}
	case "rep":
		return &ReplicableMode{}
	case "replicable":
		return &ReplicableMode{}
	case "m":
		return &MulticastMode{}
	case "mul":
		return &MulticastMode{}
	case "multicast":
		return &MulticastMode{}
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

// Infer general modality of type
func AddMissingModalities(t *SessionType, labelledTypesEnv LabelledTypesEnv) {
	if t == nil {
		return
	}

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

// Ensure that a type is constructed well, with respect to modalities
// E.g. if there is an upshift/downshift, the shift should be allowed by the mode.
// Also, the only modes allowed should be Replicable/Affine/Multicast/Linear --
// UnsetModes or InvalidModes should be flagged as an error
func (q *LabelType) checkTypeModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) error {
	_, unset := q.Mode.(*UnsetMode)
	invalidMode, invalid := q.Mode.(*InvalidMode)

	if unset || q.Mode == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	typeFound, exists := labelledTypesEnv[q.Label]

	if !exists {
		// Although this should be checked already
		return fmt.Errorf("error calling undefined label type '%s'", q.String())
	}

	if !q.Mode.Equals(currentMode) {
		return fmt.Errorf("mode of label '%s' (%s) does not match the expected mode '%s'", q.String(), q.Mode.String(), currentMode.String())
	}

	if !q.Mode.Equals(typeFound.Mode) {
		return fmt.Errorf("mode of label '%s' (%s) does not match the mode '%s' (%s)", q.String(), q.Mode.String(), typeFound.Type.String(), typeFound.Mode.String())
	}

	return nil
}

func (q *UnitType) checkTypeModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) error {
	_, unset := q.Mode.(*UnsetMode)
	invalidMode, invalid := q.Mode.(*InvalidMode)

	if unset || q.Mode == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	if !q.Mode.Equals(currentMode) {
		return fmt.Errorf("mode of unit type '%s' (%s) does not match the expected mode '%s'", q.String(), q.Mode.String(), currentMode.String())
	}

	return nil
}

func (q *SendType) checkTypeModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) error {
	_, unset := q.Mode.(*UnsetMode)
	invalidMode, invalid := q.Mode.(*InvalidMode)

	if unset || q.Mode == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	if !q.Mode.Equals(currentMode) {
		return fmt.Errorf("mode of send type '%s' (%s) does not match the expected mode '%s'", q.String(), q.Mode.String(), currentMode.String())
	}

	if err := q.Left.checkTypeModalities(labelledTypesEnv, currentMode); err != nil {
		return err
	}

	if err := q.Right.checkTypeModalities(labelledTypesEnv, currentMode); err != nil {
		return err
	}

	return nil
}
func (q *ReceiveType) checkTypeModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) error {
	_, unset := q.Mode.(*UnsetMode)
	invalidMode, invalid := q.Mode.(*InvalidMode)

	if unset || q.Mode == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	if !q.Mode.Equals(currentMode) {
		return fmt.Errorf("mode of receive type '%s' (%s) does not match the expected mode '%s'", q.String(), q.Mode.String(), currentMode.String())
	}

	if err := q.Left.checkTypeModalities(labelledTypesEnv, currentMode); err != nil {
		return err
	}

	if err := q.Right.checkTypeModalities(labelledTypesEnv, currentMode); err != nil {
		return err
	}

	return nil
}

func (q *SelectLabelType) checkTypeModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) error {
	_, unset := q.Mode.(*UnsetMode)
	invalidMode, invalid := q.Mode.(*InvalidMode)

	if unset || q.Mode == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	if !q.Mode.Equals(currentMode) {
		return fmt.Errorf("mode of select type '%s' (%s) does not match the expected mode '%s'", q.String(), q.Mode.String(), currentMode.String())
	}

	for _, j := range q.Branches {
		// Checking inside each branch
		if err := j.SessionType.checkTypeModalities(labelledTypesEnv, currentMode); err != nil {
			return err
		}
	}

	return nil
}

func (q *BranchCaseType) checkTypeModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) error {
	_, unset := q.Mode.(*UnsetMode)
	invalidMode, invalid := q.Mode.(*InvalidMode)

	if unset || q.Mode == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	if !q.Mode.Equals(currentMode) {
		return fmt.Errorf("mode of the branch type '%s' (%s) does not match the expected mode '%s'", q.String(), q.Mode.String(), currentMode.String())
	}

	for _, j := range q.Branches {
		// Checking inside each branch
		if err := j.SessionType.checkTypeModalities(labelledTypesEnv, currentMode); err != nil {
			return err
		}
	}

	return nil
}

func (q *UpType) checkTypeModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) error {
	_, unset := q.From.(*UnsetMode)
	invalidMode, invalid := q.From.(*InvalidMode)

	if unset || q.From == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	_, unset = q.To.(*UnsetMode)
	invalidMode, invalid = q.To.(*InvalidMode)

	if unset || q.To == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	if !q.To.Equals(currentMode) {
		return fmt.Errorf("mode of the upshift type '%s' (%s) does not match the expected mode '%s'", q.String(), q.From.String(), currentMode.String())
	}

	if !q.From.CanBeUpshiftedTo(q.To) {
		return fmt.Errorf("mode of the upshift type '%s': mode '%s' cannot be upshifted to mode '%s'", q.String(), q.From.String(), q.To.String())
	}

	return q.Continuation.checkTypeModalities(labelledTypesEnv, q.From)
}

func (q *DownType) checkTypeModalities(labelledTypesEnv LabelledTypesEnv, currentMode Modality) error {
	_, unset := q.From.(*UnsetMode)
	invalidMode, invalid := q.From.(*InvalidMode)

	if unset || q.From == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	_, unset = q.To.(*UnsetMode)
	invalidMode, invalid = q.To.(*InvalidMode)

	if unset || q.To == nil {
		return fmt.Errorf("type '%s' has no modality defined", q.String())
	}

	if invalid {
		return fmt.Errorf("type '%s' has an unknown modality '%s'", q.String(), invalidMode.mode)
	}

	if !q.To.Equals(currentMode) {
		return fmt.Errorf("mode of the downshift type '%s' (%s) does not match the expected mode '%s'", q.String(), q.From.String(), currentMode.String())
	}

	if !q.From.CanBeDownshiftedTo(q.To) {
		return fmt.Errorf("mode of the downshift type '%s': mode '%s' cannot be downshifted to mode '%s'", q.String(), q.From.String(), q.To.String())
	}

	return q.Continuation.checkTypeModalities(labelledTypesEnv, q.From)
}
