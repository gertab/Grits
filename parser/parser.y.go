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
}

var phiStatenames = [...]string{}

const phiEofCode = 1
const phiErrCode = 2
const phiInitialStackSize = 16

//line parser/parser.y:200

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

const phiLast = 224

var phiAct = [...]uint8{
	3, 142, 93, 137, 66, 119, 56, 138, 120, 122,
	120, 162, 160, 121, 75, 54, 74, 166, 26, 24,
	23, 9, 165, 51, 22, 106, 34, 182, 155, 154,
	107, 32, 146, 6, 52, 53, 5, 7, 8, 12,
	14, 15, 185, 25, 27, 173, 30, 16, 108, 156,
	36, 37, 38, 39, 40, 41, 21, 145, 13, 21,
	17, 10, 11, 35, 85, 81, 87, 83, 33, 57,
	126, 59, 94, 61, 102, 62, 91, 119, 112, 100,
	98, 72, 120, 67, 119, 109, 70, 71, 119, 120,
	119, 64, 76, 120, 115, 120, 57, 82, 57, 84,
	48, 123, 86, 45, 128, 129, 130, 88, 147, 4,
	101, 153, 96, 152, 133, 97, 95, 134, 67, 148,
	127, 139, 140, 141, 144, 116, 111, 67, 42, 43,
	44, 110, 149, 150, 79, 104, 157, 124, 103, 158,
	73, 68, 60, 58, 55, 31, 132, 159, 135, 117,
	113, 136, 167, 114, 125, 99, 89, 172, 80, 49,
	161, 118, 189, 174, 175, 188, 177, 176, 169, 180,
	181, 168, 183, 151, 184, 105, 65, 186, 187, 63,
	29, 9, 190, 28, 171, 163, 164, 191, 69, 192,
	193, 194, 170, 6, 143, 92, 5, 131, 8, 12,
	14, 15, 78, 178, 179, 50, 47, 16, 46, 2,
	1, 77, 24, 23, 90, 20, 19, 22, 13, 21,
	17, 10, 11, 18,
}

var phiPact = [...]int16{
	177, -1000, -1000, -1000, -1000, 14, 14, 175, 14, 133,
	27, 22, 14, 14, 14, 14, 14, 14, -16, -16,
	-16, -1000, 89, 204, 202, 84, -1000, 148, 201, -10,
	132, 14, 131, 14, 130, 14, -1000, 14, 170, 75,
	167, -1000, -1000, -1000, -1000, 14, 129, 181, 14, 14,
	65, 128, -17, -19, 14, 198, 121, 147, 14, 14,
	14, 14, -1000, 17, 14, 17, 92, 145, 191, 68,
	144, 62, 14, 17, 126, 123, 166, 12, 32, -1000,
	14, 118, -1000, 113, -1000, -1000, 61, -1000, 143, 14,
	112, 138, 151, -41, -1000, -1000, -36, -40, 68, 14,
	149, 53, 107, 17, 17, 17, -1000, 193, 14, -1000,
	-1000, -1000, -1000, 17, 68, -1000, 141, 191, 68, 68,
	68, 190, 190, 44, 15, 88, -1000, 164, 100, 98,
	-1000, 13, 11, -1000, 42, 17, 68, -1000, 136, -41,
	-41, -43, -38, 150, -39, -1000, -1000, 14, 14, -9,
	-14, 17, 162, 159, 14, 178, 17, -1000, 38, 191,
	-1000, 68, -1000, 158, 157, 14, 14, -1000, 17, 17,
	10, 17, -1000, 17, -1000, 31, 17, 17, 156, 153,
	-1000, -1000, 176, -1000, -1000, 190, -1000, -1000, 17, 17,
	17, -1000, -1000, -1000, -1000,
}

var phiPgo = [...]uint8{
	0, 109, 223, 216, 215, 0, 37, 7, 4, 6,
	214, 3, 211, 2, 1, 210, 209,
}

