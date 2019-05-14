package machine

import (
	"fmt"
	"os"
	"strconv"
)

type Runner struct {
	VM *fifVM
}

func (r *Runner) CurProgram() *Program {
	return r.VM.CurrentFrame.Func
}

func (r *Runner) Expr() string {
	return r.CurProgram().Expr()
}

func (r *Runner) Run() {
	for {
		program := r.CurProgram()
		if program.PC == len(program.Code) {
			break
		}
		r.Eval(r.Expr())
		program.PC++
	}
}

func (r *Runner) Eval(expr string) {
	switch expr {
	case "nop":
		r.push(nil)
	case "add":
		r.add()
	case "sub":
		r.sub()
	case "mul":
		r.mul()
	case "div":
		r.div()
	case "mod":
		r.mod()
	case "dup":
		r.dup()
	case "swap":
		r.swap()
	case "print":
		r.print()
	case "println":
		r.println()
	case "char":
		r.char()
	case "len":
		r.len()
	case "jmp":
		r.jmp()
	case "nequljmp":
		r.nequljmp()
	case "gtjmp":
		r.gtjmp()
	case "equljmp":
		r.equljmp()
	case "lsjmp":
		r.lsjmp()
	case "func":
		r._func()
	case "arg":
		r.arg()
	case "store":
		r.store()
	case "storev":
		r.storev()
	case "stores":
		r.stores()
	case "load":
		r.load()
	case "call":
		r.call()
	case "ret":
		r.ret()
	case "exit":
		r.exit()
	case "":
		return
	// case "push": r.push()
	// case "pop": r.pop()
	default:
		if expr[0] == '\'' || expr[0] == '"' {
			r.push(expr[1 : len(expr)-1])
		} else {
			fnum, err := strconv.ParseFloat(expr, 64)
			if err == nil {
				r.push(fnum)
				return
			}
			if fn, ok := r.CurProgram().Env.get(expr); ok {
				if ffn, ok := fn.(Program); ok {
					r._call(ffn)
					return
				}
			}
			fmt.Printf("[LOG] got unknow expr [%v]\n", expr)
		}
	}
}

func (r *Runner) pop() interface{} {
	return r.VM.Pop()
}

func (r *Runner) push(v interface{}) {
	r.VM.Push(v)
}

func (r *Runner) top() interface{} {
	return r.VM.Top()
}

func (r *Runner) dup() {
	r.push(r.top())
}

func (r *Runner) swap() {
	t := r.pop()
	l := r.pop()
	r.push(t)
	r.push(l)
}

func (r *Runner) add() {
	b := r.pop().(float64)
	a := r.pop().(float64)
	r.push(a + b)
}

func (r *Runner) sub() {
	b := r.pop().(float64)
	a := r.pop().(float64)
	r.push(a - b)
}

func (r *Runner) mul() {
	b := r.pop().(float64)
	a := r.pop().(float64)
	r.push(a * b)
}

func (r *Runner) div() {
	b := r.pop().(float64)
	a := r.pop().(float64)
	r.push(a / b)
}

func (r *Runner) mod() {
	b := r.pop().(float64)
	a := r.pop().(float64)
	r.push(float64(int(a) % int(b)))
}

func (r *Runner) char() {
	idx := int(r.pop().(float64))
	str := r.pop().(string)
	if idx >= len(str) || idx < 0 {
		r.push("")
	} else {
		r.push(float64(str[idx]))
	}
}

func (r *Runner) len() {
	str := r.pop().(string)
	r.push(float64(len(str)))
}

func (r *Runner) print() {
	fmt.Print(r.pop())
}

func (r *Runner) println() {
	fmt.Println(r.pop())
}

// jump

func (r *Runner) jmp() {
	addr := r.pop().(float64)
	r.CurProgram().PC = int(addr)
}

func (r *Runner) nequljmp() {
	addr := r.pop().(float64)
	a := r.pop().(float64)
	b := r.pop().(float64)
	if a != b {
		r.CurProgram().PC = int(addr)
	}
}

func (r *Runner) equljmp() {
	addr := r.pop().(float64)
	a := r.pop().(float64)
	b := r.pop().(float64)
	if a == b {
		r.CurProgram().PC = int(addr)
	}
}

func (r *Runner) gtjmp() {
	addr := r.pop().(float64)
	b := r.pop().(float64)
	a := r.pop().(float64)
	if a > b {
		r.CurProgram().PC = int(addr)
	}
}

func (r *Runner) lsjmp() {
	addr := r.pop().(float64)
	b := r.pop().(float64)
	a := r.pop().(float64)
	if a > b {
		r.CurProgram().PC = int(addr)
	}
}

// function
func (r *Runner) _func() {
	code := []string{}
	program := r.CurProgram()
	funcCount := 1
	for {
		program.PC++
		e := r.Expr()
		if e == "func" {
			funcCount++
		} else if e == "endfunc" {
			funcCount--
			if funcCount == 0 {
				break
			}
		}
		code = append(code, e)
	}
	r.push(Program{
		Code: code,
		Env:  NewVMEnv(r.CurProgram().Env),
		PC:   0,
	})
}

func (r *Runner) arg() {
	key := r.pop().(string)
	val := r.pop()
	r.CurProgram().Env.set(key, val, false)
}

func (r *Runner) _call(fn Program) {
	sf := NewStackFrame(r.VM.CurrentFrame, fn.Clone())
	// init
	fn.PC = 0
	r.VM.CurrentFrame = sf
}

func (r *Runner) call() {
	nameOrFn := r.pop()
	if fnName, ok := nameOrFn.(string); ok {
		switch fnName {
		case "print":
			r.print()
			return
		case "println":
			r.println()
			return
		case "len":
			r.len()
			return
		}
		if fn, ok := r.CurProgram().Env.get(fnName); ok {
			r._call(fn.(Program))
		}
	} else if fn, ok := nameOrFn.(Program); ok {
		r._call(fn)
	}
}

func (r *Runner) ret() {
	r.VM.CurrentFrame = r.VM.CurrentFrame.SuperFrame
}

func (r *Runner) store() {
	val := r.pop()
	key := r.pop().(string)
	r.CurProgram().Env.set(key, val, false)
}

func (r *Runner) storev() {
	val := r.pop()
	key := r.pop().(string)
	r.CurProgram().Env.set(key, val, true)
}

func (r *Runner) stores() {
	val := r.pop()
	key := r.pop().(string)
	r.CurProgram().Env.setSuper(key, val)
}

func (r *Runner) load() {
	key := r.pop().(string)
	val, _ := r.CurProgram().Env.get(key)
	r.VM.Push(val)
}

func (r *Runner) exit() {
	code := r.pop()
	if v, ok := code.(float32); ok {
		os.Exit(int(v))
	} else {
		os.Exit(0)
	}
}
