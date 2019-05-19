%{
package fif

import "fmt"

var IF_Label = NewLabelStack()
var WHILE_Label = NewLabelStack()

%}

%union{
	val interface{}
	code string
}

%type <val> any_T
%type <code> stmts stmt assignStmt labelStmt gotoStmt retStmt yeildStmt varStmt vars obj_propStmt callStmt inline_if_stmt inline_if_condition inline_if_then inline_if_else inline_if_end if_stmt if_BEG if_THEN then_END if_END else_BEG fact_Expr Stmts_Block while_stmt while_H while_BEG while_END breakStmt expr obj_expr get_obj_expr const_expr callExpr call_args call_arg named_func_def func_def func_args func_arg_list func_body named_gen_def gen_def fif_code fif_block S

%token LexError
%token <val> Identifier StringConstant NumConstant
%token <val> FuncDefined FuncReturn GenDefined CoroDefined
%token <val> T_IF T_ELSE T_THEN T_TRUE T_FALSE T_GOTO 
%token <val> T_FOR T_WHILE T_FIF T_BREAK
%token <val> T_EQ T_AND T_OR T_GE T_LE
%token <val> T_VAR T_NULL T_YIELD

%left '='
%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%left '<' '>'
%left T_EQ T_AND T_OR T_GE T_LE
%right UMINUS

%%

code:
		/* empty */
	|	code S						{ fmt.Fprintf(&ParserBuf,$2) }
	;

S	:	/* empty */					{ $$ = "" }
	|	stmts						{ $$ = $1 }
	|	expr						{ $$ = $1 }
	;

stmts:  /* empty */					{ $$ = "" }
	|	stmt stmts					{ $$ = fmt.Sprintf("%v %v", $1, $2) }
	;

stmt:	assignStmt					{ $$ = $1 }
	|	callStmt					{ $$ = $1 }
	|	retStmt						{ $$ = $1 }
	|	yeildStmt					{ $$ = $1 }
	|	if_stmt						{ $$ = $1 }
	|	while_stmt					{ $$ = $1 }
	|	named_func_def				{ $$ = $1 }
	|	named_gen_def				{ $$ = $1 }
	|	fif_code					{ $$ = $1 }
	|	labelStmt					{ $$ = $1 }
	|	gotoStmt					{ $$ = $1 }
	|	breakStmt					{ $$ = $1 }
	|	varStmt						{ $$ = $1 }
	|	obj_propStmt				{ $$ = $1 }
	;

assignStmt:
		Identifier '=' expr			{ $$ = fmt.Sprintf("'%v' %v store", $1, $3) }
	|	Identifier '=' T_YIELD expr { $$ = fmt.Sprintf("%v ret '%v' arg", $4, $1) }
	;

labelStmt:
		Identifier ':'				{ $$ = fmt.Sprintf("%v:", $1) }
	;

gotoStmt:
		T_GOTO Identifier			{ $$ = fmt.Sprintf("&%v jmp", $2) }	
	;

retStmt:
		FuncReturn					{ $$ = "ret" }
	|	FuncReturn expr				{ $$ = fmt.Sprintf("%v ret", $2) }
	;

yeildStmt:
		T_YIELD expr				{ $$ = fmt.Sprintf("%v ret", $2) }
	;

varStmt:
		T_VAR vars					{ $$ = $2 }
	;

vars:	/* empty */					{ $$ = "" }
	|	Identifier					{ $$ = fmt.Sprintf("'%v' nop storei", $1) }
	|	Identifier ',' vars			{ $$ = fmt.Sprintf("'%v' nop storei %v", $1, $3) }	
	;

obj_propStmt:
		Identifier '.' Identifier '=' expr
									{
										$$ = fmt.Sprintf("'set' '%v' %v '%v' load call", $3, $5, $1)
									}
	|	obj_expr '.' Identifier '=' expr
									{
										$$ = fmt.Sprintf("'set' '%v' %v '%v' load call", $3, $5, $1)
									}
	;

callStmt:
		callExpr					{ $$ = $1 }
	;

inline_if_stmt:
		inline_if_condition inline_if_then expr inline_if_else expr inline_if_end
									{
										$$ = fmt.Sprintf("%v %v %v %v %v %v", $1, $2, $3, $4, $5, $6)
									}
	;

inline_if_condition:
		expr						{
										IF_Label.BEG()
										$$ = $1
									}
	;

