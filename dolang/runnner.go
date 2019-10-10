package dodolang

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func fmtF64(f float64) string {
	if f == 0 {
		return "0"
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%.12f", f)
	ret := strings.Trim(buf.String(), "0")
	ret = strings.Trim(ret, ".")
	if ret == "" {
		return "0"
	} else {
		return ret
	}
}

type fifNumber float64

func (f *fifNumber) String() string {
	return fmtF64(float64(*f))
}

type Runner struct {
	VM *dodoVM
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
			if r.VM.CurrentFrame.SuperFrame != nil {
				// If the program does not return normally at
				// last, the stack frame is automatically judged.
				// If the execution stack is not empty,
				// the stack will continue execution.
				r.VM.CurrentFrame = r.VM.CurrentFrame.SuperFrame
				continue
			}
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
	case "equl":
		r.equl()
	case "nequl":
		r.nequl()
	case "gt":
		r.gt()
	case "ls":
		r.ls()
	case "dup":
		r.dup()
	case "not":
		r.not()
	case "neg":
		r.neg()
	case "and":
		r.and()
	case "or":
		r.or()
	case "swap":
		r.swap()
	case "stack_reverse":
		r.stackReverse()
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
	case "tjmp":
		r.tjmp()
	case "fjmp":
		r.fjmp()
	case "func":
		r._func()
	case "arg":
		r.arg()
	case "store":
		r.store()
	case "storei":
		r.storei()
	case "storeu":
		r.storeu()
	case "load":
		r.load()
	case "call":
		r.call()
	case "callr":
		r.callr()
	case "callx":
		r.callx()
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
				r.push(fifNumber(fnum))
				return
			}
			if fn, ok := r.CurProgram().Env.get(expr); ok {
				if ffn, ok := fn.(*Program); ok {
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

func (r *Runner) stackReverse() {
	r.VM.Data.Reverse()
}

func (r *Runner) add() {
	b := r.pop().(fifNumber)
	a := r.pop().(fifNumber)
	r.push(a + b)
}

func (r *Runner) sub() {
	b := r.pop().(fifNumber)
	a := r.pop().(fifNumber)
	r.push(a - b)
}

func (r *Runner) mul() {
	b := r.pop().(fifNumber)
	a := r.pop().(fifNumber)
	r.push(a * b)
}

func (r *Runner) div() {
	b := r.pop().(fifNumber)
	a := r.pop().(fifNumber)
	r.push(a / b)
}

func (r *Runner) mod() {
	b := r.pop().(fifNumber)
	a := r.pop().(fifNumber)
	r.push(fifNumber(int(a) % int(b)))
}

func (r *Runner) char() {
	idx := int(r.pop().(fifNumber))
	str := r.pop().(string)
	if idx >= len(str) || idx < 0 {
		r.push("")
	} else {
		r.push(fifNumber(str[idx]))
	}
}

func (r *Runner) len() {
	str := r.pop().(string)
	r.push(fifNumber(len(str)))
}

func (r *Runner) print() {
	v := r.pop()
	if v == nil {
		fmt.Print("null")
	} else if f, ok := v.(fifNumber); ok {
		fmt.Print(f.String())
	} else {
		fmt.Print(v)
	}
}

func (r *Runner) println() {
	v := r.pop()
	if v == nil {
		fmt.Println("null")
	} else if f, ok := v.(fifNumber); ok {
		fmt.Println(f.String())
	} else {
		fmt.Println(v)
	}
}

func isTRUE(v interface{}) bool {
	if n, ok := v.(float64); ok {
		if n > 0 {
			return true
		}
	}
	return false
}

func TRUE() fifNumber  { return fifNumber(1) }
func FALSE() fifNumber { return fifNumber(0) }

// bool
func (r *Runner) not() {
	a := r.pop()
	if isTRUE(a) {
		r.push(FALSE())
	} else {
		r.push(TRUE())
	}
}
func (r *Runner) neg() {
	a := r.pop()
	if v, ok := a.(fifNumber); ok {
		r.push(-v)
	} else {
		r.push(a)
		r.not()
	}
}
func (r *Runner) and() {
	a := r.pop()
	b := r.pop()
	if isTRUE(a) && isTRUE(b) {
		r.push(TRUE())
	} else {
		r.push(FALSE())
	}
}
func (r *Runner) or() {
	a := r.pop()
	b := r.pop()
	if isTRUE(a) || isTRUE(b) {
		r.push(TRUE())
	} else {
		r.push(FALSE())
	}
}

// fifequl
func fif_func_equl(a, b Program) bool {
	return false
}

func fif_equl(a, b interface{}) bool {
	if a == nil {
		if b == nil {
			return true
		} else {
			return false
		}
	} else if b == nil {
		return false
	}
	if va, ok := a.(fifNumber); ok {
		if vb, ok := b.(fifNumber); ok {
			if va == vb {
				return true
			}
		}
	}
	if va, ok := a.(string); ok {
		if vb, ok := b.(string); ok {
			if va == vb {
				return true
			}
		}
	}
	if va, ok := a.(Program); ok {
		if vb, ok := b.(Program); ok {
			if fif_func_equl(va, vb) {
				return true
			}
		}
	}
	return false
}

// equl nequl
// num str nil func bool(num)
func (r *Runner) nequl() {
	a := r.pop()
	b := r.pop()
	if fif_equl(a, b) {
		r.push(FALSE())
	} else {
		r.push(TRUE())
	}
}

func (r *Runner) equl() {
	a := r.pop()
	b := r.pop()
	if fif_equl(a, b) {
		r.push(TRUE())
	} else {
		r.push(FALSE())
	}
}

func (r *Runner) gt() {
	TRUE := func() { r.push(fifNumber(1)) }
	FALSE := func() { r.push(fifNumber(0)) }
	b := r.pop()
	a := r.pop()
	if va, ok := a.(fifNumber); ok {
		if vb, ok := b.(fifNumber); ok {
			if va > vb {
				TRUE()
			} else {
				FALSE()
			}
		} else {
			FALSE()
		}
		return
	}
	if va, ok := a.(string); ok {
		if vb, ok := b.(string); ok {
			if va > vb {
				TRUE()
			} else {
				FALSE()
			}
		} else {
			FALSE()
		}
		return
	}
	// if va, ok := a.(Program); ok {
	// 	if vb, ok := b.(Program); ok {
	// 		if va == vb {
	// 			TRUE()
	// 		} else {
	// 			FALSE()
	// 		}
	// 		return
	// 	} else {
	// 		FALSE()
	// 	}
	// }
	FALSE()
}

func (r *Runner) ls() {
	TRUE := func() { r.push(fifNumber(1)) }
	FALSE := func() { r.push(fifNumber(0)) }
	b := r.pop()
	a := r.pop()
	if va, ok := a.(fifNumber); ok {
		if vb, ok := b.(fifNumber); ok {
			if va < vb {
				TRUE()
			} else {
				FALSE()
			}
		} else {
			FALSE()
		}
		return
	}
	if va, ok := a.(string); ok {
		if vb, ok := b.(string); ok {
			if va < vb {
				TRUE()
			} else {
				FALSE()
			}
		} else {
			FALSE()
		}
		return
	}
	// if va, ok := a.(Program); ok {
	// 	if vb, ok := b.(Program); ok {
	// 		if va == vb {
	// 			TRUE()
	// 		} else {
	// 			FALSE()
	// 		}
	// 		return
	// 	} else {
	// 		FALSE()
	// 	}
	// }
	FALSE()
}

// jump

func (r *Runner) jmp() {
	addr := r.pop().(fifNumber)
	r.CurProgram().PC = int(addr)
}

func (r *Runner) tjmp() {
	addr := r.pop().(fifNumber)
	if r.pop().(fifNumber) == 1.0 {
		r.CurProgram().PC = int(addr)
	}
}

func (r *Runner) fjmp() {
	addr := r.pop().(fifNumber)
	if r.pop().(fifNumber) == 0.0 {
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
	r.push(&Program{
		Code: code,
		Env:  NewVMEnv(r.CurProgram().Env),
		PC:   0,
	})
}

func (r *Runner) arg() {
	key := r.pop().(string)
	val := r.pop()
	r.CurProgram().Env.set(key, val, true)
}

func (r *Runner) _call(fn *Program) {
	newFn := fn.Clone()
	sf := NewStackFrame(r.VM.CurrentFrame, newFn)
	// init
	newFn.PC = 0
	r.VM.CurrentFrame = sf
}

func (r *Runner) _callx(fn *Program) {
	sf := NewStackFrame(r.VM.CurrentFrame, fn)
	// fmt.Printf("[callx]pc:%v,code:%v", fn.PC, fn.Code)
	// init
	r.VM.CurrentFrame = sf
}

func (r *Runner) callx() {
	nameOrFn := r.pop()
	if fnName, ok := nameOrFn.(string); ok {
		if fn, ok := r.CurProgram().Env.get(fnName); ok {
			if f, ok := fn.(*Program); ok {
				r._callx(f)
			}
		}
	} else if fn, ok := nameOrFn.(*Program); ok {
		r._callx(fn)
	}
}

func (r *Runner) callr() {
	nameOrFn := r.pop()

	r.stackReverse()

	if fnName, ok := nameOrFn.(string); ok {
		if fn, ok := r.CurProgram().Env.get(fnName); ok {
			if f, ok := fn.(*Program); ok {
				r._call(f)
			}
		}
	} else if fn, ok := nameOrFn.(*Program); ok {
		r._call(fn)
	}
}

func (r *Runner) call() {
	nameOrFn := r.pop()
	if fnName, ok := nameOrFn.(string); ok {
		if fn, ok := r.CurProgram().Env.get(fnName); ok {
			if f, ok := fn.(*Program); ok {
				r._call(f)
			}
		}
	} else if fn, ok := nameOrFn.(*Program); ok {
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

func (r *Runner) storei() {
	val := r.pop()
	key := r.pop().(string)
	r.CurProgram().Env.set(key, val, true)
}

func (r *Runner) storeu() {
	val := r.pop()
	key := r.pop().(string)
	r.CurProgram().Env.setUpper(key, val)
}

func (r *Runner) load() {
	key := r.pop().(string)
	val, _ := r.CurProgram().Env.get(key)
	if val == nil {
		r.VM.Push(val)
		return
	}
	r.VM.Push(val)
}

func (r *Runner) exit() {
	code := r.pop()
	if v, ok := code.(fifNumber); ok {
		os.Exit(int(v))
	} else {
		os.Exit(0)
	}
}
