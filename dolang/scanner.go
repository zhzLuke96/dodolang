package dolang

var (
	InputContent = []byte{}
	scannerPos   = 0
	EOF          = byte('\000')
)

func nxtCh() byte {
	if scannerPos == len(InputContent) {
		return EOF
	}
	scannerPos++
	return InputContent[scannerPos-1]
}

func escape(ch byte) byte {
	switch ch {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'r':
		return '\r'
	default:
		return ch
	}
}

func Scan() []byte {
	NxtToken := []byte{}
	stringMod := false
	stringMatch := byte(0)

	for {
		ch := nxtCh()
		if ch == EOF {
			break
		}
		if !stringMod && (ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t') {
			if len(NxtToken) == 0 {
				continue
			}
			break
		} else if ch == '\\' {
			NxtToken = append(NxtToken, escape(nxtCh()))
		} else if !stringMod && (ch == '\'' || ch == '"') {
			stringMatch = ch
			stringMod = true
			NxtToken = append(NxtToken, ch)
		} else if stringMod && ch == stringMatch {
			stringMod = false
			NxtToken = append(NxtToken, ch)
		} else {
			NxtToken = append(NxtToken, ch)
		}
	}
	return NxtToken
}

func GetTokenArr() [][]byte {
	ret := [][]byte{}
	scannerPos = 0
	for {
		t := Scan()
		if len(t) == 0 {
			break
		}
		ret = append(ret, t)
	}
	return ret
}

func Tokenizer(code []byte) [][]byte {
	InputContent = code
	return GetTokenArr()
}
