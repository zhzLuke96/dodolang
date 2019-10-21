package dolang

import (
	"bytes"
	"math"
	"regexp"
	"syscall"
	"unsafe"

	"./bignum"
)

func isDoString(expr []byte) bool {
	if !(expr[0] == '"' || expr[0] == '\'') {
		return false
	}
	if bytes.HasPrefix(expr, []byte("'")) && bytes.HasSuffix(expr, []byte("'")) {
		return true
	}
	if bytes.HasPrefix(expr, []byte("\"")) && bytes.HasSuffix(expr, []byte("\"")) {
		return true
	}
	return false
}

var numberReg = regexp.MustCompile("^[+-]?(\\d+\\.\\d+|\\d+)$")

func isDoNumber(expr []byte) bool {
	return numberReg.Match(expr)
}

func bytesToNum(bs []byte) int64 {
	var nums int64
	for i, b := range bs {
		nums += int64(b) * int64(math.Pow(256, float64(i)))
	}
	return nums
}

func intPtr(n int) uintptr {
	return uintptr(n)
}

func strPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func dllProcCall(dllName, procName string, args ...uintptr) {
	dll := syscall.NewLazyDLL(dllName)
	proc := dll.NewProc(procName)
	proc.Call(args...)
}

func reverseBytes(bs *[]byte) *[]byte {
	length := len(*bs)
	if length <= 1 {
		return bs
	}
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		(*bs)[i], (*bs)[j] = (*bs)[j], (*bs)[i]
	}
	return bs
}

func doBoolean(v doValue) bool {
	if isProgram(v) {
		return true
	}
	if v.Len() == 0 {
		return false
	}
	bs := v.ToBytes()
	for _, b := range bs {
		if b != byte(0) {
			return true
		}
	}
	return false
}

func getTopAndSecond(vm *doVM) (doValue, doValue, error) {
	top, err := vm.Data.Pop()
	if err != nil {
		return nil, nil, err
	}
	second, err := vm.Data.Pop()
	if err != nil {
		return top, nil, err
	}
	return top, second, nil
}

func getTopSecondAndThrid(vm *doVM) (doValue, doValue, doValue, error) {
	top, err := vm.Data.Pop()
	if err != nil {
		return nil, nil, nil, err
	}
	second, err := vm.Data.Pop()
	if err != nil {
		return top, nil, nil, err
	}
	thrid, err := vm.Data.Pop()
	if err != nil {
		return top, second, nil, err
	}
	return top, second, thrid, nil
}

func newArrayCtx() *doCtx {
	tb := newCtxMap()
	tb.Set("length", doNum(*bignum.IntToBigNum(0)))
	tb.Set("push", doProgram{
		Native: func(v *doVM) {

		},
	})
	tb.Set("pop", doProgram{
		Native: func(v *doVM) {

		},
	})
	return &doCtx{
		Super: nil,
		Table: tb,
	}
}
