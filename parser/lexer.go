package main

//go:generate goyacc -p phi -o parser.y.go parse/parser.y

import (
	"fmt"
	"io"
	"os"
)

type phiSymType struct {
	strval string
}

// lexer for phi.
type lexer struct {
	scanner *scanner
	Errors  chan error
}

// newLexer returns a new yacc-compatible lexer.
func newLexer(r io.Reader) *lexer {
	return &lexer{scanner: newScanner(r), Errors: make(chan error, 1)}
}

// Lex is provided for yacc-compatible parser.
func (l *lexer) Lex(yylval *phiSymType) int {
	var token tok
	token, yylval.strval, _, _ = l.scanner.Scan()
	return int(token)
}

// // Error handles error.
// func (l *lexer) Error(err string) {
// 	l.Errors <- &ParseError{Err: err, Pos: l.scanner.pos}
// }

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	lexer := newLexer(file)
	val := phiSymType{}
	for {

		tok := lexer.Lex(&val)
		if tok == EOF {
			break
		}

		fmt.Printf("\t%s\t%s\n", tok, val.strval)
	}
}
