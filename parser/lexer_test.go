package parser

import (
	"strings"
	"testing"
)

// type phiSymType struct {
// 	strval string
// }

func getTokens(l *lexer) []int {
	tokens := make([]int, 0)
	val := phiSymType{}

	for {
		tok := l.Lex(&val)
		if tok == EOF {
			break
		}

		tokens = append(tokens, tok)
	}

	return tokens
}

func compareOutput(t *testing.T, got []int, expected []int) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if got[index] != expected[index] {
			// t.Errorf("got %s, expected %s\n", "sa", "de")
			t.Errorf("got %d, expected %d\n", got[index], expected[index])
		}
	}
}

func TestBasicTokens(t *testing.T) {
	input := "=abc"
	expected := []int{EQUALS, LABEL}

	reader := strings.NewReader(input)
	l := newLexer(reader)
	tokens := getTokens(l)
	compareOutput(t, tokens, expected)
}

func TestSimpleToken(t *testing.T) {
	cases := []struct {
		input    string
		expected []int
	}{
		{"().;:,[]<>", []int{LPAREN, RPAREN, DOT, SEQUENCE, COLON, COMMA, LSBRACK, RSBRACK, LANGLE, RANGLE}},
		{"==><<-<", []int{EQUALS, RIGHT_ARROW, LANGLE, LEFT_ARROW, LANGLE}},
		{"send recv receive case close wait", []int{SEND, RECEIVE, RECEIVE, CASE, CLOSE, WAIT}},
		{"cast shift accept acc acquire acq detach det", []int{CAST, SHIFT, ACCEPT, ACCEPT, ACQUIRE, ACQUIRE, DETACH, DETACH}},
		{"release rel drop split push new", []int{RELEASE, RELEASE, DROP, SPLIT, PUSH, NEW}},
		{"snew forward fwd let in end sprc prc self assuming", []int{SNEW, FORWARD, FORWARD, LET, IN, END, SPRC, PRC, SELF, ASSUMING}},
		{"print", []int{PRINT}},
		{"+-1{},()", []int{PLUS, MINUS, UNIT, LCBRACK, RCBRACK, COMMA, LPAREN, RPAREN}},
	}

	for _, c := range cases {
		reader := strings.NewReader(c.input)
		l := newLexer(reader)
		tokens := getTokens(l)
		compareOutput(t, tokens, c.expected)
	}
}

func TestIdentToken(t *testing.T) {
	cases := []struct {
		input    string
		expected []int
	}{
		{"testIdent", []int{LABEL}},
		{"ill\\egal", []int{LABEL}}, // kILLEGAL
	}

	for _, c := range cases {
		reader := strings.NewReader(c.input)
		l := newLexer(reader)
		tokens := getTokens(l)
		compareOutput(t, tokens, c.expected)
	}
}

func TestCommentToken(t *testing.T) {
	cases := []struct {
		input    string
		expected []int
	}{
		{"test/*abc*/Ident", []int{LABEL, LABEL}},
		{"il//egal", []int{LABEL}}, // kILLEGAL
	}

	for _, c := range cases {
		reader := strings.NewReader(c.input)
		l := newLexer(reader)
		tokens := getTokens(l)
		compareOutput(t, tokens, c.expected)
	}
}
