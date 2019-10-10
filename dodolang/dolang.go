package dolang

import "bytes"
import "regexp"

var ParserBuf bytes.Buffer

func ParseDolang(doCode []byte) ([]byte, error) {
	err := Parse(doCode)
	if err != nil {
		return nil, err
	}
	return ParserBuf.Bytes(), nil
}

var argRegex = regexp.MustCompile(" arg")

func argCount(t string) int {
	return len(argRegex.FindAllString(t, -1))
}
