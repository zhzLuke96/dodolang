%{
package main

import "fmt"

%}

%union{
	val interface{}
}

%token LexError
%token <val> Identifier StringConstant NumConstant
%token <val> FuncDefined FuncReturn

%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%right UMINUS
%left '='

%%

code:
		/* empty */
	|	code func_def
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
	;

assignStmt:	Identifier '=' expr		{ fmt.Printf("'%v' swap store ", $1) }
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
