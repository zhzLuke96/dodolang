%{
package fif

import "fmt"
import "strings"

var fif_code_buf = []string{}
var IF_Label = NewLabelStack()
var WHILE_Label = NewLabelStack()

%}

%union{
	val interface{}
}

%type <val> any_T

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
	|	code S
	;

S	:	/* empty */
	|	stmts
	|	expr
	;

stmts:  /* empty */
	|	stmt stmts
	;

stmt:	assignStmt
	|	callStmt
	|	retStmt
	|	yeildStmt
	|	if_stmt
	|	while_stmt
	|	named_func_def
	|	named_gen_def
	|	fif_code
	|	labelStmt
	|	gotoStmt
	|	breakStmt
	|	varStmt
	;

assignStmt:
		Identifier '=' expr			{ fmt.Fprintf(&ParserBuf,"'%v' swap store ", $1) }
	|	Identifier '=' T_YIELD expr { fmt.Fprintf(&ParserBuf,"ret '%v' arg ", $1) }
	;

labelStmt:
		Identifier ':'				{ fmt.Fprintf(&ParserBuf,"%v: ", $1) }
	;

gotoStmt:
		T_GOTO Identifier			{ fmt.Fprintf(&ParserBuf,"&%v jmp ", $2) }	
	;

retStmt:
		FuncReturn					{ fmt.Fprintf(&ParserBuf,"ret ") }
	|	FuncReturn expr				{ fmt.Fprintf(&ParserBuf,"ret ") }
	;

yeildStmt:
		T_YIELD expr				{ fmt.Fprintf(&ParserBuf,"ret ") }
	;

varStmt:
		T_VAR vars
	;

vars:	/* empty */
	|	Identifier					{ fmt.Fprintf(&ParserBuf,"'%v' nop storei ", $1) }
	|	Identifier ',' vars			{ fmt.Fprintf(&ParserBuf,"'%v' nop storei ", $1) }	
	;

callStmt:
		callExpr
	;

inline_if_stmt:
		inline_if_condition inline_if_then expr inline_if_else expr inline_if_end
	;

inline_if_condition:
		expr						{ IF_Label.BEG() }
	;

inline_if_then:
		T_THEN						{
										fmt.Fprintf(&ParserBuf,"&_THEN_END%v fjmp ", IF_Label.Topv())
									}
	;

inline_if_else:
		T_ELSE						{
										fmt.Fprintf(&ParserBuf,"&_IF_END%v jmp ", IF_Label.Topv())
										fmt.Fprintf(&ParserBuf,"_THEN_END%v: ", IF_Label.Topv())
									}
	;

inline_if_end:
		/* empty */					{
										fmt.Fprintf(&ParserBuf,"_IF_END%v: ", IF_Label.Topv())
										IF_Label.END()
									}
	;

if_stmt:
		if_BEG fact_Expr if_THEN Stmts_Block then_END if_END
	|	if_BEG fact_Expr if_THEN Stmts_Block then_END else_BEG Stmts_Block if_END
	;

if_BEG:
		T_IF						{ IF_Label.BEG() }
	;

if_THEN:
		/* empty */					{
										fmt.Fprintf(&ParserBuf,"&_THEN_END%v fjmp ", IF_Label.Topv())
									}
	;

then_END:
		/* empty */					{
										fmt.Fprintf(&ParserBuf,"&_IF_END%v jmp ", IF_Label.Topv())
										fmt.Fprintf(&ParserBuf,"_THEN_END%v: ", IF_Label.Topv())
									}
	;

if_END:
		/* empty */					{
										fmt.Fprintf(&ParserBuf,"_IF_END%v: ", IF_Label.Topv())
										IF_Label.END()
									}
	;

else_BEG:
		T_ELSE						{ 
										fmt.Fprintf(&ParserBuf,"_ELSE_BEG%v: ", IF_Label.Topv())
									}
	;

fact_Expr:
		'(' expr ')'
	|	expr
	;

Stmts_Block:
		'{' stmts '}'
	;

while_stmt:
		while_H fact_Expr while_BEG Stmts_Block while_END
	;

while_H:
		T_WHILE						{
										WHILE_Label.BEG()
										fmt.Fprintf(&ParserBuf,"_WHILE_BEG%v: ", WHILE_Label.Topv())
									}
	;

while_BEG:
		/* empty */					{
										fmt.Fprintf(&ParserBuf,"&_WHILE_END%v fjmp ", WHILE_Label.Topv())
									}
	;

while_END:
		/* empty */					{
										fmt.Fprintf(&ParserBuf,"&_WHILE_BEG%v jmp ", WHILE_Label.Topv())
										fmt.Fprintf(&ParserBuf,"_WHILE_END%v: ", WHILE_Label.Topv())
										WHILE_Label.END()
									}
	;

breakStmt:
		T_BREAK 					{
										fmt.Fprintf(&ParserBuf, "&_WHILE_END%v jmp ", WHILE_Label.Topv())
									}
	;

