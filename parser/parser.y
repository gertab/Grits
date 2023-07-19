%{
package parser

import (
	"io"
	"phi/process"
)

var processes []process.Process
var functionDefinitions []process.FunctionDefinition

%}

%union {
	strval string
	proc   process.Process
	procs []process.Process
	functions []process.FunctionDefinition
	name process.Name
	form process.Form
	branches []*process.BranchForm
}

%token LABEL LEFT_ARROW RIGHT_ARROW EQUALS DOT SEQUENCE COLON COMMA LPAREN RPAREN LSBRACK RSBRACK LANGLE RANGLE PIPE SEND RECEIVE CASE CLOSE WAIT CAST SHIFT ACCEPT ACQUIRE DETACH RELEASE DROP SPLIT PUSH NEW SNEW LET IN END SPRC PRC FORWARD
%type <strval> LABEL
%type <procs> processes 
%type <form> expression 
%type <functions> functions
%type <name> name
%type <proc> process
%type <branches> branches

%%

root : program { }
    ;

program : 
	expression { processes = append(processes, process.Process{Body:$1}) }
	 | LET functions IN processes END { 
		processes = $4
		functionDefinitions = $2
	 };

processes : process processes { $$ = append($2, $1) }
		  | process           { $$ = []process.Process{$1} }; 

process : PRC expression  { $$ = process.Process{Body:$2} }
		| SPRC expression { $$ = process.Process{Body:$2} };

functions : { $$ = nil }
		  | SPRC expression { $$ = []process.FunctionDefinition{process.FunctionDefinition{Body: $2}} };

expression : /* Send */ SEND name LANGLE name COMMA name RANGLE  
					{ $$ = process.NewSend($2, $4, $6) }
		   | /* Receive */ LANGLE name COMMA name RANGLE LEFT_ARROW RECEIVE name SEQUENCE expression 
		   			{ $$ = process.NewReceive($2, $4, $8, $10) }
		   | /* select */ name DOT LABEL LANGLE name RANGLE 
		   			{ $$ = process.NewSelect($1, process.Label{L: $3}, $5) }
		   | /* case */ CASE name LPAREN branches RPAREN 
		   			{ $$ = process.NewCase($2, $4) };

branches :   /* empty */         										 { $$ = nil }
         |               LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = []*process.BranchForm{process.NewBranch(process.Label{L: $1}, $3, $6)} }
         | branches PIPE LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = append($1, process.NewBranch(process.Label{L: $3}, $5, $8)) }
         ;

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
