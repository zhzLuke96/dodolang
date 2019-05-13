package machine

import "strconv"

func cutFuncBody(code []string) []string {
	pos := 0
	count := 0
	ret := []string{}
	for {
		if pos == len(code) {
			break
		}
		v := code[pos]
		if v == "func" {
			count++
		} else if v == "endfunc" {
			count--
			if count == 0 {
				break
			}
		}
		ret = append(ret, v)
		pos++
	}
	return ret[1:]
}

// "&end jmp 'Hello world' println end: exit"
// "4 jmp 'Hello world' println exit"
func labelLoad(text []string) []string {
	labelMap := make(map[string]int)
	ret := []string{}
	pos := 0
	for {
		if pos == len(text) {
			break
		}
		v := text[pos]
		if v[len(v)-1] == ':' {
			labelMap[v[0:len(v)-1]] = len(ret) - 1
		} else if v == "func" {
			funcbody := cutFuncBody(text[pos:])
			labeled := labelLoad(funcbody)
			ret = append(ret, "func")
			ret = append(ret, labeled...)
			ret = append(ret, "endfunc")
			pos += len(funcbody) + 1
		} else {
			ret = append(ret, v)
		}
		pos++
	}
	for i, v := range ret {
		if v[0] == '&' {
			if v, ok := labelMap[v[1:]]; ok {
				ret[i] = strconv.Itoa(v)
			} else {
				ret[i] = strconv.Itoa(i + 1)
			}
		}
	}
	return ret
}
