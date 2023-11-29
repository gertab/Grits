// Code generated by goyacc -p phi -o parser/parser.y.go parser/parser.y. DO NOT EDIT.

// Run this after each change:
// goyacc -p phi -o parser/parser.y.go parser/parser.y
//
//line parser/parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser/parser.y:4

import (
	"io"
	"phi/position"
	"phi/process"
	"phi/types"
)

//line parser/parser.y:15
type phiSymType struct {
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

var phiToknames = [...]string{
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

var phiStatenames = [...]string{}

const phiEofCode = 1
const phiErrCode = 2
const phiInitialStackSize = 16

//line parser/parser.y:261

// Parse is the entry point to the parser.
func Parse(r io.Reader) (allEnvironment, error) {
	l := newLexer(r)
	phiParse(l)
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
var phiExca = [...]int8{
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

const phiPrivate = 57344

const phiLast = 270

var phiAct = [...]uint8{
	3, 136, 147, 57, 98, 67, 115, 82, 66, 99,
	162, 103, 160, 9, 105, 56, 44, 104, 64, 52,
	130, 195, 175, 15, 171, 188, 63, 6, 139, 152,
	5, 7, 8, 10, 12, 13, 174, 31, 33, 172,
	36, 14, 39, 40, 41, 42, 43, 32, 45, 68,
	173, 141, 11, 22, 16, 29, 30, 26, 25, 51,
	98, 98, 24, 127, 93, 99, 99, 129, 128, 123,
	75, 92, 76, 69, 100, 27, 28, 78, 109, 106,
	111, 60, 112, 73, 192, 83, 53, 22, 116, 29,
	30, 168, 90, 91, 113, 118, 94, 120, 68, 84,
	68, 85, 145, 119, 131, 132, 121, 138, 107, 80,
	110, 148, 149, 133, 135, 71, 140, 89, 72, 70,
	65, 37, 144, 38, 124, 149, 153, 154, 146, 4,
	165, 157, 142, 166, 161, 143, 199, 122, 114, 158,
	83, 159, 108, 163, 87, 61, 83, 46, 47, 48,
	49, 50, 164, 116, 150, 68, 170, 169, 88, 194,
	68, 156, 167, 193, 176, 126, 179, 177, 125, 81,
	79, 77, 182, 190, 181, 180, 35, 187, 68, 189,
	151, 34, 191, 178, 101, 102, 202, 196, 86, 186,
	197, 198, 97, 137, 200, 201, 58, 9, 155, 134,
	203, 117, 96, 204, 183, 184, 185, 15, 62, 59,
	55, 6, 54, 2, 5, 1, 8, 10, 12, 13,
	23, 95, 74, 69, 21, 14, 101, 102, 20, 19,
	26, 25, 18, 73, 17, 24, 11, 22, 16, 29,
	30, 0, 0, 0, 0, 0, 0, 0, 27, 28,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 71, 0, 0, 72, 70,
}

var phiPact = [...]int16{
	193, -1000, -1000, -1000, -1000, 43, 43, 171, 43, 109,
	43, 43, 43, 43, 43, 9, 43, 20, 20, 20,
	20, 20, -1000, 15, 70, 208, 206, 192, 205, -1000,
	-1000, 63, -1000, 132, 204, -9, 106, 69, 43, -1000,
	43, 160, 59, 159, 94, 158, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 43, 85, 179, -1000, 131, 146, 103,
	43, 43, 53, 9, 43, 198, 187, -44, 219, -1000,
	-1000, -34, -37, 69, 93, 129, -1000, 9, 43, 9,
	-1000, 9, 77, 125, 192, 197, 69, 192, 69, 91,
	124, 50, 43, 157, 154, 48, 49, -15, 69, 69,
	-44, 195, 195, 177, 189, 189, 13, -1000, 43, -1000,
	32, -1000, -1000, 123, 43, 87, 115, 99, -1000, -1000,
	-1000, -1000, 43, 175, 10, 9, 9, -1000, 194, 43,
	9, -44, -44, 69, -1000, 69, -40, 122, -42, -1000,
	-1000, -1000, 9, 69, -1000, 121, 192, 74, 69, 192,
	5, 17, -1000, -1000, -1000, 18, 3, 153, -44, -44,
	-1000, 69, -1000, -1000, 174, 9, 69, -1000, 165, 112,
	-1000, -1000, 43, 43, 43, 183, 9, 12, 9, -1000,
	164, 9, 67, 152, 148, 2, 9, -1000, 189, -1000,
	9, -1000, 127, 9, 9, 180, -1000, -1000, -1000, 9,
	-1000, -1000, 9, -1000, -1000,
}

var phiPgo = [...]uint8{
	0, 129, 234, 232, 229, 228, 224, 0, 31, 3,
	7, 222, 6, 2, 15, 11, 221, 8, 1, 5,
	220, 215, 213,
}

var phiR1 = [...]int8{
	0, 21, 22, 22, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 2, 2, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 16, 16, 16, 10, 10, 11, 11, 11,
	12, 12, 12, 13, 13, 14, 14, 9, 9, 8,
	8, 8, 8, 5, 3, 3, 3, 3, 4, 17,
	17, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	18, 18, 15, 20, 20, 6,
}

var phiR2 = [...]int8{
	0, 1, 1, 1, 1, 2, 1, 2, 1, 2,
	1, 2, 1, 2, 6, 8, 7, 10, 6, 5,
	6, 8, 4, 2, 3, 10, 4, 5, 6, 4,
	3, 4, 0, 6, 8, 1, 3, 0, 1, 3,
	0, 1, 3, 0, 2, 1, 3, 1, 3, 1,
	2, 1, 2, 2, 7, 9, 8, 10, 4, 1,
	2, 1, 1, 4, 4, 3, 3, 3, 4, 4,
	3, 5, 1, 1, 1, 4,
}

var phiChk = [...]int16{
	-1000, -21, -22, -7, -1, 21, 18, -8, 23, 4,
	24, 43, 25, 26, 32, 14, 45, -2, -3, -4,
	-5, -6, 44, -20, 42, 38, 37, 55, 56, 46,
	47, -8, 4, -8, 10, 5, -8, 12, 14, -8,
	-8, -8, -8, -8, -7, -8, -1, -1, -1, -1,
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

var phiDef = [...]int8{
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

var phiTok1 = [...]int8{
	1,
}

var phiTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56,
}

var phiTok3 = [...]int8{
	0,
}

var phiErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	phiDebug        = 0
	phiErrorVerbose = false
)

type phiLexer interface {
	Lex(lval *phiSymType) int
	Error(s string)
}

type phiParser interface {
	Parse(phiLexer) int
	Lookahead() int
}

type phiParserImpl struct {
	lval  phiSymType
	stack [phiInitialStackSize]phiSymType
	char  int
}

func (p *phiParserImpl) Lookahead() int {
	return p.char
}

func phiNewParser() phiParser {
	return &phiParserImpl{}
}

const phiFlag = -1000

func phiTokname(c int) string {
	if c >= 1 && c-1 < len(phiToknames) {
		if phiToknames[c-1] != "" {
			return phiToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func phiStatname(s int) string {
	if s >= 0 && s < len(phiStatenames) {
		if phiStatenames[s] != "" {
			return phiStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func phiErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !phiErrorVerbose {
		return "syntax error"
	}

	for _, e := range phiErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + phiTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(phiPact[state])
	for tok := TOKSTART; tok-1 < len(phiToknames); tok++ {
		if n := base + tok; n >= 0 && n < phiLast && int(phiChk[int(phiAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if phiDef[state] == -2 {
		i := 0
		for phiExca[i] != -1 || int(phiExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; phiExca[i] >= 0; i += 2 {
			tok := int(phiExca[i])
			if tok < TOKSTART || phiExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if phiExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += phiTokname(tok)
	}
	return res
}

func philex1(lex phiLexer, lval *phiSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(phiTok1[0])
		goto out
	}
	if char < len(phiTok1) {
		token = int(phiTok1[char])
		goto out
	}
	if char >= phiPrivate {
		if char < phiPrivate+len(phiTok2) {
			token = int(phiTok2[char-phiPrivate])
			goto out
		}
	}
	for i := 0; i < len(phiTok3); i += 2 {
		token = int(phiTok3[i+0])
		if token == char {
			token = int(phiTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(phiTok2[1]) /* unknown char */
	}
	if phiDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", phiTokname(token), uint(char))
	}
	return char, token
}

func phiParse(philex phiLexer) int {
	return phiNewParser().Parse(philex)
}

func (phircvr *phiParserImpl) Parse(philex phiLexer) int {
	var phin int
	var phiVAL phiSymType
	var phiDollar []phiSymType
	_ = phiDollar // silence set and not used
	phiS := phircvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	phistate := 0
	phircvr.char = -1
	phitoken := -1 // phircvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		phistate = -1
		phircvr.char = -1
		phitoken = -1
	}()
	phip := -1
	goto phistack

ret0:
	return 0

ret1:
	return 1

phistack:
	/* put a state and value onto the stack */
	if phiDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", phiTokname(phitoken), phiStatname(phistate))
	}

	phip++
	if phip >= len(phiS) {
		nyys := make([]phiSymType, len(phiS)*2)
		copy(nyys, phiS)
		phiS = nyys
	}
	phiS[phip] = phiVAL
	phiS[phip].yys = phistate

phinewstate:
	phin = int(phiPact[phistate])
	if phin <= phiFlag {
		goto phidefault /* simple state */
	}
	if phircvr.char < 0 {
		phircvr.char, phitoken = philex1(philex, &phircvr.lval)
	}
	phin += phitoken
	if phin < 0 || phin >= phiLast {
		goto phidefault
	}
	phin = int(phiAct[phin])
	if int(phiChk[phin]) == phitoken { /* valid shift */
		phircvr.char = -1
		phitoken = -1
		phiVAL = phircvr.lval
		phistate = phin
		if Errflag > 0 {
			Errflag--
		}
		goto phistack
	}

phidefault:
	/* default state action */
	phin = int(phiDef[phistate])
	if phin == -2 {
		if phircvr.char < 0 {
			phircvr.char, phitoken = philex1(philex, &phircvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if phiExca[xi+0] == -1 && int(phiExca[xi+1]) == phistate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			phin = int(phiExca[xi+0])
			if phin < 0 || phin == phitoken {
				break
			}
		}
		phin = int(phiExca[xi+1])
		if phin < 0 {
			goto ret0
		}
	}
	if phin == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			philex.Error(phiErrorMessage(phistate, phitoken))
			Nerrs++
			if phiDebug >= 1 {
				__yyfmt__.Printf("%s", phiStatname(phistate))
				__yyfmt__.Printf(" saw %s\n", phiTokname(phitoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for phip >= 0 {
				phin = int(phiPact[phiS[phip].yys]) + phiErrCode
				if phin >= 0 && phin < phiLast {
					phistate = int(phiAct[phin]) /* simulate a shift of "error" */
					if int(phiChk[phistate]) == phiErrCode {
						goto phistack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if phiDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", phiS[phip].yys)
				}
				phip--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if phiDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", phiTokname(phitoken))
			}
			if phitoken == phiEofCode {
				goto ret1
			}
			phircvr.char = -1
			phitoken = -1
			goto phinewstate /* try again in the same state */
		}
	}

	/* reduction by production phin */
	if phiDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", phin, phiStatname(phistate))
	}

	phint := phin
	phipt := phip
	_ = phipt // guard against "declared and not used"

	phip -= int(phiR2[phin])
	// phip is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if phip+1 >= len(phiS) {
		nyys := make([]phiSymType, len(phiS)*2)
		copy(nyys, phiS)
		phiS = nyys
	}
	phiVAL = phiS[phip+1]

	/* consult goto table to find next state */
	phin = int(phiR1[phin])
	phig := int(phiPgo[phin])
	phij := phig + phiS[phip].yys + 1

	if phij >= phiLast {
		phistate = int(phiAct[phig])
	} else {
		phistate = int(phiAct[phij])
		if int(phiChk[phistate]) != -phin {
			phistate = int(phiAct[phig])
		}
	}
	// dummy call; replaced with literal code
	switch phint {

	case 1:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:58
		{
		}
	case 2:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:64
		{
			philex.(*lexer).processesOrFunctionsRes = append(philex.(*lexer).processesOrFunctionsRes, unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body: phiDollar[1].form, Providers: []process.Name{{Ident: "root", IsSelf: false}}}, position: phiVAL.currPosition})
		}
	case 3:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:68
		{
			philex.(*lexer).processesOrFunctionsRes = phiDollar[1].statements
		}
	case 4:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:74
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 5:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:75
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 6:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:76
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 7:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:77
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 8:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:78
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 9:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:79
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 10:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:80
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 11:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:81
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 12:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:82
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 13:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:83
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 14:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:89
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body: phiDollar[6].form, Providers: phiDollar[3].names}, position: phiVAL.currPosition}
		}
	case 15:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:91
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body: phiDollar[8].form, Type: phiDollar[6].sessionType, Providers: phiDollar[3].names}, position: phiVAL.currPosition}
		}
	case 16:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:97
		{
			phiVAL.form = process.NewSend(phiDollar[2].name, phiDollar[4].name, phiDollar[6].name)
		}
	case 17:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:101
		{
			phiVAL.form = process.NewReceive(phiDollar[2].name, phiDollar[4].name, phiDollar[8].name, phiDollar[10].form)
		}
	case 18:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:103
		{
			phiVAL.form = process.NewSelect(phiDollar[1].name, process.Label{L: phiDollar[3].strval}, phiDollar[5].name)
		}
	case 19:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:105
		{
			phiVAL.form = process.NewCase(phiDollar[2].name, phiDollar[4].branches)
		}
	case 20:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:107
		{
			phiVAL.form = process.NewNew(phiDollar[1].name, phiDollar[4].form, phiDollar[6].form)
		}
	case 21:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:109
		{
			phiVAL.form = process.NewNew(process.Name{Ident: phiDollar[1].strval, Type: phiDollar[3].sessionType, IsSelf: false}, phiDollar[6].form, phiDollar[8].form)
		}
	case 22:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:111
		{
			phiVAL.form = process.NewCall(phiDollar[1].strval, phiDollar[3].names)
		}
	case 23:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:113
		{
			phiVAL.form = process.NewClose(phiDollar[2].name)
		}
	case 24:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:115
		{
			phiVAL.form = process.NewForward(phiDollar[2].name, phiDollar[3].name)
		}
	case 25:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:117
		{
			phiVAL.form = process.NewSplit(phiDollar[2].name, phiDollar[4].name, phiDollar[8].name, phiDollar[10].form)
		}
	case 26:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:119
		{
			phiVAL.form = process.NewWait(phiDollar[2].name, phiDollar[4].form)
		}
	case 27:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:121
		{
			phiVAL.form = process.NewCast(phiDollar[2].name, phiDollar[4].name)
		}
	case 28:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:123
		{
			phiVAL.form = process.NewShift(phiDollar[1].name, phiDollar[4].name, phiDollar[6].form)
		}
	case 29:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:125
		{
			phiVAL.form = process.NewDrop(phiDollar[2].name, phiDollar[4].form)
		}
	case 30:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:127
		{
			phiVAL.form = phiDollar[2].form
		}
	case 31:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:129
		{
			phiVAL.form = process.NewPrint(phiDollar[2].name, phiDollar[4].form)
		}
	case 32:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:133
		{
			phiVAL.branches = nil
		}
	case 33:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:134
		{
			phiVAL.branches = []*process.BranchForm{process.NewBranch(process.Label{L: phiDollar[1].strval}, phiDollar[3].name, phiDollar[6].form)}
		}
	case 34:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:135
		{
			phiVAL.branches = append(phiDollar[1].branches, process.NewBranch(process.Label{L: phiDollar[3].strval}, phiDollar[5].name, phiDollar[8].form))
		}
	case 35:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:137
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 36:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:138
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 37:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:140
		{
			phiVAL.names = nil
		}
	case 38:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:141
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 39:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:142
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 40:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:145
		{
			phiVAL.names = nil
		}
	case 41:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:146
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 42:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:147
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 43:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:150
		{
			phiVAL.names = nil
		}
	case 44:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:151
		{
			phiVAL.names = phiDollar[2].names
		}
	case 45:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:155
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 46:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:156
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 47:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:161
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, IsSelf: false}
		}
	case 48:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:163
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, Type: phiDollar[3].sessionType, IsSelf: false}
		}
	case 49:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:165
		{
			phiVAL.name = process.Name{IsSelf: true}
		}
	case 50:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:167
		{
			pol := phiDollar[1].polarity
			phiVAL.name = process.Name{IsSelf: true, ExplicitPolarity: &pol}
		}
	case 51:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:169
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, IsSelf: false}
		}
	case 52:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:171
		{
			pol := phiDollar[1].polarity
			phiVAL.name = process.Name{Ident: phiDollar[2].strval, IsSelf: false, ExplicitPolarity: &pol}
		}
	case 53:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:175
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: ASSUMING_DEF, assumedFreeNameTypes: phiDollar[2].names, position: phiVAL.currPosition}
		}
	case 54:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:180
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: phiDollar[2].strval, Parameters: phiDollar[4].names, Body: phiDollar[7].form, UsesExplicitProvider: false}, position: phiVAL.currPosition}
		}
	case 55:
		phiDollar = phiS[phipt-9 : phipt+1]
