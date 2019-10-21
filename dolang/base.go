package dolang

import (
	"errors"
	"reflect"

	"./skiplist"
	"./stack"
)

const blankByte = byte(' ')

type doValue interface {
	ToBytes() []byte
	ToFunc() *doProgram
	Len() int
}

func typeName(v doValue) string {
	return reflect.TypeOf(v).Name()
}

func isProgram(dv doValue) bool {
	tp := typeName(dv)
	return tp == "doProgram" || tp == "doProgramPtr"
}

type doData []byte

func (data doData) ToBytes() []byte    { return data }
func (data doData) ToFunc() *doProgram { return nil }
func (data doData) Len() int           { return len(data) }

type doNum []byte

func (data doNum) ToBytes() []byte    { return data }
func (data doNum) ToFunc() *doProgram { return nil }
func (data doNum) Len() int           { return len(data) }

type doStr []byte

func (data doStr) ToBytes() []byte    { return data }
func (data doStr) ToFunc() *doProgram { return nil }
func (data doStr) Len() int           { return len(data) }

type nativeFunc func(*doVM)

type doProgram struct {
	Code [][]byte
	Env  *doCtx
	PC   int

	Native nativeFunc
}

func newProgram() *doProgram {
	return &doProgram{
		Env: newCtx(nil),
	}
}

func (p doProgram) Expr() []byte       { return p.Code[p.PC] }
func (p doProgram) ToFunc() *doProgram { return &p }
func (p doProgram) Len() int           { return len(p.Code) }

func (p doProgram) Clone() *doProgram {
	return &doProgram{p.Code, newCtx(p.Env.Super), 0, p.Native}
}

func (p doProgram) ToBytes() []byte {
	ret := []byte{}
	for _, v := range p.Code {
		if len(ret) != 0 {
			ret = append(ret, blankByte)
		}
		ret = append(ret, v...)
	}
	return ret
}

type doProgramPtr struct {
	p *doProgram
}

func (p doProgramPtr) ToBytes() []byte    { return p.p.ToBytes() }
func (p doProgramPtr) ToFunc() *doProgram { return p.p }
func (p doProgramPtr) Len() int           { return p.p.Len() }

// context
var errCantFoundKey = errors.New("cant found key")

type doCtx struct {
	Super *doCtx
	Table *ctxMap
}

func newCtx(super *doCtx) *doCtx {
	return &doCtx{
		Table: newCtxMap(),
		Super: super,
	}
}

func (ctx *doCtx) Get(key string) (doValue, bool) {
	if v, ok := ctx.Table.Get(key); ok {
		return v, true
	}
	if ctx.Super == nil {
		return nil, false
	}
	return ctx.Super.Get(key)
}

func (ctx *doCtx) GetTable(key string) (doValue, bool) {
	return ctx.Table.Get(key)
}

func (ctx *doCtx) GetSuperTable(key string) (doValue, bool) {
	if ctx.Super == nil {
		return nil, false
	}
	return ctx.Super.Table.Get(key)
}

func (ctx *doCtx) Set(key string, value doValue) {
	ctx.SetErr(key, value, true)
}

func (ctx *doCtx) SetErr(key string, value doValue, scoped bool) error {
	if _, ok := ctx.Table.Get(key); ok {
		ctx.Table.Set(key, value)
		return nil
	}
	if ctx.Super == nil {
		if scoped {
			ctx.Table.Set(key, value)
			return nil
		}
		return errCantFoundKey
	}
	if err := ctx.Super.SetErr(key, value, false); err != nil {
		ctx.Table.Set(key, value)
		return nil
	}
	return nil
}

func (ctx *doCtx) SetTable(key string, value doValue) {
	ctx.Table.Set(key, value)
}

func (ctx *doCtx) SetSuperTable(key string, value doValue) {
	if ctx.Super == nil {
		return
	}
	ctx.Super.Table.Set(key, value)
}

// --------------------------
// doDict
// --------------------------

type ctxMap struct {
	list *skiplist.SkipList
}

func newCtxMap() *ctxMap {
	return &ctxMap{skiplist.NewSkipList(-1)}
}

func (c ctxMap) Set(key string, val doValue) { c.list.Set(key, val) }
func (c ctxMap) Del(key string) bool         { return c.list.Del(key) }
func (c ctxMap) Len() int                    { return c.list.Len() }

func (c ctxMap) Get(key string) (doValue, bool) {
	val, ok := c.list.Get(key)
	if val == nil {
		return nil, false
	}
	return val.(doValue), ok
}

// --------------------------
// vmWordSet
// --------------------------

type vmWordSet struct {
	list *skiplist.SkipList
}

func newVmWordSet() *vmWordSet {
	return &vmWordSet{skiplist.NewSkipList(-1)}
}

func (v vmWordSet) Set(key string, val nativeFunc) { v.list.Set(key, val) }
func (v vmWordSet) Del(key string) bool            { return v.list.Del(key) }
func (v vmWordSet) Len() int                       { return v.list.Len() }

func (v vmWordSet) Get(key string) (nativeFunc, bool) {
	val, ok := v.list.Get(key)
	if val == nil {
		return nil, false
	}
	return val.(nativeFunc), ok
}

// --------------------------
// DataStack
// --------------------------

type dataStack struct {
	stack *stack.Stack
}

func newDataStack() *dataStack {
	return &dataStack{&stack.Stack{}}
}

func (d dataStack) Len() int         { return d.stack.Len() }
func (d dataStack) IsEmpty() bool    { return d.stack.IsEmpty() }
func (d dataStack) Cap() int         { return d.stack.Cap() }
func (d dataStack) Reverse()         { d.stack.Reverse() }
func (d dataStack) Push(val doValue) { d.stack.Push(val) }

func (d dataStack) Take(idx int) (doValue, bool) {
	if v, ok := d.stack.Take(idx); ok {
		if dv, ok := v.(doValue); ok {
			return dv, true
		}
	}
	return nil, false
}

func (d dataStack) Top() (doValue, error) {
	val, err := d.stack.Top()
	if val == nil {
		return nil, err
	}
	return val.(doValue), err
}
func (d dataStack) Pop() (doValue, error) {
	val, err := d.stack.Pop()
	if val == nil {
		return nil, err
	}
	return val.(doValue), err
}
