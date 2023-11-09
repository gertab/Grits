%{
// Run this after each change:
// goyacc -p phi -o parser/parser.y.go parser/parser.y
package parser

import (
	"io"
	"phi/process"
	"phi/types"
)

%}

%union {
	strval 			string
	common_type		unexpandedProcessOrFunction
	statements 		[]unexpandedProcessOrFunction
	name 			process.Name
	names 			[]process.Name
	form 			process.Form
	branches 		[]*process.BranchForm
	sessionType 	types.SessionType
	sessionTypeAlt 	[]types.BranchOption
}

%token LABEL LEFT_ARROW RIGHT_ARROW EQUALS DOT SEQUENCE COLON COMMA LPAREN RPAREN LSBRACK RSBRACK LANGLE RANGLE PIPE SEND RECEIVE CASE CLOSE WAIT CAST SHIFT ACCEPT ACQUIRE DETACH RELEASE DROP SPLIT PUSH NEW SNEW TYPE LET IN END SPRC PRC FORWARD SELF PRINT PLUS MINUS TIMES AMPERSAND UNIT LCBRACK RCBRACK LOLLI
%type <strval> LABEL
%type <statements> statements 
%type <common_type> process_def
%type <common_type> function_def
%type <common_type> type_def
%type <form> expression 
%type <name> name
%type <name> name_with_type_ann
%type <names> names
%type <names> optional_names
%type <names> optional_names_with_type_ann
%type <names> names_with_type_ann
%type <branches> branches
%type <sessionType> session_type
%type <sessionTypeAlt> session_type_alts

%left SEND
%left SEQUENCE
%right TIMES
%right LOLLI

%%

root : program { };

program : 
		/* simulate a process */
		/* todo remove */
	   expression 
		{
			philex.(*lexer).processesOrFunctionsRes = append(philex.(*lexer).processesOrFunctionsRes, unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body:$1, Providers: []process.Name{{Ident: "root", IsSelf: false}}}})
		}
	 | statements 
		{ 
			philex.(*lexer).processesOrFunctionsRes = $1
		};
/*	 | LET functions IN processes END { }; */

/* A program may consist a combination of processes, function definitions and types */
statements : process_def { $$ = []unexpandedProcessOrFunction{$1} }
		   | process_def statements { $$ = append([]unexpandedProcessOrFunction{$1}, $2...) }
		   | function_def { $$ = []unexpandedProcessOrFunction{$1} }
		   | function_def statements { $$ = append([]unexpandedProcessOrFunction{$1}, $2...) };
		   | type_def { $$ = []unexpandedProcessOrFunction{$1} }
		   | type_def statements { $$ = append([]unexpandedProcessOrFunction{$1}, $2...) };

/* A process is defined using the prc keyword */
process_def : 
			/* without type - todo remove option to force types */
		    PRC LSBRACK names RSBRACK EQUALS expression  
				{ $$ = unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body:$6, Providers: $3}} }
		  | PRC LSBRACK names RSBRACK COLON session_type EQUALS expression  
				{ $$ = unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body:$8, Type: $6, Providers: $3}} };
		/*| SPRC LSBRACK names RSBRACK COLON expression 
				{ $$ = unexpandedProcessOrFunction{kind: PROCESS, proc: incompleteProcess{Body:$6, Providers: $3}} };*/

/* Expressions form the core part of a program  */
expression : /* Send */ SEND name LANGLE name COMMA name RANGLE  
					{ $$ = process.NewSend($2, $4, $6) }
/* Send Macro */
/* | SEND name LANGLE name COMMA name RANGLE SEQUENCE expression
			{ $$ = NewSendMacroForm($2, $4, $6, $9) } */
		   | /* Receive */ LANGLE name COMMA name RANGLE LEFT_ARROW RECEIVE name SEQUENCE expression 
		   			{ $$ = process.NewReceive($2, $4, $8, $10) }
		   | /* select */ name DOT LABEL LANGLE name RANGLE 
		   			{ $$ = process.NewSelect($1, process.Label{L: $3}, $5) }
		   | /* case */ CASE name LPAREN branches RPAREN 
		   			{ $$ = process.NewCase($2, $4) }
/* new without explicit polarities */	
/* | name LEFT_ARROW NEW LPAREN expression RPAREN SEQUENCE expression 
		{ $$ = process.NewNew($1, $5, $8) } */
		   | /* new (+ve) */ name LEFT_ARROW PLUS NEW LPAREN expression RPAREN SEQUENCE expression 
					{ $$ = process.NewNew($1, $6, $9, process.POSITIVE) } 
		   | /* new (-ve) */ name LEFT_ARROW MINUS NEW LPAREN expression RPAREN SEQUENCE expression 
					{ $$ = process.NewNew($1, $6, $9, process.NEGATIVE) } 
		   | /* call (+ve) */ PLUS LABEL LPAREN optional_names RPAREN
		   			{ $$ = process.NewCall($2, $4, process.POSITIVE) }
		   | /* call (-ve) */ MINUS LABEL LPAREN optional_names RPAREN
		   			{ $$ = process.NewCall($2, $4, process.NEGATIVE) }
		   | /* close */ CLOSE name
		   			{ $$ = process.NewClose($2) }
