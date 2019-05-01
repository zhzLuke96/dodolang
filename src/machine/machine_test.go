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

	v, ok := m.Top(0).(float64)
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
	needVar := 2
	code := [...]string{"1", "num", "2", "num", "&end", "jump", "mul", "end:exit"}
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
	code := [...]string{"1", "int", "&end", "jump", "2", "int", "3", "int", "mul", "end:", "dup"}
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

func TestMachine_IfThen(t *testing.T) {
	needVar := 1
	code := [...]string{"'true'", "bool", "if", "1", "num", "exit", "then", "2", "num"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.IfThen")
	} else {
		t.Errorf("Failed Machine.IfThen need %v but %v", needVar, v)
	}
}

func TestMachine_Scoped(t *testing.T) {
	needVar := 2
	code := [...]string{"1", "num", "'var1'", "store", "2", "num", "'var1'", "load", "mul"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(float64)
	if ok && int(v) == needVar {
		t.Log("Pass Machine.Scoped")
	} else {
		t.Errorf("Failed Machine.Scoped need %v but %v", needVar, v)
	}
}

func TestMachine_Func_CaLL_Return(t *testing.T) {
	needVar := 1
	code := [...]string{"-1", "num", "&abs", "call", "exit", "abs:", "dup", "0", "num", "less", "if", "num", "mul", "return", "then", "return"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.Call")
	} else {
		t.Errorf("Failed Machine.Call need %v but %v", needVar, v)
	}

	code = [...]string{"1", "num", "&abs", "call", "exit", "abs:", "dup", "0", "num", "less", "if", "num", "mul", "return", "then", "return"}
	m = NewMachine(code[:]).Run()

	v, ok = m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.Call")
	} else {
		t.Errorf("Failed Machine.Call need %v but %v", needVar, v)
	}
}

func TestMachine_GarbageCollection(t *testing.T) {
	needVar := "undefined"
	code := [...]string{"&main", "jump", "subproc:", "'Hello'", "store_var1", "return", "main:", "'subproc'", "call", "load_var1", "exit"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(string)
	if ok && v == needVar {
		t.Log("Pass Machine.GarbageCollection")
	} else {
		t.Errorf("Failed Machine.GarbageCollection need %v but %v", needVar, v)
	}
}

func TestMachine_GlobalVar(t *testing.T) {
	needVar := 13
	code := [...]string{"13", "int", "store_gvar1", "&main", "jump", "main:", "load_gvar1", "exit"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.GlobalVar")
	} else {
		t.Errorf("Failed Machine.GlobalVar need %v but %v", needVar, v)
	}
}

func TestMachine_Closure(t *testing.T) {
	needVar := 13
	code := [...]string{"&main", "jump", "subProc:", "&Cvar1", "load", "return", "main:", "13", "int", "&Cvar1", "store", "&subProc", "call", "exit"}
	m := NewMachine(code[:]).Run()

	v, ok := m.Pop().(int)
	if ok && v == needVar {
		t.Log("Pass Machine.Closure")
	} else {
		t.Errorf("Failed Machine.Closure need %v but %v", needVar, v)
	}
}
