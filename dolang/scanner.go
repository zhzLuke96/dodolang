package dodolang

var (
	InputContent = ""
	scannerPos   = 0
	EOF          = "\000"
)

func nxtCh() string {
	if scannerPos == len(InputContent) {
		return EOF
	}
	scannerPos++
	return string(InputContent[scannerPos-1])
}

func escape(ch string) string {
	switch ch {
	case "n":
		return "\n"
	case "t":
		return "\t"
	case "r":
		return "\r"
	default:
		return ch
	}
}

func Scan() string {
	NxtToken := ""
	stringMod := false
	stringMatch := ""

	for {
		ch := nxtCh()
		if ch == EOF {
			break
		}
		if !stringMod && (ch == " " || ch == "\n" || ch == "\r" || ch == "\t") {
			if len(NxtToken) == 0 {
				continue
			}
			break
		} else if ch == "\\" {
			NxtToken += escape(nxtCh())
		} else if !stringMod && (ch == "'" || ch == "\"") {
			stringMatch = ch
			stringMod = true
			NxtToken += ch
		} else if stringMod && ch == stringMatch {
			stringMod = false
			NxtToken += ch
		} else {
			NxtToken += ch
		}
	}
	return NxtToken
}

func GetTokenArr() []string {
	ret := []string{}
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
