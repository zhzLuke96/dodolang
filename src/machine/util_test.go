package machine

import (
	"testing"
)

var testCode = `
// comment line
main:
	exit // inline comment line


	// mulit enter

subproc:
	mul "string"
	"Closure nesting:\"string\""
	"  long string  "
	return
`

func TestCodeReader(t *testing.T) {
	code := codeReader(testCode)

	needVar := "exit"
	fact := code[1]
	if fact == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, fact)
	}

	needVar = "\"string\""
	fact = code[4]
	if fact == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, fact)
	}

	needVar = "\"Closure nesting:\\\"string\\\"\""
	fact = code[5]
	if fact == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, fact)
	}

	needVar = "\"  long string  \""
	fact = code[6]
	if fact == needVar {
		t.Log("Pass Lexer.base")
	} else {
		t.Errorf("Failed Lexer.base need %v but %v", needVar, fact)
	}
}
