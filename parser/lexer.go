package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	EQUALS // =

	LABEL       // letters/digits/_
	LEFT_ARROW  // <-
	RIGHT_ARROW // =>

	DOT      // .
	SEQUENCE // ;
	COLON    // :
	COMMA    // ,
	LBRACK   // (
	RBRACK   // )
	LSBRACK  // [
	RSBRACK  // ]
	LVBRACK  // <
	RVBRACK  // >
	PIPE     // |
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	EQUALS: "=",

	// COMMENT       = /\*[\s\t\n\ra-zA-Z0-9_/]*\*/
	LEFT_ARROW:  "<-",
	RIGHT_ARROW: "=>",
	// EQUALS        = \=
	// SEND          = send
	// RECEIVE       = recv|receive
	// CASE          = case
	// CLOSE         = close
	// WAIT          = wait
	// CAST          = cast
	// SHIFT         = shift
	// ACCEPT        = accept|acc
	// ACQUIRE       = acquire|acq
	// DETACH        = detach|det
	// RELEASE       = release|rel
	// DROP          = drop
	// SPLIT         = split
	// PUSH          = push
	// NEW           = new
	// SNEW          = snew
	// FWD           = forward|fwd
	// LET           = let
	// IN            = in
	// END           = end
	// SPRC           = sprc
	// PRC           = prc
	DOT:      ".",
	SEQUENCE: ";",
	LABEL:    "LABEL",
	COLON:    ":",
	COMMA:    ",",
	LBRACK:   "(",
	RBRACK:   ")",
	LSBRACK:  "[",
	RSBRACK:  "]",
	LVBRACK:  "<",
	RVBRACK:  ">",
	PIPE:     "|",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (l *Lexer) Lex() (Position, Token, string) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}

		// update the column to the position of the newly read in rune
		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEQUENCE, ";"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, MUL, "*"
		case '/':
			return l.pos, DIV, "/"
		// case '=':
		// 	return l.pos, EQUALS, "="

		case ':':
			return l.pos, COLON, ":"
		case ',':
			return l.pos, COMMA, ","
		case '(':
			return l.pos, LBRACK, "("
		case ')':
			return l.pos, RBRACK, ")"
		case '[':
			return l.pos, LSBRACK, "["
		case ']':
			return l.pos, RSBRACK, "]"
		// case '<':
		// 	return l.pos, LVBRACK, "<"
		case '>':
			return l.pos, RVBRACK, ">"
		case '|':
			return l.pos, PIPE, "|"
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if r == '=' {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				label, lit := l.lexEquals()
				return startPos, label, lit
			} else if r == '<' {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				label, lit := l.lexLVBrackets()
				return startPos, label, lit
			} else if isLabel(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				lit := l.lexLabel()
				return startPos, LABEL, lit
				// } else if unicode.IsDigit(r) {
				// 	// backup and let lexInt rescan the beginning of the int
				// 	startPos := l.pos
				// 	l.backup()
				// 	lit := l.lexInt()
				// 	return startPos, INT, lit
				// } else if unicode.IsLetter(r) {
				// 	// backup and let lexIdent rescan the beginning of the ident
				// 	startPos := l.pos
				// 	l.backup()
				// 	lit := l.lexIdent()
				// 	return startPos, IDENT, lit
			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

// lexInt scans the input until the end of an integer and then returns the
// literal.
func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return lit
		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexEquals() (Token, string) {
	r, _, err := l.reader.ReadRune()
	var lit = string(r)

	r, _, err = l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			// at the end of the identifier
			return EQUALS, lit
		}
	}

	l.pos.column++
	if r == '>' {
		// lit + string(r)
		return RIGHT_ARROW, "=>"

	} else {
		// scanned something not in the identifier
		l.backup()
		return EQUALS, lit

	}
}

func (l *Lexer) lexLVBrackets() (Token, string) {
	r, _, err := l.reader.ReadRune()
	var lit = string(r)

	r, _, err = l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			// at the end of the identifier
			return LVBRACK, lit
		}
	}

	l.pos.column++
	if r == '=' {
		// return lit + string(r)
		return LEFT_ARROW, "<="

	} else {
		// scanned something not in the identifier
		l.backup()
		return EQUALS, lit
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexLabel() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.column++
		if isLabel(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}

func isLabel(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
func main() {
	file, err := os.Open("parser/input.test")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
}
