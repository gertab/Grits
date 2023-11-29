package parser

import (
	"bufio"
	"bytes"
	"io"
)

const (
	EOF = iota
	// kILLEGAL
)

// 	LABEL // letters/digits/_

// scanner is a lexical scanner.
type scanner struct {
	r   *bufio.Reader
	pos TokenPos
}

// newScanner returns a new instance of Scanner.
func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r), pos: TokenPos{Char: 0, Lines: []int{}}}
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if reached the end or error occurs.
func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	if ch == '\n' {
		s.pos.Lines = append(s.pos.Lines, s.pos.Char)
		s.pos.Char = 0
	} else {
		s.pos.Char++
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *scanner) unread() {
	_ = s.r.UnreadRune()
	if s.pos.Char == 0 {
		s.pos.Char = s.pos.Lines[len(s.pos.Lines)-1]
		s.pos.Lines = s.pos.Lines[:len(s.pos.Lines)-1]
	} else {
		s.pos.Char--
	}
}

// Scan returns the next token and parsed value.
func (s *scanner) Scan() (token tok, value string, startPos, endPos TokenPos) {
	ch := s.read()

	if isWhitespace(ch) {
		s.skipWhitespace()
		ch = s.read()
	}

	// Track token positions.
	startPos = s.pos
	defer func() { endPos = s.pos }()

	switch ch {
	case eof:
		return 0, "", startPos, endPos
	case '>':
		return RANGLE, string(ch), startPos, endPos
	case '(':
		return LPAREN, string(ch), startPos, endPos
	case ')':
		return RPAREN, string(ch), startPos, endPos
	case '[':
		return LSBRACK, string(ch), startPos, endPos
	case ']':
		return RSBRACK, string(ch), startPos, endPos
	case '{':
		return LCBRACK, string(ch), startPos, endPos
	case '}':
		return RCBRACK, string(ch), startPos, endPos
	case '.':
		return DOT, string(ch), startPos, endPos
	case ';':
		return SEQUENCE, string(ch), startPos, endPos
	case ':':
		return COLON, string(ch), startPos, endPos
	case '|':
		return PIPE, string(ch), startPos, endPos
	case ',':
		return COMMA, string(ch), startPos, endPos
	case '+':
		return PLUS, string(ch), startPos, endPos
	case '*':
		return TIMES, string(ch), startPos, endPos
	case '&':
		return AMPERSAND, string(ch), startPos, endPos
	case '%':
		return PERCENTAGE, string(ch), startPos, endPos
	}

	if s.consumeIfComment(ch) {
		return s.Scan()
	}

	if isSpecialSymbol(ch) {
		// s.unread()
		return s.scanSpecialSymbol(ch)
	}

	if isAlphaNum(ch) || isUnderscore(ch) {
		// s.unread()
		return s.scanLabel(ch)
	}

	return kILLEGAL, string(ch), startPos, endPos
}

