

%token NUMBER TOKHEAT STATE TOKTARGET TOKTEMPERATURE


%union {   
value string
}
%token <value> NUMBER
%token <value> STATE

%%

list	: /* empty */
	| list '\n'
	;
