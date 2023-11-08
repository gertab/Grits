package types

import "fmt"

func SanityChecksTypeDefinitions(typesDefs []SessionTypeDefinition) error {
	for _, j := range typesDefs {
		// Avoid labelled types at the top level, e.g. type A = B
		switch j.SessionType.(type) {
		case *LabelType:
			return fmt.Errorf("session type definition for %s cannot use contractive type '%s'", j.Name, j.SessionType.String())
		}
	}

	// Check that all labelled reference point to a defined type
	typeDefNames := make(map[string]bool)
	for _, j := range typesDefs {
		_, exists := typeDefNames[j.Name] // check for existence

		if exists {
			return fmt.Errorf("error redefinition of the same type called '%s'", j.Name)
		} else {
			// add to our map
			typeDefNames[j.Name] = true
		}

	}

	for _, j := range typesDefs {
		err := j.SessionType.checkLabelledTypes(typeDefNames)

		if err != nil {
			return err
		}
	}

	return nil
}

// Performs similar check to the preceding function
func SanityChecksType(types []SessionType, typesDefs []SessionTypeDefinition) error {
	// Check that all labelled reference point to a defined type
	typeDefNames := make(map[string]bool)
	for _, j := range typesDefs {
		typeDefNames[j.Name] = true
	}

	for _, j := range types {
		err := j.checkLabelledTypes(typeDefNames)

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
func (q *LabelType) checkLabelledTypes(typeDefNames map[string]bool) error {
	_, exists := typeDefNames[q.Label]

	if !exists {
		return fmt.Errorf("error calling undefined label type '%s'", q.String())
	}

	return nil
}
func (q *WIPType) checkLabelledTypes(typeDefNames map[string]bool) error {
	return nil
}
func (q *UnitType) checkLabelledTypes(typeDefNames map[string]bool) error {
	return nil
}
func (q *SendType) checkLabelledTypes(typeDefNames map[string]bool) error {
	err := q.Left.checkLabelledTypes(typeDefNames)

	if err != nil {
		return err
	}

	err = q.Right.checkLabelledTypes(typeDefNames)

	if err != nil {
		return err
	}

	return nil
}
func (q *ReceiveType) checkLabelledTypes(typeDefNames map[string]bool) error {
	err := q.Left.checkLabelledTypes(typeDefNames)

	if err != nil {
		return err
	}

	err = q.Right.checkLabelledTypes(typeDefNames)

	if err != nil {
		return err
	}

	return nil
}
func (q *SelectLabelType) checkLabelledTypes(typeDefNames map[string]bool) error {
	for _, j := range q.Branches {
		err := j.Session_type.checkLabelledTypes(typeDefNames)

		if err != nil {
			return err
		}
	}

	return nil
}
func (q *BranchCaseType) checkLabelledTypes(typeDefNames map[string]bool) error {
	for _, j := range q.Branches {
		err := j.Session_type.checkLabelledTypes(typeDefNames)

		if err != nil {
			return err
		}
	}

	return nil
}

// func (q *LabelType) String() string {}
// func (q *WIPType) String() string {}
// func (q *UnitType) String() string {}
// func (q *SendType) String() string {}
// func (q *ReceiveType) String() string {}
// func (q *SelectLabelType) String() string {}
// func (q *BranchCaseType) String() string {}
