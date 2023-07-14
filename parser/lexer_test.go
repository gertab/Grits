package main

import (
	"strings"
	"testing"
)

func getTokens(l *Lexer) []Token {
	tokens := make([]Token, 0)
	for {
		_, tok, _ := l.Lex()
		if tok == EOF {
			break
		}

		tokens = append(tokens, tok)
	}

	return tokens
}

func compareOutput(t *testing.T, got []Token, expected []Token) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if got[index] != expected[index] {
			t.Errorf("got %s, expected %s\n", got[index], expected[index])
		}
	}
}

func TestBasicTokens(t *testing.T) {
	input := "+-*/="
	expected := []Token{ADD, SUB, MUL, DIV}

	reader := strings.NewReader(input)
	l := NewLexer(reader)
	tokens := getTokens(l)
	compareOutput(t, tokens, expected)
}

func TestIntToken(t *testing.T) {
	cases := []struct {
		input    string
		expected []Token
	}{
		{"123+23", []Token{LEFT_ARROW, ADD, LEFT_ARROW}},
		{"11111111010100-", []Token{LEFT_ARROW, SUB}},
		{"24593753790175972954 5439574375348", []Token{LEFT_ARROW, LEFT_ARROW}},
		{"213one", []Token{LEFT_ARROW, RELEASE}},
	}

	for _, c := range cases {
		reader := strings.NewReader(c.input)
		l := NewLexer(reader)
		tokens := getTokens(l)
		compareOutput(t, tokens, c.expected)
	}
}

func TestIdentToken(t *testing.T) {
	cases := []struct {
		input    string
		expected []Token
	}{
		{"testIdent", []Token{LABEL}},
		{"ill)egal", []Token{LABEL, ILLEGAL, LABEL}},
		{"ill.egal", []Token{LABEL, ILLEGAL, LABEL}},
		{"one two", []Token{LABEL, LABEL}},
	}

	for _, c := range cases {
		reader := strings.NewReader(c.input)
		l := NewLexer(reader)
		tokens := getTokens(l)
		compareOutput(t, tokens, c.expected)
	}
}
