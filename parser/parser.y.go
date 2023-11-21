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
	"phi/process"
	"phi/types"
)

//line parser/parser.y:14
type phiSymType struct {
	yys            int
	strval         string
	common_type    unexpandedProcessOrFunction
	statements     []unexpandedProcessOrFunction
	name           process.Name
	names          []process.Name
	form           process.Form
	branches       []*process.BranchForm
	sessionType    types.SessionType
	sessionTypeAlt []types.BranchOption
	polarity       process.Polarity
}

const LABEL = 57346
const LEFT_ARROW = 57347
const RIGHT_ARROW = 57348
const EQUALS = 57349
const DOT = 57350
const SEQUENCE = 57351
const COLON = 57352
const COMMA = 57353
const LPAREN = 57354
const RPAREN = 57355
const LSBRACK = 57356
const RSBRACK = 57357
const LANGLE = 57358
const RANGLE = 57359
const PIPE = 57360
const SEND = 57361
const RECEIVE = 57362
const CASE = 57363
const CLOSE = 57364
const WAIT = 57365
const CAST = 57366
const SHIFT = 57367
const ACCEPT = 57368
const ACQUIRE = 57369
const DETACH = 57370
const RELEASE = 57371
const DROP = 57372
const SPLIT = 57373
const PUSH = 57374
const NEW = 57375
const SNEW = 57376
const TYPE = 57377
const LET = 57378
const IN = 57379
const END = 57380
const SPRC = 57381
const PRC = 57382
const FORWARD = 57383
const SELF = 57384
const PRINT = 57385
const PLUS = 57386
const MINUS = 57387
const TIMES = 57388
const AMPERSAND = 57389
const UNIT = 57390
const LCBRACK = 57391
const RCBRACK = 57392
const LOLLI = 57393
const PERCENTAGE = 57394

var phiToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"LABEL",
	"LEFT_ARROW",
	"RIGHT_ARROW",
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
}

var phiStatenames = [...]string{}

const phiEofCode = 1
const phiErrCode = 2
const phiInitialStackSize = 16

//line parser/parser.y:239

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
}

const phiPrivate = 57344

const phiLast = 271

var phiAct = [...]uint8{
	3, 120, 169, 153, 102, 133, 73, 187, 63, 57,
	170, 88, 154, 171, 148, 184, 10, 42, 135, 86,
	87, 146, 90, 123, 89, 88, 26, 25, 116, 7,
	145, 24, 82, 179, 194, 27, 29, 28, 32, 22,
	23, 37, 38, 39, 40, 41, 87, 43, 54, 35,
	87, 88, 87, 87, 81, 88, 87, 88, 88, 55,
	87, 88, 58, 165, 64, 88, 66, 53, 67, 96,
	62, 98, 91, 99, 94, 21, 159, 74, 22, 23,
	78, 79, 138, 111, 113, 83, 36, 106, 103, 114,
	160, 126, 122, 108, 164, 64, 95, 118, 119, 97,
	124, 161, 60, 117, 115, 61, 59, 80, 129, 69,
	109, 139, 50, 141, 22, 23, 190, 144, 156, 100,
	75, 47, 76, 74, 130, 125, 92, 71, 149, 65,
	4, 74, 33, 56, 34, 134, 135, 136, 150, 174,
	158, 163, 155, 131, 157, 143, 167, 107, 103, 44,
	45, 46, 172, 151, 162, 127, 152, 168, 128, 200,
	101, 93, 173, 176, 51, 147, 132, 182, 203, 192,
	191, 183, 186, 166, 185, 140, 189, 112, 188, 110,
	72, 70, 195, 68, 196, 175, 197, 77, 199, 198,
	177, 178, 201, 202, 180, 9, 204, 31, 121, 181,
	30, 206, 205, 16, 207, 208, 137, 6, 104, 193,
	5, 142, 8, 11, 13, 14, 105, 85, 52, 49,
	48, 15, 2, 1, 84, 20, 26, 25, 19, 9,
	18, 24, 12, 21, 17, 22, 23, 16, 0, 0,
	0, 6, 0, 0, 5, 0, 8, 11, 13, 14,
	0, 0, 0, 0, 0, 15, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 12, 21, 17, 22,
	23,
}

