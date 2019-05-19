package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"./fif"
	"./machine"
)

func main() {
	// REPL()
	fifmod := flag.Bool("fif", false, "just parse to fifcode. not exec code.")
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
			fmt.Println(fif.ParseFifth(code))
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
	machine.InputContent = code
	codearr := machine.GetTokenArr()

	vm := machine.NewFifVM(labelLoad(codearr))
	runner := machine.Runner{vm}
	runner.Run()
}

func ExecFifthCode(code []byte) error {
	fifcode, err := fif.ParseFifth(code)
	// fmt.Printf("[LOG] pcode = %v\n", pcode)
	if err != nil {
		return err
	}
	RunfifCode(fifcode)
	return nil
}

// func ParseFifthCode(code string) (string, error) {
// 	cmd := exec.Command("./fif_parser.exe")
// 	cmd.Stdin = strings.NewReader(code + "\000")

// 	var out bytes.Buffer
// 	var serr bytes.Buffer
// 	cmd.Stdout = &out
// 	cmd.Stderr = &serr

// 	err := cmd.Run()
// 	if err != nil {
// 		return "", err
// 	}
// 	if serr.Len() != 0 {
// 		return out.String(), fmt.Errorf(serr.String())
// 	}
// 	return out.String(), nil
// }
