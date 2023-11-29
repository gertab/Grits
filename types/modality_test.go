package types

import (
	"testing"
)

func TestEqualModality(t *testing.T) {
	label1 := NewLabelType("abc", NewUnrestrictedMode())
	label2 := NewLabelType("def", NewUnrestrictedMode())

	unit := NewUnitType(NewAffineMode())
	send := NewSendType(label1, label2, NewLinearMode())
	receive := NewReceiveType(label1, label2, NewAffineMode())
	sel_opt := []Option{{Label: "a", SessionType: label1}}
	sel := NewSelectLabelType(sel_opt, NewLinearMode())
	branch_opt := []Option{{Label: "a", SessionType: label1}, {Label: "bb", SessionType: label2}}
	branch := NewBranchCaseType(branch_opt, NewLinearMode())
	branch_opt2 := []Option{{Label: "bb", SessionType: label2}, {Label: "a", SessionType: label1}}
	branch2 := NewBranchCaseType(branch_opt2, NewLinearMode())

	cases := []struct {
		input    SessionType
		expected SessionType
	}{
		{CopyType(unit), unit},
		{CopyType(label1), label1},
		{CopyType(send), send},
		{CopyType(receive), receive},
		{CopyType(sel), sel},
		{CopyType(branch), branch},
		{branch, branch2},
	}

	for i, c := range cases {
		if !c.input.Modality().Equals(c.expected.Modality()) {
			t.Errorf("error (Modality Equals) in case #%d: Got %s, expected %s\n", i, c.input.String(), c.expected.String())
		}

		if !EqualType(c.input, c.expected, make(LabelledTypesEnv)) {
			t.Errorf("error (EqualType) in case #%d: Got %s, expected %s\n", i, c.input.String(), c.expected.String())
		}
	}
}

func TestNotEqualModality(t *testing.T) {
	label1 := NewLabelType("abc", NewUnrestrictedMode())
	label2 := NewLabelType("def", NewUnrestrictedMode())

	unit := NewUnitType(NewAffineMode())
	unit2 := NewUnitType(NewReplicableMode())
	send := NewSendType(label1, label2, NewLinearMode())
	send2 := NewSendType(label1, label2, NewAffineMode())
	receive := NewReceiveType(label1, label2, NewAffineMode())
	receive2 := NewReceiveType(label1, label2, NewUnrestrictedMode())
	sel_opt := []Option{{Label: "a", SessionType: label1}}
	sel := NewSelectLabelType(sel_opt, NewLinearMode())
	sel2 := NewSelectLabelType(sel_opt, NewAffineMode())
	branch_opt := []Option{{Label: "a", SessionType: label1}, {Label: "bb", SessionType: label2}}
	branch := NewBranchCaseType(branch_opt, NewLinearMode())
	branch_opt2 := []Option{{Label: "bb", SessionType: label2}, {Label: "a", SessionType: label1}}
	branch2 := NewBranchCaseType(branch_opt2, NewUnrestrictedMode())

	cases := []struct {
		input    SessionType
		expected SessionType
	}{
		{unit, unit2},
		{label1, label2},
		{send, send2},
		{receive, receive2},
		{sel, sel2},
		{branch, branch2},
	}

	for i, c := range cases {
		if EqualType(c.input, c.expected, make(LabelledTypesEnv)) {
			t.Errorf("error (EqualType) in case #%d: Got %s, expected %s\n", i, c.input.String(), c.expected.String())
		}
	}
}
