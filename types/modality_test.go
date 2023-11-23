package types

import (
	"testing"
)

func TestModality(t *testing.T) {

	unit := NewUnitType()
	cases := []struct {
		input    SessionType
		expected SessionType
	}{
		{CopyType(unit), unit},
	}

	for i, c := range cases {
		if !EqualType(c.input, c.expected, make(LabelledTypesEnv)) {
			t.Errorf("error (EqualType) in case #%d: Got %s, expected %s\n", i, c.input.String(), c.expected.String())
		}
	}
}
