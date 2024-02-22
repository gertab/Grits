// Code generated by goyacc -p grits -o parser/parser.y.go parser/parser.y. DO NOT EDIT.

// Run this after each change:
// goyacc -p grits -o parser/parser.y.go parser/parser.y
//
//line parser/parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser/parser.y:4

import (
	"grits/position"
	"grits/process"
	"grits/types"
	"io"
)

//line parser/parser.y:15
type gritsSymType struct {
	yys                   int
	strval                string
	currPosition          position.Position
	common_type           unexpandedProcessOrFunction
	statements            []unexpandedProcessOrFunction
	name                  process.Name
	names                 []process.Name
	form                  process.Form
	branches              []*process.BranchForm
	sessionType           types.SessionType
	sessionTypeInitial    types.SessionTypeInitial
	sessionTypeAltInitial []types.OptionInitial
	polarity              types.Polarity
}

const LABEL = 57346
const LEFT_ARROW = 57347
const RIGHT_ARROW = 57348
const UP_ARROW = 57349
const DOWN_ARROW = 57350
const EQUALS = 57351
const DOT = 57352
const SEQUENCE = 57353
const COLON = 57354
const COMMA = 57355
const LPAREN = 57356
const RPAREN = 57357
const LSBRACK = 57358
const RSBRACK = 57359
const LANGLE = 57360
const RANGLE = 57361
const PIPE = 57362
const SEND = 57363
const RECEIVE = 57364
const CASE = 57365
const CLOSE = 57366
const WAIT = 57367
const CAST = 57368
const SHIFT = 57369
const ACCEPT = 57370
const ACQUIRE = 57371
const DETACH = 57372
const RELEASE = 57373
const DROP = 57374
const SPLIT = 57375
const PUSH = 57376
const NEW = 57377
const SNEW = 57378
const TYPE = 57379
const LET = 57380
const IN = 57381
const END = 57382
const SPRC = 57383
const PRC = 57384
const FORWARD = 57385
const SELF = 57386
const PRINT = 57387
const PLUS = 57388
const MINUS = 57389
const TIMES = 57390
const AMPERSAND = 57391
const UNIT = 57392
const LCBRACK = 57393
const RCBRACK = 57394
const LOLLI = 57395
const PERCENTAGE = 57396
const ASSUMING = 57397
const EXEC = 57398

var gritsToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"LABEL",
	"LEFT_ARROW",
	"RIGHT_ARROW",
	"UP_ARROW",
	"DOWN_ARROW",
	"EQUALS",
	"DOT",
	"SEQUENCE",
	"COLON",
	"COMMA",
	"LPAREN",
	"RPAREN",
	"LSBRACK",
	"RSBRACK",
	"LANGLE",
	"RANGLE",
	"PIPE",
	"SEND",
	"RECEIVE",
	"CASE",
	"CLOSE",
	"WAIT",
	"CAST",
	"SHIFT",
	"ACCEPT",
	"ACQUIRE",
	"DETACH",
	"RELEASE",
	"DROP",
	"SPLIT",
	"PUSH",
	"NEW",
	"SNEW",
	"TYPE",
	"LET",
	"IN",
	"END",
	"SPRC",
	"PRC",
	"FORWARD",
	"SELF",
	"PRINT",
	"PLUS",
	"MINUS",
	"TIMES",
	"AMPERSAND",
	"UNIT",
	"LCBRACK",
	"RCBRACK",
	"LOLLI",
	"PERCENTAGE",
	"ASSUMING",
	"EXEC",
}

var gritsStatenames = [...]string{}

const gritsEofCode = 1
const gritsErrCode = 2
const gritsInitialStackSize = 16

//line parser/parser.y:261

// Parse is the entry point to the parser.
func Parse(r io.Reader) (allEnvironment, error) {
	l := newLexer(r)
	gritsParse(l)
	allEnvironment := allEnvironment{}
	select {
	case err := <-l.Errors:
		return allEnvironment, err
	default:
		// allEnvironment := l
		allEnvironment.procsAndFuns = l.processesOrFunctionsRes
		return allEnvironment, nil
	}
}

//line yacctab:1
var gritsExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 69,
	4, 72,
	7, 72,
	8, 72,
	14, 72,
	46, 72,
	49, 72,
	50, 72,
	-2, 61,
}

const gritsPrivate = 57344

const gritsLast = 270