//line parser/parser.y:182
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: phiDollar[2].strval, Parameters: phiDollar[4].names, Body: phiDollar[9].form, Type: phiDollar[7].sessionType, UsesExplicitProvider: false}, position: phiVAL.currPosition}
		}
	case 56:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:185
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{
				FunctionName:         phiDollar[2].strval,
				Parameters:           phiDollar[5].names,
				Body:                 phiDollar[8].form,
				UsesExplicitProvider: true,
				ExplicitProvider:     process.Name{Ident: phiDollar[4].strval, IsSelf: true},
				// Type: $6,
			}, position: phiVAL.currPosition}
		}
	case 57:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:196
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{
				FunctionName:         phiDollar[2].strval,
				Parameters:           phiDollar[7].names,
				Body:                 phiDollar[10].form,
				UsesExplicitProvider: true,
				ExplicitProvider:     process.Name{Ident: phiDollar[4].strval, IsSelf: true},
				Type:                 phiDollar[6].sessionType}, position: phiVAL.currPosition}
		}
	case 58:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:206
		{
			phiVAL.common_type = unexpandedProcessOrFunction{
				kind:         TYPE_DEF,
				session_type: types.SessionTypeDefinition{Name: phiDollar[2].strval, SessionType: phiDollar[4].sessionType},
				position:     phiVAL.currPosition}
		}
	case 59:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:213
		{
			phiVAL.sessionType = types.ConvertSessionTypeInitialToSessionType(phiDollar[1].sessionTypeInitial)
		}
	case 60:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:215
		{
			mode := types.StringToMode(phiDollar[1].strval)
			phiVAL.sessionType = types.ConvertSessionTypeInitialToSessionType(types.NewExplicitModeTypeInitial(mode, phiDollar[2].sessionTypeInitial))
		}
	case 61:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:221
		{
			phiVAL.sessionTypeInitial = types.NewLabelTypeInitial(phiDollar[1].strval)
		}
	case 62:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:223
		{
			phiVAL.sessionTypeInitial = types.NewUnitTypeInitial()
		}
	case 63:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:225
		{
			phiVAL.sessionTypeInitial = types.NewSelectLabelTypeInitial(phiDollar[3].sessionTypeAltInitial)
		}
	case 64:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:227
		{
			phiVAL.sessionTypeInitial = types.NewBranchCaseTypeInitial(phiDollar[3].sessionTypeAltInitial)
		}
	case 65:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:229
		{
			phiVAL.sessionTypeInitial = types.NewSendTypeInitial(phiDollar[1].sessionTypeInitial, phiDollar[3].sessionTypeInitial)
		}
	case 66:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:231
		{
			phiVAL.sessionTypeInitial = types.NewReceiveTypeInitial(phiDollar[1].sessionTypeInitial, phiDollar[3].sessionTypeInitial)
		}
	case 67:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:233
		{
			phiVAL.sessionTypeInitial = phiDollar[2].sessionTypeInitial
		}
	case 68:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:235
		{
			modeFrom := types.StringToMode(phiDollar[1].strval)
			modeTo := types.StringToMode(phiDollar[3].strval)
			phiVAL.sessionTypeInitial = types.NewUpTypeInitial(modeFrom, modeTo, phiDollar[4].sessionTypeInitial)
		}
	case 69:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:239
		{
			modeFrom := types.StringToMode(phiDollar[1].strval)
			modeTo := types.StringToMode(phiDollar[3].strval)
			phiVAL.sessionTypeInitial = types.NewDownTypeInitial(modeFrom, modeTo, phiDollar[4].sessionTypeInitial)
		}
	case 70:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:245
		{
			phiVAL.sessionTypeAltInitial = []types.OptionInitial{*types.NewOptionInitial(phiDollar[1].strval, phiDollar[3].sessionTypeInitial)}
		}
	case 71:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:247
		{
			phiVAL.sessionTypeAltInitial = append([]types.OptionInitial{*types.NewOptionInitial(phiDollar[1].strval, phiDollar[3].sessionTypeInitial)}, phiDollar[5].sessionTypeAltInitial...)
		}
	case 72:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:249
		{
			phiVAL.strval = phiDollar[1].strval
		}
	case 73:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:251
		{
			phiVAL.polarity = types.POSITIVE
		}
	case 74:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:252
		{
			phiVAL.polarity = types.NEGATIVE
		}
	case 75:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:256
		{
			phiVAL.common_type = unexpandedProcessOrFunction{
				kind:     EXEC_DEF,
				proc:     incompleteProcess{Body: process.NewCall(phiDollar[2].strval, []process.Name{})},
				position: phiVAL.currPosition}
		}
	}
	goto phistack /* stack new state and value */
}
