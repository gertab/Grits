%{
package parser

import (
	"io"
	"phi/process"
)

var processes []process.Process
// funcDefinitions type

%}

%union {
	strval string
	proc   process.Process
}

%token LABEL LEFT_ARROW RIGHT_ARROW EQUALS DOT SEQUENCE COLON COMMA LPAREN RPAREN LSBRACK RSBRACK LANGLE RANGLE PIPE SEND RECEIVE CASE CLOSE WAIT CAST SHIFT ACCEPT ACQUIRE DETACH RELEASE DROP SPLIT PUSH NEW SNEW LET IN END SPRC PRC FORWARD
%type <proc> proc 
%type <strval> LABEL

%%

top : proc { 
	__yyfmt__.Println($1)
	// processes = $1 
	}
    ;

proc : SEND LABEL LANGLE LABEL COMMA LABEL RANGLE  { 
	__yyfmt__.Println("send")
	__yyfmt__.Println($2)
	__yyfmt__.Println("<")
	__yyfmt__.Println($4)
	__yyfmt__.Println(",")
	__yyfmt__.Println($6)
	__yyfmt__.Println(">")
	
	}
     | proc  LANGLE { }
     ;

%%

// Parse is the entry point to the parser.
func Parse(r io.Reader) ([]process.Process, error) {
	l := newLexer(r)
	phiParse(l)
	select {
	case err := <-l.Errors:
		return nil, err
	default:
		return processes, nil
	}
}
