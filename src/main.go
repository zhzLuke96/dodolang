package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"./machine"
)

func main() {
	// REPL()
	if len(os.Args) == 1 {
		log.Fatal("need code file")
	}
	codefile, err := ReadAll(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	err = ExecFifthCode(string(codefile))
	if err != nil {
		log.Fatal(err)
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

func ExecFifthCode(code string) error {
	pcode, err := ParseFifthCode(code)
	// fmt.Printf("[LOG] pcode = %v\n", pcode)
	if err != nil {
		return err
	}
	RunfifCode(pcode)
	return nil
}

func ParseFifthCode(code string) (string, error) {
	cmd := exec.Command("./fif_parser.exe")
	cmd.Stdin = strings.NewReader(code + "\000")

	var out bytes.Buffer
	var serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	if serr.Len() != 0 {
		return out.String(), fmt.Errorf(serr.String())
	}
	return out.String(), nil
}
