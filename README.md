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
10 store_var1
1 store_var2
1 store_var3
fib:
load_var3 println
load_var2 load_var3 plus
load_var3 store_var2 store_var3
load_var1 1 sub dup store_var1
0 greater if &fib call return then return
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