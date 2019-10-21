package dolang

import (
	"fmt"
	"math/rand"
	"time"

	"./bignum"
)

var (
	valueTRUE  = doData{1}
	valueFALSE = doData{0}

	gDictionary = newVmWordSet()
)

func defineVMInsts() {
	defineSysWords()
	defineStackWords()
	defineBoolWords()
	defineMathWords()
	defineDataWords()
	defineIOWords()
	defineCtxWords()

	defineFuncWords()
	defineObjectWords()
}

func defineStackWords() {
	gDictionary.Set("nop", func(vm *doVM) {
		// 存入空元素
		vm.Data.Push(nil)
	})

	gDictionary.Set("depth", func(vm *doVM) {
		// 获取堆栈长度
		depth := vm.Data.Len()
		vm.Data.Push(doData{byte(depth)})
	})
	gDictionary.Set("dup", func(vm *doVM) {
		// 重复top
		if v, err := vm.Data.Top(); err == nil {
			vm.Data.Push(v)
		}
	})
	gDictionary.Set("over", func(vm *doVM) {
		// 重复second
		length := vm.Data.Len()
		if length <= 1 {
			return
		}
		vm.Data.Push((*vm.Data.stack)[length-2].(doValue))
	})
	gDictionary.Set("swap", func(vm *doVM) {
		// 交换top和second
		if top, err := vm.Data.Pop(); err == nil {
			if second, err := vm.Data.Pop(); err == nil {
				vm.Data.Push(top)
				vm.Data.Push(second)
			} else {
				vm.Data.Push(top)
			}
		}
	})
	gDictionary.Set("drop", func(vm *doVM) {
		// 丢弃顶部堆栈
		vm.Data.Pop()
	})
	gDictionary.Set("clear", func(vm *doVM) {
		// 清空堆栈
		for {
			if _, err := vm.Data.Pop(); err != nil {
				return
			}
		}
	})
	gDictionary.Set("pick", func(vm *doVM) {
		// 丢弃顶部堆栈
		v, err := vm.Data.Pop()
		if err != nil {
			return
		}
		if isProgram(v) {
			vm.Data.Push(v)
			return
		}
		idx := bytesToNum(v.ToBytes())
		if iv, ok := vm.Data.Take(int(idx)); ok {
			vm.Data.Push(iv)
		} else {
			// [TODO] 报错
		}
	})
}

func defineBoolWords() {
	gDictionary.Set("not", func(vm *doVM) {
		// 判断顶部元素为 TRUE 或 FALSE
		if vm.Data.IsEmpty() {
			vm.Data.Push(valueTRUE)
			return
		}
		top, err := vm.Data.Top()
		if err != nil {
			vm.Data.Push(valueTRUE)
			return
		}
		if doBoolean(top) {
			vm.Data.Push(valueFALSE)
			return
		}
		vm.Data.Push(valueTRUE)

	})
	gDictionary.Set("and", func(vm *doVM) {
		// 与运算
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		if doBoolean(top) && doBoolean(sed) {
			vm.Data.Push(valueTRUE)
		} else {
			vm.Data.Push(valueFALSE)
		}
	})
	gDictionary.Set("or", func(vm *doVM) {
		// 或运算
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		if doBoolean(top) || doBoolean(sed) {
			vm.Data.Push(valueTRUE)
		} else {
			vm.Data.Push(valueFALSE)
		}
	})
	gDictionary.Set("xor", func(vm *doVM) {
		// 异或运算
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		if doBoolean(top) != doBoolean(sed) {
			vm.Data.Push(valueTRUE)
		} else {
			vm.Data.Push(valueFALSE)
		}
	})
	gDictionary.Set("xnor", func(vm *doVM) {
		// 同或运算
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		if doBoolean(top) == doBoolean(sed) {
			vm.Data.Push(valueTRUE)
		} else {
			vm.Data.Push(valueFALSE)
		}
	})

	// jump
	gDictionary.Set("jump", func(vm *doVM) {
		// 跳转
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		addr := bignum.BigNum(top.ToBytes())
		vm.CurProgram().PC = addr.ToInt()
	})
	gDictionary.Set("tjump", func(vm *doVM) {
		// 为真跳转
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		if !doBoolean(sed) {
			return
		}
		addr := bignum.BigNum(top.ToBytes())
		vm.CurProgram().PC = addr.ToInt()
	})
	gDictionary.Set("fjump", func(vm *doVM) {
		// 为假跳转
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		if doBoolean(sed) {
			return
		}
		addr := bignum.BigNum(top.ToBytes())
		vm.CurProgram().PC = addr.ToInt()
	})

	// skip
	gDictionary.Set("skip", func(vm *doVM) {
		// 跳过
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		addr := bignum.BigNum(top.ToBytes())
		vm.CurProgram().PC += addr.ToInt()

	})
	gDictionary.Set("tskip", func(vm *doVM) {
		// 为真跳过
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		if !doBoolean(sed) {
			return
		}
		addr := bignum.BigNum(top.ToBytes())
		vm.CurProgram().PC += addr.ToInt()

	})
	gDictionary.Set("fskip", func(vm *doVM) {
		// 为假跳过
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		if doBoolean(sed) {
			return
		}
		addr := bignum.BigNum(top.ToBytes())
		vm.CurProgram().PC += addr.ToInt()
	})
}

