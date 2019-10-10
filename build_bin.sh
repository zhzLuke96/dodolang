# go get dodolang
# go get github.com/zhzluke96/dodolang/dolang
# go get github.com/zhzluke96/dodolang/dodolang

cd $(go env GOPATH)/src/github.com/zhzluke96/dodolang/dodolang/
./build.sh

cd $(go env GOPATH)/src/github.com/zhzluke96/dodolang/tools/exec
go build -o ./dodo.exe .
