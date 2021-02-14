# TODO

### 语法

```js
var Y = f => (x => x(x))(x => f(y => x(y)));
var F = g => n => n == 0 ? 1 : n * g(n-1);

var FACT = Y(F)
FACT(5) // =>  120
```

之前的语法太golang了，还是js这种简洁的看着舒服，从语法层面简化逻辑

### 宏系统

大概两种方案：
1. 宏生成的是代码文本，等于是一个 `() => String` 函数，需要预处理
2. 直接返回对象，操作起来和函数一样，但是接受的不是参数是exprs语句，属于运行时


```js
macro vec {
    [ $:expr:* ] >>= {
        ...
    }
    [] >>= {
        return '';
    }
}

```

### 系统调用
比较难，感觉是golang没法做，需要系统级语言，比如rust c才可以直接做系统调用...

需要调研


### 类型检查

类似于python的运行时类型检查玩法


```js
var fn = (x:Str):Str => x;
// 调用的时候
fn('x');
// 先检查变量
if (typeof x != 'string') throw err;
// 返回的时候再检查
```

### 重载模式匹配退版本

```js
var fact = (n = 0):number => 1;
var fact = (n: number):number => n * fact(n-1);
```
根据顺序，如果第一个匹配类型错误，自动用下一个同名重载调用

类似于Haskell的模式匹配

### 公式推导

x = [1,3,5...100];


### 虚拟机级别调试
这个应该好弄

