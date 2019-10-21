# do-code
> 本身作为底层语言，forth已经足够了，dolang在此基础上提供了更多的函数支持和抽象

## Index
- [do-code](#do-code)
	- [Index](#index)
	- [Usage](#usage)
	- [Closure](#closure)
	- [Immediately Invoked Function Expression](#immediately-invoked-function-expression)
	- [oop](#oop)

## Usage
```
Hit CTRL+C or type "exit" or "quit" to quit.
>>> "hello world"
>>> print
hello world
>>>
```

```
>>> "square" func dup mul ret endfunc set
>>> 12 square print
144
```

## Closure
```
"main" func 
	"count1"
	"counter" call
	set
	"count2"
	"counter" call
	set

	"count1" call println
	"count1" call println
	"count2" call println
endfunc set

"counter" func
	"count" 0 set
	func "count" get 1 add dup "count" swap set ret endfunc
	ret
endfunc set

"main" call
```

## Immediately Invoked Function Expression
```
"main" func
	&end jmp

	func
		&end jmp
		"foo" println
		end: 
		"bar" println
		ret
	endfunc call

	
	"sub" func 
		"hello world" println ret
	endfunc set
	&& jmp
	"sub" call
	ret
endfunc set

"main" call
```

## oop
```
"newMap" func 
	func 
		"val" arg "key" arg "opt" arg
		"opt" get "get" strEqul 0 &set equljmp
		"key" get get ret
		set:
		"opt" get "set" strEqul 0 &end equljmp
		"key" get "val" get set
		end: ret
	endfunc
	ret 
endfunc set

"main" func
	"map1" newMap set
	"set" "name" "alice" "map1" get call
	"set" "age" 20 "map1" call
	"get" "name" nop "map1" call println
	"get" "age" nop "map1" call println
endfunc set

"main" call
```