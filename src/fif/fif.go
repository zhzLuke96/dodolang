package fif

import "bytes"
import "regexp"

var ParserBuf bytes.Buffer

func ParseFifth(fifthCode []byte) (string, error) {
	err := Parse(fifthCode)
	if err != nil {
		return "", err
	}
	return ParserBuf.String(), nil
}

var argRegex = regexp.MustCompile(" arg")

func argCount(t string) int {
	return len(argRegex.FindAllString(t, -1))
}
