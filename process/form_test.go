package process

import (
	"testing"
)

// type phiSymType struct {
// 	strval string
// }

func compareOutput(t *testing.T, got []string, expected []string) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if got[index] != expected[index] {
			// t.Errorf("got %s, expected %s\n", "sa", "de")
			t.Errorf("got %s, expected %s\n", got[index], expected[index])
		}
	}
}

func TestBasicTokens(t *testing.T) {
	expected, output := []string{}, []string{}
	to_c := Name{Ident: "to_c"}
	pay_c := Name{Ident: "pay_c"}
	cont_c := Name{Ident: "cont_c"}
	from_c := Name{Ident: "from_c"}
	end := NewClose(Name{Ident: "self"})

	// Send
	input := NewSend(to_c, pay_c, cont_c)
	output = append(output, input.String())
	expected = append(expected, "send to_c<pay_c,cont_c>")

	// Receive
	input2 := NewReceive(pay_c, cont_c, from_c, end)
	output = append(output, input2.String())
	expected = append(expected, "<pay_c,cont_c> <- recv from_c; close self")

	// Case
	input4 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "label1"}, pay_c, end)})
	output = append(output, input4.String())
	expected = append(expected, "case from_c ( \n   label1<pay_c> => close self\n)")

	// Case (multiple cases)
	input5 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "label1"}, pay_c, end), NewBranch(Label{L: "label2"}, pay_c, end)})
	output = append(output, input5.String())
	expected = append(expected, "case from_c ( \n   label1<pay_c> => close self\n | label2<pay_c> => close self\n)")

	// Select
	input6 := NewSelect(to_c, Label{L: "label1"}, cont_c)
	output = append(output, input6.String())
	expected = append(expected, "to_c.label1<cont_c>")

	// New
	input7 := NewNew(cont_c, end, end)
	output = append(output, input7.String())
	expected = append(expected, "cont_c <- new (close self); close self")

	// Close
	input8 := NewClose(from_c)
	output = append(output, input8.String())
	expected = append(expected, "close from_c")

	compareOutput(t, output, expected)
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
