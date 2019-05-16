package machine

import (
	"../stack"
)

type Program struct {
	Code []string
	Env  *VMEnv
	PC   int
}

func (p Program) Expr() string {
	return p.Code[p.PC]
}

func (p Program) Clone() *Program {
	return &Program{p.Code, NewVMEnv(p.Env.Super), 0}
}

type StackFrame struct {
	SuperFrame *StackFrame
	Func       *Program
}

func NewStackFrame(sf *StackFrame, fn *Program) *StackFrame {
	return &StackFrame{
		SuperFrame: sf,
		Func:       fn,
	}
}

type VMEnv struct {
	Table map[string]interface{}
	Super *VMEnv
}

func NewVMEnv(super *VMEnv) *VMEnv {
	return &VMEnv{
		Table: make(map[string]interface{}),
		Super: super,
	}
}

func (v *VMEnv) get(key string) (interface{}, bool) {
	if _, ok := v.Table[key]; !ok {
		if v.Super == nil {
			return nil, false
		}
		return v.Super.get(key)
	}
	return v.Table[key], true
}

func (v *VMEnv) set(key string, val interface{}, init bool) {
	if init {
		v.Table[key] = val
		return
	}
	if _, ok := v.Table[key]; ok {
		v.Table[key] = val
	} else if !v.setScope(key, val) {
		v.Table[key] = val
	}
}

func (v *VMEnv) setScope(key string, val interface{}) bool {
	if _, ok := v.Table[key]; ok {
		v.Table[key] = val
		return true
	} else if v.Super != nil {
		return v.Super.setScope(key, val)
	}
	return false
}

type fifVM struct {
	Data         stack.Stack
	CurrentFrame *StackFrame
}

func NewFifVM(code []string) *fifVM {
	return &fifVM{
		Data:         stack.Stack{},
		CurrentFrame: NewStackFrame(nil, &Program{code, NewVMEnv(nil), 0}),
	}
}

func (f *fifVM) Pop() interface{} {
	v, _ := f.Data.Pop()
	return v
}

func (f *fifVM) Push(v interface{}) {
	f.Data.Push(v)
}

func (f *fifVM) Top() interface{} {
	val, _ := f.Data.Top()
	return val
}

func (f *fifVM) Run() {
	r := Runner{f}
	r.Run()
}
