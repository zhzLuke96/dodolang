./goyacc.exe -o fif_parser.go -p Fif fif_parser.y
go build fif_parser.go fif.go lex.go labelStack.go