var gritsAct = [...]uint8{
	3, 136, 147, 57, 98, 67, 115, 82, 66, 99,
	162, 103, 160, 9, 105, 56, 44, 104, 64, 52,
	130, 195, 175, 15, 171, 188, 63, 6, 139, 152,
	5, 7, 8, 10, 12, 13, 174, 31, 33, 172,
	36, 14, 39, 40, 41, 42, 43, 32, 141, 68,
	173, 123, 11, 22, 16, 29, 30, 26, 25, 51,
	98, 98, 24, 127, 93, 99, 99, 129, 128, 92,
	75, 78, 76, 69, 100, 27, 28, 60, 109, 106,
	111, 192, 112, 73, 168, 83, 53, 22, 116, 29,
	30, 113, 90, 91, 145, 118, 94, 120, 68, 84,
	68, 85, 121, 119, 131, 132, 107, 138, 80, 37,
	110, 38, 89, 133, 135, 71, 140, 65, 72, 70,
	148, 149, 144, 149, 124, 165, 153, 154, 166, 4,
	146, 157, 142, 122, 114, 143, 199, 108, 87, 158,
	83, 159, 61, 163, 161, 88, 83, 46, 47, 48,
	49, 50, 164, 116, 150, 68, 170, 169, 194, 193,
	68, 156, 167, 176, 126, 125, 179, 177, 81, 79,
	77, 35, 182, 190, 181, 180, 34, 187, 68, 189,
	178, 86, 191, 101, 102, 202, 186, 196, 151, 97,
	197, 198, 137, 58, 200, 201, 155, 9, 134, 117,
	203, 96, 62, 204, 183, 184, 185, 15, 59, 55,
	54, 6, 45, 2, 5, 1, 8, 10, 12, 13,
	23, 95, 74, 69, 21, 14, 101, 102, 20, 19,
	26, 25, 18, 73, 17, 24, 11, 22, 16, 29,
	30, 0, 0, 0, 0, 0, 0, 0, 27, 28,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 71, 0, 0, 72, 70,
}

var gritsPact = [...]int16{
	193, -1000, -1000, -1000, -1000, 43, 43, 166, 43, 97,
	43, 43, 43, 43, 43, 9, 208, 20, 20, 20,
	20, 20, -1000, 15, 70, 206, 205, 189, 204, -1000,
	-1000, 59, -1000, 129, 198, -9, 103, 69, 43, -1000,
	43, 159, 53, 158, 93, 157, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 43, 85, 172, -1000, 125, 133, 98,
	43, 43, 51, 9, 43, 197, 184, -44, 219, -1000,
	-1000, -34, -37, 69, 91, 124, -1000, 9, 43, 9,
	-1000, 9, 74, 121, 189, 195, 69, 189, 69, 87,
	120, 32, 43, 154, 153, 48, 49, -15, 69, 69,
	-44, 194, 194, 176, 188, 188, 13, -1000, 43, -1000,
	29, -1000, -1000, 123, 43, 79, 117, 108, -1000, -1000,
	-1000, -1000, 43, 183, 10, 9, 9, -1000, 192, 43,
	9, -44, -44, 69, -1000, 69, -40, 132, -42, -1000,
	-1000, -1000, 9, 69, -1000, 116, 189, 67, 69, 189,
	5, 17, -1000, -1000, -1000, 18, 3, 152, -44, -44,
	-1000, 69, -1000, -1000, 171, 9, 69, -1000, 165, 110,
	-1000, -1000, 43, 43, 43, 180, 9, 12, 9, -1000,
	164, 9, 64, 148, 147, 2, 9, -1000, 188, -1000,
	9, -1000, 127, 9, 9, 179, -1000, -1000, -1000, 9,
	-1000, -1000, 9, -1000, -1000,
}

var gritsPgo = [...]uint8{
	0, 129, 234, 232, 229, 228, 224, 0, 31, 3,
	7, 222, 6, 2, 15, 11, 221, 8, 1, 5,
	220, 215, 213,
}

var gritsR1 = [...]int8{
	0, 21, 22, 22, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 2, 2, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 16, 16, 16, 10, 10, 11, 11, 11,
	12, 12, 12, 13, 13, 14, 14, 9, 9, 8,
	8, 8, 8, 5, 3, 3, 3, 3, 4, 17,
	17, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	18, 18, 15, 20, 20, 6,
}

