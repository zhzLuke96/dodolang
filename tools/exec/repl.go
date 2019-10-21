package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zhzluke96/dodolang/parser"
)

func REPL() {
	fmt.Print("[dodolang code]\n")
	fmt.Print("Hit CTRL+C or type \"quit\" to quit.")
	vm := dolang.NewVM([]string{})
	for {
		fmt.Print("\n>>> ")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		if strings.TrimSpace(string(data)) == "quit" {
			return
		}

		dolangcode, err := parser.Parse(data)
		if err != nil {
			fmt.Println(err)
		}

		code := dolang.Tokenizer(dolangcode)
		runner.VM.CurrentFrame.Func.Code = labelLoad(code)
		runner.VM.CurrentFrame.Func.PC = 0
		runner.Run()
	}
}
