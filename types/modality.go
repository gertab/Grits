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

// Unrestricted => {W, C}
type UnrestrictedMode struct{}

func NewUnrestrictedMode() *UnrestrictedMode {
	return &UnrestrictedMode{}
}

func (q *UnrestrictedMode) String() string {
	return "U"
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
	return "R"
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
	return "A"
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
	return "L"
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
	return q.mode
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
	return ""
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
