# Macros system

# Range
```go
gen lib_range(start,end,step){
    step = if (step == nop) then 1 else step ;
    count = start
    onup = start < end
    while(true){
        count = count + step
        yield count
        if enup & count > end {
            break
        } else {
            if count < end {
                break
            }
        }
    }
}

macro range{
    (\d+) => {
        return lib_range(0,$1)
    }
    (\d+)\.\.\.(\d+) => {
        return lib_range($1,$2)
    }
    (\d+),(\d+)\.\.\.(\d+) => {
        return lib_range($1,$3,$2-$1)
    }
}

range![5] // => [0,1,2,3,4,5]
range![5...10] // => [5,6,7,8,9,10]
range![5,10...20] // => [5,10,15,20]
range![50,40...0] // => [50,40,30,20,10,0]
```

# more

```
'numis' macro
'0' tpl 'is zero' print tplend
'[13579]' tpl "'is odd num'" "print" tplend
'[2468]|10' tpl "'is even num'" "print" tplend
'' tpl "'is big'" "print" tplend
macroend
```

```go
macro swap{
    (\w+),(\w+) => {
        load!($1)
        set!($1,load!($2))
        store!($2)
    }
    _ => {}
}
```

```go
macro vac{
    // (\w+,?)+
    $items:($x:(\w+),?)+ => {
        arr = newArr()
        index = 0
        $item {
            arr[index] = x
            index = index + 1
        }
        return arr
    }
    (\w+) => {
        arr = newArr()
        arr[0] = load!($1)
        return arr
    }
}
```
