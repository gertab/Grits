%{
// Run this after each change:
// goyacc -p phi -o parser/parser.y.go parser/parser.y
package parser

import (
	"io"
	"phi/process"
)

%}

%union {
	strval string
	proc   incompleteProcess
	procs []incompleteProcess
	functions []process.FunctionDefinition
	name process.Name
	names []process.Name
	form process.Form
	branches []*process.BranchForm
}

%token LABEL LEFT_ARROW RIGHT_ARROW EQUALS DOT SEQUENCE COLON COMMA LPAREN RPAREN LSBRACK RSBRACK LANGLE RANGLE PIPE SEND RECEIVE CASE CLOSE WAIT CAST SHIFT ACCEPT ACQUIRE DETACH RELEASE DROP SPLIT PUSH NEW SNEW LET IN END SPRC PRC FORWARD SELF PRINT POL_POSITIVE POL_NEGATIVE
%type <strval> LABEL
%type <procs> processes 
%type <form> expression 
%type <functions> functions
%type <name> name
%type <names> names
%type <proc> process
%type <branches> branches

%left SEND
%left SEQUENCE

%%

root : program { };

program : 
		/* collect results in processesRes and functionDefinitionsRes */
	expression 
		{
			philex.(*lexer).processesRes = append(philex.(*lexer).processesRes, incompleteProcess{Body:$1, Providers: []process.Name{{Ident: "root", IsSelf: false}}})
		}
	 | processes 
		{ 
			philex.(*lexer).processesRes = $1
		}
	 | LET functions IN processes END 
		{ 
			philex.(*lexer).processesRes = $4
			philex.(*lexer).functionDefinitionsRes = $2
		};

processes : process processes { $$ = append([]incompleteProcess{$1}, $2...) }
		  | process           { $$ = []incompleteProcess{$1} }; 

process : PRC LSBRACK names RSBRACK COLON expression  { $$ = incompleteProcess{Body:$6, Providers: $3} }
		| SPRC LSBRACK names RSBRACK COLON expression { $$ = incompleteProcess{Body:$6, Providers: $3} };

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
		   | /* new (+ve) */ name LEFT_ARROW POL_POSITIVE NEW LPAREN expression RPAREN SEQUENCE expression 
					{ $$ = process.NewNew($1, $6, $9, process.POSITIVE) } 
		   | /* new (-ve) */ name LEFT_ARROW POL_NEGATIVE NEW LPAREN expression RPAREN SEQUENCE expression 
					{ $$ = process.NewNew($1, $6, $9, process.NEGATIVE) } 
		   | /* call (+ve) */ POL_POSITIVE LABEL LPAREN names RPAREN
		   			{ $$ = process.NewCall($2, $4, process.POSITIVE) }
		   | /* call (+ve) */ POL_POSITIVE LABEL LPAREN RPAREN
		   			{ $$ = process.NewCall($2, []process.Name{}, process.POSITIVE) }
		   | /* call (-ve) */ POL_NEGATIVE LABEL LPAREN names RPAREN
		   			{ $$ = process.NewCall($2, $4, process.NEGATIVE) }
		   | /* call (-ve) */ POL_NEGATIVE LABEL LPAREN RPAREN
		   			{ $$ = process.NewCall($2, []process.Name{}, process.NEGATIVE) }
		   | /* close */ CLOSE name
		   			{ $$ = process.NewClose($2) }
/* forward without explitit polarities */
/* | FORWARD name name
		{ $$ = process.NewForward($2, $3) } */
		   | /* forward (+ve) */ POL_POSITIVE FORWARD name name
		   			{ $$ = process.NewForward($3, $4, process.POSITIVE) }
		   | /* forward (-ve) */ POL_NEGATIVE FORWARD name name
		   			{ $$ = process.NewForward($3, $4, process.NEGATIVE) }
/* split without explicit polarities */
/* | LANGLE name COMMA name RANGLE LEFT_ARROW SPLIT name SEQUENCE expression
	{ $$ = process.NewSplit($2, $4, $8, $10) } */
		   | /* split (+ve) */ LANGLE name COMMA name RANGLE LEFT_ARROW POL_POSITIVE SPLIT name SEQUENCE expression
		   			{ $$ = process.NewSplit($2, $4, $9, $11, process.POSITIVE) }
		   | /* split (-ve) */ LANGLE name COMMA name RANGLE LEFT_ARROW POL_NEGATIVE SPLIT name SEQUENCE expression
		   			{ $$ = process.NewSplit($2, $4, $9, $11, process.NEGATIVE) }
		   | /* wait */ WAIT name SEQUENCE expression
		   			{ $$ = process.NewWait($2, $4) }
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
 	  | name COMMA names { $$ = append($3, $1) }

name : SELF { $$ = process.Name{IsSelf: true} };
name : LABEL { $$ = process.Name{Ident: $1, IsSelf: false} };

functions : /* empty */         										 
				{ $$ = nil }
		  | LABEL LPAREN RPAREN EQUALS expression 
				{ $$ = []process.FunctionDefinition{{FunctionName: $1, Parameters: []process.Name{}, Body: $5}} }
		  | LABEL LPAREN names RPAREN EQUALS expression functions
		  		{ $$ = append($7, process.FunctionDefinition{FunctionName: $1, Parameters: $3, Body: $6}) };

%%

// Parse is the entry point to the parser.
func Parse(r io.Reader) (unexpandedProcesses, error) {
	l := newLexer(r)
	phiParse(l)
	select {
	case err := <-l.Errors:
		return unexpandedProcesses{}, err
	default:
		unexpandedProcesses := unexpandedProcesses{procs: l.processesRes, functions: l.functionDefinitionsRes}
		return unexpandedProcesses, nil	
	}
}
