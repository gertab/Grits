%{
// Run this after each change:
// goyacc -p phi -o parser/parser.y.go parser/parser.y
package parser

import (
	"io"
	"phi/process"
	"phi/types"
	"phi/position"
)

%}

%union {
	strval 			      string
	currPosition 	      position.Position
	common_type		      unexpandedProcessOrFunction
	statements 		      []unexpandedProcessOrFunction
	name 			      process.Name
	names 			      []process.Name
	form 			      process.Form
	branches 		      []*process.BranchForm
	sessionType 	      types.SessionType
	sessionTypeInitial 	  types.SessionTypeInitial
	sessionTypeAltInitial []types.OptionInitial
	polarity 		      types.Polarity
}

%token LABEL LEFT_ARROW RIGHT_ARROW UP_ARROW DOWN_ARROW  EQUALS DOT SEQUENCE COLON COMMA LPAREN RPAREN LSBRACK RSBRACK LANGLE RANGLE PIPE SEND RECEIVE CASE CLOSE WAIT CAST SHIFT ACCEPT ACQUIRE DETACH RELEASE DROP SPLIT PUSH NEW SNEW TYPE LET IN END SPRC PRC FORWARD SELF PRINT PRINTL PLUS MINUS TIMES AMPERSAND UNIT LCBRACK RCBRACK LOLLI PERCENTAGE ASSUMING EXEC
%type <strval> LABEL
%type <statements> statements 
%type <common_type> process_def
%type <common_type> function_def
%type <common_type> type_def
%type <common_type> assuming_def
%type <common_type> exec_def
%type <form> expression 
%type <name> name
%type <name> name_with_type_ann
%type <names> names
%type <names> optional_names
%type <names> optional_names_with_type_ann
%type <names> comma_optional_names_with_type_ann
%type <names> names_with_type_ann
%type <strval> modality
%type <branches> branches
%type <sessionType> session_type
%type <sessionTypeAltInitial> session_type_options_init
%type <sessionTypeInitial> session_type_init
%type <polarity> polarity

%left SEQUENCE RANGLE
%right TIMES LOLLI UP_ARROW DOWN_ARROW

%%

root : program { };

program : 
		/* simulate a process */
		/* todo remove */
	   expression 
		{
			philex.(*lexer).processesOrFunctionsRes = append(philex.(*lexer).processesOrFunctionsRes, unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body:$1, Providers: []process.Name{{Ident: "root" , IsSelf: false}}}, position: phiVAL.currPosition})
		}
	 | statements 
		{ 
			philex.(*lexer).processesOrFunctionsRes = $1
		};
/*	 | LET functions IN processes END { }; */

/* A program may consist a combination of processes, function definitions and types */
statements : process_def             { $$ = []unexpandedProcessOrFunction{$1} }
		   | process_def statements  { $$ = append([]unexpandedProcessOrFunction{$1}, $2...) }
		   | function_def            { $$ = []unexpandedProcessOrFunction{$1} }
		   | function_def statements { $$ = append([]unexpandedProcessOrFunction{$1}, $2...) }
		   | type_def 				 { $$ = []unexpandedProcessOrFunction{$1} }
		   | type_def statements 	 { $$ = append([]unexpandedProcessOrFunction{$1}, $2...) }
		   | assuming_def 			 { $$ = []unexpandedProcessOrFunction{$1} }
		   | assuming_def statements { $$ = append([]unexpandedProcessOrFunction{$1}, $2...) }
		   | exec_def 			 	 { $$ = []unexpandedProcessOrFunction{$1} }
		   | exec_def statements 	 { $$ = append([]unexpandedProcessOrFunction{$1}, $2...) };

/* A process is defined using the prc keyword */
process_def : 
			/* without type - todo remove option to force types */
		    PRC LSBRACK names RSBRACK EQUALS expression 
				{ $$ = unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body:$6, Providers: $3}, position: phiVAL.currPosition} }
		  | PRC LSBRACK names RSBRACK COLON session_type EQUALS expression 
				{ $$ = unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body:$8, Type: $6, Providers: $3}, position: phiVAL.currPosition} };
		/*| SPRC LSBRACK names RSBRACK COLON expression
				{ $$ = unexpandedProcessOrFunction{kind: PROCESS_DEF, proc: incompleteProcess{Body:$6, Providers: $3}, position: phiVAL.currPosition} };*/