var phiR1 = [...]int8{
	0, 15, 16, 16, 1, 1, 1, 1, 1, 1,
	2, 2, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 12, 12, 12, 8, 8, 9,
	9, 9, 10, 10, 10, 11, 11, 7, 7, 6,
	6, 3, 3, 4, 13, 13, 13, 13, 13, 13,
	13, 14, 14,
}

var phiR2 = [...]int8{
	0, 1, 1, 1, 1, 2, 1, 2, 1, 2,
	6, 8, 7, 10, 6, 5, 8, 9, 9, 4,
	5, 5, 2, 3, 4, 4, 10, 11, 11, 4,
	5, 6, 4, 2, 0, 6, 8, 1, 3, 0,
	1, 3, 0, 1, 3, 1, 3, 1, 3, 1,
	1, 7, 9, 4, 1, 1, 4, 4, 3, 3,
	3, 3, 5,
}

var phiChk = [...]int16{
	-1000, -15, -16, -5, -1, 19, 16, -6, 21, 4,
	44, 45, 22, 41, 23, 24, 30, 43, -2, -3,
	-4, 42, 40, 36, 35, -6, 4, -6, 8, 5,
	-6, 12, 4, 41, 4, 41, -6, -6, -6, -6,
	-6, -6, -1, -1, -1, 14, 4, 4, 16, 11,
	4, 33, 44, 45, 25, 12, -9, -6, 12, -6,
	12, -6, -6, 9, 16, 9, -8, -6, 12, 7,
	-6, -6, 16, 12, 33, 33, -6, -12, 4, 13,
	11, -9, -6, -9, -6, -5, -6, -5, 15, 11,
	-10, -7, 4, -13, 4, 48, 44, 47, 12, 11,
	17, -6, -5, 12, 12, 9, 13, 18, 16, -8,
	13, 13, 17, 7, 10, -8, 13, 11, 10, 46,
	51, 49, 49, -13, -6, 5, 17, 13, -5, -5,
	-5, 4, -6, -5, -13, 7, 10, -11, -7, -13,
	-13, -13, -14, 4, -14, 13, 17, 20, 31, 44,
	45, 9, 13, 13, 16, 17, 7, -5, -13, 11,
	50, 10, 50, -6, -6, 31, 31, -5, 9, 9,
	-6, 6, -5, 7, -11, -13, 9, 9, -6, -6,
	-5, -5, 17, -5, -5, 11, -5, -5, 9, 9,
	6, -14, -5, -5, -5,
}

