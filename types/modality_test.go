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

func TestCheckModality(t *testing.T) {
	// var stI SessionTypeInitial
	// stI = NewUnitTypeInitial()

	unsetMode := NewUnsetMode()
	linearMode := NewLinearMode()
	affineMode := NewAffineMode()
	unrestrictedMode := NewUnrestrictedMode()
	replicableMode := NewReplicableMode()

	unitType := NewUnitTypeInitial()
	labelType := NewLabelTypeInitial("A")
	sendType := NewSendTypeInitial(labelType, unitType)
	receiveType := NewReceiveTypeInitial(labelType, unitType)
	selectType := NewSelectLabelTypeInitial([]OptionInitial{*NewOptionInitial("opt1", unitType)})
	branchType := NewBranchCaseTypeInitial([]OptionInitial{*NewOptionInitial("opt1", unitType)})
	upType := NewUpTypeInitial(linearMode, affineMode, unitType)
	downType := NewDownTypeInitial(unrestrictedMode, affineMode, unitType)

	unitNormalType := NewUnitType(unsetMode)
	labelNormalType := NewLabelType("A", unsetMode)
	sendNormalType := NewSendType(labelNormalType, unitNormalType, unsetMode)
	receiveNormalType := NewReceiveType(labelNormalType, unitNormalType, unsetMode)
	selectNormalType := NewSelectLabelType([]Option{*NewOption("opt1", unitNormalType)}, unsetMode)
	branchNormalType := NewBranchCaseType([]Option{*NewOption("opt1", unitNormalType)}, unsetMode)
	upNormalType := NewUpType(linearMode, affineMode, unitNormalType)
	downNormalType := NewDownType(unrestrictedMode, affineMode, unitNormalType)

	cases := []struct {
		stI              SessionTypeInitial
		expectedSt       SessionType
		expectedModality Modality
	}{
		{NewExplicitModeTypeInitial(linearMode, unitType), unitNormalType, linearMode},
		{NewExplicitModeTypeInitial(affineMode, unitType), unitNormalType, affineMode},
		{NewExplicitModeTypeInitial(unrestrictedMode, unitType), unitNormalType, unrestrictedMode},
		{NewExplicitModeTypeInitial(replicableMode, unitType), unitNormalType, replicableMode},
		{labelType, labelNormalType, unsetMode},
		{unitType, unitNormalType, unsetMode},
		{sendType, sendNormalType, unsetMode},
		{receiveType, receiveNormalType, unsetMode},
		{selectType, selectNormalType, unsetMode},
		{branchType, branchNormalType, unsetMode},
		{upType, upNormalType, affineMode},
		{downType, downNormalType, affineMode},
	}

	for i, c := range cases {

		// Convert to normalized session type
		st := ConvertSessionTypeInitialToSessionType(cases[i].stI)

		if st.String() != cases[i].expectedSt.String() {
			t.Errorf("error (in comparison of st) in case #%d: Got %s, expected %s\n", i, st.String(), cases[i].expectedSt.String())
		}

		if !st.Modality().Equals(cases[i].expectedModality) {
			t.Errorf("error (Modality equal) in case #%d: Got %s, expected %s\n", i, st.Modality().String(), c.expectedModality.String())
		}
	}
}
