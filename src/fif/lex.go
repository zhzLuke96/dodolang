package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

//go:generate goyacc -l -o parser.go parser.y

func Parse(input []byte) error {
	l := newLex(input)
	_ = FifParse(l)
	return l.err
}

type lex struct {
	input   []byte
	pos     int
	linepos int
	line    int
	err     error
}

func newLex(input []byte) *lex {
	return &lex{
		input: input,
	}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *FifSymType) int {
	return l.scanNormal(lval)
}

var idRegex = regexp.MustCompile("[\\w\\$_]")

func (l *lex) scanNormal(lval *FifSymType) int {
	for b := l.next(); b != 0; b = l.next() {
		switch {
		case b == '\n':
			l.line++
			l.linepos = 0
			continue
		case unicode.IsSpace(rune(b)):
			continue
		case b == '"' || b == '\'':
			return l.scanString(lval, b)
		case unicode.IsDigit(rune(b)):
			l.backup()
			return l.scanNum(lval)
		// case unicode.IsLetter(rune(b)):
		case idRegex.Match([]byte{b}):
			l.backup()
			return l.scanIdentifier(lval)
		default:
			return int(b)
		}
	}
	return 0
}

var escape = map[byte]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func (l *lex) scanString(lval *FifSymType, match byte) int {
	buf := bytes.NewBuffer(nil)
	for b := l.next(); b != 0; b = l.next() {
		switch b {
		case '\\':
			// TODO(sougou): handle \uxxxx construct.
			b2 := escape[l.next()]
			if b2 == 0 {
				return LexError
			}
			buf.WriteByte(b2)
		case match:
			lval.val = buf.String()
			return StringConstant
		default:
			buf.WriteByte(b)
		}
	}
	return LexError
}

func (l *lex) scanNum(lval *FifSymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case unicode.IsDigit(rune(b)):
			buf.WriteByte(b)
		case strings.IndexByte(".+-eE", b) != -1:
			buf.WriteByte(b)
		default:
			l.backup()
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return LexError
			}
			lval.val = val
			return NumConstant
		}
	}
}

var reserved_words = map[string]int{
	"func":        FuncDefined,
	"return":      FuncReturn,
	"if":          T_IF,
	"else":        T_ELSE,
	"for":         T_FOR,
	"while":       T_WHILE,
	"__fifcode__": T_FIF,
}

func (l *lex) scanIdentifier(lval *FifSymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case idRegex.Match([]byte{b}):
			buf.WriteByte(b)
		// case unicode.IsSpace(rune(b)):
		default:
			l.backup()
			lval.val = buf.String()
			if ty, ok := reserved_words[lval.val.(string)]; ok {
				// fmt.Printf("\nty:[%v] [%v]\n", ty, lval.val.(string))
				return ty
			}
			return Identifier
		}
	}
}

func (l *lex) backup() {
	if l.pos == -1 {
		return
	}
	l.pos--
	l.linepos--
}

func (l *lex) next() byte {
	if l.pos >= len(l.input) || l.pos == -1 {
		l.pos = -1
		l.linepos = -1
		return 0
	}
	l.pos++
	l.linepos++
	return l.input[l.pos-1]
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
	fmt.Printf("\n[ERROR]:%v pos:[L:%v,P:%v] \n", s, l.line, l.linepos)
}
