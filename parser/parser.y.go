// Code generated by goyacc -p phi -o parser/parser.y.go parser/parser.y. DO NOT EDIT.

//line parser/parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser/parser.y:2

import (
	"io"
	"phi/process"
)

var processes []incompleteProcess
var functionDefinitions []process.FunctionDefinition

//line parser/parser.y:14
type phiSymType struct {
	yys       int
	strval    string
	proc      incompleteProcess
	procs     []incompleteProcess
	functions []process.FunctionDefinition
	name      process.Name
	names     []process.Name
	form      process.Form
	branches  []*process.BranchForm
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
const LET = 57377
const IN = 57378
const END = 57379
const SPRC = 57380
const PRC = 57381
const FORWARD = 57382

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
	"LET",
	"IN",
	"END",
	"SPRC",
	"PRC",
	"FORWARD",
}

var phiStatenames = [...]string{}

const phiEofCode = 1
const phiErrCode = 2
const phiInitialStackSize = 16

// Parse is the entry point to the parser.
//
//line parser/parser.y:89
func Parse(r io.Reader) (unexpandedProcesses, error) {
	l := newLexer(r)
	phiParse(l)
	select {
	case err := <-l.Errors:
		return unexpandedProcesses{}, err
	default:
		unexpandedProcesses := unexpandedProcesses{procs: processes, functions: functionDefinitions}
		// todo: not sure if copy is needed
		processes = nil
		functionDefinitions = nil
		return unexpandedProcesses, nil
	}
}

//line yacctab:1
var phiExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const phiPrivate = 57344

const phiLast = 99

var phiAct = [...]int8{
	3, 32, 31, 11, 13, 40, 52, 21, 26, 66,
	80, 36, 69, 29, 22, 6, 64, 65, 5, 49,
	8, 9, 7, 11, 50, 57, 45, 37, 14, 15,
	68, 18, 19, 20, 51, 6, 35, 47, 5, 10,
	8, 9, 23, 28, 41, 11, 33, 34, 62, 59,
	54, 43, 42, 58, 4, 27, 63, 6, 46, 10,
	5, 44, 8, 9, 24, 53, 53, 55, 74, 72,
	71, 77, 70, 78, 61, 79, 67, 81, 48, 17,
	82, 10, 16, 83, 84, 76, 53, 56, 11, 73,
	60, 75, 39, 25, 2, 1, 38, 30, 12,
}

var phiPact = [...]int16{
	19, -1000, -1000, -1000, -34, 84, 84, 74, 84, 84,
	84, -1000, -29, 41, 26, 53, 89, -25, 43, -1000,
	84, -37, -1000, 84, 84, 20, -1, 88, -1000, -32,
	-37, 38, 37, 50, 9, 84, 41, 69, 6, 18,
	-1000, -1000, 84, 84, 84, 82, 8, 40, 41, -1000,
	86, 84, 33, 45, 1, 0, -11, -1000, 67, -1000,
	14, -5, 62, 84, 59, -1000, 84, 41, 84, 79,
	41, -1000, 41, 66, -1000, -7, 41, -1000, -1000, 41,
	77, -1000, -1000, 41, -1000,
}

var phiPgo = [...]int8{
	0, 13, 0, 98, 22, 6, 97, 96, 95, 94,
}

var phiR1 = [...]int8{
	0, 8, 9, 9, 1, 1, 6, 6, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2, 7, 7,
	7, 5, 5, 4,
}

var phiR2 = [...]int8{
	0, 1, 1, 5, 2, 1, 6, 6, 0, 2,
	7, 10, 6, 5, 8, 6, 2, 3, 0, 6,
	8, 1, 3, 1,
}

var phiChk = [...]int16{
	-1000, -8, -9, -2, 35, 19, 16, -4, 21, 22,
	40, 4, -3, 38, -4, -4, 8, 5, -4, -4,
	-4, 36, -2, 16, 11, 4, 33, 12, -4, -1,
	-6, 39, 38, -4, -4, 16, 12, -2, -7, 4,
	37, -1, 14, 14, 11, 17, -4, -2, 9, 13,
	18, 16, -5, -4, -5, -4, 5, 17, 13, -2,
	4, -4, 15, 11, 15, 17, 20, 9, 16, 17,
	10, -5, 10, -4, -2, -4, 6, -2, -2, 9,
	17, -2, -2, 6, -2,
}

