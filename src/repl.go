package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"./fif"
	"./machine"
)

func fifREPL() {
	fmt.Print("[fif code]\n")
	fmt.Print("Hit CTRL+C or type \"exit\" or \"quit\" to quit.")
	vm := machine.NewFifVM([]string{})
	runner := machine.Runner{vm}
	for {
		fmt.Print("\n>>> ")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		machine.InputContent = string(data)
		if strings.Trim(string(data), " ") == "quit" {
			return
		}
		code := machine.GetTokenArr()
		runner.VM.CurrentFrame.Func.Code = labelLoad(code)
		runner.VM.CurrentFrame.Func.PC = 0
		runner.Run()
	}
}

func fifthREPL() {
	fmt.Print("[fifth code]\n")
	fmt.Print("Hit CTRL+C or type \"quit\" to quit.")
	vm := machine.NewFifVM([]string{})
	runner := machine.Runner{vm}
	for {
		fmt.Print("\n>>> ")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		if strings.Trim(string(data), " ") == "quit" {
			return
		}

		fifcode, err := fif.ParseFifth(data)
		// fmt.Printf("[LOG] pcode = %v\n", pcode)
		if err != nil {
			log.Fatalln(err)
		}

		machine.InputContent = fifcode
		code := machine.GetTokenArr()
		runner.VM.CurrentFrame.Func.Code = labelLoad(code)
		runner.VM.CurrentFrame.Func.PC = 0
		runner.Run()
	}
}
