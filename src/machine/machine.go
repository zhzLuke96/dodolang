package machine

import (
	"fmt"
	"os"
	"strings"

	"../stack"
)

type Machine struct {
	code            []string
	dataStack       *stack.Stack
	returnAddrStack *stack.Stack
	programCounter  int
	labelMap        map[string]int
	dispatchMap     map[string]func(args ...string)
	scopedVars      map[string]interface{}
}

func NewMachine(code []string) *Machine {
	m := new(Machine)
	m.code = code
	m.dataStack = new(stack.Stack)
	m.returnAddrStack = new(stack.Stack)
	m.programCounter = 0
	m.labelMap = make(map[string]int)
	m.scopedVars = make(map[string]interface{})
	// m.dispatchMap = make(map[string]func())
	m.reload()
	return m
}

func (m *Machine) reload() {
	labels, clearCode := cutLabelInCode(m.code)
	m.code = clearCode
	for _, v := range labels {
		m.labelMap[v[0].(string)] = v[1].(int)
	}
	m.dispatchMap = map[string]func(args ...string){
		"%": func(args ...string) {
			m.push(m.pop().(int) % m.pop().(int))
		},
		"*": func(args ...string) {
			m.push(m.pop().(int) * m.pop().(int))
		},
		"+": func(args ...string) {
			m.push(m.pop().(int) + m.pop().(int))
		},
		"-": func(args ...string) {
			m.push(m.pop().(int) - m.pop().(int))
		},
		"/": func(args ...string) {
			m.push(m.pop().(int) / m.pop().(int))
		},
		"==": func(args ...string) {
			m.push(m.pop().(int) == m.pop().(int))
		},
		"num": func(args ...string) {
			m.push(str2num(m.pop()))
		},
		"bool": func(args ...string) {
			dataText := strings.ToUpper(m.pop().(string))
			m.push(dataText == "TRUE")
		},
		"drop": func(args ...string) {
			m.pop()
		},
		"dup": func(args ...string) {
			if len(args) == 0 {
				m.dup(0)
			} else {
				idx := str2num(args[0])
				m.dup(idx)
			}
		},
		"swap": func(args ...string) {
			var idx int
			top := m.pop().(int)
			if len(args) == 0 {
				idx = 1
				m.dup(1 - 1) // => dup_1
			} else {
				idx = str2num(args[0])
				if idx == 0 {
					return
				}
				m.dup(idx - 1)
			}
			(*m.dataStack)[m.dataStack.Len()-idx] = top
		},
		"read": func(args ...string) {
			var input string
			fmt.Scanln(&input)
			m.push(input)
		},
		"print": func(args ...string) {
			fmt.Print(m.pop())
		},
		"println": func(args ...string) {
			fmt.Println(m.pop())
		},
		"call": func(args ...string) {
			m.returnAddrStack.Push(m.programCounter)
			m.jump()
		},
		"return": func(args ...string) {
			addr := m.pop().(int)
			m.programCounter = addr
		},
		"exit": func(args ...string) {
			if len(args) == 0 {
				os.Exit(0)
			} else {
				retCode := str2num(args[0])
				os.Exit(retCode)
			}
		},
		"load": func(args ...string) {
			if len(args) != 0 {
				key := args[0]
				if v, ok := m.scopedVars[key]; ok {
					m.push(v)
				} else {
					m.push("0")
				}
			} else {
				// ERROR CALL
			}
		},
		"store": func(args ...string) {
			if len(args) == 0 {
				key := args[0]
				m.scopedVars[key] = m.pop()
			} else {
				// ERROR CALL
			}
		},
		"jump": func(args ...string) {
			m.jump()
		},
	}
}

func (m *Machine) dup(idx int) {
	m.push(m.top(idx))
}

func (m *Machine) jump() interface{} {
	jumpLabel := m.pop().(string)
	if isStringType(jumpLabel) {
		if addr, ok := m.labelMap[jumpLabel]; ok {
			m.programCounter = addr
		}
	}
	// Error("JMP address must be a valid label")
	return nil
}

func (m *Machine) pop() interface{} {
	if v, err := m.dataStack.Pop(); err == nil {
		return v
	} else {
		// ERROR
		fmt.Println(err)
		return 0
	}
}

func (m *Machine) push(v interface{}) {
	m.dataStack.Push(v)
}

func (m *Machine) top(idx int) interface{} {
	if idx < 0 || idx >= m.dataStack.Len() {
		return "null"
	}
	return (*m.dataStack)[m.dataStack.Len()-1-idx]
}

func (m *Machine) Run() *Machine {
	for {
		if m.programCounter < len(m.code) {
			opt := m.code[m.programCounter]
			m.programCounter++
			m.dispatch(opt)
		} else {
			break
		}
	}
	return m
}

func (m *Machine) dispatch(opt string) {
	tokenType, arg := GetTokenTypeName(opt)
	switch tokenType {
	case "Operator":
		fallthrough
	case "Instruction":
		// m.dispatchMap[opt]()
		fallthrough
	case "Instruction_Args":
		m.dispatchMap[opt](arg)
	case "Number":
		m.push(opt)
	case "String":
		m.push(opt[1 : len(opt)-1])
	default:
		// error
		fmt.Printf("[ERROR] UNKNOW TOKEN: %s\n", opt)
	}
}
