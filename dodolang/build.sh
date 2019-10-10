./goyacc.exe -o dolang_parser.go -p Dolang dolang_parser.y
go build dolang_parser.go dolang.go lex.go