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
}

%token LABEL LEFT_ARROW RIGHT_ARROW EQUALS DOT SEQUENCE COLON COMMA LPAREN RPAREN LSBRACK RSBRACK LANGLE RANGLE PIPE SEND RECEIVE CASE CLOSE WAIT CAST SHIFT ACCEPT ACQUIRE DETACH RELEASE DROP SPLIT PUSH NEW SNEW LET IN END SPRC PRC FORWARD
%type <strval> LABEL
%type <procs> processes 
%type <form> expression 
%type <functions> functions
%type <name> name
%type <proc> process

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

expression : SEND name LANGLE name COMMA name RANGLE  { $$ = process.NewSend($2, $4, $6) }
		   | LANGLE name COMMA name RANGLE LEFT_ARROW RECEIVE name SEQUENCE expression { $$ = process.NewReceive($2, $4, $8, $10) };

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
