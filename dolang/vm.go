package dolang

import (
	"fmt"
	"time"

	"./bignum"
)

type stackFrame struct {
	SuperFrame *stackFrame
	Function   *doProgram
}

func newStackFrame(ssf *stackFrame, fn *doProgram) *stackFrame {
	return &stackFrame{
		SuperFrame: ssf,
		Function:   fn,
	}
}

type doVM struct {
	Data         *dataStack
	CurrentFrame *stackFrame
	Scheduler    *doScheduler
	WordSet      *vmWordSet
}

func NewVM(code [][]byte) *doVM {
	return &doVM{
		Data:         newDataStack(),
		CurrentFrame: newStackFrame(nil, &doProgram{code, newCtx(nil), 0, nil}),
		Scheduler:    newDoScheduler(),
		WordSet:      gDictionary,
	}
}

func (vm *doVM) CurProgram() *doProgram {
	return vm.CurrentFrame.Function
}

func (vm *doVM) CurExpr() []byte {
	return vm.CurProgram().Expr()
}

func (vm *doVM) Run() {
	for {
		if vm.CurrentFrame == nil {
			if !vm.Scheduler.Empty() {
				vm.CurrentFrame = vm.Scheduler.Dequeue()
				continue
			}
			if vm.Scheduler.BlockCounter != 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		}
		program := vm.CurProgram()
		if program.PC == len(program.Code) {
			if !vm.Scheduler.Empty() {
				vm.CurrentFrame = vm.Scheduler.Dequeue()
				continue
			}
			if vm.Scheduler.BlockCounter != 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			if vm.CurrentFrame.SuperFrame != nil {
				// If the program does not return normally at
				// last, the stack frame is automatically judged.
				// If the execution stack is not empty,
				// the stack will continue execution.
				vm.CurrentFrame = vm.CurrentFrame.SuperFrame
				continue
			}
			break
		}
		// eval expr
		vm.Eval(vm.CurExpr())
		program.PC++
	}
	if !vm.Scheduler.Empty() {
		vm.Run()
	}
}

func (vm *doVM) CallProgram(fn *doProgram) {
	if fn == nil {
		return
	}
	if fn.Native != nil {
		fn.Native(vm)
	}
	vm.CurrentFrame = newStackFrame(vm.CurrentFrame, fn)
}

func (vm *doVM) matchMethod(expr []byte) (nativeFunc, bool) {
	return vm.WordSet.Get(string(expr))
}

func (vm *doVM) Eval(expr []byte) {
	if isDoString(expr) {
		vm.Data.Push(doStr(expr[1 : len(expr)-1]))
		return
	}
	if isDoNumber(expr) {
		vm.Data.Push(doNum(*bignum.Eval(string(expr))))
		return
	}
	if m, ok := vm.matchMethod(expr); ok {
		m(vm)
		return
	}
	fmt.Printf("[LOG] Unkonw expr [%v]\n", expr)
	// vm.Data.Push(doData(expr))
}

func (vm *doVM) DispatchTask(sf *stackFrame) {
	if vm.Scheduler.Empty() {
		vm.CurrentFrame = sf
		return
	}
	vm.CurrentFrame = vm.Scheduler.Dequeue()
	vm.Scheduler.Eequeue(sf)
}
