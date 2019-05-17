# fifth
🛸Coding for code.

# Program <==> Data
> Program Is Data, Program Operating Data

# build
```
bash build_parser.sh
bash build.sh
```

# fifth
```go
func print(text){
    __fifcode__ {
        'text' load print
    }
    return
}

func sum(a,b){
    return a+b 
}

func main(){
    a = 10
    b = -8.5
    res = sum(a,b)
    print(res)
}

main()
```

### generator
```go
gen counter1(num){
    while(1){
        getv = yield num
        
        if getv == null {
            num = num + 1
        } else {
            num = num + getv
        }
    }
}

gen counter2(){
    num = 0
    while true{
        yield num
        num = num + 1
    }
}
```

# fif-code

### Usage
```
Hit CTRL+C or type "exit" or "quit" to quit.
>>> "hello world"
>>> print
hello world
>>>
```

```
>>> "square" func dup mul ret endfunc store
>>> 12 square print
144
```

### Closure
```
"main" func 
	"count1"
	"counter" call
	store
	"count2"
	"counter" call
	store

	"count1" call println
	"count1" call println
	"count2" call println
endfunc store

"counter" func
	"count" 0 store
	func "count" load 1 add dup "count" swap storev ret endfunc
	ret
endfunc store

"main" call
```

### Immediately Invoked Function Expression
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
	endfunc store
	&& jmp
	"sub" call
	ret
endfunc store

"main" call
```

### oop
```
"newMap" func 
	func 
		"val" arg "key" arg "opt" arg
		"opt" load "get" strEqul 0 &set equljmp
		"key" load load ret
		set:
		"opt" load "set" strEqul 0 &end equljmp
		"key" load "val" load stores
		end: ret
	endfunc
	ret 
endfunc store

"main" func
	"map1" newMap store
	"set" "name" "alice" "map1" load call
	"set" "age" 20 "map1" call
	"get" "name" nop "map1" call println
	"get" "age" nop "map1" call println
endfunc store

"main" call
```

# LICENSE
GPL-3.0