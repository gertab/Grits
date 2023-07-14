package main

//go:generate goyacc -p main -o parser.y.go parser.y

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

	// KEYWORDS
	SEND
	RECEIVE
	CASE
	CLOSE
	WAIT
	CAST
	SHIFT
	ACCEPT
	ACQUIRE
	DETACH
	RELEASE
	DROP
	SPLIT
	PUSH
	NEW
	SNEW
	FWD
	LET
	IN
	END
	SPRC
	PRC
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	EQUALS: "=",

	// COMMENT       = /\*[\s\t\n\ra-zA-Z0-9_/]*\*/
	LEFT_ARROW:  "<-",
	RIGHT_ARROW: "=>",
	DOT:         ".",
	SEQUENCE:    ";",
	LABEL:       "LABEL",
	COLON:       ":",
	COMMA:       ",",
	LBRACK:      "(",
	RBRACK:      ")",
	LSBRACK:     "[",
	RSBRACK:     "]",
	LVBRACK:     "<",
	RVBRACK:     ">",
	PIPE:        "|",

	SEND:    "SEND",
	RECEIVE: "RECEIVE",
	CASE:    "CASE",
	CLOSE:   "CLOSE",
	WAIT:    "WAIT",
	CAST:    "CAST",
	SHIFT:   "SHIFT",
	ACCEPT:  "ACCEPT",
	ACQUIRE: "ACQUIRE",
	DETACH:  "DETACH",
	RELEASE: "RELEASE",
	DROP:    "DROP",
	SPLIT:   "SPLIT",
	PUSH:    "PUSH",
	NEW:     "NEW",
	SNEW:    "SNEW",
	FWD:     "FWD",
	LET:     "LET",
	IN:      "IN",
	END:     "END",
	SPRC:    "SPRC",
	PRC:     "PRC",
}

var keywords = map[string]Token{
	"send":    SEND,
	"recv":    RECEIVE,
	"receive": RECEIVE,
	"case":    CASE,
	"close":   CLOSE,
	"wait":    WAIT,
	"cast":    CAST,
	"shift":   SHIFT,
	"accept":  ACCEPT,
	"acc":     ACCEPT,
	"acquire": ACQUIRE,
	"acq":     ACQUIRE,
	"detach":  DETACH,
	"det":     DETACH,
	"release": RELEASE,
	"rel":     RELEASE,
	"drop":    DROP,
	"split":   SPLIT,
	"push":    PUSH,
	"new":     NEW,
	"snew":    SNEW,
	"forward": FWD,
	"fwd":     FWD,
	"let":     LET,
	"in":      IN,
	"end":     END,
	"sprc":    SPRC,
	"prc":     PRC,
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

type yySymType struct {
	value string
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (l *Lexer) Lex(lval *yySymType) Token {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return EOF
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
			return SEQUENCE
		case '+':
			return ADD
		case '-':
			return SUB
		case '*':
			return MUL
		// case '/':
		// 	return l.pos, DIV, "/"
		case ':':
			return COLON
		case ',':
			return COMMA
		case '(':
			return LBRACK
		case ')':
			return RBRACK
		case '[':
			return LSBRACK
		case ']':
			return RSBRACK
		case '>':
			return RVBRACK
		case '|':
			return PIPE
		case '/':
			// Potential comment, move to next line
			r, _, err = l.reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					// at the end of the file
					return EOF
				}
			}

			if r == '/' {
				l.readUntilNextLine()
				return l.Lex(lval)

			} else {
				return ILLEGAL
			}

		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if r == '=' {
				// backup and let lexInt rescan the beginning of the int
				// startPos := l.pos
				l.backup()
				label, lit := l.lexEquals()
				lval.value = lit
				return label
			} else if r == '<' {
				// backup and let lexInt rescan the beginning of the int
				// startPos := l.pos
				l.backup()
				label, lit := l.lexLVBrackets()
				lval.value = lit
				return label
			} else if isAlphaNumeric(r) {
				// backup and let lexInt rescan the beginning of the int
				// startPos := l.pos
				l.backup()
				label, lit := l.lexLabel()
				lval.value = lit
				return label
			} else {
				return ILLEGAL
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

func (l *Lexer) readUntilNextLine() {
	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				// at the end of the int
				l.backup()
				return
			}
		}

		if r == '\n' {
			break
		}
	}

	l.resetPosition()
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
	if r == '-' {
		return LEFT_ARROW, "<-"
	} else {
		// scanned something not in the identifier
		l.backup()
		return LVBRACK, lit
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexLabel() (Token, string) {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return getLabelOrKeyword(lit)
			}
		}

		l.pos.column++
		if isAlphaNumeric(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return getLabelOrKeyword(lit)
		}
	}
}

func getLabelOrKeyword(lit string) (Token, string) {
	val, ok := keywords[lit]

	if ok {
		return val, lit
	}
	return LABEL, lit
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
func main() {
	file, err := os.Open("parser/input.test")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	val := yySymType{}
	for {

		tok := lexer.Lex(&val)
		if tok == EOF {
			break
		}

		fmt.Printf("\t%s\t%s\n", tok, val.value)
	}
}