var gritsR2 = [...]int8{
	0, 1, 1, 1, 1, 2, 1, 2, 1, 2,
	1, 2, 1, 2, 6, 8, 7, 10, 6, 5,
	6, 8, 4, 2, 3, 10, 4, 5, 6, 4,
	3, 4, 0, 6, 8, 1, 3, 0, 1, 3,
	0, 1, 3, 0, 2, 1, 3, 1, 3, 1,
	2, 1, 2, 2, 7, 9, 8, 10, 4, 1,
	2, 1, 1, 4, 4, 3, 3, 3, 4, 4,
	3, 5, 1, 1, 1, 4,
}

var gritsChk = [...]int16{
	-1000, -21, -22, -7, -1, 21, 18, -8, 23, 4,
	24, 43, 25, 26, 32, 14, 45, -2, -3, -4,
	-5, -6, 44, -20, 42, 38, 37, 55, 56, 46,
	47, -8, 4, -8, 10, 5, -8, 12, 14, -8,
	-8, -8, -8, -8, -7, 4, -1, -1, -1, -1,
	-1, 44, 4, 16, 4, 4, -14, -9, 4, 4,
	18, 13, 4, 35, 27, 14, -17, -19, -15, 4,
	50, 46, 49, 14, -11, -8, -8, 11, 18, 11,
	15, 11, -10, -8, 14, 16, 9, 13, 12, 14,
	-8, -8, 18, -7, -8, -16, 4, 5, 48, 53,
	-19, 7, 8, -15, 51, 51, -19, 15, 13, -7,
	-8, -7, -7, 17, 13, -12, -9, 4, -17, -14,
	-17, 15, 13, 19, -8, 11, 11, 15, 20, 18,
	35, -19, -19, -15, 4, -15, -18, 4, -18, 15,
	-10, 19, 9, 12, -10, 15, 13, -13, 12, 13,
	-8, 5, 19, -7, -7, 4, -8, -7, -19, -19,
	52, 12, 52, -7, -17, 9, 12, -14, 17, -17,
	-12, 19, 22, 33, 18, 19, 11, -19, 9, -7,
	-17, 9, -13, -8, -8, -8, 6, -7, 13, -7,
	9, -7, 17, 11, 11, 19, -7, -18, -7, 9,
	-7, -7, 6, -7, -7,
}

var gritsDef = [...]int8{
	0, -2, 1, 2, 3, 0, 0, 0, 0, 51,
	0, 0, 0, 0, 0, 0, 0, 4, 6, 8,
	10, 12, 49, 0, 0, 0, 0, 0, 0, 73,
	74, 0, 51, 0, 0, 0, 0, 0, 37, 23,
	0, 0, 0, 0, 0, 0, 5, 7, 9, 11,
	13, 50, 52, 0, 0, 0, 53, 45, 47, 0,
	0, 0, 0, 0, 0, 32, 0, 59, 0, -2,
	62, 0, 0, 0, 0, 38, 24, 0, 0, 0,
	30, 0, 0, 35, 40, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	60, 0, 0, 0, 0, 0, 0, 22, 0, 26,
	0, 29, 31, 0, 0, 0, 41, 43, 58, 46,
	48, 75, 0, 0, 0, 0, 0, 19, 0, 0,
	0, 65, 66, 0, 72, 0, 0, 0, 0, 67,
	39, 27, 0, 0, 36, 0, 0, 0, 0, 40,
	0, 0, 18, 20, 28, 0, 0, 0, 68, 69,
	63, 0, 64, 14, 0, 0, 0, 42, 0, 43,
	44, 16, 0, 0, 0, 0, 0, 70, 0, 54,
	0, 0, 0, 0, 0, 0, 0, 21, 0, 15,
	0, 56, 0, 0, 0, 0, 33, 71, 55, 0,
	17, 25, 0, 57, 34,
}

var gritsTok1 = [...]int8{
	1,
}

var gritsTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56,
}

var gritsTok3 = [...]int8{
	0,
}

var gritsErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	gritsDebug        = 0
	gritsErrorVerbose = false
)

type gritsLexer interface {
	Lex(lval *gritsSymType) int
	Error(s string)
}

type gritsParser interface {
	Parse(gritsLexer) int
	Lookahead() int
}

type gritsParserImpl struct {
	lval  gritsSymType
	stack [gritsInitialStackSize]gritsSymType
	char  int
}

