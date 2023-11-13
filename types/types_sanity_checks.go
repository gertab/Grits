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
	}

	// Mapping of  labels to their session type
	labelledTypesEnv := ProduceLabelledSessionTypeEnvironment(typesDefs)

	// Check that all labelled reference point to a defined type
	for _, j := range typesDefs {
		err := j.SessionType.checkLabelledTypes(labelledTypesEnv)

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

		err := j.SessionType.checkLabelledTypes(labelledTypesEnv)

		if err != nil {
			return err
		}

		// }
	}

	return nil
}

// Performs similar check to the preceding function,
// however, this deals with plain types directly (not labelled)
func SanityChecksType(types []SessionType, typesDefs []SessionTypeDefinition) error {
	// Check that all labelled reference point to a defined type

	labelledTypesEnv := ProduceLabelledSessionTypeEnvironment(typesDefs)

	for _, j := range types {
		err := j.checkLabelledTypes(labelledTypesEnv)

		if err != nil {
			return err
		}
	}

	return nil
}

// Check whether a reference to a session type label exists
// Example:
//
//	type A = 1			[correct]
//	type B = A -o 1  	[correct]
//	type C = A -o D		[incorrect, because D is undefined]
func (q *LabelType) checkLabelledTypes(labelledTypesEnv LabelledTypesEnv) error {
	if !LabelledTypedExists(labelledTypesEnv, q.Label) {
		return fmt.Errorf("error calling undefined label type '%s'", q.String())
	}

	return nil
}
func (q *WIPType) checkLabelledTypes(labelledTypesEnv LabelledTypesEnv) error {
	return nil
}
func (q *UnitType) checkLabelledTypes(labelledTypesEnv LabelledTypesEnv) error {
	return nil
}
func (q *SendType) checkLabelledTypes(labelledTypesEnv LabelledTypesEnv) error {
	err := q.Left.checkLabelledTypes(labelledTypesEnv)

	if err != nil {
		return err
	}

	err = q.Right.checkLabelledTypes(labelledTypesEnv)

	if err != nil {
		return err
	}

	return nil
}
func (q *ReceiveType) checkLabelledTypes(labelledTypesEnv LabelledTypesEnv) error {
	err := q.Left.checkLabelledTypes(labelledTypesEnv)

	if err != nil {
		return err
	}

	err = q.Right.checkLabelledTypes(labelledTypesEnv)

	if err != nil {
		return err
	}

	return nil
}
func (q *SelectLabelType) checkLabelledTypes(labelledTypesEnv LabelledTypesEnv) error {
	for _, j := range q.Branches {
		err := j.Session_type.checkLabelledTypes(labelledTypesEnv)

		if err != nil {
			return err
		}
	}

	return nil
}
func (q *BranchCaseType) checkLabelledTypes(labelledTypesEnv LabelledTypesEnv) error {
	for _, j := range q.Branches {
		err := j.Session_type.checkLabelledTypes(labelledTypesEnv)

		if err != nil {
			return err
		}
	}

	return nil
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

func (q *WIPType) isContractive(labelledTypesEnv LabelledTypesEnv, snapshots map[string]bool) bool {
	return true
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

// func (q *LabelType) String() string {}
// func (q *WIPType) String() string {}
// func (q *UnitType) String() string {}
// func (q *SendType) String() string {}
// func (q *ReceiveType) String() string {}
// func (q *SelectLabelType) String() string {}
// func (q *BranchCaseType) String() string {}
