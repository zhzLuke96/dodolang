package main

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
