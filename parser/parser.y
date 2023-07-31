%{
package parser

import (
	"io"
	"phi/process"
)

var processes []incompleteProcess
var functionDefinitions []process.FunctionDefinition

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

%token LABEL LEFT_ARROW RIGHT_ARROW EQUALS DOT SEQUENCE COLON COMMA LPAREN RPAREN LSBRACK RSBRACK LANGLE RANGLE PIPE SEND RECEIVE CASE CLOSE WAIT CAST SHIFT ACCEPT ACQUIRE DETACH RELEASE DROP SPLIT PUSH NEW SNEW LET IN END SPRC PRC FORWARD SELF
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
	expression { processes = append(processes, incompleteProcess{Body:$1, Names: []process.Name{{Ident: "root", IsSelf: false}}}) }
	 | LET functions IN processes END { 
		processes = $4
		functionDefinitions = $2
	 };

processes : process processes { $$ = append([]incompleteProcess{$1}, $2...) }
		  | process           { $$ = []incompleteProcess{$1} }; 

process : PRC LSBRACK names RSBRACK COLON expression  { $$ = incompleteProcess{Body:$6, Names: $3} }
		| SPRC LSBRACK names RSBRACK COLON expression { $$ = incompleteProcess{Body:$6, Names: $3} };

functions : { $$ = nil }
		  | SPRC expression { $$ = []process.FunctionDefinition{{Body: $2}} };

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
		   | /* close */ CLOSE name
		   			{ $$ = process.NewClose($2) };
		   | /* forward */ FORWARD name name
		   			{ $$ = process.NewForward($2, $3) };

/* remaining expressions:
Call, Split, Drop, Snew
Wait, Close, Cast, Shift
Acquire, Accept, Push, Detach, Release
*/

branches :   /* empty */         										 { $$ = nil }
         |               LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = []*process.BranchForm{process.NewBranch(process.Label{L: $1}, $3, $6)} }
         | branches PIPE LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = append($1, process.NewBranch(process.Label{L: $3}, $5, $8)) }
         ;

names : name { $$ = []process.Name{$1} }
 	  | name COMMA names { $$ = append($3, $1) }

name : SELF { $$ = process.Name{IsSelf: true} };
name : LABEL { $$ = process.Name{Ident: $1, IsSelf: false} };

%%

// Parse is the entry point to the parser.
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
		return unexpandedProcesses, nil	}
}