func defineMathWords() {
	gDictionary.Set("neg", func(vm *doVM) {
		// 变号
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		if typeName(top) == "doNum" {
			num := bignum.BigNum(top.ToBytes())
			num.Neg()
			vm.Data.Push(doNum(num))
		} else {
			vm.Data.Push(top)
		}
	})
	gDictionary.Set("succ", func(vm *doVM) {
		// ++
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		if typeName(top) == "doNum" {
			num := bignum.BigNum(top.ToBytes())
			num.Succ()
			vm.Data.Push(doNum(num))
		} else {
			vm.Data.Push(top)
		}
	})
	gDictionary.Set("pred", func(vm *doVM) {
		// --
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		if typeName(top) == "doNum" {
			num := bignum.BigNum(top.ToBytes())
			num.Pred()
			vm.Data.Push(doNum(num))
		} else {
			vm.Data.Push(top)
		}
	})
	gDictionary.Set("add", func(vm *doVM) {
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		num1 := bignum.BigNum(sed.ToBytes())
		num2 := bignum.BigNum(top.ToBytes())
		num1.Plus(&num2)
		vm.Data.Push(doNum(num1))
	})
	gDictionary.Set("sub", func(vm *doVM) {
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		num1 := bignum.BigNum(sed.ToBytes())
		num2 := bignum.BigNum(top.ToBytes())
		num1.Minus(&num2)
		vm.Data.Push(doNum(num1))
	})
	gDictionary.Set("mul", func(vm *doVM) {
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		num1 := bignum.BigNum(sed.ToBytes())
		num2 := bignum.BigNum(top.ToBytes())
		num1.Mul(&num2)
		vm.Data.Push(doNum(num1))

	})
	gDictionary.Set("div", func(vm *doVM) {
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		num1 := bignum.BigNum(sed.ToBytes())
		num2 := bignum.BigNum(top.ToBytes())
		num1.Div(&num2)
		vm.Data.Push(doNum(num1))

	})
	gDictionary.Set("mod", func(vm *doVM) {
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		num1 := bignum.BigNum(sed.ToBytes())
		num2 := bignum.BigNum(top.ToBytes())
		num1.Mod(&num2)
		vm.Data.Push(doNum(num1))

	})
	// push bool
	gDictionary.Set("equl", func(vm *doVM) {
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		num1 := bignum.BigNum(sed.ToBytes())
		num2 := bignum.BigNum(top.ToBytes())
		if num1.Equle(&num2) {
			vm.Data.Push(valueTRUE)
		} else {
			vm.Data.Push(valueFALSE)
		}
	})
	gDictionary.Set("less", func(vm *doVM) {
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		num1 := bignum.BigNum(sed.ToBytes())
		num2 := bignum.BigNum(top.ToBytes())
		if num1.Less(&num2) {
			vm.Data.Push(valueTRUE)
		} else {
			vm.Data.Push(valueFALSE)
		}
	})
	gDictionary.Set("great", func(vm *doVM) {
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		num1 := bignum.BigNum(sed.ToBytes())
		num2 := bignum.BigNum(top.ToBytes())
		if num1.Great(&num2) {
			vm.Data.Push(valueTRUE)
		} else {
			vm.Data.Push(valueFALSE)
		}
	})

	// rand
	gDictionary.Set("rand", func(vm *doVM) {
		// 随机数
		r := rand.Intn(INT_MAX)
		num := bignum.IntToBigNum(int64(r))
		vm.Data.Push(doNum(*num))
	})
	gDictionary.Set("randx", func(vm *doVM) {
		// 指定大小的随机数
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		max := bignum.BigNum(top.ToBytes())
		r := rand.Intn(max.ToInt())
		num := bignum.IntToBigNum(int64(r))
		vm.Data.Push(doNum(*num))
	})

	// trigo
	gDictionary.Set("cos", func(vm *doVM) {
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		p := bignum.BigNum(top.ToBytes())
		num := bignum.Cos(&p)
		vm.Data.Push(doNum(*num))
	})
	gDictionary.Set("sin", func(vm *doVM) {
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		p := bignum.BigNum(top.ToBytes())
		num := bignum.Sin(&p)
		vm.Data.Push(doNum(*num))
	})
	gDictionary.Set("tan", func(vm *doVM) {
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		p := bignum.BigNum(top.ToBytes())
		num := bignum.Tan(&p)
		vm.Data.Push(doNum(*num))
	})
	// 指定精度 泰勒逼近
	gDictionary.Set("cosx", func(vm *doVM) {

	})
	gDictionary.Set("sinx", func(vm *doVM) {

	})
	gDictionary.Set("tanx", func(vm *doVM) {

	})

	// 阶乘
	gDictionary.Set("fact", func(vm *doVM) {
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		num := bignum.BigNum(top.ToBytes())
		num.Factorial()
		vm.Data.Push(doNum(num))
	})
}

