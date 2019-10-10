![dodo bird](https://upload.wikimedia.org/wikipedia/commons/9/9b/Frohawk_Dodo.png)
# dodolang
![LICENSE badge](https://img.shields.io/badge/license-GPL3.0-blue)
![build badge](https://img.shields.io/badge/build-error-red)
> Program <==> Data

🛸Do What U Want To Do.


# Overview
📑Game Engine Internal Script.

# Index
- [dodolang](#dodolang)
- [Overview](#overview)
- [Index](#index)
- [Install](#install)
- [Build](#build)
- [Usage](#usage)
- [Changelog](#changelog)
- [LICENSE](#license)

# Install
download form [releases](#) or build form source.

***build requirement***
- [goyacc](https://godoc.org/golang.org/x/tools/cmd/goyacc)

# Build
```
# go get dodolang
go get github.com/zhzluke96/dodolang/dolang
go get github.com/zhzluke96/dodolang/dodolang

cd $(go env GOPATH)/src/github.com/zhzluke96/dodolang/dodolang/
./build.sh

# build tools...
```

# Usage
```go
func print(text){
    __do__ {
        'text' load print
    }
    return
}

func add(a,b){
    return a+b 
}

func main(){
    a = 10
    b = -8.5
    res = add(a,b)
    print(res)
}

main()
```

# Changelog
- rename to dolang & dodolang
- Reorganization working directory

# LICENSE
GPL-3.0