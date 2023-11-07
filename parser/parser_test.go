package parser

import (
	"phi/process"
	"phi/types"
	"testing"
)

func compareOutputProgram(t *testing.T, got []process.Form, expected []process.Form) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if !process.EqualForm(got[index], expected[index]) {
			t.Errorf("[%d] got %s, expected %s\n", index, got[index].String(), expected[index].String())
		}
	}
}

func compareOutputTypes(t *testing.T, got []types.SessionType, expected []types.SessionType) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if !types.EqualType(got[index], expected[index]) {
			t.Errorf("[%d] got %s, expected %s\n", index, got[index].String(), expected[index].String())
		}
	}
}

func compareOutputType(t *testing.T, got types.SessionType, expected types.SessionType) bool {
	if !types.EqualType(got, expected) {
		t.Errorf("got %s, expected %s\n", got.String(), expected.String())
		return false
	}
	return true
}

func assertEqual(t *testing.T, i1, i2 process.Form) {
	if !process.EqualForm(i1, i2) {
		t.Errorf("got %s, expected %s\n", i1.String(), i2.String())
	}
}

func parseGetBody(input string) process.Form {
	body, _, err := ParseString(input)

	if err != nil {
		return nil
	}

	return body[0].Body
}

func parseGetEnvironment(input string) *process.GlobalEnvironment {
	_, globalEnv, err := ParseString(input)

	if err != nil {
		return nil
	}

	return globalEnv
}

func TestBasicForms(t *testing.T) {
	var output, expected []process.Form
	to_c := process.Name{Ident: "to_c", IsSelf: false}
	pay_c := process.Name{Ident: "pay_c", IsSelf: false}
	cont_c := process.Name{Ident: "cont_c", IsSelf: false}
	from_c := process.Name{Ident: "from_c", IsSelf: false}
	end := process.NewClose(process.Name{IsSelf: true})

	input := "send to_c<pay_c,cont_c>"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewSend(to_c, pay_c, cont_c))

	input = "<pay_c,cont_c> <- recv from_c; close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewReceive(pay_c, cont_c, from_c, end))
	i1 := parseGetBody(input)
	i2 := parseGetBody(input)
	assertEqual(t, i1, i2)

	input = "case from_c ( \n   label1<pay_c> => close self\n)"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewCase(from_c, []*process.BranchForm{process.NewBranch(process.Label{L: "label1"}, pay_c, end)}))

	input = "case from_c ( \n   label1<pay_c> => close self\n | label2<pay_c> => close self\n)"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewCase(from_c, []*process.BranchForm{process.NewBranch(process.Label{L: "label1"}, pay_c, end), process.NewBranch(process.Label{L: "label2"}, pay_c, end)}))

	input = "to_c.label1<cont_c>"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewSelect(to_c, process.Label{L: "label1"}, cont_c))

	input = "cont_c <- +new (close self); close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewNew(cont_c, end, end, process.POSITIVE))

	input = "close from_c"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewClose(from_c))

	input = "+fwd to_c from_c"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewForward(to_c, from_c, process.POSITIVE))

	input = "<pay_c,cont_c> <- +split from_c; close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewSplit(pay_c, cont_c, from_c, end, process.POSITIVE))

	compareOutputProgram(t, output, expected)
}
func TestSimpleTypes(t *testing.T) {
	unitType := types.NewUnitType()
	labelType1 := types.NewLabelType("abc")
	labelType2 := types.NewLabelType("def")

	cases := []struct {
		input    string
		expected types.SessionType
	}{
		{"type A = 1", unitType},
		{"type A = abc", labelType1},
		{"type B = abc * def", types.NewSendType(labelType1, labelType2)},
		{"type C = abc -o def", types.NewReceiveType(labelType1, labelType2)},
		{"type C = +{a : abc}", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", labelType1)})},
		{"type D = &{a : abc}", types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", labelType1)})},
	}

	for i, c := range cases {
		output := *parseGetEnvironment(c.input).Types
		outputST := output[0].SessionType
		if !compareOutputType(t, outputST, c.expected) {
			t.Errorf("error in case #%d\n", i)
		}
	}
}

func TestTypes(t *testing.T) {
	unitType := types.NewUnitType()
	labelType1 := types.NewLabelType("abc")
	labelType2 := types.NewLabelType("def")

	cases := []struct {
		input    string
		expected types.SessionType
	}{
		{"type A = (1)", unitType},
		{"type A = (abc)", labelType1},
		{"type B = (abc * def)", types.NewSendType(labelType1, labelType2)},
		{"type C = (abc -o def)", types.NewReceiveType(labelType1, labelType2)},
		{"type C = +{a : abc, bb : def}", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})},
		{"type D = &{a : abc, bb : def}", types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})},
		{"type D = &{a : +{a : abc, bb : def}, bb : def}", types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})), *types.NewBranchOption("bb", labelType2)})},
		{"type D = &{a : +{a : abc, bb : def}, bb : def}", types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})), *types.NewBranchOption("bb", labelType2)})},
		{"type E = (abc -o (abc -o def))", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2))},
		{"type E = abc -o (abc -o def)", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2))},
		{"type E = abc -o abc -o def", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2))},
		{"type E = +{a : (abc -o (abc -o def)), bb : def}", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2))), *types.NewBranchOption("bb", labelType2)})},
		{"type E = +{a : (abc -o (abc -o &{a : abc, bb : def})), bb : def}", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})))), *types.NewBranchOption("bb", labelType2)})},
	}

	for i, c := range cases {
		output := *parseGetEnvironment(c.input).Types
		outputST := output[0].SessionType
		if !compareOutputType(t, outputST, c.expected) {
			t.Errorf("error in case #%d\n", i)
		}
	}
}

func TestSimpleTypesStrings(t *testing.T) {

	cases := []struct {
		input    string
		expected string
	}{
		{"type A = 1", "1"},
		{"type A = abc", "abc"},
		{"type B = abc * def", "abc * def"},
		{"type C = abc -o def", "abc -o def"},
		{"type C = +{a : abc}", "+{a : abc}"},
		{"type D = &{a : abc}", "&{a : abc}"},
		{"type E = +{a : (abc -o (abc -o &{a : abc, bb : def})), bb : def}", "+{a : abc -o abc -o &{a : abc, bb : def}, bb : def}"},
	}

	for i, c := range cases {
		output := *parseGetEnvironment(c.input).Types
		outputST := output[0].SessionType
		if c.expected != outputST.String() {
			t.Errorf("error in case #%d: Got %s, expected %s\n", i, outputST.String(), c.expected)
		}
	}
}