var phiPact = [...]int16{
	191, -1000, -1000, -1000, -1000, 33, 33, 192, 33, 122,
	45, 33, 33, 33, 33, 33, 225, 33, -9, -9,
	-9, -1000, -1000, -1000, 107, 216, 215, 96, -1000, 153,
	214, 34, 121, 58, 33, 117, 33, -1000, 33, 174,
	93, 172, 114, 171, -1000, -1000, -1000, 33, 108, 180,
	33, 33, 91, 225, -1, 33, 213, 14, -1000, -1000,
	-25, -27, 58, 113, 150, 33, 33, -1000, 225, 33,
	225, -1000, 225, 104, 149, 204, 212, 58, 136, 76,
	33, 170, 225, 168, 71, 88, -5, 58, 58, 194,
	194, 10, -1000, 33, 112, -1000, -1000, 74, -1000, -1000,
	148, 33, 111, 132, 156, 125, -26, 33, 201, 65,
	225, 166, 225, -1000, 207, 33, 225, -3, -26, -40,
	-29, 155, -36, -1000, -1000, -1000, -1000, 225, 58, -1000,
	146, 204, 58, 103, 58, 204, 59, 70, -1000, -1000,
	225, -1000, 78, 46, 164, 225, -1000, 58, -1000, -42,
	6, 225, 58, -1000, 128, -26, 178, 7, -1000, -1000,
	33, 33, 2, -1000, 33, 193, 225, 162, 4, -1000,
	204, 225, -1000, 0, 204, 225, 101, 161, 160, 33,
	17, 225, -1000, 225, 194, -1000, -42, 225, -1000, -1000,
	152, 225, 225, 159, 190, -1000, -1000, -1000, -1000, -42,
	225, -1000, -1000, 225, 225, -1000, -1000, -1000, -1000,
}

var phiPgo = [...]uint8{
	0, 130, 230, 228, 225, 0, 29, 12, 6, 8,
	4, 5, 3, 2, 224, 9, 1, 16, 223, 222,
}

var phiR1 = [...]int8{
	0, 18, 19, 19, 1, 1, 1, 1, 1, 1,
	2, 2, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 14, 14, 14, 8, 8, 9, 9,
	9, 10, 10, 10, 11, 11, 12, 12, 7, 7,
	6, 6, 13, 13, 3, 3, 3, 3, 4, 15,
	15, 15, 15, 15, 15, 15, 16, 16, 17, 17,
}

var phiR2 = [...]int8{
	0, 1, 1, 1, 1, 2, 1, 2, 1, 2,
	7, 9, 7, 10, 6, 5, 6, 8, 7, 9,
	4, 5, 2, 3, 4, 10, 11, 4, 5, 6,
	4, 3, 4, 0, 6, 8, 1, 3, 0, 1,
	3, 0, 1, 3, 0, 2, 1, 3, 1, 3,
	1, 1, 0, 2, 7, 10, 8, 10, 4, 1,
	1, 4, 4, 3, 3, 3, 3, 5, 1, 1,
}

var phiChk = [...]int16{
	-1000, -18, -19, -5, -1, 19, 16, -6, 21, 4,
	-17, 22, 41, 23, 24, 30, 12, 43, -2, -3,
	-4, 42, 44, 45, 40, 36, 35, -6, 4, -6,
	8, 5, -6, 10, 12, 4, 41, -6, -6, -6,
	-6, -6, -5, -6, -1, -1, -1, 14, 4, 4,
	16, 11, 4, 33, -17, 25, 12, -15, 4, 48,
	44, 47, 12, -9, -6, 12, -6, -6, 9, 16,
	9, 13, 9, -8, -6, 12, 14, 7, -6, -6,
	16, -5, 33, -6, -14, 4, 5, 46, 51, 49,
	49, -15, 13, 11, -9, -6, -5, -6, -5, -5,
	15, 11, -10, -7, 4, 4, -15, 11, 17, -6,
	9, -5, 9, 13, 18, 16, 33, -17, -15, -15,
	-16, 4, -16, 13, -8, 13, 17, 7, 10, -8,
	13, 11, 10, -11, 10, 11, -6, 5, 17, -5,
	9, -5, 4, -6, -5, 33, 50, 10, 50, -5,
	-15, 7, 10, -12, -7, -15, 15, -15, -10, 17,
	20, 31, -17, -5, 16, 17, 9, -5, -15, -13,
	52, 7, -5, -15, 11, 7, -11, -6, -6, 31,
	-6, 6, -5, 9, 11, -12, -5, 7, -12, -5,
	15, 9, 9, -6, 17, -5, -5, -16, -13, -5,
	7, -5, -5, 9, 6, -13, -5, -5, -5,
}