func defineDataWords() {
	gDictionary.Set("concat", func(vm *doVM) {
		b, err := vm.Data.Pop()
		if err != nil || isProgram(b) {
			vm.Data.Push(b)
			return
		}
		a, err := vm.Data.Pop()
		if err != nil || isProgram(a) {
			vm.Data.Push(a)
			vm.Data.Push(b)
			return
		}
		ab := append(a.ToBytes(), b.ToBytes()...)
		vm.Data.Push(doData(ab))
	})

	gDictionary.Set("head", func(vm *doVM) {
		// 取出第一个Byte
		top, err := vm.Data.Top()
		if err != nil || isProgram(top) {
			return
		}
		bs := top.ToBytes()
		if len(bs) == 0 {
			vm.Data.Push(nil)
			return
		}
		vm.Data.Push(doData(bs[:1]))
	})

	gDictionary.Set("tail", func(vm *doVM) {
		// 取出除了第一个的Byte
		top, err := vm.Data.Top()
		if err != nil || isProgram(top) {
			return
		}
		bs := top.ToBytes()
		if len(bs) == 0 {
			vm.Data.Push(nil)
			return
		}
		vm.Data.Push(doData(bs[1:]))
	})

	gDictionary.Set("last", func(vm *doVM) {
		// 取出末尾的Byte
		top, err := vm.Data.Top()
		if err != nil || isProgram(top) {
			return
		}
		bs := top.ToBytes()
		if len(bs) == 0 {
			vm.Data.Push(nil)
			return
		}
		vm.Data.Push(doData(bs[len(bs)-1:]))
	})

	gDictionary.Set("init", func(vm *doVM) {
		// 取出除了末尾的Byte
		top, err := vm.Data.Top()
		if err != nil || isProgram(top) {
			return
		}
		bs := top.ToBytes()
		if len(bs) == 0 {
			vm.Data.Push(nil)
			return
		}
		vm.Data.Push(doData(bs[:len(bs)-1]))
	})

	gDictionary.Set("length", func(vm *doVM) {
		// 取出上一个元素的长度
		top, err := vm.Data.Top()
		if err != nil {
			vm.Data.Push(doData{0})
			return
		}
		vm.Data.Push(doData{byte(top.Len())})
	})

	gDictionary.Set("null", func(vm *doVM) {
		// 判断上一个是不是空
		top, err := vm.Data.Top()
		if err != nil || top == nil || len(top.ToBytes()) == 0 {
			vm.Data.Push(valueTRUE)
			return
		}
		if isProgram(top) {
			vm.Data.Push(valueFALSE)
			return
		}
		vm.Data.Push(valueFALSE)
	})

	gDictionary.Set("reverse", func(vm *doVM) {
		// 反转
		top, err := vm.Data.Top()
		if err != nil || isProgram(top) {
			vm.Data.Push(valueFALSE)
			return
		}
		bs := top.ToBytes()
		reverseBytes(&bs)
		vm.Data.Push(doData(bs))
	})

	gDictionary.Set("maximum", func(vm *doVM) {
		// 取最大值

	})
	gDictionary.Set("minimum", func(vm *doVM) {
		// 取最小值

	})

	gDictionary.Set("sum", func(vm *doVM) {
		// 计算和

	})

	gDictionary.Set("take", func(vm *doVM) {
		// 取前几个元素

	})
	gDictionary.Set("drop", func(vm *doVM) {
		// 取后几个元素

	})

}

