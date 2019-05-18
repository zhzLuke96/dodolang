package main

import "regexp"

var lineComment = regexp.MustCompile("\\/\\/[^\\n]*\\n")
var mulitLineComment = regexp.MustCompile("\\/\\*[\\s\\S]*\\*\\/")

func clearComment(code []byte) []byte {
	ret := lineComment.ReplaceAll(code, []byte(""))
	ret = mulitLineComment.ReplaceAll(ret, []byte(""))
	return ret
}