func (p *gritsParserImpl) Lookahead() int {
	return p.char
}

func gritsNewParser() gritsParser {
	return &gritsParserImpl{}
}

const gritsFlag = -1000

func gritsTokname(c int) string {
	if c >= 1 && c-1 < len(gritsToknames) {
		if gritsToknames[c-1] != "" {
			return gritsToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func gritsStatname(s int) string {
	if s >= 0 && s < len(gritsStatenames) {
		if gritsStatenames[s] != "" {
			return gritsStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func gritsErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !gritsErrorVerbose {
		return "syntax error"
	}

	for _, e := range gritsErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + gritsTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(gritsPact[state])
	for tok := TOKSTART; tok-1 < len(gritsToknames); tok++ {
		if n := base + tok; n >= 0 && n < gritsLast && int(gritsChk[int(gritsAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if gritsDef[state] == -2 {
		i := 0
		for gritsExca[i] != -1 || int(gritsExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; gritsExca[i] >= 0; i += 2 {
			tok := int(gritsExca[i])
			if tok < TOKSTART || gritsExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if gritsExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += gritsTokname(tok)
	}
	return res
}

func gritslex1(lex gritsLexer, lval *gritsSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(gritsTok1[0])
		goto out
	}
	if char < len(gritsTok1) {
		token = int(gritsTok1[char])
		goto out
	}
	if char >= gritsPrivate {
		if char < gritsPrivate+len(gritsTok2) {
			token = int(gritsTok2[char-gritsPrivate])
			goto out
		}
	}
	for i := 0; i < len(gritsTok3); i += 2 {
		token = int(gritsTok3[i+0])
		if token == char {
			token = int(gritsTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(gritsTok2[1]) /* unknown char */
	}
	if gritsDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", gritsTokname(token), uint(char))
	}
	return char, token
}

func gritsParse(gritslex gritsLexer) int {
	return gritsNewParser().Parse(gritslex)
}

func (gritsrcvr *gritsParserImpl) Parse(gritslex gritsLexer) int {
	var gritsn int
	var gritsVAL gritsSymType
	var gritsDollar []gritsSymType
	_ = gritsDollar // silence set and not used
	gritsS := gritsrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	gritsstate := 0
	gritsrcvr.char = -1
	gritstoken := -1 // gritsrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		gritsstate = -1
		gritsrcvr.char = -1
		gritstoken = -1
	}()
	gritsp := -1
	goto gritsstack

ret0:
	return 0

ret1:
	return 1

gritsstack:
	/* put a state and value onto the stack */
	if gritsDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", gritsTokname(gritstoken), gritsStatname(gritsstate))
	}

	gritsp++
	if gritsp >= len(gritsS) {
		nyys := make([]gritsSymType, len(gritsS)*2)
		copy(nyys, gritsS)
		gritsS = nyys
	}
	gritsS[gritsp] = gritsVAL
	gritsS[gritsp].yys = gritsstate

gritsnewstate:
	gritsn = int(gritsPact[gritsstate])
	if gritsn <= gritsFlag {
		goto gritsdefault /* simple state */
	}
	if gritsrcvr.char < 0 {
		gritsrcvr.char, gritstoken = gritslex1(gritslex, &gritsrcvr.lval)
	}
	gritsn += gritstoken
	if gritsn < 0 || gritsn >= gritsLast {
		goto gritsdefault
	}
	gritsn = int(gritsAct[gritsn])
	if int(gritsChk[gritsn]) == gritstoken { /* valid shift */
		gritsrcvr.char = -1
		gritstoken = -1
		gritsVAL = gritsrcvr.lval
		gritsstate = gritsn
		if Errflag > 0 {
			Errflag--
		}
		goto gritsstack
	}

gritsdefault:
	/* default state action */
	gritsn = int(gritsDef[gritsstate])
	if gritsn == -2 {
		if gritsrcvr.char < 0 {
			gritsrcvr.char, gritstoken = gritslex1(gritslex, &gritsrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if gritsExca[xi+0] == -1 && int(gritsExca[xi+1]) == gritsstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			gritsn = int(gritsExca[xi+0])
			if gritsn < 0 || gritsn == gritstoken {
				break
			}
		}
		gritsn = int(gritsExca[xi+1])
		if gritsn < 0 {
			goto ret0
		}
	}
	if gritsn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			gritslex.Error(gritsErrorMessage(gritsstate, gritstoken))
			Nerrs++
			if gritsDebug >= 1 {
				__yyfmt__.Printf("%s", gritsStatname(gritsstate))
				__yyfmt__.Printf(" saw %s\n", gritsTokname(gritstoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for gritsp >= 0 {
				gritsn = int(gritsPact[gritsS[gritsp].yys]) + gritsErrCode
				if gritsn >= 0 && gritsn < gritsLast {
					gritsstate = int(gritsAct[gritsn]) /* simulate a shift of "error" */
					if int(gritsChk[gritsstate]) == gritsErrCode {
						goto gritsstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if gritsDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", gritsS[gritsp].yys)
				}
				gritsp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if gritsDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", gritsTokname(gritstoken))
			}
			if gritstoken == gritsEofCode {
				goto ret1
			}
			gritsrcvr.char = -1
			gritstoken = -1
			goto gritsnewstate /* try again in the same state */
		}
	}

	/* reduction by production gritsn */
	if gritsDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", gritsn, gritsStatname(gritsstate))
	}

	gritsnt := gritsn
	gritspt := gritsp
	_ = gritspt // guard against "declared and not used"

	gritsp -= int(gritsR2[gritsn])
	// gritsp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if gritsp+1 >= len(gritsS) {
		nyys := make([]gritsSymType, len(gritsS)*2)
		copy(nyys, gritsS)
		gritsS = nyys
	}
	gritsVAL = gritsS[gritsp+1]

	/* consult goto table to find next state */
	gritsn = int(gritsR1[gritsn])
	gritsg := int(gritsPgo[gritsn])
	gritsj := gritsg + gritsS[gritsp].yys + 1

	if gritsj >= gritsLast {
		gritsstate = int(gritsAct[gritsg])
	} else {
		gritsstate = int(gritsAct[gritsj])
		if int(gritsChk[gritsstate]) != -gritsn {
			gritsstate = int(gritsAct[gritsg])
		}
	}
	// dummy call; replaced with literal code
	switch gritsnt {

	case 1:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:58
		{
		}
	case 2:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:64
		{
			gritslex.(*lexer).processesOrFunctionsRes = append(gritslex.(*lexer).processesOrFunctionsRes, unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body: gritsDollar[1].form, Providers: []process.Name{{Ident: "root", IsSelf: false}}}, position: gritsVAL.currPosition})
		}
	case 3:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:68
		{
			gritslex.(*lexer).processesOrFunctionsRes = gritsDollar[1].statements
		}
	case 4:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:74
		{
			gritsVAL.statements = []unexpandedProcessOrFunction{gritsDollar[1].common_type}
		}
	case 5:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:75
		{
			gritsVAL.statements = append([]unexpandedProcessOrFunction{gritsDollar[1].common_type}, gritsDollar[2].statements...)
		}
	case 6:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:76
		{
			gritsVAL.statements = []unexpandedProcessOrFunction{gritsDollar[1].common_type}
		}
	case 7:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:77
		{
			gritsVAL.statements = append([]unexpandedProcessOrFunction{gritsDollar[1].common_type}, gritsDollar[2].statements...)
		}
	case 8:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:78
		{
			gritsVAL.statements = []unexpandedProcessOrFunction{gritsDollar[1].common_type}
		}
	case 9:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:79
		{
			gritsVAL.statements = append([]unexpandedProcessOrFunction{gritsDollar[1].common_type}, gritsDollar[2].statements...)
		}
	case 10:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:80
		{
			gritsVAL.statements = []unexpandedProcessOrFunction{gritsDollar[1].common_type}
		}
	case 11:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:81
		{
			gritsVAL.statements = append([]unexpandedProcessOrFunction{gritsDollar[1].common_type}, gritsDollar[2].statements...)
		}
	case 12:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:82
		{
			gritsVAL.statements = []unexpandedProcessOrFunction{gritsDollar[1].common_type}
		}
	case 13:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:83
		{
			gritsVAL.statements = append([]unexpandedProcessOrFunction{gritsDollar[1].common_type}, gritsDollar[2].statements...)
		}
	case 14:
		gritsDollar = gritsS[gritspt-6 : gritspt+1]
//line parser/parser.y:89
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body: gritsDollar[6].form, Providers: gritsDollar[3].names}, position: gritsVAL.currPosition}
		}
	case 15:
		gritsDollar = gritsS[gritspt-8 : gritspt+1]
//line parser/parser.y:91
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body: gritsDollar[8].form, Type: gritsDollar[6].sessionType, Providers: gritsDollar[3].names}, position: gritsVAL.currPosition}
		}
	case 16:
		gritsDollar = gritsS[gritspt-7 : gritspt+1]
//line parser/parser.y:97
		{
			gritsVAL.form = process.NewSend(gritsDollar[2].name, gritsDollar[4].name, gritsDollar[6].name)
		}
	case 17:
		gritsDollar = gritsS[gritspt-10 : gritspt+1]
//line parser/parser.y:101
		{
			gritsVAL.form = process.NewReceive(gritsDollar[2].name, gritsDollar[4].name, gritsDollar[8].name, gritsDollar[10].form)
		}
	case 18:
		gritsDollar = gritsS[gritspt-6 : gritspt+1]
//line parser/parser.y:103
		{
			gritsVAL.form = process.NewSelect(gritsDollar[1].name, process.Label{L: gritsDollar[3].strval}, gritsDollar[5].name)
		}
	case 19:
		gritsDollar = gritsS[gritspt-5 : gritspt+1]
//line parser/parser.y:105
		{
			gritsVAL.form = process.NewCase(gritsDollar[2].name, gritsDollar[4].branches)
		}
	case 20:
		gritsDollar = gritsS[gritspt-6 : gritspt+1]
//line parser/parser.y:107
		{
			gritsVAL.form = process.NewNew(gritsDollar[1].name, gritsDollar[4].form, gritsDollar[6].form)
		}
	case 21:
		gritsDollar = gritsS[gritspt-8 : gritspt+1]
//line parser/parser.y:109
		{
			gritsVAL.form = process.NewNew(process.Name{Ident: gritsDollar[1].strval, Type: gritsDollar[3].sessionType, IsSelf: false}, gritsDollar[6].form, gritsDollar[8].form)
		}
	case 22:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:111
		{
			gritsVAL.form = process.NewCall(gritsDollar[1].strval, gritsDollar[3].names)
		}
	case 23:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:113
		{
			gritsVAL.form = process.NewClose(gritsDollar[2].name)
		}
	case 24:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:115
		{
			gritsVAL.form = process.NewForward(gritsDollar[2].name, gritsDollar[3].name)
		}
	case 25:
		gritsDollar = gritsS[gritspt-10 : gritspt+1]
//line parser/parser.y:117
		{
			gritsVAL.form = process.NewSplit(gritsDollar[2].name, gritsDollar[4].name, gritsDollar[8].name, gritsDollar[10].form)
		}
	case 26:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:119
		{
			gritsVAL.form = process.NewWait(gritsDollar[2].name, gritsDollar[4].form)
		}
	case 27:
		gritsDollar = gritsS[gritspt-5 : gritspt+1]
//line parser/parser.y:121
		{
			gritsVAL.form = process.NewCast(gritsDollar[2].name, gritsDollar[4].name)
		}
	case 28:
		gritsDollar = gritsS[gritspt-6 : gritspt+1]
//line parser/parser.y:123
		{
			gritsVAL.form = process.NewShift(gritsDollar[1].name, gritsDollar[4].name, gritsDollar[6].form)
		}
	case 29:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:125
		{
			gritsVAL.form = process.NewDrop(gritsDollar[2].name, gritsDollar[4].form)
		}
	case 30:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:127
		{
			gritsVAL.form = gritsDollar[2].form
		}
	case 31:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:129
		{
			gritsVAL.form = process.NewPrint(process.Label{L: gritsDollar[2].strval}, gritsDollar[4].form)
		}
	case 32:
		gritsDollar = gritsS[gritspt-0 : gritspt+1]
//line parser/parser.y:133
		{
			gritsVAL.branches = nil
		}
	case 33:
		gritsDollar = gritsS[gritspt-6 : gritspt+1]
//line parser/parser.y:134
		{
			gritsVAL.branches = []*process.BranchForm{process.NewBranch(process.Label{L: gritsDollar[1].strval}, gritsDollar[3].name, gritsDollar[6].form)}
		}
	case 34:
		gritsDollar = gritsS[gritspt-8 : gritspt+1]
//line parser/parser.y:135
		{
			gritsVAL.branches = append(gritsDollar[1].branches, process.NewBranch(process.Label{L: gritsDollar[3].strval}, gritsDollar[5].name, gritsDollar[8].form))
		}
	case 35:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:137
		{
			gritsVAL.names = []process.Name{gritsDollar[1].name}
		}
	case 36:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:138
		{
			gritsVAL.names = append([]process.Name{gritsDollar[1].name}, gritsDollar[3].names...)
		}
	case 37:
		gritsDollar = gritsS[gritspt-0 : gritspt+1]
//line parser/parser.y:140
		{
			gritsVAL.names = nil
		}
	case 38:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:141
		{
			gritsVAL.names = []process.Name{gritsDollar[1].name}
		}
	case 39:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:142
		{
			gritsVAL.names = append([]process.Name{gritsDollar[1].name}, gritsDollar[3].names...)
		}
	case 40:
		gritsDollar = gritsS[gritspt-0 : gritspt+1]
//line parser/parser.y:145
		{
			gritsVAL.names = nil
		}
	case 41:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:146
		{
			gritsVAL.names = []process.Name{gritsDollar[1].name}
		}
	case 42:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:147
		{
			gritsVAL.names = append([]process.Name{gritsDollar[1].name}, gritsDollar[3].names...)
		}
	case 43:
		gritsDollar = gritsS[gritspt-0 : gritspt+1]
//line parser/parser.y:150
		{
			gritsVAL.names = nil
		}
	case 44:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:151
		{
			gritsVAL.names = gritsDollar[2].names
		}
	case 45:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:155
		{
			gritsVAL.names = []process.Name{gritsDollar[1].name}
		}
	case 46:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:156
		{
			gritsVAL.names = append([]process.Name{gritsDollar[1].name}, gritsDollar[3].names...)
		}
	case 47:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:161
		{
			gritsVAL.name = process.Name{Ident: gritsDollar[1].strval, IsSelf: false}
		}
	case 48:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:163
		{
			gritsVAL.name = process.Name{Ident: gritsDollar[1].strval, Type: gritsDollar[3].sessionType, IsSelf: false}
		}
	case 49:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:165
		{
			gritsVAL.name = process.Name{IsSelf: true}
		}
	case 50:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:167
		{
			pol := gritsDollar[1].polarity
			gritsVAL.name = process.Name{IsSelf: true, ExplicitPolarity: &pol}
		}
	case 51:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:169
		{
			gritsVAL.name = process.Name{Ident: gritsDollar[1].strval, IsSelf: false}
		}
	case 52:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:171
		{
			pol := gritsDollar[1].polarity
			gritsVAL.name = process.Name{Ident: gritsDollar[2].strval, IsSelf: false, ExplicitPolarity: &pol}
		}
	case 53:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:175
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{kind: ASSUMING_DEF, assumedFreeNameTypes: gritsDollar[2].names, position: gritsVAL.currPosition}
		}
	case 54:
		gritsDollar = gritsS[gritspt-7 : gritspt+1]
//line parser/parser.y:180
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: gritsDollar[2].strval, Parameters: gritsDollar[4].names, Body: gritsDollar[7].form, UsesExplicitProvider: false}, position: gritsVAL.currPosition}
		}
	case 55:
		gritsDollar = gritsS[gritspt-9 : gritspt+1]
