if [ ! -f "./goyacc.exe" ];then
    go build -o goyacc.exe $(go env GOPATH)/src/golang.org/x/tools/cmd/goyacc
fi

./goyacc.exe -o parser.go -p Dolang parser.y
# go build parser.go main.go lex.go utils.go