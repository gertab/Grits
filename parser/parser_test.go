package parser

import (
	"grits/process"
	"grits/types"
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
		t.Errorf("got %s, expected %s\n", got.StringWithModality(), expected.StringWithModality())
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

	input = "cont_c <- new (close self); close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewNew(cont_c, end, end))

	input = "close from_c"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewClose(from_c))

	input = "fwd to_c from_c"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewForward(to_c, from_c))

	input = "<pay_c,cont_c> <- split from_c; close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewSplit(pay_c, cont_c, from_c, end))

	input = "wait from_c; close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewWait(from_c, end))

	input = "drop from_c; close self"
	output = append(output, parseGetBody(input))
	expected = append(expected, process.NewDrop(from_c, end))

	compareOutputProgram(t, output, expected)
}
func TestSimpleTypes(t *testing.T) {
	u := types.NewReplicableMode()

	unitType := types.NewUnitType(u)
	labelType1 := types.NewLabelType("abc", u)
	labelType2 := types.NewLabelType("def", u)

	cases := []struct {
		input    string
		expected types.SessionType
	}{
		{"type A = 1", unitType},
		{"type A = abc", labelType1},
		{"type B = abc * def", types.NewSendType(labelType1, labelType2, u)},
		{"type C = abc -* def", types.NewReceiveType(labelType1, labelType2, u)},
		{"type C = +{a : abc}", types.NewSelectLabelType([]types.Option{*types.NewOption("a", labelType1)}, u)},
		{"type D = &{a : abc}", types.NewBranchCaseType([]types.Option{*types.NewOption("a", labelType1)}, u)},
	}

	for i, c := range cases {
		globalEnv := parseGetEnvironment(c.input)
		typeDefs := *globalEnv.Types
		outputST := typeDefs[0].SessionType
		if !compareOutputType(t, outputST, c.expected, make(types.LabelledTypesEnv)) {
			t.Errorf("error in case #%d\n", i)
		}
	}
}