var phiDef = [...]int8{
	0, -2, 1, 2, 3, 0, 0, 0, 0, 50,
	0, 0, 0, 0, 0, 0, 0, 0, 4, 6,
	8, 49, 0, 0, 0, 0, 50, 0, 0, 0,
	0, 39, 0, 0, 0, 0, 22, 0, 0, 0,
	0, 33, 5, 7, 9, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 34, 0, 40, 39, 0,
	39, 0, 23, 0, 0, 0, 0, 37, 42, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 19,
	0, 0, 24, 0, 25, 29, 0, 32, 0, 0,
	0, 43, 47, 53, 54, 55, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 15, 0, 0, 41,
	20, 21, 30, 0, 0, 38, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 14, 0, 0, 0,
	31, 0, 0, 10, 0, 0, 0, 44, 45, 48,
	58, 59, 0, 0, 0, 60, 12, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 51, 0, 0,
	56, 0, 57, 0, 0, 0, 0, 16, 0, 0,
	0, 0, 11, 0, 46, 61, 0, 0, 0, 0,
	17, 18, 0, 35, 52, 0, 13, 26, 0, 0,
	0, 62, 27, 28, 36,
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
//line parser/parser.y:50
		{
		}
	case 2:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:56
		{
			philex.(*lexer).processesOrFunctionsRes = append(philex.(*lexer).processesOrFunctionsRes, unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body: phiDollar[1].form, Providers: []process.Name{{Ident: "root", IsSelf: false}}}})
		}
	case 3:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:60
		{
			philex.(*lexer).processesOrFunctionsRes = phiDollar[1].statements
		}
	case 4:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:66
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 5:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:67
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 6:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:68
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 7:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:69
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 8:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:70
		{
			phiVAL.statements = []unexpandedProcessOrFunction{phiDollar[1].common_type}
		}
	case 9:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:71
		{
			phiVAL.statements = append([]unexpandedProcessOrFunction{phiDollar[1].common_type}, phiDollar[2].statements...)
		}
	case 10:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:77
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body: phiDollar[6].form, Providers: phiDollar[3].names}}
		}
	case 11:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:79
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body: phiDollar[8].form, Type: phiDollar[6].sessionType, Providers: phiDollar[3].names}}
		}
	case 12:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:85
		{
			phiVAL.form = process.NewSend(phiDollar[2].name, phiDollar[4].name, phiDollar[6].name)
		}
	case 13:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:90
		{
			phiVAL.form = process.NewReceive(phiDollar[2].name, phiDollar[4].name, phiDollar[8].name, phiDollar[10].form)
		}
	case 14:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:92
		{
			phiVAL.form = process.NewSelect(phiDollar[1].name, process.Label{L: phiDollar[3].strval}, phiDollar[5].name)
		}
	case 15:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:94
		{
			phiVAL.form = process.NewCase(phiDollar[2].name, phiDollar[4].branches)
		}
	case 16:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:96
		{
			phiVAL.form = process.NewNew(phiDollar[1].name, phiDollar[5].form, phiDollar[8].form, process.UNKNOWN)
		}
	case 17:
		phiDollar = phiS[phipt-9 : phipt+1]
//line parser/parser.y:98
		{
			phiVAL.form = process.NewNew(phiDollar[1].name, phiDollar[6].form, phiDollar[9].form, process.POSITIVE)
		}
	case 18:
		phiDollar = phiS[phipt-9 : phipt+1]
//line parser/parser.y:100
		{
			phiVAL.form = process.NewNew(phiDollar[1].name, phiDollar[6].form, phiDollar[9].form, process.NEGATIVE)
		}
	case 19:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:102
		{
			phiVAL.form = process.NewCall(phiDollar[1].strval, phiDollar[3].names, process.UNKNOWN)
		}
	case 20:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:104
		{
			phiVAL.form = process.NewCall(phiDollar[2].strval, phiDollar[4].names, process.POSITIVE)
		}
	case 21:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:106
		{
			phiVAL.form = process.NewCall(phiDollar[2].strval, phiDollar[4].names, process.NEGATIVE)
		}
	case 22:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:108
		{
			phiVAL.form = process.NewClose(phiDollar[2].name)
		}
	case 23:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:110
		{
			phiVAL.form = process.NewForward(phiDollar[2].name, phiDollar[3].name, process.UNKNOWN)
		}
	case 24:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:112
		{
			phiVAL.form = process.NewForward(phiDollar[3].name, phiDollar[4].name, process.POSITIVE)
		}
	case 25:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:114
		{
			phiVAL.form = process.NewForward(phiDollar[3].name, phiDollar[4].name, process.NEGATIVE)
		}
	case 26:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:116
		{
			phiVAL.form = process.NewSplit(phiDollar[2].name, phiDollar[4].name, phiDollar[8].name, phiDollar[10].form, process.UNKNOWN)
		}
	case 27:
		phiDollar = phiS[phipt-11 : phipt+1]
//line parser/parser.y:118
		{
			phiVAL.form = process.NewSplit(phiDollar[2].name, phiDollar[4].name, phiDollar[9].name, phiDollar[11].form, process.POSITIVE)
		}
	case 28:
		phiDollar = phiS[phipt-11 : phipt+1]