inline_if_then:
		T_THEN						{
										$$ = fmt.Sprintf("&_THEN_END%v fjmp", IF_Label.Topv())
									}
	;

inline_if_else:
		T_ELSE						{
										$$ = fmt.Sprintf("&_IF_END%v jmp", IF_Label.Topv()) + fmt.Sprintf("_THEN_END%v:", IF_Label.Topv())
									}
	;

inline_if_end:
		/* empty */					{
										$$ = fmt.Sprintf("_IF_END%v:", IF_Label.Topv())
										IF_Label.END()
									}
	;

if_stmt:
		if_BEG fact_Expr if_THEN Stmts_Block then_END if_END
									{
										$$ = fmt.Sprintf("%v %v %v %v %v", $2,$3,$4,$5,$6)
									}
	|	if_BEG fact_Expr if_THEN Stmts_Block then_END else_BEG Stmts_Block if_END
									{
										$$ = fmt.Sprintf("%v %v %v %v %v %v %v", $2,$3,$4,$5,$6,$7,$8)
									}
	;

if_BEG:
		T_IF						{ IF_Label.BEG() }
	;

if_THEN:
		/* empty */					{
										$$ = fmt.Sprintf("&_THEN_END%v fjmp", IF_Label.Topv())
									}
	;

then_END:
		/* empty */					{
										$$ = fmt.Sprintf("&_IF_END%v jmp ", IF_Label.Topv()) + fmt.Sprintf("_THEN_END%v:", IF_Label.Topv())
									}
	;

if_END:
		/* empty */					{
										$$ = fmt.Sprintf("_IF_END%v:", IF_Label.Topv())
										IF_Label.END()
									}
	;

else_BEG:
		T_ELSE						{ 
										$$ = fmt.Sprintf("_ELSE_BEG%v:", IF_Label.Topv())
									}
	;

fact_Expr:
		'(' expr ')'				{ $$ = $2 }
	|	expr						{ $$ = $1 }
	;

Stmts_Block:
		'{' stmts '}'				{ $$ = $2 }
	;

while_stmt:
		while_H fact_Expr while_BEG Stmts_Block while_END
									{
										$$ = fmt.Sprintf("%v %v %v %v %v", $1,$2,$3,$4,$5)
									}
	;

while_H:
		T_WHILE						{
										WHILE_Label.BEG()
										$$ = fmt.Sprintf("_WHILE_BEG%v:", WHILE_Label.Topv())
									}
	;

while_BEG:
		/* empty */					{
										$$ = fmt.Sprintf("&_WHILE_END%v fjmp", WHILE_Label.Topv())
									}
	;

while_END:
		/* empty */					{
										$$ = fmt.Sprintf("&_WHILE_BEG%v jmp ", WHILE_Label.Topv()) + fmt.Sprintf("_WHILE_END%v:", WHILE_Label.Topv())
										WHILE_Label.END()
									}
	;

breakStmt:
		T_BREAK 					{
										$$ = fmt.Sprintf("&_WHILE_END%v jmp", WHILE_Label.Topv())
									}
	;

expr:   expr '+' expr               { $$ = fmt.Sprintf("%v %v add", $1, $3) }
	|   expr '-' expr               { $$ = fmt.Sprintf("%v %v sub", $1, $3) }
	|   expr '*' expr               { $$ = fmt.Sprintf("%v %v mul", $1, $3) }
	|   expr '/' expr               { $$ = fmt.Sprintf("%v %v div", $1, $3) }
	|   expr '&' expr               { $$ = fmt.Sprintf("%v %v and", $1, $3) }
	|   expr '|' expr               { $$ = fmt.Sprintf("%v %v or", $1, $3) }
	|   expr '%' expr               { $$ = fmt.Sprintf("%v %v mod", $1, $3) }
	|   expr '>' expr               { $$ = fmt.Sprintf("%v %v gt", $1, $3) }
	|   expr '<' expr               { $$ = fmt.Sprintf("%v %v ls", $1, $3) }
	|   expr T_EQ expr           	{ $$ = fmt.Sprintf("%v %v equl", $1, $3) }
	|  	'!' expr           			{ $$ = fmt.Sprintf("%v not", $2) }
	|   expr T_AND expr           	{ $$ = fmt.Sprintf("%v %v and_b", $1, $3) }
	|   expr T_OR expr           	{ $$ = fmt.Sprintf("%v %v or_b", $1, $3) }
	|	T_TRUE						{ $$ = "1" }
	|	T_FALSE						{ $$ = "0" }
	|	T_NULL						{ $$ = "nop" }
	|   const_expr					{ $$ = $1 }
	|	obj_expr					{ $$ = $1 }
	| 	callExpr          			{ $$ = $1 }
	|	func_def					{ $$ = $1 }
	|	gen_def						{ $$ = $1 }
	|	inline_if_stmt				{ $$ = $1 }
	|	'(' expr ')'                { $$ = $2 }
	;

