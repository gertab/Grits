package types

import (
	"testing"
)

func TestSimpleStrings(t *testing.T) {
	label1 := NewLabelType("abc", NewReplicableMode())
	label2 := NewLabelType("def", NewReplicableMode())

	cases := []struct {
		input    SessionType
		expected string
	}{
		{NewUnitType(NewReplicableMode()), "1"},
		{label1, "abc"},
		{NewSendType(label1, label2, NewReplicableMode()), "abc * def"},
		{NewReceiveType(label1, label2, NewReplicableMode()), "abc -* def"},
		{NewSelectLabelType([]Option{{Label: "a", SessionType: label1}}, NewReplicableMode()), "+{a : abc}"},
		{NewBranchCaseType([]Option{{Label: "a", SessionType: label1}}, NewReplicableMode()), "&{a : abc}"},
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
	label1 := NewLabelType("abc", NewReplicableMode())
	label2 := NewLabelType("def", NewReplicableMode())

	unit := NewUnitType(NewReplicableMode())
	send := NewSendType(label1, label2, NewReplicableMode())
	receive := NewReceiveType(label1, label2, NewReplicableMode())
	sel_opt := []Option{{Label: "a", SessionType: label1}}
	sel := NewSelectLabelType(sel_opt, NewReplicableMode())
	branch_opt := []Option{{Label: "a", SessionType: label1}, {Label: "bb", SessionType: label2}}
	branch := NewBranchCaseType(branch_opt, NewReplicableMode())
	branch_opt2 := []Option{{Label: "bb", SessionType: label2}, {Label: "a", SessionType: label1}}
	branch2 := NewBranchCaseType(branch_opt2, NewReplicableMode())

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
	label1 := NewLabelType("abc", NewReplicableMode())
	label2 := NewLabelType("def", NewReplicableMode())

	unit := NewUnitType(NewReplicableMode())
	send := NewSendType(label1, label2, NewReplicableMode())
	receive := NewReceiveType(label1, label2, NewReplicableMode())
	sel_opt := []Option{{Label: "a", SessionType: label1}}
	sel := NewSelectLabelType(sel_opt, NewReplicableMode())
	branch_opt := []Option{{Label: "bb", SessionType: label1}}
	branch := NewBranchCaseType(branch_opt, NewReplicableMode())
	branch_opt2 := []Option{{Label: "cc", SessionType: label2}, {Label: "a", SessionType: label1}}
	branch2 := NewBranchCaseType(branch_opt2, NewReplicableMode())
	branch_opt3 := []Option{{Label: "a", SessionType: label1}, {Label: "bb", SessionType: label2}}
	branch3 := NewBranchCaseType(branch_opt3, NewReplicableMode())

	cases := []struct {
		input    SessionType
		expected SessionType
	}{
		{CopyType(unit), label1},
		{CopyType(label1), NewLabelType("x", NewReplicableMode())},
		{CopyType(send), receive},
		{CopyType(receive), NewReceiveType(label1, NewLabelType("l", NewReplicableMode()), NewReplicableMode())},
		{CopyType(sel), NewSelectLabelType([]Option{{Label: "a", SessionType: label2}}, NewReplicableMode())},
		{CopyType(sel), NewSelectLabelType([]Option{{Label: "a", SessionType: label1}, {Label: "b", SessionType: label2}}, NewReplicableMode())},
		{CopyType(branch), NewBranchCaseType([]Option{{Label: "ff", SessionType: label1}}, NewReplicableMode())},
		{CopyType(branch), NewBranchCaseType([]Option{{Label: "ff", SessionType: label1}}, NewReplicableMode())},
		{branch2, branch3},
	}

	for i, c := range cases {
		if EqualType(c.input, c.expected, make(LabelledTypesEnv)) {
			t.Errorf("error (EqualType) in case #%d: Got %s, expected %s\n", i, c.input.String(), c.expected.String())
		}
	}
}

func TestCopy(t *testing.T) {
	label1 := NewLabelType("abc", NewReplicableMode())
	label2 := NewLabelType("def", NewReplicableMode())

	unit := NewUnitType(NewReplicableMode())
	send := NewSendType(label1, label2, NewReplicableMode())
	receive := NewReceiveType(label1, label2, NewReplicableMode())
	sel_opt := []Option{{Label: "a", SessionType: label1}}
	sel := NewSelectLabelType(sel_opt, NewReplicableMode())
	branch_opt := []Option{{Label: "bb", SessionType: label1}}
	branch := NewBranchCaseType(branch_opt, NewReplicableMode())

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

	label1 = NewLabelType("a", NewReplicableMode())
	label1.Label = "new name"
	label2.Label = "new name2"

	sel_opt[0].Label = "new select branch label"
	sel_opt[0].SessionType = NewUnitType(NewReplicableMode())

	branch_opt[0].Label = "new branch label"
	branch_opt[0].SessionType = NewUnitType(NewReplicableMode())

	for i, c := range cases {
		output := c.input.String()
		if c.expected != output {
			t.Errorf("error in case #%d: Got %s, expected %s\n", i, output, c.expected)
		}
	}
}
