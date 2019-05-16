package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// parseFile()
	// repl()
	stdInParse()
}

func stdInParse() {
	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadBytes('\000')
	if err != nil {
		log.Fatal(err)
	}
	Parse(code)
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

func fmtFloat64(f float64) string {
	s := strconv.FormatFloat(f, 'f', 10, 64)
	return strings.Trim(s, ".0")
}
