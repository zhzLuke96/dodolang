package main

import (
	"bufio"
	"fmt"
	"os"

	"./machine"
)

func REPL() {
	fmt.Print("Hit CTRL+C or type \"exit\" or \"quit\" to quit.")
	vm := machine.NewFifVM([]string{})
	runner := machine.Runner{vm}
	for {
		fmt.Print("\n>>> ")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		machine.InputContent = string(data)
		code := machine.GetTokenArr()
		runner.VM.CurrentFrame.Func.Code = code
		runner.VM.CurrentFrame.Func.PC = 0
		runner.Run()
	}
}