//line parser/parser.y:182
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: gritsDollar[2].strval, Parameters: gritsDollar[4].names, Body: gritsDollar[9].form, Type: gritsDollar[7].sessionType, UsesExplicitProvider: false}, position: gritsVAL.currPosition}
		}
	case 56:
		gritsDollar = gritsS[gritspt-8 : gritspt+1]
//line parser/parser.y:185
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{
				FunctionName:         gritsDollar[2].strval,
				Parameters:           gritsDollar[5].names,
				Body:                 gritsDollar[8].form,
				UsesExplicitProvider: true,
				ExplicitProvider:     process.Name{Ident: gritsDollar[4].strval, IsSelf: true},
				// Type: $6,
			}, position: gritsVAL.currPosition}
		}
	case 57:
		gritsDollar = gritsS[gritspt-10 : gritspt+1]
//line parser/parser.y:196
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{
				FunctionName:         gritsDollar[2].strval,
				Parameters:           gritsDollar[7].names,
				Body:                 gritsDollar[10].form,
				UsesExplicitProvider: true,
				ExplicitProvider:     process.Name{Ident: gritsDollar[4].strval, IsSelf: true},
				Type:                 gritsDollar[6].sessionType}, position: gritsVAL.currPosition}
		}
	case 58:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:206
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{
				kind:         TYPE_DEF,
				session_type: types.SessionTypeDefinition{Name: gritsDollar[2].strval, SessionType: gritsDollar[4].sessionType},
				position:     gritsVAL.currPosition}
		}
	case 59:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:213
		{
			gritsVAL.sessionType = types.ConvertSessionTypeInitialToSessionType(gritsDollar[1].sessionTypeInitial)
		}
	case 60:
		gritsDollar = gritsS[gritspt-2 : gritspt+1]