var phiDef = [...]int8{
	0, -2, 1, 2, 8, 0, 0, 0, 0, 0,
	0, 23, 0, 0, 0, 0, 0, 0, 0, 16,
	0, 0, 9, 0, 0, 0, 0, 18, 17, 0,
	5, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	3, 4, 0, 0, 0, 0, 0, 0, 0, 13,
	0, 0, 0, 21, 0, 0, 0, 12, 0, 15,
	0, 0, 0, 0, 0, 10, 0, 0, 0, 0,
	0, 22, 0, 0, 14, 0, 0, 6, 7, 0,
	0, 19, 11, 0, 20,
}

var phiTok1 = [...]int8{
	1,
}

var phiTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40,
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
//line parser/parser.y:37
		{
		}
	case 2:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:41
		{
			processes = append(processes, incompleteProcess{Body: phiDollar[1].form, Names: []process.Name{{Ident: "root"}}})
		}
	case 3:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:42
		{
			processes = phiDollar[4].procs
			functionDefinitions = phiDollar[2].functions
		}
	case 4:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:47
		{
			phiVAL.procs = append(phiDollar[2].procs, phiDollar[1].proc)
		}
	case 5:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:48
		{
			phiVAL.procs = []incompleteProcess{phiDollar[1].proc}
		}
	case 6:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:50
		{
			phiVAL.proc = incompleteProcess{Body: phiDollar[6].form, Names: phiDollar[3].names}
		}
	case 7:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:51
		{
			phiVAL.proc = incompleteProcess{Body: phiDollar[6].form, Names: phiDollar[3].names}
		}
	case 8:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:53
		{
			phiVAL.functions = nil
		}
	case 9:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:54
		{
			phiVAL.functions = []process.FunctionDefinition{{Body: phiDollar[2].form}}
		}
	case 10:
		phiDollar = phiS[phipt-7 : phipt+1]
//line parser/parser.y:57
		{
			phiVAL.form = process.NewSend(phiDollar[2].name, phiDollar[4].name, phiDollar[6].name)
		}
	case 11:
		phiDollar = phiS[phipt-10 : phipt+1]
//line parser/parser.y:59
		{
			phiVAL.form = process.NewReceive(phiDollar[2].name, phiDollar[4].name, phiDollar[8].name, phiDollar[10].form)
		}
	case 12:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:61
		{
			phiVAL.form = process.NewSelect(phiDollar[1].name, process.Label{L: phiDollar[3].strval}, phiDollar[5].name)
		}
	case 13:
		phiDollar = phiS[phipt-5 : phipt+1]
//line parser/parser.y:63
		{
			phiVAL.form = process.NewCase(phiDollar[2].name, phiDollar[4].branches)
		}
	case 14:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:65
		{
			phiVAL.form = process.NewNew(phiDollar[1].name, phiDollar[5].form, phiDollar[8].form)
		}
	case 15:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:67
		{
			phiVAL.form = process.NewNew(phiDollar[1].name, phiDollar[4].form, phiDollar[6].form)
		}
	case 16:
		phiDollar = phiS[phipt-2 : phipt+1]
//line parser/parser.y:69
		{
			phiVAL.form = process.NewClose(phiDollar[2].name)
		}
	case 17:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:71
		{
			phiVAL.form = process.NewForward(phiDollar[2].name, phiDollar[3].name)
		}
	case 18:
		phiDollar = phiS[phipt-0 : phipt+1]
//line parser/parser.y:79
		{
			phiVAL.branches = nil
		}
	case 19:
		phiDollar = phiS[phipt-6 : phipt+1]
//line parser/parser.y:80
		{
			phiVAL.branches = []*process.BranchForm{process.NewBranch(process.Label{L: phiDollar[1].strval}, phiDollar[3].name, phiDollar[6].form)}
		}
	case 20:
		phiDollar = phiS[phipt-8 : phipt+1]
//line parser/parser.y:81
		{
			phiVAL.branches = append(phiDollar[1].branches, process.NewBranch(process.Label{L: phiDollar[3].strval}, phiDollar[5].name, phiDollar[8].form))
		}
	case 21:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:84
		{
			phiVAL.names = []process.Name{phiDollar[1].name}
		}
	case 22:
		phiDollar = phiS[phipt-3 : phipt+1]
//line parser/parser.y:85
		{
			phiVAL.names = append(phiDollar[3].names, phiDollar[1].name)
		}
	case 23:
		phiDollar = phiS[phipt-1 : phipt+1]
//line parser/parser.y:87
		{
			phiVAL.name = process.Name{Ident: phiDollar[1].strval}
		}
	}
	goto phistack /* stack new state and value */
}