obj_expr:
		Identifier					{ $$ = fmt.Sprintf("'%v' load",$1) }
	|	Identifier '.' get_obj_expr	{ $$ = fmt.Sprintf("'get' '%v' '%v' load call", $3, $1) }
	|	obj_expr '.' get_obj_expr	{ $$ = fmt.Sprintf("'get' '%v' '%v' call", $3, $1) }
	;

get_obj_expr:
		Identifier					{ $$ = $1.(string) }
	;

const_expr:
		NumConstant					{ $$ = fmt.Sprintf("%v",$1) }
	|   '-' NumConstant %prec UMINUS{ $$ = fmt.Sprintf("-%v",$2) }
	|   '-' Identifier %prec UMINUS { $$ = fmt.Sprintf("'%v' laod neg",$2) }
	|	StringConstant				{ $$ = "'" + $1.(string) + "'" }
	;

callExpr:
		expr '(' call_args ')'		{ $$ = fmt.Sprintf("%v %v call", $3, $1) }
	|	callExpr '(' call_args ')'	{ $$ = fmt.Sprintf("%v %v call", $3, $1) }
	;
		
call_args:
		/* empty */					{ $$ = "" }
	|	call_arg					{ $$ = $1 }
	;

call_arg:
		expr						{ $$ = $1 }
	|	expr ',' call_arg			{ $$ = fmt.Sprintf("%v %v", $1, $3) }
	;

named_func_def:
		FuncDefined Identifier '(' func_args ')' '{' func_body '}'
									{
										sr := "stack_reverse"
										if argCount($4) <= 1 {
											sr = ""
										}
										$$ = fmt.Sprintf("'%v' func %v %v %v endfunc storei", $2, sr,$4, $7)
									}
	;

func_def:
		FuncDefined '(' func_args ')' '{' func_body '}'
									{
										sr := "stack_reverse"
										if argCount($3) <= 1 {
											sr = ""
										}
										$$ = fmt.Sprintf("func %v %v %v endfunc", sr, $3 ,$6)
									}
	;

func_args:
		/* empty */					{ $$ = "" }
	|	func_arg_list				{ $$ = $1 }
	;

func_arg_list:
		Identifier					{ $$ = fmt.Sprintf("'%v' arg", $1) }
	|	func_arg_list ',' Identifier{ $$ = fmt.Sprintf("%v '%v' arg", $1, $3) }
	;

func_body:
		S							{ $$ = $1 }
	;

named_gen_def:
		GenDefined Identifier '(' func_args ')' '{' func_body '}'
									{
										$$ = fmt.Sprintf(
											"'%v' func '__gen__' func %v %v endfunc storei func '__gen__' load callx endfunc ret endfunc storei",
											$2, $4, $7)
									}
	;

gen_def:
		GenDefined '(' func_args ')' '{' func_body '}'
									{
										/* 	
											func
												'__gen__' func 
												//...
												endfunc store
												func
													'__gen__' load callx
												endfunc ret
											endfunc
										*/
										$$ = fmt.Sprintf(
											"func '__gen__' func %v %v endfunc storei func '__gen__' load callx endfunc ret endfunc storei",
											$1, $3, $6)
									}
	;

fif_code:
		T_FIF '{' fif_block '}'		{
										$$ = $3
									}
	;

fif_block:
		/* empty */					{ $$ = "" }
	|	any_T fif_block				{ $$ = fmt.Sprintf("%v %v", $1, $2) }
	;

any_T:
		Identifier					{ $$ = $1.(string) }
	|	StringConstant				{ $$ = "'" + $1.(string) + "'" }
	|	NumConstant					{ $$ = fmtFloat64($1.(float64)) }
	;
