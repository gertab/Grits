package types

import (
	"testing"
)

func TestSimpleStrings(t *testing.T) {
	label1 := NewLabelType("abc", NewUnrestrictedMode())
	label2 := NewLabelType("def", NewUnrestrictedMode())

	cases := []struct {
		input    SessionType
		expected string
	}{
		{NewUnitType(NewUnrestrictedMode()), "1"},
		{label1, "abc"},
		{NewSendType(label1, label2, NewUnrestrictedMode()), "abc * def"},
		{NewReceiveType(label1, label2, NewUnrestrictedMode()), "abc -* def"},
		{NewSelectLabelType([]Option{{Label: "a", SessionType: label1}}, NewUnrestrictedMode()), "+{a : abc}"},
		{NewBranchCaseType([]Option{{Label: "a", SessionType: label1}}, NewUnrestrictedMode()), "&{a : abc}"},
		// {NewUpType(), ""},
	}

	for i, c := range cases {
		output := c.input.String()
		if c.expected != output {
			t.Errorf("error in case #%d: Got %s, expected %s\n", i, output, c.expected)
		}
	}
}

func TestEqualType(t *testing.T) {
	label1 := NewLabelType("abc", NewUnrestrictedMode())
	label2 := NewLabelType("def", NewUnrestrictedMode())

	unit := NewUnitType(NewUnrestrictedMode())
	send := NewSendType(label1, label2, NewUnrestrictedMode())
	receive := NewReceiveType(label1, label2, NewUnrestrictedMode())
	sel_opt := []Option{{Label: "a", SessionType: label1}}
	sel := NewSelectLabelType(sel_opt, NewUnrestrictedMode())
	branch_opt := []Option{{Label: "a", SessionType: label1}, {Label: "bb", SessionType: label2}}
	branch := NewBranchCaseType(branch_opt, NewUnrestrictedMode())
	branch_opt2 := []Option{{Label: "bb", SessionType: label2}, {Label: "a", SessionType: label1}}
	branch2 := NewBranchCaseType(branch_opt2, NewUnrestrictedMode())

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
		if !EqualType(c.input, c.expected, make(LabelledTypesEnv)) {
			t.Errorf("error (EqualType) in case #%d: Got %s, expected %s\n", i, c.input.String(), c.expected.String())
		}
	}
}

func TestNotEqualType(t *testing.T) {
	label1 := NewLabelType("abc", NewUnrestrictedMode())
	label2 := NewLabelType("def", NewUnrestrictedMode())

	unit := NewUnitType(NewUnrestrictedMode())
	send := NewSendType(label1, label2, NewUnrestrictedMode())
	receive := NewReceiveType(label1, label2, NewUnrestrictedMode())
	sel_opt := []Option{{Label: "a", SessionType: label1}}
	sel := NewSelectLabelType(sel_opt, NewUnrestrictedMode())
	branch_opt := []Option{{Label: "bb", SessionType: label1}}
	branch := NewBranchCaseType(branch_opt, NewUnrestrictedMode())
	branch_opt2 := []Option{{Label: "cc", SessionType: label2}, {Label: "a", SessionType: label1}}
	branch2 := NewBranchCaseType(branch_opt2, NewUnrestrictedMode())
	branch_opt3 := []Option{{Label: "a", SessionType: label1}, {Label: "bb", SessionType: label2}}
	branch3 := NewBranchCaseType(branch_opt3, NewUnrestrictedMode())

	cases := []struct {
		input    SessionType
		expected SessionType
	}{
		{CopyType(unit), label1},
		{CopyType(label1), NewLabelType("x", NewUnrestrictedMode())},
		{CopyType(send), receive},
		{CopyType(receive), NewReceiveType(label1, NewLabelType("l", NewUnrestrictedMode()), NewUnrestrictedMode())},
		{CopyType(sel), NewSelectLabelType([]Option{{Label: "a", SessionType: label2}}, NewUnrestrictedMode())},
		{CopyType(sel), NewSelectLabelType([]Option{{Label: "a", SessionType: label1}, {Label: "b", SessionType: label2}}, NewUnrestrictedMode())},
		{CopyType(branch), NewBranchCaseType([]Option{{Label: "ff", SessionType: label1}}, NewUnrestrictedMode())},
		{CopyType(branch), NewBranchCaseType([]Option{{Label: "ff", SessionType: label1}}, NewUnrestrictedMode())},
		{branch2, branch3},
	}

	for i, c := range cases {
		if EqualType(c.input, c.expected, make(LabelledTypesEnv)) {
			t.Errorf("error (EqualType) in case #%d: Got %s, expected %s\n", i, c.input.String(), c.expected.String())
		}
	}
}

func TestCopy(t *testing.T) {
	label1 := NewLabelType("abc", NewUnrestrictedMode())
	label2 := NewLabelType("def", NewUnrestrictedMode())

	unit := NewUnitType(NewUnrestrictedMode())
	send := NewSendType(label1, label2, NewUnrestrictedMode())
	receive := NewReceiveType(label1, label2, NewUnrestrictedMode())
	sel_opt := []Option{{Label: "a", SessionType: label1}}
	sel := NewSelectLabelType(sel_opt, NewUnrestrictedMode())
	branch_opt := []Option{{Label: "bb", SessionType: label1}}
	branch := NewBranchCaseType(branch_opt, NewUnrestrictedMode())

	cases := []struct {
		input    SessionType
		expected string
	}{
		{CopyType(unit), "1"},
		{CopyType(label1), "abc"},
		{CopyType(send), "abc * def"},
		{CopyType(receive), "abc -* def"},
		{CopyType(sel), "+{a : abc}"},
		{CopyType(branch), "&{bb : abc}"},
	}

	label1 = NewLabelType("a", NewUnrestrictedMode())
	label1.Label = "new name"
	label2.Label = "new name2"

	sel_opt[0].Label = "new select branch label"
	sel_opt[0].SessionType = NewUnitType(NewUnrestrictedMode())

	branch_opt[0].Label = "new branch label"
	branch_opt[0].SessionType = NewUnitType(NewUnrestrictedMode())

	for i, c := range cases {
		output := c.input.String()
		if c.expected != output {
			t.Errorf("error in case #%d: Got %s, expected %s\n", i, output, c.expected)
		}
	}
}
