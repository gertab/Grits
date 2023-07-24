package parser

import (
	"phi/process"
	"testing"
)

// type phiSymType struct {
// 	strval string
// }

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

// func compareOutputString(t *testing.T, got []string, expected []string) {
// 	if len(got) != len(expected) {
// 		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
// 		return
// 	}

// 	for index := range got {
// 		if got[index] != expected[index] {
// 			t.Errorf("got %s, expected %s\n", got[index], expected[index])
// 		}
// 	}
// }

func TestBasicForms(t *testing.T) {
	var output, expected []process.Form
	to_c := process.Name{Ident: "to_c"}
	pay_c := process.Name{Ident: "pay_c"}
	cont_c := process.Name{Ident: "cont_c"}
	from_c := process.Name{Ident: "from_c"}
	end := process.NewClose(process.Name{Ident: "self"})

	input := "send to_c<pay_c,cont_c>"
	output = append(output, ParseString(input)[0].Body)
	expected = append(expected, process.NewSend(to_c, pay_c, cont_c))

	input = "<pay_c,cont_c> <- recv from_c; close self"
	output = append(output, ParseString(input)[1].Body)
	expected = append(expected, process.NewReceive(pay_c, cont_c, from_c, end))

	// input = "case from_c ( \n   label1<pay_c> => close self\n)"
	// output = append(output, ParseString(input)[0].Body)
	// expected = append(expected, process.NewBranch(process.Label{L: "label1"}, pay_c, end))

	// input = "case from_c ( \n   label1<pay_c> => close self\n | label2<pay_c> => close self\n)"
	// output = append(output, ParseString(input)[0].Body)
	// expected = append(expected, process.NewCase(from_c, []*process.BranchForm{process.NewBranch(process.Label{L: "label1"}, pay_c, end), process.NewBranch(process.Label{L: "label2"}, pay_c, end)}))

	// input = "to_c.label1<cont_c>"
	// output = append(output, ParseString(input)[0].Body)
	// expected = append(expected, process.NewSelect(to_c, process.Label{L: "label1"}, cont_c))

	// input = "cont_c <- new (close self); close self"
	// output = append(output, ParseString(input)[0].Body)
	// expected = append(expected, process.NewNew(cont_c, end, end))

	// input = "close from_c"
	// output = append(output, ParseString(input)[0].Body)
	// expected = append(expected, process.NewClose(from_c))

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