//line parser/parser.y:120
		{
			phiVAL.form = process.NewSplit(phiDollar[2].name, phiDollar[4].name, phiDollar[9].name, phiDollar[11].form, process.NEGATIVE)
		}
	case 29:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:122
		{
			phiVAL.form = process.NewWait(phiDollar[2].name, phiDollar[4].form)
		}
	case 30:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:124
		{
			phiVAL.form = process.NewCast(phiDollar[2].name, phiDollar[4].name)
		}
	case 31:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:126
		{
			phiVAL.form = process.NewShift(phiDollar[1].name, phiDollar[4].name, phiDollar[6].form)
		}
	case 32:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:128
		{
			phiVAL.form = process.NewDrop(phiDollar[2].name, phiDollar[4].form)
		}
	case 33:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:137
		{
			phiVAL.form = process.NewPrint(phiDollar[2].name)
		}
	case 34:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:139
		{
			phiVAL.branches = nil
		}
	case 35:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:140
		{
			phiVAL.branches = []*process.BranchForm{process.NewBranch(process.Label{L: phiDollar[1].strval}, phiDollar[3].name, phiDollar[6].form)}
		}
	case 36:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:141
		{
			phiVAL.branches = append(phiDollar[1].branches, process.NewBranch(process.Label{L: phiDollar[3].strval}, phiDollar[5].name, phiDollar[8].form))
		}
	case 37:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:144
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 38:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:145
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 39:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:147
		{
			phiVAL.names = nil
		}
	case 40:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:148
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 41:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:149
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 42:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:152
		{
			phiVAL.names = nil
		}
	case 43:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:153
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 44:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:154
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 45:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:157
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 46:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:158
		{
			phiVAL.names = append([]process.Name{phiDollar[1].name}, phiDollar[3].names...)
		}
	case 47:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:163
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, IsSelf: false}
		}
	case 48:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:165
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, Type: phiDollar[3].sessionType, IsSelf: false}
		}
	case 49:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:167
		{
			phiVAL.name = process.Name{IsSelf: true}
		}
	case 50:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:168
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval, IsSelf: false}
		}
	case 51:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:173
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: phiDollar[2].strval, Parameters: phiDollar[4].names, Body: phiDollar[7].form}}
		}
	case 52:
		phiDollar = phiS[phipt-9 : phipt+1]
//line parser/parser.y:175
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: phiDollar[2].strval, Parameters: phiDollar[4].names, Body: phiDollar[9].form, Type: phiDollar[7].sessionType}}
		}
	case 53:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:178
		{
			phiVAL.common_type = unexpandedProcessOrFunction{kind: TYPE_DEF, session_type: types.SessionTypeDefinition{Name: phiDollar[2].strval, SessionType: phiDollar[4].sessionType}}
		}
	case 54:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:182
		{
			phiVAL.sessionType = types.NewLabelType(phiDollar[1].strval)
		}
	case 55:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:184
		{
			phiVAL.sessionType = types.NewUnitType()
		}
	case 56:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:186
		{
			phiVAL.sessionType = types.NewSelectType(phiDollar[3].sessionTypeAlt)
		}
	case 57:
		phiDollar = phiS[phipt-4 : phipt+1]
//line parser/parser.y:188
		{
			phiVAL.sessionType = types.NewBranchCaseType(phiDollar[3].sessionTypeAlt)
		}
	case 58:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:190
		{
			phiVAL.sessionType = types.NewSendType(phiDollar[1].sessionType, phiDollar[3].sessionType)
		}
	case 59:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:192
		{
			phiVAL.sessionType = types.NewReceiveType(phiDollar[1].sessionType, phiDollar[3].sessionType)
		}
	case 60:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:194
		{
			phiVAL.sessionType = phiDollar[2].sessionType
		}
	case 61:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:197
		{
			phiVAL.sessionTypeAlt = []types.BranchOption{*types.NewBranchOption(phiDollar[1].strval, phiDollar[3].sessionType)}
		}
	case 62:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:198
		{
			phiVAL.sessionTypeAlt = append([]types.BranchOption{*types.NewBranchOption(phiDollar[1].strval, phiDollar[3].sessionType)}, phiDollar[5].sessionTypeAlt...)
		}
	}
	goto phistack /* stack new state and value */
}
