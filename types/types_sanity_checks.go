package types

import "fmt"

func SanityChecksTypeDefinitions(typesDefs []SessionTypeDefinition) error {
	// Check for redeclaration of the same name
	typeDefNames := make(map[string]bool)
	for _, j := range typesDefs {
		_, exists := typeDefNames[j.Name]

		// Check for duplicate name
		if exists {
			return fmt.Errorf("error redefinition of the same type called '%s'", j.Name)
		}

		typeDefNames[j.Name] = true
	}

	// Mapping of labels to their session type
	labelledTypesEnv := ProduceLabelledSessionTypeEnvironment(typesDefs)

	// Check that all labelled reference point to a defined type
	for _, j := range typesDefs {
		err := CheckTypeWellFormedness(j.SessionType, labelledTypesEnv)

		if err != nil {
			return err
		}
	}

	// Ensures that the labelled types are contractive
	// For example, the following definition is not allowed:
	// -> type C = D
	// -> type D = E
	// -> type E = C
	for _, j := range typesDefs {
		ok := j.SessionType.isContractive(labelledTypesEnv, make(map[string]bool))

		if !ok {
			return fmt.Errorf("session type definition for %s (= %s) is not contractive", j.Name, j.SessionType.String())
		}

		err := CheckTypeWellFormedness(j.SessionType, labelledTypesEnv)

		if err != nil {
			return err
		}
	}

	return nil
}

// Performs similar check to the preceding function,
// however, this deals with plain types directly (not labelled)
func SanityChecksType(types []SessionType, typesDefs []SessionTypeDefinition) error {
	// Check that all labelled reference point to a defined type

	labelledTypesEnv := ProduceLabelledSessionTypeEnvironment(typesDefs)

	for _, j := range types {
		err := CheckTypeWellFormedness(j, labelledTypesEnv)

		if err != nil {
			return err
		}
	}

	return nil
}

// Run some checks on the construction of the type
func CheckTypeWellFormedness(t SessionType, labelledTypesEnv LabelledTypesEnv) error {
	err := t.checkTypeLabels(labelledTypesEnv)
	if err != nil {
		return err
	}

	err = t.checkTypeModalities(labelledTypesEnv, t.Modality())

	return err
}

// Check whether a reference to a session type label exists
//
// Example:
//
//	type A = 1			[correct]
//	type B = A -* 1  	[correct]
//	type C = A -* D		[incorrect, because D is undefined]
//
// Ensures also the branches are made up of unique labels
func (q *LabelType) checkTypeLabels(labelledTypesEnv LabelledTypesEnv) error {
	if !LabelledTypedExists(labelledTypesEnv, q.Label) {
		return fmt.Errorf("error calling undefined label type '%s'", q.String())
	}

	return nil
}
func (q *UnitType) checkTypeLabels(labelledTypesEnv LabelledTypesEnv) error {
	return nil
}
func (q *SendType) checkTypeLabels(labelledTypesEnv LabelledTypesEnv) error {
	err := q.Left.checkTypeLabels(labelledTypesEnv)

	if err != nil {
		return err
	}

	err = q.Right.checkTypeLabels(labelledTypesEnv)

	if err != nil {
		return err
	}

	return nil
}
func (q *ReceiveType) checkTypeLabels(labelledTypesEnv LabelledTypesEnv) error {
	err := q.Left.checkTypeLabels(labelledTypesEnv)

	if err != nil {
		return err
	}

	err = q.Right.checkTypeLabels(labelledTypesEnv)

	if err != nil {
		return err
	}

	return nil
}
func (q *SelectLabelType) checkTypeLabels(labelledTypesEnv LabelledTypesEnv) error {

	existingLabels := make(map[string]bool)

	for _, j := range q.Branches {
		// Check for unique labels
		_, exists := existingLabels[j.Label]

		if exists {
			return fmt.Errorf("duplicate label '%s' found in type '%s'", j.Label, q.String())
		}

		// Checking inside the branch
		err := j.SessionType.checkTypeLabels(labelledTypesEnv)

		if err != nil {
			return err
		}
	}

	return nil
}
func (q *BranchCaseType) checkTypeLabels(labelledTypesEnv LabelledTypesEnv) error {
	existingLabels := make(map[string]bool)

	for _, j := range q.Branches {
		// Check for unique labels
		_, exists := existingLabels[j.Label]

		if exists {
			return fmt.Errorf("duplicate label '%s' found in type '%s'", j.Label, q.String())
		}

		// Checking inside the branch
		err := j.SessionType.checkTypeLabels(labelledTypesEnv)

		if err != nil {
			return err
		}
	}

	return nil
}

func (q *UpType) checkTypeLabels(labelledTypesEnv LabelledTypesEnv) error {
	return q.Continuation.checkTypeLabels(labelledTypesEnv)
}

func (q *DownType) checkTypeLabels(labelledTypesEnv LabelledTypesEnv) error {
	return q.Continuation.checkTypeLabels(labelledTypesEnv)
}

// Ensures that the labelled types are contractive
// For example, the following definitions is not allowed:
// -> type C = D
// -> type D = E
// -> type E = C
// where the the D is non-contractive
func (q *LabelType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {

	presentSnapshot := q.String()

	// Cycle reached, so type is not contractive
	_, exists := snapshots[presentSnapshot]
	if exists {
		return false
	}

	snapshots[q.Label] = true

	// This succeeds since we already checked that all labels map to some type
	unfoldedType := labelledTypesEnv[q.Label].Type

	return unfoldedType.isContractive(labelledTypesEnv, snapshots)
}

func (q *UnitType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {
	return true
}

func (q *SendType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {
	return true
}

func (q *ReceiveType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {
	return true
}

func (q *SelectLabelType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {
	return true
}

func (q *BranchCaseType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {
	return true
}

func (q *UpType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {
	// not entirely sure about shifting
	return true
}

func (q *DownType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {
	return true
}
