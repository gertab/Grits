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

func compareOutputType(t *testing.T, got types.SessionType, expected types.SessionType, labelledTypesEnv types.LabelledTypesEnv) bool {
	if !types.EqualType(got, expected, labelledTypesEnv) {
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
	body, _, _, err := ParseString(input)

	if err != nil {
		return nil
	}

	return body[0].Body
}

func parseGetEnvironment(input string) *process.GlobalEnvironment {
	_, _, globalEnv, err := ParseString(input)

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

	input = "wait from_c; close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewWait(from_c, end))

	input = "drop from_c; close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewDrop(from_c, end))

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
		{"type C = abc -* def", types.NewReceiveType(labelType1, labelType2)},
		{"type C = +{a : abc}", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", labelType1)})},
		{"type D = &{a : abc}", types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", labelType1)})},
	}

	for i, c := range cases {
		output := *parseGetEnvironment(c.input).Types
		outputST := output[0].SessionType
		if !compareOutputType(t, outputST, c.expected, make(types.LabelledTypesEnv)) {
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
		{"type C = (abc -* def)", types.NewReceiveType(labelType1, labelType2)},
		{"type C = +{a : abc, bb : def}", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})},
		{"type D = &{a : abc, bb : def}", types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})},
		{"type D = &{a : +{a : abc, bb : def}, bb : def}", types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})), *types.NewBranchOption("bb", labelType2)})},
		{"type D = &{a : +{a : abc, bb : def}, bb : def}", types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})), *types.NewBranchOption("bb", labelType2)})},
		{"type E = (abc -* (abc -* def))", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2))},
		{"type E = abc -* (abc -* def)", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2))},
		{"type E = abc -* abc -* def", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2))},
		{"type E = +{a : (abc -* (abc -* def)), bb : def}", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2))), *types.NewBranchOption("bb", labelType2)})},
		{"type E = +{a : (abc -* (abc -* &{a : abc, bb : def})), bb : def}", types.NewSelectType([]types.BranchOption{*types.NewBranchOption("a", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, types.NewBranchCaseType([]types.BranchOption{*types.NewBranchOption("a", labelType1), *types.NewBranchOption("bb", labelType2)})))), *types.NewBranchOption("bb", labelType2)})},
	}

	for i, c := range cases {
		output := *parseGetEnvironment(c.input).Types
		outputST := output[0].SessionType
		if !compareOutputType(t, outputST, c.expected, make(types.LabelledTypesEnv)) {
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
		{"type C = abc -* def", "abc -* def"},
		{"type C = +{a : abc}", "+{a : abc}"},
		{"type D = &{a : abc}", "&{a : abc}"},
		{"type E = +{a : (abc -* (abc -* &{a : abc, bb : def})), bb : def}", "+{a : abc -* abc -* &{a : abc, bb : def}, bb : def}"},
	}

	for i, c := range cases {
		output := *parseGetEnvironment(c.input).Types
		outputST := output[0].SessionType
		if c.expected != outputST.String() {
			t.Errorf("error in case #%d: Got %s, expected %s\n", i, outputST.String(), c.expected)
		}
	}
}

func TestEqualType(t *testing.T) {

	commonProgram :=
		`type A = 1 -* 1
		type B = &{a : A}
		type C = +{a : D}
		type D = 1 * C
		type E = F // these should be avoided
		type F = E`

	// Ensures that input1 and input2 are equivalent
	cases := []struct {
		input1 string
		input2 string
	}{
		{"abc", "abc"},
		{"A", "A"},
		{"&{a : (1 -* 1)}", "B"},
		{"D", "1 * C"},
		{"1 * +{a : D}", "1 * C"},
		{"1 * +{a : D}", "D"},
		{"1 * +{a : 1 * +{a : 1 * +{a : 1 * +{a : 1 * +{a : D}}}}}", "D"},
		{"1 * +{a : 1 * +{a : 1 * +{a : 1 * +{a : 1 * +{a : D}}}}}", "1 * C"},
		{"E", "F"},
		{"&{a : (1 -* 1), b : 1}", "&{b : 1, a : (1 -* 1)}"},
	}

	sessionTypeDefinitions := *parseGetEnvironment(commonProgram).Types
	labelledTypesEnv := types.ProduceLabelledSessionTypeEnvironment(sessionTypeDefinitions)

	for i, c := range cases {
		input1 := *parseGetEnvironment("type X = " + c.input1).Types
		input2 := *parseGetEnvironment("type X = " + c.input2).Types
		input1ST := input1[0].SessionType
		input2ST := input2[0].SessionType
		if !compareOutputType(t, input1ST, input2ST, labelledTypesEnv) {
			t.Errorf("error in case #%d\n", i)
		}
	}
}

func TestNotEqualType(t *testing.T) {

	commonProgram :=
		`type A = 1 -* 1
		type B = &{a : A}
		type C = +{a : D}
		type D = 1 * C`

	cases := []struct {
		input1 string
		input2 string
	}{
		{"abc", "dd"},
		{"1 * +{a : D}", "C"},
		{"1 * +{a : C}", "D"},
	}

	sessionTypeDefinitions := *parseGetEnvironment(commonProgram).Types
	labelledTypesEnv := types.ProduceLabelledSessionTypeEnvironment(sessionTypeDefinitions)
	for i, c := range cases {
		input1 := *parseGetEnvironment("type X = " + c.input1).Types
		input2 := *parseGetEnvironment("type X = " + c.input2).Types
		input1ST := input1[0].SessionType
		input2ST := input2[0].SessionType

		if types.EqualType(input1ST, input2ST, labelledTypesEnv) {
			t.Errorf("error (EqualType) in case #%d: Got %s == %s, Expected %s != %s\n", i, c.input1, c.input2, c.input1, c.input2)
		}
	}
}

func TestSimpleFunctionDefinitions(t *testing.T) {

	cases := []struct {
		input          string
		expectedNumber int
	}{
		{"type A = 1", 0},
		{"prc[a] = close self", 0},
		{"let f() : 1= close self", 1},
		// {"let f(a, b, c) = close self", 1},
		// {"let f(a : 1, b : 1, c : 1) = close self", 1},
		{"let f(a : 1) : 1 = wait a; close self", 1},
		{`type A = 1
		let f(a : 1) : A = drop a; close self`, 1},
	}

	for i, c := range cases {
		processes, processesFreeNames, globalEnv, err := ParseString(c.input)

		if err != nil {
			t.Errorf("compilation error in case #%d: %s\n", i, err.Error())
		}

		err = process.Typecheck(processes, processesFreeNames, globalEnv)

		if err != nil {
			t.Errorf("type error in case #%d: %s\n", i, err.Error())
		}

		if len(*globalEnv.FunctionDefinitions) != c.expectedNumber {
			t.Errorf("error in case #%d: Got %d, expected %d\n", i, len(*globalEnv.FunctionDefinitions), c.expectedNumber)
		}
	}
}
