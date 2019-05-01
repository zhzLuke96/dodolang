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
	labelMap        map[string]*label_body
	dispatchMap     map[string]func(args ...string)
	scopedVars      map[string]interface{}
}

func NewMachine(code []string) *Machine {
	m := new(Machine)
	m.code = code
	m.dataStack = new(stack.Stack)
	m.returnAddrStack = new(stack.Stack)
	m.programCounter = 0
	// m.labelMap = make(map[string]int)
	// m.scopedVars = make(map[string]interface{})
	// m.dispatchMap = make(map[string]func())
	m.reload()
	return m
}

func (m *Machine) reload() {
	labels, clearCode := cutLabelInCode(m.code)
	m.code = clearCode
	m.labelMap = labels
	m.dispatchMap = map[string]func(args ...string){
		"null": func(args ...string) {
			return
		},
		"mul": func(args ...string) {
			m.Push(num2float(m.Pop()) * num2float(m.Pop()))
		},
		"plus": func(args ...string) {
			m.Push(num2float(m.Pop()) + num2float(m.Pop()))
		},
		"equal": func(args ...string) {
			m.Push(num2float(m.Pop()) == num2float(m.Pop()))
		},
		"mod": func(args ...string) {
			numb := num2float(m.Pop())
			numa := num2float(m.Pop())
			m.Push(int(numa) % int(numb))
		},
		"sub": func(args ...string) {
			numb := num2float(m.Pop())
			numa := num2float(m.Pop())
			m.Push(numa - numb)
		},
		"div": func(args ...string) {
			numb := num2float(m.Pop())
			numa := num2float(m.Pop())
			m.Push(numa / numb)
		},
		"greater": func(args ...string) {
			numb := num2float(m.Pop())
			numa := num2float(m.Pop())
			m.Push(numa > numb)
		},
		"less": func(args ...string) {
			numb := num2float(m.Pop())
			numa := num2float(m.Pop())
			m.Push(numa < numb)
		},
		"num": func(args ...string) {
			m.Push(str2num(m.Pop()))
		},
		"int": func(args ...string) {
			m.Push(str2int(m.Pop()))
		},
		"float": func(args ...string) {
			m.Push(str2float(m.Pop()))
		},
		"bool": func(args ...string) {
			dataText := strings.ToUpper(m.Pop().(string))
			m.Push(dataText == "TRUE")
		},
		"drop": func(args ...string) {
			m.Pop()
		},
		"dup": func(args ...string) {
			if len(args) == 0 {
				m.Dup(0)
			} else {
				idx := str2int(args[0])
				m.Dup(idx)
			}
		},
		"swap": func(args ...string) {
			var idx int
			top := m.Pop().(int)
			if len(args) == 0 {
				idx = 1
				m.Dup(1 - 1) // => dup_1
			} else {
				idx = str2int(args[0])
				if idx == 0 {
					return
				}
				m.Dup(idx - 1)
			}
			(*m.dataStack)[m.dataStack.Len()-idx] = top
		},
		"if": func(args ...string) {
			falseStatement := m.Pop().(string)
			trueStatement := m.Pop().(string)
			condition := m.Pop().(bool)
			if condition {
				m.Push(trueStatement)
			} else {
				m.Push(falseStatement)
			}
			m.Jump()
		},
		"read": func(args ...string) {
			var input string
			fmt.Scanln(&input)
			m.Push(input)
		},
		"print": func(args ...string) {
			fmt.Print(m.Pop())
		},
		"println": func(args ...string) {
			fmt.Println(m.Pop())
		},
		"call": func(args ...string) {
			m.returnAddrStack.Push(m.programCounter)
			m.Jump()
		},
		"return": func(args ...string) {
			if addr, err := m.returnAddrStack.Pop(); err != nil {
				fmt.Printf("[ERROR] return addr error")
			} else {
				m.programCounter = addr.(int)
			}
		},
		"exit": func(args ...string) {
			if len(args) == 0 {
				os.Exit(0)
			} else {
				retCode := str2int(args[0])
				os.Exit(retCode)
			}
		},
		"load": func(args ...string) {
			// if len(args) != 0 {
			// 	key := args[0]
			// 	if v, ok := m.scopedVars[key]; ok {
			// 		m.Push(v)
			// 	} else {
			// 		m.Push("0")
			// 	}
			// } else {
			// 	// ERROR CALL
			// }
			labelKey := m.Pop().(string)

			m.Push(m.labelMap[labelKey].Value)
		},
		"store": func(args ...string) {
			// if len(args) == 0 {
			// 	key := args[0]
			// 	m.scopedVars[key] = m.Pop()
			// } else {
			// 	// ERROR CALL
			// }
			labelKey := m.Pop().(string)
			m.labelMap[labelKey].Value = m.Pop()
		},
		"jump": func(args ...string) {
			m.Jump()
		},
	}
}

func (m *Machine) Dup(idx int) {
	m.Push(m.Top(idx))
}

func (m *Machine) Jump() {
	if jumpLabel, ok := m.Pop().(string); ok {
		// if isStringType(jumpLabel) {
		jumpLabel = strings.ToLower(jumpLabel)
		if addrLabel, ok := m.labelMap[jumpLabel]; ok {
			m.programCounter = addrLabel.Idx
		}
		// }
	}
	// Error("JMP address must be a valid label")
	// return nil
}

func (m *Machine) Pop() interface{} {
	if v, err := m.dataStack.Pop(); err == nil {
		return v
	} else {
		// ERROR
		fmt.Println(err)
	}
	return 0
}

func (m *Machine) Push(v interface{}) {
	m.dataStack.Push(v)
}

func (m *Machine) Top(idx int) interface{} {
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
	opt = strings.ToLower(opt)
	tokenType, arg := GetTokenTypeName(opt)
	switch tokenType {
	case "Operator":
		fallthrough
	case "Instruction":
		// m.dispatchMap[opt]()
		fallthrough
	case "Instruction_Args":
		m.dispatchMap[opt](arg)
	case "Label_Pointer":
		m.Push(opt[1:])
	case "Number":
		m.Push(opt)
	case "String":
		m.Push(opt[1 : len(opt)-1])
	default:
		// error
		fmt.Printf("[ERROR] UNKNOW TOKEN: %s\n", opt)
	}
}