/* Expressions form the core part of a program  */
expression : /* Send */ SEND name LANGLE name COMMA name RANGLE  
					{ $$ = process.NewSend($2, $4, $6) }
		   /* Send Macro */ /* SEND name LANGLE name COMMA name RANGLE SEQUENCE expression
			{ $$ = NewSendMacro($2, $4, $6, $9) }*/
		   | /* Receive */ LANGLE name COMMA name RANGLE LEFT_ARROW RECEIVE name SEQUENCE expression 
		   			{ $$ = process.NewReceive($2, $4, $8, $10) }
		   | /* Select */ name DOT LABEL LANGLE name RANGLE 
		   			{ $$ = process.NewSelect($1, process.Label{L: $3}, $5) }
		   | /* Case */ CASE name LPAREN branches RPAREN 
		   			{ $$ = process.NewCase($2, $4) }
		   | /* New */ name LEFT_ARROW NEW expression SEQUENCE expression 
					{ $$ = process.NewNew($1, $4, $6) } 
		   | /* New */ LABEL COLON session_type LEFT_ARROW NEW expression SEQUENCE expression 
					{ $$ = process.NewNew(process.Name{Ident: $1, Type: $3, IsSelf: false}, $6, $8) } 		   
		   | /* Call */ LABEL LPAREN optional_names RPAREN
		   			{ $$ = process.NewCall($1, $3) }
		   | /* Close */ CLOSE name
		   			{ $$ = process.NewClose($2) }
		   | /* Forward */ FORWARD name name
				{ $$ = process.NewForward($2, $3) } 
		   | /* Split */ LANGLE name COMMA name RANGLE LEFT_ARROW SPLIT name SEQUENCE expression
		   			{ $$ = process.NewSplit($2, $4, $8, $10) }
		   | /* Wait */ WAIT name SEQUENCE expression
		   			{ $$ = process.NewWait($2, $4) }
		   | /* Cast */ CAST name LANGLE name RANGLE  
					{ $$ = process.NewCast($2, $4) }
		   | /* Shift */ name LEFT_ARROW SHIFT name SEQUENCE expression 
		   			{ $$ = process.NewShift($1, $4, $6) }
		   | /* Drop */ DROP name SEQUENCE expression
					{ $$ = process.NewDrop($2, $4) }
		   | /* Brackets */ LPAREN expression RPAREN
					{ $$ = $2 }
		   | /* Print - for output */ PRINT name SEQUENCE expression
		   			{ $$ = process.NewPrint($2, $4) }
		   | /* PrintL - for output */ PRINTL LABEL SEQUENCE expression
		   			{ $$ = process.NewPrintL(process.Label{L: $2}, $4) };
/* remaining expressions - used for shared processes
	SNew, Acquire, Accept, Push, Detach, Release*/
 
branches :   /* empty */         										 { $$ = nil }
         |               LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = []*process.BranchForm{process.NewBranch(process.Label{L: $1}, $3, $6)} }
         | branches PIPE LABEL LANGLE name RANGLE RIGHT_ARROW expression { $$ = append($1, process.NewBranch(process.Label{L: $3}, $5, $8)) };

names : name { $$ = []process.Name{$1} }
 	  | name COMMA names { $$ = append([]process.Name{$1}, $3...) };

optional_names : /* empty */ { $$ = nil }
		| name { $$ = []process.Name{$1} }
		| name COMMA names { $$ = append([]process.Name{$1}, $3...) };

optional_names_with_type_ann : 
			/* empty */ { $$ = nil }
		| name_with_type_ann { $$ = []process.Name{$1} }
		| name_with_type_ann COMMA names_with_type_ann { $$ = append([]process.Name{$1}, $3...) };

comma_optional_names_with_type_ann : 
			/* empty */ { $$ = nil }
		| COMMA optional_names_with_type_ann { $$ = $2 };


names_with_type_ann : 
	name_with_type_ann { $$ = []process.Name{$1} }
	| name_with_type_ann COMMA names_with_type_ann { $$ = append([]process.Name{$1}, $3...) };

name_with_type_ann : 
			/* without type - todo remove option to force types */
			LABEL
					{ $$ = process.Name{Ident: $1, IsSelf: false} }
			| LABEL COLON session_type 
			 		{ $$ = process.Name{Ident: $1, Type: $3, IsSelf: false} };

name : SELF { $$ = process.Name{IsSelf: true} }
	 | polarity SELF  
		{ pol := $1
		  $$ = process.Name{IsSelf: true, ExplicitPolarity: &pol} }
	 | LABEL { $$ = process.Name{Ident: $1, IsSelf: false} }
	 | polarity LABEL
		{ pol := $1
		  $$ = process.Name{Ident: $2, IsSelf: false, ExplicitPolarity: &pol} };

assuming_def : ASSUMING names_with_type_ann
			{ $$ = unexpandedProcessOrFunction{kind: ASSUMING_DEF, assumedFreeNameTypes: $2, position: phiVAL.currPosition} };

