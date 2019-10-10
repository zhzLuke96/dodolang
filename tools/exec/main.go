package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/zhzluke96/dodolang/dodolang"
	"github.com/zhzluke96/dodolang/dolang"
)

func main() {
	// REPL()
	fifmod := flag.Bool("do", false, "just parse to fifcode. not exec code.")
	debugmod := flag.Bool("debug", false, "run one line fifthcode.")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		if *fifmod {
			fifREPL()
		} else {
			fifthREPL()
		}
	} else {
		var code []byte
		var err error
		if *debugmod {
			code = []byte(args[0])
		} else {
			code, err = ReadAll(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}
		// clear comment
		code = clearComment(code)
		if *fifmod {
			fmt.Println(dolang.ParseFifth(code))
		} else {
			err = ExecFifthCode(code)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func RunfifCode(code string) {
	dodolang.InputContent = code
	codearr := dodolang.GetTokenArr()

	vm := dodolang.NewFifVM(labelLoad(codearr))
	runner := dodolang.Runner{vm}
	runner.Run()
}

func ExecFifthCode(code []byte) error {
	fifcode, err := dolang.ParseFifth(code)
	// fmt.Printf("[LOG] pcode = %v\n", pcode)
	if err != nil {
		return err
	}
	RunfifCode(fifcode)
	return nil
}
