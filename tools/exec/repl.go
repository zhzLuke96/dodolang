package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zhzluke96/dodolang/dodolang"
	"github.com/zhzluke96/dodolang/dolang"
)

func dodolangEPL() {
	fmt.Print("[dodolang code]\n")
	fmt.Print("Hit CTRL+C or type \"exit\" or \"quit\" to quit.")
	vm := dodolang.NewDodoVM([]string{})
	runner := dodolang.Runner{vm}
	for {
		fmt.Print("\n>>> ")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		dodolang.InputContent = string(data)
		if strings.Trim(string(data), " ") == "quit" {
			return
		}
		code := dodolang.GetTokenArr()
		runner.VM.CurrentFrame.Func.Code = labelLoad(code)
		runner.VM.CurrentFrame.Func.PC = 0
		runner.Run()
	}
}

func dolangREPL() {
	fmt.Print("[dolang code]\n")
	fmt.Print("Hit CTRL+C or type \"quit\" to quit.")
	vm := dodolang.NewDodoVM([]string{})
	runner := dodolang.Runner{vm}
	for {
		fmt.Print("\n>>> ")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		if strings.Trim(string(data), " ") == "quit" {
			return
		}

		dolangcode, err := dolang.ParseDolang(data)
		// fmt.Printf("[LOG] pcode = %v\n", pcode)
		if err != nil {
			log.Fatalln(err)
		}

		dodolang.InputContent = dolangcode
		code := dodolang.GetTokenArr()
		runner.VM.CurrentFrame.Func.Code = labelLoad(code)
		runner.VM.CurrentFrame.Func.PC = 0
		runner.Run()
	}
}