var phiDef = [...]int8{
	0, -2, 1, 2, 3, 0, 0, 0, 0, 51,
	0, 0, 0, 0, 0, 0, 0, 0, 4, 6,
	8, 50, 68, 69, 0, 0, 0, 0, 51, 0,
	0, 0, 0, 0, 38, 0, 0, 22, 0, 0,
	0, 0, 0, 0, 5, 7, 9, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 33, 0, 59, 60,
	0, 0, 0, 0, 39, 38, 0, 23, 0, 0,
	0, 31, 0, 0, 36, 41, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 20, 0, 0, 24, 27, 0, 30, 32,
	0, 0, 0, 42, 48, 44, 58, 0, 0, 0,
	0, 0, 0, 15, 0, 0, 0, 0, 63, 64,
	0, 0, 0, 65, 40, 21, 28, 0, 0, 37,
	0, 0, 0, 0, 0, 41, 0, 0, 14, 16,
	0, 29, 0, 0, 0, 0, 61, 0, 62, 52,
	0, 0, 0, 43, 46, 49, 0, 44, 45, 12,
	0, 0, 0, 18, 0, 0, 0, 0, 66, 10,
	0, 0, 54, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 17, 0, 0, 53, 52, 0, 47, 56,
	0, 0, 0, 0, 0, 34, 19, 67, 11, 52,
	0, 13, 25, 0, 0, 55, 57, 26, 35,
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
	52,
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
//line parser/parser.y:56
		{
		}
	case 2:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:62
		{
			philex.(*lexer).processesOrFunctionsRes = append(philex.(*lexer).processesOrFunctionsRes, unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body: phiDollar[1].form, Providers: []process.Name{{Ident: "root", IsSelf: false}}}})
		}
	case 3:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:66
		{
			philex.(*lexer).processesOrFunctionsRes = phiDollar[1].statements
		}
	case 4:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:72
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 5:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:73
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 6:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:74
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 7:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:75
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 8:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:76
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 9:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:77
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 10:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:83
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body: phiDollar[6].form, Providers: phiDollar[3].names}, freeNamesWithType: phiDollar[7].names}
		}
	case 11:
		phiDollar = phiS[phipt-9 : phipt+1]
//line parser/parser.y:85
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body: phiDollar[8].form, Type: phiDollar[6].sessionType, Providers: phiDollar[3].names}, freeNamesWithType: phiDollar[9].names}
		}
	case 12:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:91
		{
			phiVAL.form = process.NewSend(phiDollar[2].name, phiDollar[4].name, phiDollar[6].name)
		}
	case 13:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:96
		{
			phiVAL.form = process.NewReceive(phiDollar[2].name, phiDollar[4].name, phiDollar[8].name, phiDollar[10].form)
		}
	case 14:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:98
		{
			phiVAL.form = process.NewSelect(phiDollar[1].name, process.Label{L: phiDollar[3].strval}, phiDollar[5].name)
		}
	case 15:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:100
		{
			phiVAL.form = process.NewCase(phiDollar[2].name, phiDollar[4].branches)
		}
	case 16:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:102
		{
			phiVAL.form = process.NewNew(phiDollar[1].name, phiDollar[4].form, phiDollar[6].form, process.UNKNOWN)
		}
	case 17:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:104
		{
			phiVAL.form = process.NewNew(process.Name{Ident: phiDollar[1].strval, Type: phiDollar[3].sessionType, IsSelf: false}, phiDollar[6].form, phiDollar[8].form, process.UNKNOWN)
		}
	case 18:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:106
		{
			phiVAL.form = process.NewNew(phiDollar[1].name, phiDollar[5].form, phiDollar[7].form, phiDollar[3].polarity)
		}
	case 19:
		phiDollar = phiS[phipt-9 : phipt+1]
//line parser/parser.y:108
		{
			phiVAL.form = process.NewNew(process.Name{Ident: phiDollar[1].strval, Type: phiDollar[3].sessionType, IsSelf: false}, phiDollar[7].form, phiDollar[9].form, phiDollar[5].polarity)
		}
	case 20:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:110
		{
			phiVAL.form = process.NewCall(phiDollar[1].strval, phiDollar[3].names, process.UNKNOWN)
		}
	case 21:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:112
		{
			phiVAL.form = process.NewCall(phiDollar[2].strval, phiDollar[4].names, phiDollar[1].polarity)
		}
	case 22:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:114
		{
			phiVAL.form = process.NewClose(phiDollar[2].name)
		}
	case 23:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:116
		{
			phiVAL.form = process.NewForward(phiDollar[2].name, phiDollar[3].name, process.UNKNOWN)
		}
	case 24:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:118
		{
			phiVAL.form = process.NewForward(phiDollar[3].name, phiDollar[4].name, phiDollar[1].polarity)
		}
	case 25:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:120
		{
			phiVAL.form = process.NewSplit(phiDollar[2].name, phiDollar[4].name, phiDollar[8].name, phiDollar[10].form, process.UNKNOWN)
		}
	case 26:
		phiDollar = phiS[phipt-11 : phipt+1]
