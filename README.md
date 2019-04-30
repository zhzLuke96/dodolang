# fifth
🛸Coding for code.

# Program <==> Data
> Program Is Data, Program Operating Data


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

func inner_fib(a,b){
    var temp = a + b
    a = b
    b = temp
    print(a)
    count = count + 1
    if(count == 10){
        return
    }
    inner_fib(a,b)
}
inner_fib(1,1)
```

# LICENSE
GPL-3.0