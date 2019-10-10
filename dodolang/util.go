package dolang

import (
	"strconv"
	"strings"
)

func fmtFloat64(f float64) string {
	s := strconv.FormatFloat(f, 'f', 10, 64)
	return strings.Trim(s, ".0")
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

type LabelStack struct {
	Top   int
	Count int
	S     [100]int
}

func NewLabelStack() *LabelStack {
	return &LabelStack{
		Top:   0,
		Count: 0,
		S:     [100]int{},
	}
}

func (l *LabelStack) Topv() int {
	return l.S[l.Top]
}

func (l *LabelStack) BEG() int {
	l.Top++
	l.Count++
	l.S[l.Top] = l.Count
	return l.Topv()
}

func (l *LabelStack) END() {
	l.Top--
}
