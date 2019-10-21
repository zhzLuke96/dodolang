package main

import "io/ioutil"

func execDoByteCode(fn string) error {
	// [TODO]
	return nil
}

func ececDo(fn string) error {
	fbs, err := loadFile(fn)
	if err != nil {
		return err
	}
	vm := dolang.NewVm(fbs)
	vm.Run()
}

func execDoDo(fn string) error {
	fbs, err := loadFile(fn)
	if err != nil {
		return err
	}
	docode, err := parser.Parse(fbs)
	if err != nil {
		return err
	}
	vm := dolang.NewVm(docode)
	vm.Run()
}

func loadFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}
