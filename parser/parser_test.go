package parser

import (
	"phi/process"
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

// func TestSimpleToken(t *testing.T) {
// 	cases := []struct {
// 		input    string
// 		expected []int
// 	}{
// 		{"().;:,[]<>", []int{LPAREN, RPAREN, DOT, SEQUENCE, COLON, COMMA, LSBRACK, RSBRACK, LANGLE, RANGLE}},
// 		{"==><<-<", []int{EQUALS, RIGHT_ARROW, LANGLE, LEFT_ARROW, LANGLE}},
// 		{"send recv receive case close wait", []int{SEND, RECEIVE, RECEIVE, CASE, CLOSE, WAIT}},
// 		{"cast shift accept acc acquire acq detach det", []int{CAST, SHIFT, ACCEPT, ACCEPT, ACQUIRE, ACQUIRE, DETACH, DETACH}},
// 		{"release rel drop split push new", []int{RELEASE, RELEASE, DROP, SPLIT, PUSH, NEW}},
// 		{"snew forward fwd let in end sprc prc", []int{SNEW, FORWARD, FORWARD, LET, IN, END, SPRC, PRC}},
// 	}

// 	for _, c := range cases {
// 		reader := strings.NewReader(c.input)
// 		l := newLexer(reader)
// 		tokens := getTokens(l)
// 		compareOutput(t, tokens, c.expected)
// 	}
// }

// func TestIdentToken(t *testing.T) {
// 	cases := []struct {
// 		input    string
// 		expected []int
// 	}{
// 		{"testIdent", []int{LABEL}},
// 		{"ill\\egal", []int{LABEL}}, // kILLEGAL
// 	}

// 	for _, c := range cases {
// 		reader := strings.NewReader(c.input)
// 		l := newLexer(reader)
// 		tokens := getTokens(l)
// 		compareOutput(t, tokens, c.expected)
// 	}
// }

// func TestCommentToken(t *testing.T) {
// 	cases := []struct {
// 		input    string
// 		expected []int
// 	}{
// 		{"test/*abc*/Ident", []int{LABEL, LABEL}},
// 		{"il//egal", []int{LABEL}}, // kILLEGAL
// 	}

// 	for _, c := range cases {
// 		reader := strings.NewReader(c.input)
// 		l := newLexer(reader)
// 		tokens := getTokens(l)
// 		compareOutput(t, tokens, c.expected)
// 	}
// }
