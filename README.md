# fifth
🛸Coding for code.

# Program <==> Data
> Program Is Data, Program Operating Data

# build
```
bash build.sh
```

# Usage
```
Hit CTRL+C or type "exit" or "quit" to quit.
>>> "hello world"
>>> print
hello world
>>>
```

```
>>> square: dup mul return
>>> 12 num
>>> &square call print
144
```

# fif-code
```
0 store_1
1 store_2
1 store_3
func fib
load_3 println
load_2 load_3 +
load_3 store_2 store_3
load_1 1 + dup store_1
10 == "end" true_jump return
end:"fib" call
func_end
```

# fifth
```js
var count = 0

func fib(a,b){
    var temp = a + b
    a = b
    b = temp
    print(a)
    count = count + 1
    if(count == 10){
        return
    }
    fib(a,b)
}
fib(1,1)
```

# LICENSE
GPL-3.0