func defineSysWords() {
	gDictionary.Set("syscall", func(vm *doVM) {
		// 系统调用
		proc, err := vm.Data.Pop()
		if err != nil {
			return
		}
		procName := string(proc.ToBytes())
		dll, err := vm.Data.Pop()
		if err != nil {
			return
		}
		dllName := string(dll.ToBytes())
		ptrs := []uintptr{}
		for {
			arg, err := vm.Data.Pop()
			if err != nil {
				break
			}
			var ptr uintptr
			if typeName(arg) == "doStr" {
				str := string(arg.ToBytes())
				ptr = strPtr(str)
			} else {
				bn := bignum.BigNum(arg.ToBytes())
				num := bn.ToInt()
				ptr = intPtr(num)
			}
			ptrs = append([]uintptr{ptr}, ptrs...)
		}
		dllProcCall(dllName, procName, ptrs...)
	})

	gDictionary.Set("const", func(vm *doVM) {
		// 定义静态变量
		key, err := vm.Data.Pop()
		if err != nil {
			return
		}
		value, err := vm.Data.Pop()
		if err != nil {
			return
		}
		gDictionary.Set(string(key.ToBytes()), func(vm *doVM) {
			vm.Data.Push(value)
		})
	})
	gDictionary.Set("get", func(vm *doVM) {
		key, err := vm.Data.Pop()
		if err != nil {
			return
		}
		ctx := vm.CurProgram().Env
		if v, ok := ctx.Get(string(key.ToBytes())); ok {
			vm.Data.Push(v)
		}
	})

	gDictionary.Set("set", func(vm *doVM) {
		val, err := vm.Data.Pop()
		if err != nil {
			return
		}
		key, err := vm.Data.Pop()
		if err != nil {
			return
		}
		ctx := vm.CurProgram().Env
		ctx.Set(string(key.ToBytes()), val)
	})
}

func defineIOWords() {
	gDictionary.Set("println", func(vm *doVM) {
		if b, err := vm.Data.Pop(); err == nil {
			tp := typeName(b)
			if tp == "doNum" {
				num := bignum.BigNum(b.ToBytes())
				fmt.Println(num.String())
			} else {
				fmt.Println(string(b.ToBytes()))
			}
		}
	})
	gDictionary.Set("inputln", func(vm *doVM) {
		// 获取输入
		inp := ""
		fmt.Scanln(&inp)
		if isDoNumber([]byte(inp)) {
			vm.Data.Push(doNum(*bignum.Eval(inp)))
		} else {
			vm.Data.Push(doStr(inp))
			fmt.Printf("[LOG] inp = %v\n", inp)
		}
	})
	gDictionary.Set("now", func(vm *doVM) {
		// 获取UTC时间
		t := time.Now().UnixNano()
		num := bignum.FloatToBigNum(float64(t) / 1000000000)
		vm.Data.Push(doNum(*num))
	})
}

func defineCtxWords() {

}

