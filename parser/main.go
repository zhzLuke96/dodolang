package parser

import "bytes"

var ParserBuf bytes.Buffer

func Parse(code []byte) ([]byte, error) {
	err := StartParse(code)
	if err != nil {
		return nil, err
	}
	return ParserBuf.Bytes(), nil
}

func StartParse(input []byte) error {
	l := newLex(input)
	_ = DolangParse(l)
	return l.err
}
