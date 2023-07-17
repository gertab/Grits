package main

import (
	"bufio"
	"bytes"
	"io"
)

const (
	EOF = iota
	// ILLEGAL

	kNIL
	kNEW
	kNAME

	EQUALS // =

	// LABEL       // letters/digits/_
	LEFT_ARROW  // <-
	RIGHT_ARROW // =>

	DOT      // .
	SEQUENCE // ;
	COLON    // :
	COMMA    // ,
	LPAREN   // (
	RPAREN   // )
	LSBRACK  // [
	RSBRACK  // ]
	LANGLE   // <
	RANGLE   // >
	PIPE     // |

	// // KEYWORDS
	// SEND
	// RECEIVE
	// CASE
	// CLOSE
	// WAIT
	// CAST
	// SHIFT
	// ACCEPT
	// ACQUIRE
	// DETACH
	// RELEASE
	// DROP
	// SPLIT
	// PUSH
	// NEW
	// SNEW
	// FWD
	// LET
	// IN
	// END
	// SPRC
	// PRC
)

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
	}

	if s.consumeIfComment(ch) {
		return s.Scan()
	}

	if isSpecialSymbol(ch) {
		s.unread()
		return s.scanSpecialSymbol()
	}

	if isAlphaNum(ch) || isUnderscore(ch) {
		s.unread()
		return s.scanName()
	}

	return kILLEGAL, string(ch), startPos, endPos
}

func (s *scanner) scanName() (token tok, value string, startPos, endPos TokenPos) {
	var buf bytes.Buffer
	startPos = s.pos
	defer func() { endPos = s.pos }()
	buf.WriteRune(s.read())

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
	case "0":
		return kNIL, buf.String(), startPos, endPos
	case "new":
		return kNEW, buf.String(), startPos, endPos
	}
	return kNAME, buf.String(), startPos, endPos
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

func (s *scanner) consumeIfComment(ch rune) bool {
	if ch == '/' {
		if ch = s.read(); ch == '/' {
			s.skipToEOL()
			return true
		} else if ch == '*' {
			// todo implement /* ... */
			s.skipToEndOfComment()
			return true
		}
		s.unread()
	}
	// Not a comment, so undo changes
	s.unread()
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
	return ch == '=' || ch == '<'
}

func (s *scanner) scanSpecialSymbol() (token tok, value string, startPos, endPos TokenPos) {
	startPos = s.pos
	defer func() { endPos = s.pos }()
	ch := s.read()

	// for {
	// 	if ch := s.read(); ch == eof {
	// 		break
	// 	} else if !isAlphaNum(ch) && !isNameSymbols(ch) {
	// 		s.unread()
	// 		break
	// 	} else {
	// 		_, _ = buf.WriteRune(ch)
	// 	}
	// }

	// switch buf.String() {
	// case "0":
	// 	return kNIL, buf.String(), startPos, endPos
	// case "new":
	// 	return kNEW, buf.String(), startPos, endPos
	// }
	// return kNAME, buf.String(), startPos, endPos

	switch ch {
	case '=':
		// Can be = or =>

		if ch2 := s.read(); ch2 == '>' {
			// is =>
			return RIGHT_ARROW, "=>", startPos, endPos
		} else {
			// is just =
			s.unread()
			return EQUALS, "=", startPos, endPos
		}
	case '<':
		// Can be < or <-

		if ch2 := s.read(); ch2 == '-' {
			// is <-
			return LEFT_ARROW, "<-", startPos, endPos
		} else {
			// is just =
			s.unread()
			return LANGLE, "<", startPos, endPos
		}
	}
	// Not one of the special commands
	return kILLEGAL, string(ch), startPos, endPos
}