/* forward without explitit polarities */
/* | FORWARD name name
		{ $$ = process.NewForward($2, $3) } */
		   | /* forward (+ve) */ PLUS FORWARD name name
		   			{ $$ = process.NewForward($3, $4, process.POSITIVE) }
		   | /* forward (-ve) */ MINUS FORWARD name name
		   			{ $$ = process.NewForward($3, $4, process.NEGATIVE) }
/* split without explicit polarities */
/* | LANGLE name COMMA name RANGLE LEFT_ARROW SPLIT name SEQUENCE expression
	{ $$ = process.NewSplit($2, $4, $8, $10) } */
		   | /* split (+ve) */ LANGLE name COMMA name RANGLE LEFT_ARROW PLUS SPLIT name SEQUENCE expression
		   			{ $$ = process.NewSplit($2, $4, $9, $11, process.POSITIVE) }
		   | /* split (-ve) */ LANGLE name COMMA name RANGLE LEFT_ARROW MINUS SPLIT name SEQUENCE expression
		   			{ $$ = process.NewSplit($2, $4, $9, $11, process.NEGATIVE) }
		   | /* wait */ WAIT name SEQUENCE expression
		   			{ $$ = process.NewWait($2, $4) }
		   | /* Cast */ CAST name LANGLE name RANGLE  
					{ $$ = process.NewCast($2, $4) }
		   | /* Shift */ name LEFT_ARROW SHIFT name SEQUENCE expression 
		   			{ $$ = process.NewShift($1, $4, $6) }
					
					/* used for shared processes */
					/* for debugging */
/* remaining expressions:
 Drop, 
	Snew, Cast, Shift
	Acquire, Accept, Push, Detach, Release
*/
		   | /* print - for debugging */ PRINT name
		   			{ $$ = process.NewPrint($2) };
 
branches :   /* empty */         										 { $$ = nil }
         |               LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = []*process.BranchForm{process.NewBranch(process.Label{L: $1}, $3, $6)} }
         | branches PIPE LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = append($1, process.NewBranch(process.Label{L: $3}, $5, $8)) }
         ;

names : name { $$ = []process.Name{$1} }
 	  | name COMMA names { $$ = append([]process.Name{$1}, $3...) }

optional_names : /* empty */ { $$ = nil }
		| name { $$ = []process.Name{$1} }
		| name COMMA names { $$ = append([]process.Name{$1}, $3...) }

optional_names_with_type_ann : 
			/* empty */ { $$ = nil }
		| name_with_type_ann { $$ = []process.Name{$1} }
		| name_with_type_ann COMMA names_with_type_ann { $$ = append([]process.Name{$1}, $3...) }

names_with_type_ann : 
	name_with_type_ann { $$ = []process.Name{$1} }
	| name_with_type_ann COMMA names_with_type_ann { $$ = append([]process.Name{$1}, $3...) }

name_with_type_ann : 
			/* without type - todo remove option to force types */
			LABEL
					{ $$ = process.Name{Ident: $1, IsSelf: false} } | 
			LABEL COLON session_type 
			 		{ $$ = process.Name{Ident: $1, Type: $3, IsSelf: false} };

name : SELF { $$ = process.Name{IsSelf: true} }
	 | LABEL { $$ = process.Name{Ident: $1, IsSelf: false} };

function_def : 
			/* without type - todo remove option to force types */
			 LET LABEL LPAREN optional_names_with_type_ann RPAREN EQUALS expression
			{ $$ = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: $2, Parameters: $4, Body: $7}} } | 
			/* with type annotation */ LET LABEL LPAREN optional_names_with_type_ann RPAREN COLON session_type EQUALS expression
			{ $$ = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: $2, Parameters: $4, Body: $9, Type: $7}} };

type_def : TYPE LABEL EQUALS session_type
			{ $$ = unexpandedProcessOrFunction{kind: TYPE_DEF, session_type: types.SessionTypeDefinition{Name: $2, SessionType: $4}} };

session_type :  
			/* label */ LABEL
				{ $$ = types.NewLabelType($1) }
		   | /* unit */ UNIT
		   		{ $$ = types.NewUnitType()}
		   | /* select +{ } */ PLUS LCBRACK session_type_alts RCBRACK  
		   		{ $$ = types.NewSelectType($3)}
		   | /* branch &{ } */ AMPERSAND LCBRACK session_type_alts RCBRACK  
		   		{ $$ = types.NewBranchCaseType($3)}
		   | /* send A * B */ session_type TIMES session_type
		   		{ $$ = types.NewSendType($1, $3)}
		   | /* receive A -o B */ session_type LOLLI session_type
		   		{ $$ = types.NewReceiveType($1, $3)}
		   | /* brackets (A) */ LPAREN session_type RPAREN
		   		{ $$ = $2};

session_type_alts : 
			LABEL COLON session_type { $$ = []types.BranchOption{*types.NewBranchOption($1, $3)}} 
	 	  | LABEL COLON session_type COMMA session_type_alts { $$ = append([]types.BranchOption{*types.NewBranchOption($1, $3)}, $5...) };

%%

// Parse is the entry point to the parser.
func Parse(r io.Reader) (allEnvironment, error) {
	l := newLexer(r)
	phiParse(l)
	allEnvironment := allEnvironment{}
	select {
	case err := <-l.Errors:
		return  allEnvironment, err
	default:
		// allEnvironment := l
		allEnvironment.procsAndFuns = l.processesOrFunctionsRes
		return allEnvironment, nil
	}
}