function_def : 
			/* without type - todo remove option to force types */
			 LET LABEL LPAREN optional_names_with_type_ann RPAREN EQUALS expression
					{ $$ = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: $2, Parameters: $4, Body: $7, UsesExplicitProvider: false}, position: phiVAL.currPosition} }
			| /* with type annotation */ LET LABEL LPAREN optional_names_with_type_ann RPAREN COLON session_type EQUALS expression
					{ $$ = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: process.FunctionDefinition{FunctionName: $2, Parameters: $4, Body: $9, Type: $7, UsesExplicitProvider: false}, position: phiVAL.currPosition} }
			| /* explicit provider name : without type - todo remove option to force types */
			 LET LABEL LSBRACK LABEL comma_optional_names_with_type_ann RSBRACK EQUALS expression
					{ $$ = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: 
							process.FunctionDefinition{
								FunctionName: $2, 
								Parameters: $5, 
								Body: $8, 
								UsesExplicitProvider: true, 
								ExplicitProvider: process.Name{Ident: $4, IsSelf: true}, 
								// Type: $6,
								}, position: phiVAL.currPosition} }
			| /* explicit provider name :  with type annotation */
			 LET LABEL LSBRACK LABEL COLON session_type comma_optional_names_with_type_ann RSBRACK EQUALS expression 
					{ $$ = unexpandedProcessOrFunction{kind: FUNCTION_DEF, function: 
							process.FunctionDefinition{
								FunctionName: $2, 
								Parameters: $7, 
								Body: $10, 
								UsesExplicitProvider: true, 
								ExplicitProvider: process.Name{Ident: $4, IsSelf: true}, 
								Type: $6}, position: phiVAL.currPosition} };

type_def : TYPE LABEL EQUALS session_type
			{ $$ = unexpandedProcessOrFunction{
						kind: TYPE_DEF, 
						session_type: types.SessionTypeDefinition{Name: $2, SessionType: $4},
						position: phiVAL.currPosition} };

/* Returns a SessionType struct */
session_type : /* no explicit mode */ session_type_init
					{ $$ = types.ConvertSessionTypeInitialToSessionType($1)}
			 | /* explicit mode */ modality session_type_init
					{ mode := types.StringToMode($1)
					  $$ = types.ConvertSessionTypeInitialToSessionType(types.NewExplicitModeTypeInitial(mode, $2))};

/* Returns a SessionTypeInitial struct */
session_type_init : 
			/* label */ LABEL
				{ $$ = types.NewLabelTypeInitial($1) }
		   | /* unit */ UNIT
		   		{ $$ = types.NewUnitTypeInitial() }
		   | /* select +{ } */ PLUS LCBRACK session_type_options_init RCBRACK  
		   		{ $$ = types.NewSelectLabelTypeInitial($3) }
		   | /* branch &{ } */ AMPERSAND LCBRACK session_type_options_init RCBRACK  
		   		{ $$ = types.NewBranchCaseTypeInitial($3) }
		   | /* send A * B */ session_type_init TIMES session_type_init
		   		{ $$ = types.NewSendTypeInitial($1, $3) }
		   | /* receive A -o B */ session_type_init LOLLI session_type_init
		   		{ $$ = types.NewReceiveTypeInitial($1, $3) }
		   | /* brackets (A) */ LPAREN session_type_init RPAREN
		   		{ $$ = $2 }
		   | /* upshift mode /\ model type */ modality UP_ARROW modality session_type_init
		   		{ modeFrom := types.StringToMode($1)
				  modeTo := types.StringToMode($3)
				  $$ = types.NewUpTypeInitial(modeFrom, modeTo, $4) }
		   | /* downshift mode /\ model type */ modality DOWN_ARROW modality session_type_init
		   		{ modeFrom := types.StringToMode($1)
				  modeTo := types.StringToMode($3)
				  $$ = types.NewDownTypeInitial(modeFrom, modeTo, $4) };

session_type_options_init : 
            LABEL COLON session_type_init 
				{ $$ = []types.OptionInitial{*types.NewOptionInitial($1, $3)}} 
	 	  | LABEL COLON session_type_init COMMA session_type_options_init 
		  { $$ = append([]types.OptionInitial{*types.NewOptionInitial($1, $3)}, $5...) };

modality : LABEL { $$ = $1 };

polarity : PLUS { $$ = types.POSITIVE }
	     | MINUS { $$ = types.NEGATIVE };

/* execute function definitions directly */
exec_def : EXEC LABEL LPAREN RPAREN
			{ $$ = unexpandedProcessOrFunction{
				kind: EXEC_DEF, 
				proc: incompleteProcess{Body: process.NewCall($2, []process.Name{})},
				position: phiVAL.currPosition}};

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