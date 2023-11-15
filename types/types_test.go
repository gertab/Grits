package types

import (
	"testing"
)

func TestSimpleStrings(t *testing.T) {
	label1 := NewLabelType("abc")
	label2 := NewLabelType("def")

	cases := []struct {
		input    SessionType
		expected string
	}{
		{NewUnitType(), "1"},
		{label1, "abc"},
		{NewSendType(label1, label2), "abc * def"},
		{NewReceiveType(label1, label2), "abc -* def"},
		{NewSelectType([]BranchOption{{Label: "a", Session_type: label1}}), "+{a : abc}"},
		{NewBranchCaseType([]BranchOption{{Label: "a", Session_type: label1}}), "&{a : abc}"},
	}

	for i, c := range cases {
		output := c.input.String()
		if c.expected != output {
			t.Errorf("error in case #%d: Got %s, expected %s\n", i, output, c.expected)
		}
	}
}

func TestEqualType(t *testing.T) {
	label1 := NewLabelType("abc")
	label2 := NewLabelType("def")

	unit := NewUnitType()
	send := NewSendType(label1, label2)
	receive := NewReceiveType(label1, label2)
	sel_opt := []BranchOption{{Label: "a", Session_type: label1}}
	sel := NewSelectType(sel_opt)
	branch_opt := []BranchOption{{Label: "a", Session_type: label1}, {Label: "bb", Session_type: label2}}
	branch := NewBranchCaseType(branch_opt)
	branch_opt2 := []BranchOption{{Label: "bb", Session_type: label2}, {Label: "a", Session_type: label1}}
	branch2 := NewBranchCaseType(branch_opt2)

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
	label1 := NewLabelType("abc")
	label2 := NewLabelType("def")

	unit := NewUnitType()
	send := NewSendType(label1, label2)
	receive := NewReceiveType(label1, label2)
	sel_opt := []BranchOption{{Label: "a", Session_type: label1}}
	sel := NewSelectType(sel_opt)
	branch_opt := []BranchOption{{Label: "bb", Session_type: label1}}
	branch := NewBranchCaseType(branch_opt)
	branch_opt2 := []BranchOption{{Label: "cc", Session_type: label2}, {Label: "a", Session_type: label1}}
	branch2 := NewBranchCaseType(branch_opt2)
	branch_opt3 := []BranchOption{{Label: "a", Session_type: label1}, {Label: "bb", Session_type: label2}}
	branch3 := NewBranchCaseType(branch_opt3)

	cases := []struct {
		input    SessionType
		expected SessionType
	}{
		{CopyType(unit), label1},
		{CopyType(label1), NewLabelType("x")},
		{CopyType(send), receive},
		{CopyType(receive), NewReceiveType(label1, NewLabelType("l"))},
		{CopyType(sel), NewSelectType([]BranchOption{{Label: "a", Session_type: label2}})},
		{CopyType(sel), NewSelectType([]BranchOption{{Label: "a", Session_type: label1}, {Label: "b", Session_type: label2}})},
		{CopyType(branch), NewBranchCaseType([]BranchOption{{Label: "ff", Session_type: label1}})},
		{CopyType(branch), NewBranchCaseType([]BranchOption{{Label: "ff", Session_type: label1}})},
		{branch2, branch3},
	}

	for i, c := range cases {
		if EqualType(c.input, c.expected, make(LabelledTypesEnv)) {
			t.Errorf("error (EqualType) in case #%d: Got %s, expected %s\n", i, c.input.String(), c.expected.String())
		}
	}
}

func TestCopy(t *testing.T) {
	label1 := NewLabelType("abc")
	label2 := NewLabelType("def")

	unit := NewUnitType()
	send := NewSendType(label1, label2)
	receive := NewReceiveType(label1, label2)
	sel_opt := []BranchOption{{Label: "a", Session_type: label1}}
	sel := NewSelectType(sel_opt)
	branch_opt := []BranchOption{{Label: "bb", Session_type: label1}}
	branch := NewBranchCaseType(branch_opt)

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

	label1 = NewLabelType("a")
	label1.Label = "new name"
	label2.Label = "new name2"

	sel_opt[0].Label = "new select branch label"
	sel_opt[0].Session_type = NewUnitType()

	branch_opt[0].Label = "new branch label"
	branch_opt[0].Session_type = NewUnitType()

	for i, c := range cases {
		output := c.input.String()
		if c.expected != output {
			t.Errorf("error in case #%d: Got %s, expected %s\n", i, output, c.expected)
		}
	}
}
