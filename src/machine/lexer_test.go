package machine

import (
	"testing"
)

func TestLexer_Base(t *testing.T) {
	token := "1"
	needVar := "Number"
	name, _ := GetTokenTypeName(token)
	if name == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, name)
	}

	token = "-1"
	needVar = "Number"
	name, _ = GetTokenTypeName(token)
	if name == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, name)
	}

	token = "'name'"
	needVar = "String"
	name, _ = GetTokenTypeName(token)
	if name == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, name)
	}

	token = "mul"
	needVar = "Operator"
	name, _ = GetTokenTypeName(token)
	if name == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, name)
	}

	token = "return"
	needVar = "Instruction"
	name, _ = GetTokenTypeName(token)
	if name == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, name)
	}
}

func TestLexer_InsWithArgs(t *testing.T) {
	tokenHead := "dup"
	tokenTail := "5"
	token := tokenHead + "_" + tokenTail
	needVar := "Instruction_Args"
	name, arg := GetTokenTypeName(token)
	if name == needVar && arg == tokenTail {
		t.Log("Pass Lexer[Instruction_Args]")
	} else {
		t.Errorf("Failed Lexer[Instruction_Args] need %v,%v but %v,%v", needVar, tokenTail, name, arg)
	}

	tokenHead = "load"
	tokenTail = "some_vars"
	token = tokenHead + "_" + tokenTail
	needVar = "Instruction_Args"
	name, arg = GetTokenTypeName(token)
	if name == needVar && arg == tokenTail {
		t.Log("Pass Lexer[Instruction_Args]")
	} else {
		t.Errorf("Failed Lexer[Instruction_Args] need %v,%v but %v,%v", needVar, tokenTail, name, arg)
	}
}

func TestLexer_Unknow(t *testing.T) {
	token := "unknow_token"
	needVar := "UNKNOW TOKEN"
	name, _ := GetTokenTypeName(token)
	if name == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, name)
	}
}
