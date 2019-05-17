package fif

import "bytes"

var ParserBuf bytes.Buffer

func ParseFifth(fifthCode []byte) (string, error) {
	err := Parse(fifthCode)
	if err != nil {
		return "", err
	}
	return ParserBuf.String(), nil
}