expr:   expr '+' expr               { fmt.Fprintf(&ParserBuf,"add ") }
	|   expr '-' expr               { fmt.Fprintf(&ParserBuf,"sub ") }
	|   expr '*' expr               { fmt.Fprintf(&ParserBuf,"mul ") }
	|   expr '/' expr               { fmt.Fprintf(&ParserBuf,"div ") }
	|   expr '&' expr               { fmt.Fprintf(&ParserBuf,"and ") }
	|   expr '|' expr               { fmt.Fprintf(&ParserBuf,"or ") }
	|   expr '%' expr               { fmt.Fprintf(&ParserBuf,"mod ") }
	|   expr '>' expr               { fmt.Fprintf(&ParserBuf,"gt ") }
	|   expr '<' expr               { fmt.Fprintf(&ParserBuf,"ls ") }
	|   expr T_EQ expr           	{ fmt.Fprintf(&ParserBuf,"equl ") }
	|  	'!' expr           			{ fmt.Fprintf(&ParserBuf,"not ") }
	|   expr T_AND expr           	{ fmt.Fprintf(&ParserBuf,"and_b ") }
	|   expr T_OR expr           	{ fmt.Fprintf(&ParserBuf,"or_b ") }
	|	T_TRUE						{ fmt.Fprintf(&ParserBuf,"1 ") }
	|	T_FALSE						{ fmt.Fprintf(&ParserBuf,"0 ") }
	|	T_NULL						{ fmt.Fprintf(&ParserBuf,"nop ") }
	|   Identifier					{ fmt.Fprintf(&ParserBuf,"'%v' load ", $1) }
	|	NumConstant					{ fmt.Fprintf(&ParserBuf,"%v ", $1) }
	|   '-' NumConstant %prec UMINUS{ fmt.Fprintf(&ParserBuf,"-%v ",$2) }
	|   '-' Identifier %prec UMINUS { fmt.Fprintf(&ParserBuf,"'%v' load neg ",$2) }
	|	StringConstant				{ fmt.Fprintf(&ParserBuf,"'%v' ", $1) }
	| 	callExpr          			{ /* empty */ }
	|	func_def					{ /* empty */ }
	|	gen_def						{ /* empty */ }
	|	inline_if_stmt				{ /* empty */ }
	|	'(' expr ')'                { /* empty */ }
	;

callExpr:
		Identifier '(' call_args ')'{ fmt.Fprintf(&ParserBuf,"'%v' call ", $1) }
	;

call_args:
		/* empty */					{ /* empty */ }
	|	expr call_arg				{ /* empty */ }
	;

call_arg:
		/* empty */					{ /* empty */ }
	|	call_arg ',' expr			{ /* empty */ }
	;

named_func_def:
		named_func_h '(' func_args ')' '{' func_body '}'
									{ fmt.Fprintf(&ParserBuf,"endfunc storei ") }
	;

named_func_h:
		FuncDefined Identifier		{ fmt.Fprintf(&ParserBuf,"'%v' func ",$2) }
	;

func_def:
		func_def_h '(' func_args ')' '{' func_body '}'
									{ fmt.Fprintf(&ParserBuf,"endfunc ") }
	;

func_def_h:
		FuncDefined					{ fmt.Fprintf(&ParserBuf,"func ") }
	;

func_args:
		/* empty */					{ /* empty */ }
	|	Identifier func_arg			{ fmt.Fprintf(&ParserBuf,"'%v' arg ", $1) }
	;

func_arg:
		/* empty */					{ /* empty */ }
	|	func_arg ',' Identifier		{ fmt.Fprintf(&ParserBuf,"'%v' arg ", $3) }
	;

func_body:
		S
	;

named_gen_def:
		named_gen_h '(' func_args ')' '{' func_body '}'
									{ fmt.Fprintf(&ParserBuf,"endfunc storei func '__gen__' load callx endfunc ret endfunc storei ") }
	;

named_gen_h:
		GenDefined Identifier		{ fmt.Fprintf(&ParserBuf,"'%v' func '__gen__' func ",$2) }
	;

gen_def:
		gen_def_h '(' func_args ')' '{' func_body '}'
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
										fmt.Fprintf(&ParserBuf,"endfunc storei func '__gen__' load callx endfunc ret endfunc")
									}
	;

gen_def_h:
		GenDefined					{ fmt.Fprintf(&ParserBuf,"func '__gen__' func ") }
	;

fif_code:
		T_FIF '{' fif_block '}'		{
										fmt.Fprintf(&ParserBuf,"%v ",
											strings.Join(
												reverse(fif_code_buf), " "))
										fif_code_buf = []string{}
									}
	;

fif_block:
		/* empty */					
	|	any_T fif_block				{ fif_code_buf = append(fif_code_buf,$1.(string)) }
	;

any_T:
		Identifier					{ $$ = $1.(string) }
	|	StringConstant				{ $$ = "'" + $1.(string) + "'" }
	|	NumConstant					{ $$ = fmtFloat64($1.(float64)) }
	;
