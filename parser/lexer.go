package parser

//go:generate goyacc -p phi -o parser/parser.y.go parser/parser.y

import (
	"io"
)

// Generated from goyacc
// type phiSymType struct {
// 	strval string
// }

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

// Error handles error.
func (l *lexer) Error(err string) {
	l.Errors <- &ParseError{Err: err, Pos: l.scanner.pos}
}

// func main() {
// 	file, err := os.Open("input.test")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// prc, err := Parse(file)

// 	// fmt.Println(prc)

// 	lexer := newLexer(file)
// 	val := phiSymType{}
// 	for {

// 		tok := lexer.Lex(&val)
// 		if tok == EOF {
// 			break
// 		}

// 		fmt.Printf("\t%d\t%s\n", tok, val.strval)
// 	}
// }