// Scan label or keyword
func (s *scanner) scanLabel(ch rune) (token tok, value string, startPos, endPos TokenPos) {
	var buf bytes.Buffer
	startPos = s.pos
	defer func() { endPos = s.pos }()
	buf.WriteRune(ch)

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isAlphaNum(ch) && !isNameSymbols(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch buf.String() {
	case "send":
		return SEND, buf.String(), startPos, endPos
	case "recv":
		return RECEIVE, buf.String(), startPos, endPos
	case "receive":
		return RECEIVE, buf.String(), startPos, endPos
	case "case":
		return CASE, buf.String(), startPos, endPos
	case "close":
		return CLOSE, buf.String(), startPos, endPos
	case "wait":
		return WAIT, buf.String(), startPos, endPos
	case "cast":
		return CAST, buf.String(), startPos, endPos
	case "shift":
		return SHIFT, buf.String(), startPos, endPos
	case "accept":
		return ACCEPT, buf.String(), startPos, endPos
	case "acc":
		return ACCEPT, buf.String(), startPos, endPos
	case "acquire":
		return ACQUIRE, buf.String(), startPos, endPos
	case "acq":
		return ACQUIRE, buf.String(), startPos, endPos
	case "detach":
		return DETACH, buf.String(), startPos, endPos
	case "det":
		return DETACH, buf.String(), startPos, endPos
	case "release":
		return RELEASE, buf.String(), startPos, endPos
	case "rel":
		return RELEASE, buf.String(), startPos, endPos
	case "drop":
		return DROP, buf.String(), startPos, endPos
	case "split":
		return SPLIT, buf.String(), startPos, endPos
	case "push":
		return PUSH, buf.String(), startPos, endPos
	case "new":
		return NEW, buf.String(), startPos, endPos
	case "snew":
		return SNEW, buf.String(), startPos, endPos
	case "forward":
		return FORWARD, buf.String(), startPos, endPos
	case "fwd":
		return FORWARD, buf.String(), startPos, endPos
	case "type":
		return TYPE, buf.String(), startPos, endPos
	case "let":
		return LET, buf.String(), startPos, endPos
	case "in":
		return IN, buf.String(), startPos, endPos
	case "end":
		return END, buf.String(), startPos, endPos
	case "sprc":
		return SPRC, buf.String(), startPos, endPos
	case "prc":
		return PRC, buf.String(), startPos, endPos
	case "self":
		return SELF, buf.String(), startPos, endPos
	case "assuming":
		return ASSUMING, buf.String(), startPos, endPos
	case "exec":
		return EXEC, buf.String(), startPos, endPos
	case "print":
		// Debug keyword
		return PRINT, buf.String(), startPos, endPos
	}
	return LABEL, buf.String(), startPos, endPos
}

func (s *scanner) skipWhitespace() {
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
	}
}

// Consumes line comments (//...) or multiline comments (/*...*/)
func (s *scanner) consumeIfComment(ch rune) bool {
	if ch == '/' {
		if ch = s.read(); ch == '/' {
			s.skipToEOL()
			return true
		} else if ch == '*' {
			s.skipToEndOfComment()
			return true
		} else {
			s.unread()
		}
		// s.unread()
	}
	// Not a comment, so do nothing
	return false
}

func (s *scanner) skipToEndOfComment() {
	for {
		if ch := s.read(); ch == '*' {
			for {
				if ch := s.read(); ch == '/' {
					return
				}
			}
		}
	}
}

func (s *scanner) skipToEOL() {
	for {
		if ch := s.read(); ch == '\n' || ch == eof {
			break
		}
	}
}

// Some commands are multi-character. So, they have to be check explicitly
func isSpecialSymbol(ch rune) bool {
	return ch == '=' || ch == '<' || ch == '-' || ch == '1' || ch == '/' || ch == '\\'
}

func (s *scanner) scanSpecialSymbol(ch rune) (token tok, value string, startPos, endPos TokenPos) {
	startPos = s.pos
	defer func() { endPos = s.pos }()
	// ch := s.read()
	ch2 := s.read()

	switch ch {
	case '=':
		// Can be = or =>
		if ch2 == '>' {
			// is =>
			return RIGHT_ARROW, "=>", startPos, endPos
		} else {
			// is just =
			s.unread()
			return EQUALS, "=", startPos, endPos
		}
	case '<':
		// Can be < or <-
		if ch2 == '-' {
			// is <-
			return LEFT_ARROW, "<-", startPos, endPos
		} else {
			// is just <
			s.unread()
			return LANGLE, "<", startPos, endPos
		}
	case '-':
		// Can be - or -* (or -o)
		if ch2 == '*' {
			// is -o
			return LOLLI, "-*", startPos, endPos
		} else if ch2 == 'o' {
			// is -o
			return LOLLI, "-o", startPos, endPos
		} else {
			// is just -
			s.unread()
			return MINUS, "-", startPos, endPos
		}
	case '1':
		// Can be 1 or label starting with 1
		if isAlphaNum(ch2) || isUnderscore(ch2) {
			// is a label
			s.unread()
			return s.scanLabel(ch)
		} else {
			// is just 1
			s.unread()
			return UNIT, "1", startPos, endPos
		}
	case '\\':
		// Should be \/
		if ch2 == '/' {
			return DOWN_ARROW, "\\/", startPos, endPos
		}
	case '/':
		// Should be /\
		if ch2 == '\\' {
			return UP_ARROW, "\\/", startPos, endPos
		}
	}
	// Not one of the special commands
	return kILLEGAL, string(ch), startPos, endPos
}