//line parser/parser.y:215
		{
			mode := types.StringToMode(gritsDollar[1].strval)
			gritsVAL.sessionType = types.ConvertSessionTypeInitialToSessionType(types.NewExplicitModeTypeInitial(mode, gritsDollar[2].sessionTypeInitial))
		}
	case 61:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:221
		{
			gritsVAL.sessionTypeInitial = types.NewLabelTypeInitial(gritsDollar[1].strval)
		}
	case 62:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:223
		{
			gritsVAL.sessionTypeInitial = types.NewUnitTypeInitial()
		}
	case 63:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:225
		{
			gritsVAL.sessionTypeInitial = types.NewSelectLabelTypeInitial(gritsDollar[3].sessionTypeAltInitial)
		}
	case 64:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:227
		{
			gritsVAL.sessionTypeInitial = types.NewBranchCaseTypeInitial(gritsDollar[3].sessionTypeAltInitial)
		}
	case 65:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:229
		{
			gritsVAL.sessionTypeInitial = types.NewSendTypeInitial(gritsDollar[1].sessionTypeInitial, gritsDollar[3].sessionTypeInitial)
		}
	case 66:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:231
		{
			gritsVAL.sessionTypeInitial = types.NewReceiveTypeInitial(gritsDollar[1].sessionTypeInitial, gritsDollar[3].sessionTypeInitial)
		}
	case 67:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:233
		{
			gritsVAL.sessionTypeInitial = gritsDollar[2].sessionTypeInitial
		}
	case 68:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:235
		{
			modeFrom := types.StringToMode(gritsDollar[1].strval)
			modeTo := types.StringToMode(gritsDollar[3].strval)
			gritsVAL.sessionTypeInitial = types.NewUpTypeInitial(modeFrom, modeTo, gritsDollar[4].sessionTypeInitial)
		}
	case 69:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:239
		{
			modeFrom := types.StringToMode(gritsDollar[1].strval)
			modeTo := types.StringToMode(gritsDollar[3].strval)
			gritsVAL.sessionTypeInitial = types.NewDownTypeInitial(modeFrom, modeTo, gritsDollar[4].sessionTypeInitial)
		}
	case 70:
		gritsDollar = gritsS[gritspt-3 : gritspt+1]
