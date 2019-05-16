%{
package main

import "fmt"
import "strings"

var fif_code_buf = []string{}
var IF_Label = NewLabelStack()
var WHILE_Label = NewLabelStack()

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

%}

%union{
	val interface{}
}

%type <val> any_T

%token LexError
%token <val> Identifier StringConstant NumConstant
%token <val> FuncDefined FuncReturn
%token <val> T_IF T_ELSE T_THEN T_TRUE T_FALSE
%token <val> T_FOR T_WHILE T_FIF
%token <val> T_EQ T_AND T_OR T_GE T_LE

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
	|	if_stmt
	|	named_func_def
	|	fif_code
	;

assignStmt:	assignId '=' expr		{ fmt.Printf("store ") }
	;

assignId:
		Identifier					{ fmt.Printf("'%v' ", $1) }
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
										fmt.Printf("&_THEN_END%v fjmp ", IF_Label.Topv())
									}
	;

inline_if_else:
		T_ELSE						{
										fmt.Printf("&_IF_END%v jmp ", IF_Label.Topv())
										fmt.Printf("_THEN_END%v: ", IF_Label.Topv())
									}
	;

inline_if_end:
		/* empty */					{
										fmt.Printf("_IF_END%v: ", IF_Label.Topv())
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
										fmt.Printf("&_THEN_END%v fjmp ", IF_Label.Topv())
									}
	;

then_END:
		/* empty */					{
										fmt.Printf("&_IF_END%v jmp ", IF_Label.Topv())
										fmt.Printf("_THEN_END%v: ", IF_Label.Topv())
									}
	;

if_END:
		/* empty */					{
										fmt.Printf("_IF_END%v: ", IF_Label.Topv())
										IF_Label.END()
									}
	;

else_BEG:
		T_ELSE						{ 
										fmt.Printf("_ELSE_BEG%v: ", IF_Label.Topv())
									}
	;

fact_Expr:
		'(' expr ')'
	|	expr
	;

Stmts_Block:
		'{' stmts '}'
	;

expr:   '(' expr ')'                { /* empty */ }
	|   expr '+' expr               { fmt.Print("add ") }
	|   expr '-' expr               { fmt.Print("sub ") }
	|   expr '*' expr               { fmt.Print("mul ") }
	|   expr '/' expr               { fmt.Print("div ") }
	|   expr '&' expr               { fmt.Print("and ") }
	|   expr '|' expr               { fmt.Print("or ") }
	|   expr '%' expr               { fmt.Print("mod ") }
	|   expr '>' expr               { fmt.Print("gt ") }
	|   expr '<' expr               { fmt.Print("ls ") }
	|   expr T_EQ expr           	{ fmt.Print("equl ") }
	|  	'!' expr           			{ fmt.Print("not ") }
	|   expr T_AND expr           	{ fmt.Print("and_b ") }
	|   expr T_OR expr           	{ fmt.Print("or_b ") }
	|	T_TRUE						{ fmt.Printf("1 ") }
	|	T_FALSE						{ fmt.Printf("0 ") }
	|   Identifier					{ fmt.Printf("'%v' load ", $1) }
	|	NumConstant					{ fmt.Printf("%v ", $1) }
	|   '-' NumConstant %prec UMINUS{ fmt.Printf("-%v ",$2) }
	|   '-' Identifier %prec UMINUS { fmt.Printf("'%v' load neg ",$2) }
	|	StringConstant				{ fmt.Printf("'%v' ", $1) }
	| 	callExpr          			{ /* empty */ }
	|	func_def					{ /* empty */ }
	|	inline_if_stmt				{ /* empty */ }
	;

callExpr:
		Identifier '(' call_args ')'{ fmt.Printf("'%v' call ", $1) }
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
									{ fmt.Printf("endfunc store ") }
	;

named_func_h:
		FuncDefined Identifier		{ fmt.Printf("'%v' func ",$2) }
	;

func_def:
		func_def_h '(' func_args ')' '{' func_body '}'
									{ fmt.Printf("endfunc ") }
	;

func_def_h:
		FuncDefined					{ fmt.Printf("func ") }
	;

func_args:
		/* empty */					{ /* empty */ }
	|	Identifier func_arg			{ fmt.Printf("'%v' arg ", $1) }
	;

func_arg:
		/* empty */					{ /* empty */ }
	|	func_arg ',' Identifier		{ fmt.Printf("'%v' arg ", $3) }
	;

func_body:
		S
	;

retStmt:
		FuncReturn					{ fmt.Printf("ret ") }
	|	FuncReturn expr				{ fmt.Printf("ret ") }
	;

fif_code:
		T_FIF '{' fif_block '}'		{
										fmt.Printf("%v ",
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
	|	NumConstant					{ $$ = $1.(string) }
	;
