package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"./machine"
)

var procRegex = regexp.MustCompile("^\\s*.+:.+ return\\s*$")

func ReplFifCode() {
	fmt.Println("Hit CTRL+C or type \"exit\" or \"quit\" to quit.")
	m := machine.NewMachine([]string{})
	for {
		fmt.Print("\n>>> ")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input := string(data)
		// fmt.Println("input: ", input)
		if input == "quit" {
			break
		}
		if procRegex.MatchString(input) {
			m.EvalProc(input)
		} else {
			m.Eval(input)
		}
	}
	return
}