func defineFuncWords() {
	gDictionary.Set("call", func(vm *doVM) {
		// 调用函数
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		if isProgram(top) {
			fn := top.ToFunc().Clone()
			// init
			fn.PC = 0
			vm.CallProgram(fn)
		} else {
			// [TODO] ERROR catch
			vm.Data.Push(top)
		}
	})
	gDictionary.Set("callr", func(vm *doVM) {
		vm.Data.Reverse()
		vm.Eval([]byte("call"))
	})
	gDictionary.Set("callx", func(vm *doVM) {
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		if isProgram(top) {
			fn := top.ToFunc().Clone()
			vm.CallProgram(fn)
		} else {
			// [TODO] ERROR catch
			vm.Data.Push(top)
		}
	})
	gDictionary.Set("func", func(vm *doVM) {
		code := [][]byte{}
		program := vm.CurProgram()
		funcCount := 1
		for {
			program.PC++
			e := vm.CurExpr()
			estr := string(e)
			if estr == "func" {
				funcCount++
			} else if estr == "endfunc" {
				funcCount--
				if funcCount == 0 {
					break
				}
			}
			code = append(code, e)
		}
		vm.Data.Push(doProgramPtr{&doProgram{
			Code: code,
			Env:  newCtx(program.Env),
		}})
	})

	gDictionary.Set("arg", func(vm *doVM) {
		key, err := vm.Data.Pop()
		if err != nil {
			return
		}
		keys := string(key.ToBytes())
		val, err := vm.Data.Pop()
		if err != nil {
			return
		}
		vm.CurProgram().Env.SetTable(keys, val)
	})

	gDictionary.Set("ret", func(vm *doVM) {
		if !vm.Scheduler.Empty() {
			vm.CurrentFrame = vm.Scheduler.Dequeue()
			return
		}
		vm.CurrentFrame = vm.CurrentFrame.SuperFrame
	})

	// coroutine
	gDictionary.Set("do", func(vm *doVM) {
		// 将函数注册为协程
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		if !isProgram(top) {
			// [TODO] ERROR
			vm.Data.Push(top)
			return
		}
		fn := top.ToFunc().Clone()
		sf := newStackFrame(vm.CurrentFrame, fn)
		vm.Scheduler.Eequeue(sf)
	})
	gDictionary.Set("yield", func(vm *doVM) {
		// 协程出让
		vm.Scheduler.Eequeue(vm.CurrentFrame)
		vm.CurrentFrame = vm.Scheduler.Dequeue()
	})
	gDictionary.Set("block", func(vm *doVM) {
		callback := vm.Scheduler.Block(vm.CurrentFrame)
		p := newProgram()
		p.Native = func(v *doVM) {
			callback()
		}
		vm.Data.Push(p)
	})
	gDictionary.Set("sleep", func(vm *doVM) {
		// 等待
		top, err := vm.Data.Pop()
		if err != nil {
			return
		}
		if typeName(top) != "doNum" {
			// [TODO] ERROR
			return
		}
		num := bignum.BigNum(top.ToBytes())
		callback := vm.Scheduler.Block(vm.CurrentFrame)
		go func() {
			time.Sleep(time.Second * time.Duration(num.ToInt()))
			callback()
		}()
		// 出让
		vm.CurrentFrame = vm.Scheduler.Dequeue()
	})
}

func defineObjectWords() {
	gDictionary.Set("newMap", func(vm *doVM) {
		// 创建对象
		vm.Data.Push(doProgram{Env: newCtx(nil)})
	})
	gDictionary.Set("getv", func(vm *doVM) {
		// 获取对象属性
		top, sed, err := getTopAndSecond(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			return
		}
		key := string(top.ToBytes())
		if !isProgram(sed) {
			vm.Data.Push(sed)
			return
		}
		fn := sed.ToFunc()
		value, _ := fn.Env.GetTable(key)
		vm.Data.Push(value)
	})
	gDictionary.Set("setv", func(vm *doVM) {
		// 设置对象属性
		top, sed, thr, err := getTopSecondAndThrid(vm)
		if err != nil {
			if top != nil {
				vm.Data.Push(top)
			}
			if sed != nil {
				vm.Data.Push(sed)
			}
			return
		}
		key := string(sed.ToBytes())
		if !isProgram(thr) {
			vm.Data.Push(thr)
			return
		}
		thr.ToFunc().Env.SetTable(key, top)
	})

	gDictionary.Set("newArr", func(vm *doVM) {
		// 创建数组对象

		vm.Data.Push(doValue(doProgram{Env: newCtx(nil)}))
	})
}
