package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	parseFile()
	// repl()
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func parseFile() {
	if len(os.Args) == 1 {
		fmt.Println("need input file")
		return
	}
	if code, err := ReadAll(os.Args[1]); err != nil {
		fmt.Println(err)
		return
	} else {
		Parse(code)
	}
}

func repl() {
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()

		Parse(data)
		fmt.Print("\n")
	}
}
