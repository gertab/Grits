package process

import "testing"

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

func compareOutputProgram(t *testing.T, got []Form, expected []Form) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if !EqualForm(got[index], expected[index]) {
			t.Errorf("[%d] got %s, expected %s\n", index, got[index].String(), expected[index].String())
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
	expected = append(expected, "case from_c (label1<pay_c> => close self)")

	// Case (multiple cases)
	input5 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "label1"}, pay_c, end), NewBranch(Label{L: "label2"}, pay_c, end)})
	output = append(output, input5.String())
	expected = append(expected, "case from_c (label1<pay_c> => close self | label2<pay_c> => close self)")

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

	// Forward
	input9 := NewForward(to_c, from_c)
	output = append(output, input9.String())
	expected = append(expected, "fwd to_c from_c")

	compareOutput(t, output, expected)
}

func TestSubstitutions(t *testing.T) {
	expected, output := []Form{}, []Form{}
	to_c := Name{Ident: "to_c"}
	new_to_c := Name{Ident: "new_to_c"}
	pay_c := Name{Ident: "pay_c"}
	new_pay_c := Name{Ident: "new_pay_c"}
	cont_c := Name{Ident: "cont_c"}
	new_cont_c := Name{Ident: "new_cont_c"}
	from_c := Name{Ident: "from_c"}
	new_from_c := Name{Ident: "new_from_c"}
	self := Name{Ident: "from_c"}
	new_self := Name{Ident: "new_from_c"}
	end := NewClose(self)
	new_end := NewClose(new_self)

	// Send
	input := NewSend(to_c, pay_c, cont_c)
	input.Substitute(to_c, new_to_c)
	input.Substitute(pay_c, new_pay_c)
	input.Substitute(cont_c, new_cont_c)
	output = append(output, input)
	result := NewSend(new_to_c, new_pay_c, new_cont_c)
	expected = append(expected, result)

	// Receive
	input2 := NewReceive(pay_c, cont_c, from_c, end)
	input2.Substitute(cont_c, new_cont_c)
	input2.Substitute(pay_c, new_pay_c)
	input2.Substitute(cont_c, new_cont_c)
	input2.Substitute(self, new_self)
	result2 := NewReceive(pay_c, cont_c, new_from_c, new_end)
	output = append(output, input2)
	expected = append(expected, result2)

	// Case
	input3 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "from_c"}, pay_c, end)})
	input3.Substitute(from_c, new_from_c)
	input3.Substitute(pay_c, new_pay_c)
	input3.Substitute(cont_c, new_cont_c)
	input3.Substitute(self, new_self)
	result3 := NewCase(new_from_c, []*BranchForm{NewBranch(Label{L: "from_c"}, pay_c, new_end)})
	output = append(output, input3)
	expected = append(expected, result3)

	// Select
	input4 := NewSelect(to_c, Label{L: "label1"}, cont_c)
	input4.Substitute(to_c, new_to_c)
	input4.Substitute(cont_c, new_cont_c)
	result4 := NewSelect(new_to_c, Label{L: "label1"}, new_cont_c)
	output = append(output, input4)
	expected = append(expected, result4)

	// New
	input5 := NewNew(cont_c, end, end)
	input5.Substitute(cont_c, new_cont_c)
	input5.Substitute(self, new_self)
	result5 := NewNew(cont_c, end, end)
	output = append(output, input5)
	expected = append(expected, result5)

	// Close
	input6 := NewClose(from_c)
	input6.Substitute(from_c, new_from_c)
	input6.Substitute(self, new_self)
	result6 := NewClose(new_from_c)
	output = append(output, input6)
	expected = append(expected, result6)

	// Forward
	input7 := NewForward(to_c, from_c)
	input7.Substitute(from_c, new_from_c)
	input7.Substitute(to_c, new_to_c)
	result7 := NewForward(new_to_c, new_from_c)
	output = append(output, input7)
	expected = append(expected, result7)

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