//line parser/parser.y:245
		{
			gritsVAL.sessionTypeAltInitial = []types.OptionInitial{*types.NewOptionInitial(gritsDollar[1].strval, gritsDollar[3].sessionTypeInitial)}
		}
	case 71:
		gritsDollar = gritsS[gritspt-5 : gritspt+1]
//line parser/parser.y:247
		{
			gritsVAL.sessionTypeAltInitial = append([]types.OptionInitial{*types.NewOptionInitial(gritsDollar[1].strval, gritsDollar[3].sessionTypeInitial)}, gritsDollar[5].sessionTypeAltInitial...)
		}
	case 72:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:249
		{
			gritsVAL.strval = gritsDollar[1].strval
		}
	case 73:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:251
		{
			gritsVAL.polarity = types.POSITIVE
		}
	case 74:
		gritsDollar = gritsS[gritspt-1 : gritspt+1]
//line parser/parser.y:252
		{
			gritsVAL.polarity = types.NEGATIVE
		}
	case 75:
		gritsDollar = gritsS[gritspt-4 : gritspt+1]
//line parser/parser.y:256
		{
			gritsVAL.common_type = unexpandedProcessOrFunction{
				kind:     EXEC_DEF,
				proc:     incompleteProcess{Body: process.NewCall(gritsDollar[2].strval, []process.Name{})},
				position: gritsVAL.currPosition}
		}
	}
	goto gritsstack /* stack new state and value */
}