//line parser/parser.y:122
		{
			phiVAL.form = process.NewSplit(phiDollar[2].name, phiDollar[4].name, phiDollar[9].name, phiDollar[11].form, phiDollar[7].polarity)
		}
	case 27:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:124
		{
			phiVAL.form = process.NewWait(phiDollar[2].name, phiDollar[4].form)
		}
	case 28:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:126
		{
			phiVAL.form = process.NewCast(phiDollar[2].name, phiDollar[4].name)
		}
	case 29:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:128
		{
			phiVAL.form = process.NewShift(phiDollar[1].name, phiDollar[4].name, phiDollar[6].form)
		}
	case 30:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:130
		{
			phiVAL.form = process.NewDrop(phiDollar[2].name, phiDollar[4].form)
		}
	case 31:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:132
		{
			phiVAL.form = phiDollar[2].form
		}
	case 32:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:141
		{
			phiVAL.form = process.NewPrint(phiDollar[2].name, phiDollar[4].form)
		}
	case 33:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:143
		{
			phiVAL.branches = nil
		}
	case 34:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:144
		{
			phiVAL.branches = []*process.BranchForm{process.NewBranch(process.Label{L: phiDollar[1].strval}, phiDollar[3].name, phiDollar[6].form)}
		}
	case 35:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:145
		{
			phiVAL.branches = append(phiDollar[1].branches, process.NewBranch(process.Label{L: phiDollar[3].strval}, phiDollar[5].name, phiDollar[8].form))
		}
	case 36:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:148
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 37:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:149
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 38:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:151
		{
			phiVAL.names = nil
		}
	case 39:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:152
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 40:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:153
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 41:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:156
		{
			phiVAL.names = nil
		}
	case 42:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:157
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 43:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:158
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 44:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:161
		{
			phiVAL.names = nil
		}
	case 45:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:162
		{
			phiVAL.names = phiDollar[2].names
		}
	case 46:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:166
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 47:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:167
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 48:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:172
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, IsSelf: false}
		}
	case 49:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:174
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, Type: phiDollar[3].sessionType, IsSelf: false}
		}
	case 50:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:176
		{
			phiVAL.name = process.Name{IsSelf: true}
		}
	case 51:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:177
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, IsSelf: false}
		}
	case 52:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:181
		{
			phiVAL.names = nil
		}
	case 53:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:182
		{
			phiVAL.names = phiDollar[2].names
		}
	case 54:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:188
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: phiDollar[2].strval, Parameters: phiDollar[4].names, Body: phiDollar[7].form, UsesExplicitProvider: false}}
		}
	case 55:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:190
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: phiDollar[2].strval, Parameters: phiDollar[4].names, Body: phiDollar[9].form, Type: phiDollar[7].sessionType, UsesExplicitProvider: false}}
		}
	case 56:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:193
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{
				FunctionName:         phiDollar[2].strval,
				Parameters:           phiDollar[5].names,
				Body:                 phiDollar[8].form,
				UsesExplicitProvider: true,
				ExplicitProvider:     process.Name{Ident: phiDollar[4].strval, IsSelf: true},
				// Type: $6,
			}}
		}
	case 57:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:204
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{
				FunctionName:         phiDollar[2].strval,
				Parameters:           phiDollar[7].names,
				Body:                 phiDollar[10].form,
				UsesExplicitProvider: true,
				ExplicitProvider:     process.Name{Ident: phiDollar[4].strval, IsSelf: true},
				Type:                 phiDollar[6].sessionType}}
		}
	case 58:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:214
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: TYPE_DEF, session_type: types.SessionTypeDefinition{Name: phiDollar[2].strval, SessionType: phiDollar[4].sessionType}}
		}
	case 59:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:218
		{
			phiVAL.sessionType = types.NewLabelType(phiDollar[1].strval)
		}
	case 60:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:220
		{
			phiVAL.sessionType = types.NewUnitType()
		}
	case 61:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:222
		{
			phiVAL.sessionType = types.NewSelectType(phiDollar[3].sessionTypeAlt)
		}
	case 62:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:224
		{
			phiVAL.sessionType = types.NewBranchCaseType(phiDollar[3].sessionTypeAlt)
		}
	case 63:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:226
		{
			phiVAL.sessionType = types.NewSendType(phiDollar[1].sessionType, phiDollar[3].sessionType)
		}
	case 64:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:228
		{
			phiVAL.sessionType = types.NewReceiveType(phiDollar[1].sessionType, phiDollar[3].sessionType)
		}
	case 65:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:230
		{
			phiVAL.sessionType = phiDollar[2].sessionType
		}
	case 66:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:233
		{
			phiVAL.sessionTypeAlt = []types.BranchOption{*types.NewBranchOption(phiDollar[1].strval, phiDollar[3].sessionType)}
		}
	case 67:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:234
		{
			phiVAL.sessionTypeAlt = append([]types.BranchOption{*types.NewBranchOption(phiDollar[1].strval, phiDollar[3].sessionType)}, phiDollar[5].sessionTypeAlt...)
		}
	case 68:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:237
		{
			phiVAL.polarity = process.POSITIVE
		}
	case 69:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:238
		{
			phiVAL.polarity = process.NEGATIVE
		}
	}
	goto phistack /* stack new state and value */
}
