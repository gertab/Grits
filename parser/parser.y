%{
package parser

import (
	"io"
	"phi/process"
)

var processes []earlyProcess
var functionDefinitions []process.FunctionDefinition

%}

%union {
	strval string
	proc   earlyProcess
	procs []earlyProcess
	functions []process.FunctionDefinition
	name process.Name
	names []process.Name
	form process.Form
	branches []*process.BranchForm
}

%token LABEL LEFT_ARROW RIGHT_ARROW EQUALS DOT SEQUENCE COLON COMMA LPAREN RPAREN LSBRACK RSBRACK LANGLE RANGLE PIPE SEND RECEIVE CASE CLOSE WAIT CAST SHIFT ACCEPT ACQUIRE DETACH RELEASE DROP SPLIT PUSH NEW SNEW LET IN END SPRC PRC FORWARD
%type <strval> LABEL
%type <procs> processes 
%type <form> expression 
%type <functions> functions
%type <name> name
%type <names> names
%type <proc> process
%type <branches> branches

%%

root : program { }
    ;

program : 
	expression { processes = append(processes, earlyProcess{Body:$1, Names: []process.Name{process.Name{Ident: "root"}}}) }
	 | LET functions IN processes END { 
		processes = $4
		functionDefinitions = $2
	 };

processes : process processes { $$ = append($2, $1) }
		  | process           { $$ = []earlyProcess{$1} }; 

process : PRC LSBRACK names RSBRACK COLON expression  { $$ = earlyProcess{Body:$6, Names: $3} }
		| SPRC LSBRACK names RSBRACK COLON expression { $$ = earlyProcess{Body:$6, Names: $3} };

functions : { $$ = nil }
		  | SPRC expression { $$ = []process.FunctionDefinition{process.FunctionDefinition{Body: $2}} };

expression : /* Send */ SEND name LANGLE name COMMA name RANGLE  
					{ $$ = process.NewSend($2, $4, $6) }
		   | /* Receive */ LANGLE name COMMA name RANGLE LEFT_ARROW RECEIVE name SEQUENCE expression 
		   			{ $$ = process.NewReceive($2, $4, $8, $10) }
		   | /* select */ name DOT LABEL LANGLE name RANGLE 
		   			{ $$ = process.NewSelect($1, process.Label{L: $3}, $5) }
		   | /* case */ CASE name LPAREN branches RPAREN 
		   			{ $$ = process.NewCase($2, $4) }
		   | /* new */ name LEFT_ARROW NEW LPAREN expression RPAREN SEQUENCE expression 
		   			{ $$ = process.NewNew($1, $5, $8) }
		   | /* new */ name LEFT_ARROW NEW expression SEQUENCE expression 
		   			{ $$ = process.NewNew($1, $4, $6) };

/* remaining expressions:
Call, Split, Forward, Drop, Snew
Wait, Close, Cast, Shift
Acquire, Accept, Push, Detach, Release
*/

branches :   /* empty */         										 { $$ = nil }
         |               LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = []*process.BranchForm{process.NewBranch(process.Label{L: $1}, $3, $6)} }
         | branches PIPE LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = append($1, process.NewBranch(process.Label{L: $3}, $5, $8)) }
         ;

names : name { $$ = []process.Name{$1} }
 	  | name COMMA names { $$ = append($3, $1) }

name : LABEL { $$ = process.Name{Ident: $1} };

%%

// Parse is the entry point to the parser.
func Parse(r io.Reader) (unexpandedProcesses, error) {
	l := newLexer(r)
	phiParse(l)
	select {
	case err := <-l.Errors:
		return unexpandedProcesses{}, err
	default:
		return unexpandedProcesses{procs: processes, functions: functionDefinitions}, nil
	}
}
