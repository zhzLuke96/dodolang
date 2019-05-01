package machine

import (
	"testing"
)

const (
	TestSymbol = "Machine_Test_Symbol"
)

func TestMachine_Run(t *testing.T) {
	code := [...]string{"1", "num", "2", "num", "plus"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Top(0).(int)
	if ok && v == 3 {
		t.Log("Pass Machine.Run")
	} else {
		t.Errorf("Failed Machine.Run need %v but %v", 3, v)
	}
}

func TestMachine_Push(t *testing.T) {
	code := [...]string{"1"}
	m := NewMachine(code[:])
	m.Push(TestSymbol)

	v, ok := m.Top(0).(string)
	if ok && v == TestSymbol {
		t.Log("Pass Machine.Push")
	} else {
		t.Errorf("Failed Machine.Push need %v but %v", TestSymbol, v)
	}
}

func TestMachine_Pop(t *testing.T) {
	code := [...]string{"1"}
	m := NewMachine(code[:])
	m.Push(TestSymbol)

	v, ok := m.Pop().(string)
	if ok && v == TestSymbol {
		t.Log("Pass Machine.Pop")
	} else {
		t.Errorf("Failed Machine.Pop need %v but %v", TestSymbol, v)
	}
}

func TestMachine_Label(t *testing.T) {
	needVar := 11
	code := [...]string{"var1:", "null", "11", "int", "'var1'", "store", "'var1'", "load"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.Label")
	} else {
		t.Errorf("Failed Machine.Label need %v but %v", needVar, v)
	}
}

func TestMachine_Jump(t *testing.T) {
	needVar := 1
	code := [...]string{"1", "int", "'end'", "jump", "2", "int", "3", "int", "mul", "end:", "dup"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.Jump")
	} else {
		t.Errorf("Failed Machine.Jump need %v but %v", needVar, v)
	}
}

func TestMachine_Dup(t *testing.T) {
	needVar := 1
	code := [...]string{"1", "int", "dup"}
	m := NewMachine(code[:]).Run()

	m.Pop()
	v, ok := m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.Dup")
	} else {
		t.Errorf("Failed Machine.Dup need %v but %v", needVar, v)
	}
}

func TestMachine_Func_CaLL_Return(t *testing.T) {
	needVar := 1
	code := [...]string{"-1", "num", "&abs", "call", "exit", "abs:", "dup", "0", "num", "greater", "&true_s", "&false_s", "if", "true_s:return", "false_s:-1", "num", "mul", "return"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.Call")
	} else {
		t.Errorf("Failed Machine.Call need %v but %v", needVar, v)
	}

	code = [...]string{"1", "num", "&abs", "call", "exit", "abs:", "dup", "0", "num", "greater", "&true_s", "&false_s", "if", "true_s:return", "false_s:-1", "num", "mul", "return"}
	m = NewMachine(code[:]).Run()

	v, ok = m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.Call")
	} else {
		t.Errorf("Failed Machine.Call need %v but %v", needVar, v)
	}
}