func TestTypes(t *testing.T) {
	u := types.NewReplicableMode()

	unitType := types.NewUnitType(u)
	labelType1 := types.NewLabelType("abc", u)
	labelType2 := types.NewLabelType("def", u)

	cases := []struct {
		input    string
		expected types.SessionType
	}{
		{"type A = (1)", unitType},
		{"type A = (abc)", labelType1},
		{"type B = (abc * def)", types.NewSendType(labelType1, labelType2, u)},
		{"type C = (abc -* def)", types.NewReceiveType(labelType1, labelType2, u)},
		{"type C = +{a : abc, bb : def}", types.NewSelectLabelType([]types.Option{*types.NewOption("a", labelType1), *types.NewOption("bb", labelType2)}, u)},
		{"type D = &{a : abc, bb : def}", types.NewBranchCaseType([]types.Option{*types.NewOption("a", labelType1), *types.NewOption("bb", labelType2)}, u)},
		{"type D = &{a : +{a : abc, bb : def}, bb : def}", types.NewBranchCaseType([]types.Option{*types.NewOption("a", types.NewSelectLabelType([]types.Option{*types.NewOption("a", labelType1), *types.NewOption("bb", labelType2)}, u)), *types.NewOption("bb", labelType2)}, u)},
		{"type D = &{a : +{a : abc, bb : def}, bb : def}", types.NewBranchCaseType([]types.Option{*types.NewOption("a", types.NewSelectLabelType([]types.Option{*types.NewOption("a", labelType1), *types.NewOption("bb", labelType2)}, u)), *types.NewOption("bb", labelType2)}, u)},
		{"type E = (abc -* (abc -* def))", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2, u), u)},
		{"type E = abc -* (abc -* def)", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2, u), u)},
		{"type E = abc -* abc -* def", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2, u), u)},
		{"type E = +{a : (abc -* (abc -* def)), bb : def}", types.NewSelectLabelType([]types.Option{*types.NewOption("a", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, labelType2, u), u)), *types.NewOption("bb", labelType2)}, u)},
		{"type E = +{a : (abc -* (abc -* &{a : abc, bb : def})), bb : def}", types.NewSelectLabelType([]types.Option{*types.NewOption("a", types.NewReceiveType(labelType1, types.NewReceiveType(labelType1, types.NewBranchCaseType([]types.Option{*types.NewOption("a", labelType1), *types.NewOption("bb", labelType2)}, u), u), u)), *types.NewOption("bb", labelType2)}, u)},
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
		type F = E
		//type G = linear +{a : 1}`

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
		// {"linear +{a : 1}", "G"},
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
		type D = 1 * C
		//type E = linear +{a : 1}`

	cases := []struct {
		input1 string
		input2 string
	}{
		{"abc", "dd"},
		{"1 * +{a : D}", "C"},
		{"1 * +{a : C}", "D"},
		// {"affine +{a : 1}", "E"},
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

func TestSimpleFunctionDefinitionsWithoutTypechecking(t *testing.T) {

	cases := []struct {
		input          string
		expectedNumber int
	}{
		{"type A = 1", 0},
		{"prc[a] : 1 = close self", 0},
		{"let f() : 1 = close self", 1},
		{"let f(a, b, c) : 1 = close self", 1},
		{"let f(a, b, c) = close self", 1},
		{"let f(a : 1, b : 1, c : 1) : 1 = close self", 1},
		{"let f(a : 1) : 1 = wait a; close self", 1},
		{`type A = 1
		let f(a : 1) : A = drop a; close self`, 1},
	}

	for i, c := range cases {
		// processes, assumedFreeNames, globalEnv, err := ParseString(c.input)
		_, _, globalEnv, err := ParseString(c.input)

		if err != nil {
			t.Errorf("compilation error in case #%d: %s\n", i, err.Error())
		}

		// err = process.Typecheck(processes, assumedFreeNames, globalEnv)

		// if err != nil {
		// 	t.Errorf("type error in case #%d: %s\n", i, err.Error())
		// }

		if len(*globalEnv.FunctionDefinitions) != c.expectedNumber {
			t.Errorf("error in case #%d: Got %d, expected %d\n", i, len(*globalEnv.FunctionDefinitions), c.expectedNumber)
		}
	}
}

func TestProcessesWithoutTypechecking(t *testing.T) {
	cases := []struct {
		input                   string
		expectedNumberProcesses int
		expectedNumberFreeNames int
	}{
		{"type A = 1", 0, 0},
		{"prc[a] = close self", 1, 0},
		{"prc[a] : 1 * 1 = close self", 1, 0},
		{"prc[a, b, c] : 1 * 1 = close self", 1, 0},
		{"prc[a, b, c] = close self", 1, 0},
		{`assuming x : 1
		  prc[a, b, c] = close self`, 1, 1},
		{`assuming x, y, z
		  prc[a, b, c] = close self`, 1, 3},
		{`assuming x, y : 1 * 1, z
		  prc[a, b, c] = close self`, 1, 3},
	}

	for i, c := range cases {
		processes, assumedFreeNames, _, err := ParseString(c.input)

		if err != nil {
			t.Errorf("compilation error in case #%d: %s\n", i, err.Error())
		}

		for i := range processes {

			if len(processes) != c.expectedNumberProcesses {
				t.Errorf("error in case #%d: Got %d processes, expected %d\n", i, len(processes), c.expectedNumberProcesses)
			}

			if len(assumedFreeNames) != c.expectedNumberFreeNames {
				t.Errorf("error in case #%d: Got %d free names, expected %d\n", i, len(processes), c.expectedNumberFreeNames)
			}
		}
	}
}

// Modalities

func TestTypeDefinitionModes(t *testing.T) {

	inputProgram :=
		`type A = 1 -* 1
		 type B = 1
		 type C = linear 1 * 1
		 type D = &{a : 1, b : 1}
		 type E = +{a : 1, b : 1}
		 type F = &{a : A}
		 type G = +{a : D}
		 type H = 1 * C
		 type I = J
		 type J = I
		 type I2 = J2
		 type J2 = affine I2
		 type K = linear +{a : 1}
		 type L = linear 1 -* 1
		 type M = affine 1 -* 1
		 type N = multicast 1 -* 1
		 type O = replicable 1 -* 1
		 type P = 1 -* (1 * +{ab : 1, cd : 1})
		 type Q = affine 1 -* (1 * +{ab : 1, cd : 1})
		 type R =  1 * (1 * linear /\ affine +{ab : 1, cd : 1})
		 type S =  1 -* (1 * affine /\ linear &{ab : 1, cd : 1})`

	cases := []struct {
		input1 string
	}{
		{"[rep]1 [rep]-* [rep]1"},        // A
		{"[rep]1"},                       // B
		{"[lin]1 [lin]* [lin]1"},         // C
		{"rep&{a : [rep]1, b : [rep]1}"}, // D
		{"rep+{a : [rep]1, b : [rep]1}"}, // E
		{"rep&{a : [rep]A}"},             // F
		{"rep+{a : [rep]D}"},             // G
		{"[lin]1 [lin]* [lin]C"},         // H
		{"[rep]J"},                       // I
		{"[rep]I"},                       // J
		{"[aff]J2"},                      // I2
		{"[aff]I2"},                      // J2
		{"lin+{a : [lin]1}"},             // K
		{"[lin]1 [lin]-* [lin]1"},        // L
		{"[aff]1 [aff]-* [aff]1"},        // M
		{"[mul]1 [mul]-* [mul]1"},        // N
		{"[rep]1 [rep]-* [rep]1"},        // O
		{"[rep]1 [rep]-* [rep]1 [rep]* rep+{ab : [rep]1, cd : [rep]1}"},          // P
		{"[aff]1 [aff]-* [aff]1 [aff]* aff+{ab : [aff]1, cd : [aff]1}"},          // Q
		{`[aff]1 [aff]* [aff]1 [aff]* lin/\aff lin+{ab : [lin]1, cd : [lin]1}`},  // R
		{`[lin]1 [lin]-* [lin]1 [lin]* aff/\lin aff&{ab : [aff]1, cd : [aff]1}`}, // S
	}

	sessionTypeDefinitions := *parseGetEnvironment(inputProgram).Types
	// labelledTypesEnv := types.ProduceLabelledSessionTypeEnvironment(sessionTypeDefinitions)

	if len(sessionTypeDefinitions) != len(cases) {
		t.Errorf("number of cases do not match with the type definitions\n")
	}
	for i, c := range sessionTypeDefinitions {
		if c.SessionType.StringWithModality() != cases[i].input1 {
			t.Errorf("error in case #%d: Got %s, but expected %s\n", i, c.SessionType.StringWithModality(), cases[i].input1)
		}
	}
}
