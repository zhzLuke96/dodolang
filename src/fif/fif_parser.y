%{
package main

import "fmt"
import "strings"

var fif_code_buf = []string{}

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
%token <val> T_IF T_ELSE T_FOR T_WHILE T_FIF

%left '='
%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
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

expr:   '(' expr ')'                { /* empty */ }
	|   expr '+' expr               { fmt.Print("add ") }
	|   expr '-' expr               { fmt.Print("sub ") }
	|   expr '*' expr               { fmt.Print("mul ") }
	|   expr '/' expr               { fmt.Print("div ") }
	|   expr '&' expr               { fmt.Print("and ") }
	|   expr '|' expr               { fmt.Print("or ") }
	|   expr '%' expr               { fmt.Print("mod ") }
	|   Identifier					{ fmt.Printf("'%v' load ", $1) }
	|	NumConstant					{ fmt.Printf("%v ", $1) }
	|   '-' NumConstant %prec UMINUS{ fmt.Printf("-%v ",$2) }
	|	StringConstant				{ fmt.Printf("'%v' ", $1) }
	| 	callExpr          			{ /* empty */ }
	|	func_def					{ /* empty */ }
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
