package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	var err error
	args := os.Args[1:]
	if len(args) == 0 {
		REPL()
	}
	filename := args[0]
	switch path.Ext(filename) {
	case ".dbc", ".dobc":
		err = execDoByteCode(filename)
	case ".d", ".do":
		err = ececDo(filename)
	case ".dd", ".dodo":
		err = execDoDo(filename)
	}
	fmt.Println(err)
}
