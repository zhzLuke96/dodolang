package bignum

import (
	"strconv"
	"testing"
)

func TestBigIntString(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected string
	}{
		{BigInt{10}, "10"},
		{BigInt{10, 10}, "1010"},
		{BigInt{0, 10}, "1000"},
		{BigInt{0}, "0"},
		{BigInt{0, 1}, "100"},
		{BigInt{77, 88, 99}, "998877"},
		{BigInt{99, 88, 77}, "778899"},
		{BigInt{}, ""},
	}
	for _, test := range tests {
		actual := test.bigint.String()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] String() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
	}
}

func bigintEquleInt(b *BigInt, n int) bool {
	return b.String() == strconv.Itoa(n)
}

func TestBigIntToInt(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected int
	}{
		{BigInt{10}, 10},
	}
	for _, test := range tests {
		actual := test.bigint.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] ToInt() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
	}
}

func TestEvalInt(t *testing.T) {
	tests := []struct {
		expr     string
		expected int
	}{
		{"1", 1},
		{"11", 11},
		{"111", 111},
		{"1111", 1111},
		{"0", 0},
		{"00", 0},
		{"000", 0},
	}
	for _, test := range tests {
		actual := EvalInt(test.expr).ToInt()
		if actual != test.expected {
			t.Fatalf("EvalInt [%v] need %v, but %v", test.expr, test.expected, actual)
		}
	}
}

func TestEvalDec(t *testing.T) {
	tests := []struct {
		expr     string
		expected string
	}{
		{"1", "10"},
		{"10", "10"},
		{"100", "1000"},
		{"12", "12"},
		{"123", "1230"},
	}
	for _, test := range tests {
		actual := EvalDec(test.expr).String()
		if actual != test.expected {
			t.Fatalf("EvalDec [%v] need %v, but %v", test.expr, test.expected, actual)
		}
	}
}

func TestIntToBigInt(t *testing.T) {
	tests := []struct {
		n        int64
		expected string
	}{
		{1, "1"},
		{0, "0"},
		{100, "100"},
		{10000, "10000"},
		{1000000, "1000000"},
	}
	for _, test := range tests {
		num := intToBigInt(test.n)
		actual := num.String()
		if actual != test.expected {
			t.Fatalf("intToBigInt [%v] need %v, but %v", test.n, test.expected, actual)
		}
	}
}

func TestBigIntLen(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected int
	}{
		{BigInt{10}, 1},
		{BigInt{10, 1}, 2},
		{BigInt{0, 1}, 2},
		{BigInt{0, 1, 0}, 2},
		{BigInt{0, 1, 0, 0}, 2},
		{BigInt{0, 0, 0, 0}, 1},
		{BigInt{0}, 1},
	}
	for _, test := range tests {
		actual := test.bigint.Len()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Len() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
	}
}

func TestBigIntIsZero(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected bool
	}{
		{BigInt{}, true},
		{BigInt{0}, true},
		{BigInt{0, 0}, true},
		{BigInt{0, 0, 0}, true},
		{BigInt{1}, false},
	}
	for _, test := range tests {
		actual := test.bigint.IsZero()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] IsZero() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		b1       BigInt
		b2       BigInt
		expected bool
	}{
		{BigInt{1}, BigInt{0}, false},
		{BigInt{0}, BigInt{1}, true},
		{BigInt{0, 0}, BigInt{1}, true},
		{BigInt{}, BigInt{1}, true},
		{BigInt{0, 0, 1}, BigInt{1}, false},
	}
	for _, test := range tests {
		actual := test.b1.Less(&test.b2)
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Less(%v) need %v, but %v", test.b1.String(), test.b2.String(), test.expected, actual)
		}
	}
}

func TestGreat(t *testing.T) {
	tests := []struct {
		b1       BigInt
		b2       BigInt
		expected bool
	}{
		{BigInt{1}, BigInt{0}, true},
		{BigInt{0}, BigInt{1}, false},
		{BigInt{0, 0}, BigInt{1}, false},
		{BigInt{}, BigInt{1}, false},
		{BigInt{0, 0, 1}, BigInt{1}, true},
	}
	for _, test := range tests {
		actual := test.b1.Great(&test.b2)
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Great(%v) need %v, but %v", test.b1.String(), test.b2.String(), test.expected, actual)
		}
	}
}

func TestEqule(t *testing.T) {
	tests := []struct {
		b1       BigInt
		b2       BigInt
		expected bool
	}{
		{BigInt{0}, BigInt{0}, true},
		{BigInt{0}, BigInt{1}, false},
		{BigInt{}, BigInt{0}, false},
		{BigInt{}, BigInt{}, true},
		{BigInt{0, 0}, BigInt{0}, true},
		{BigInt{0, 1}, BigInt{0}, false},
		{BigInt{0, 1}, BigInt{1}, false},
	}
	for _, test := range tests {
		actual := test.b1.Equle(&test.b2)
		if actual != test.expected {
			t.Fatalf("Equle [%v] [%v] need %v, but %v", test.b1.String(), test.b2.String(), test.expected, actual)
		}
	}
}

func TestBigIntSucc(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected int
	}{
		{BigInt{0}, 1},
		{BigInt{1}, 2},
		{BigInt{99}, 100},
		{BigInt{99, 9}, 1000},
		{BigInt{99, 99}, 10000},
		{BigInt{}, 1},
	}
	for _, test := range tests {
		num := test.bigint.Clone()
		num.Succ()
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Succ() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
	}
}

func TestBigIntPred(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected int
	}{
		{BigInt{0}, 0},
		{BigInt{1}, 0},
		{BigInt{99}, 98},
		{BigInt{99, 9}, 998},
		{BigInt{}, 0},
		{BigInt{0, 1}, 99},
		{BigInt{0, 10}, 999},
		{BigInt{0, 0, 1}, 9999},
	}
	for _, test := range tests {
		num := test.bigint.Clone()
		num.Pred()
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Fail() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
	}
}

func TestBigIntPlus(t *testing.T) {
	tests := []struct {
		b1       BigInt
		b2       BigInt
		expected int
	}{
		{BigInt{0}, BigInt{0}, 0},
		{BigInt{1}, BigInt{0}, 1},
		{BigInt{0}, BigInt{1}, 1},
		{BigInt{99}, BigInt{1}, 100},
		{BigInt{99, 99}, BigInt{1}, 10000},
		{BigInt{1}, BigInt{99, 99}, 10000},
		{BigInt{99, 99}, BigInt{2}, 10001},
		{BigInt{99, 99}, BigInt{0}, 9999},
		{BigInt{99, 99}, BigInt{0, 0}, 9999},
	}
	for _, test := range tests {
		num := test.b1.Clone()
		num.Plus(&test.b2)
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Plus(%v) need %v, but %v", test.b1.String(), test.b2.String(), test.expected, actual)
		}
	}
}

func TestBigIntMinus(t *testing.T) {
	tests := []struct {
		b1       BigInt
		b2       BigInt
		expected int
	}{
		{BigInt{1}, BigInt{0}, 1},
		{BigInt{1}, BigInt{1}, 0},
		{BigInt{0, 1}, BigInt{1}, 99},
		{BigInt{0, 0, 1}, BigInt{1}, 9999},
		{BigInt{1}, BigInt{2}, 1},
		{BigInt{0, 1}, BigInt{2}, 98},
		{BigInt{0, 0, 1}, BigInt{2}, 9998},
	}
	for _, test := range tests {
		num := test.b1.Clone()
		num.Minus(&test.b2)
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Minus(%v) need %v, but %v", test.b1.String(), test.b2.String(), test.expected, actual)
		}
	}
}

func TestBigIntMulByte(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		by       int
		expected int
	}{
		{BigInt{1}, 2, 2},
		{BigInt{1}, 20, 20},
		{BigInt{0}, 20, 0},
		{BigInt{20}, 0, 0},
		{BigInt{5}, 4, 20},
		{BigInt{5}, 40, 200},
		{BigInt{50}, 40, 2000},
	}
	for _, test := range tests {
		actual := test.bigint.MulByte(byte(test.by)).ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] MulByte(%v) need %v, but %v", test.bigint.String(), test.by, test.expected, actual)
		}
	}
}

func TestBigIntLeftShift(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected int
	}{
		{BigInt{0}, 0},
		{BigInt{1}, 0},
		{BigInt{12}, 1},
		{BigInt{23, 1}, 12},
		{BigInt{34, 12}, 123},
		{BigInt{56, 34, 12}, 12345},
	}
	for _, test := range tests {
		num := test.bigint.Clone()
		num.LeftShift()
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] LeftShift() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
	}
}

func TestBigIntRightShift(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected int
	}{
		{BigInt{0}, 0},
		{BigInt{1}, 10},
		{BigInt{12}, 120},
		{BigInt{23, 1}, 1230},
		{BigInt{34, 12}, 12340},
		{BigInt{56, 34, 12}, 1234560},
	}
	for _, test := range tests {
		num := test.bigint.Clone()
		num.RightShift()
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] RightShift() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
	}
}

func TestBigIntMul(t *testing.T) {
	tests := []struct {
		b1       BigInt
		b2       BigInt
		expected int
	}{
		{BigInt{10}, BigInt{0}, 0},
		{BigInt{0}, BigInt{10}, 0},
		{BigInt{10}, BigInt{10}, 100},
		{BigInt{0, 1}, BigInt{0, 1}, 10000},
		{BigInt{0, 10}, BigInt{0, 10}, 1000000},
		{BigInt{0, 0, 1}, BigInt{0, 0, 1}, 100000000},
	}
	for _, test := range tests {
		num := test.b1.Clone()
		num.Mul(&test.b2)
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Mul(%v) need %v, but %v", test.b1.String(), test.b2.String(), test.expected, actual)
		}
	}
}

func TestBigIntDiv(t *testing.T) {
	tests := []struct {
		b1       BigInt
		b2       BigInt
		expected int
	}{
		{BigInt{0}, BigInt{12}, 0},
		{BigInt{2}, BigInt{2}, 1},
		{BigInt{10}, BigInt{2}, 5},
		{BigInt{0, 1}, BigInt{2}, 50},
		{BigInt{0, 10}, BigInt{2}, 500},
		{BigInt{55}, BigInt{2}, 27},
	}
	for _, test := range tests {
		num := test.b1.Clone()
		num.Div(&test.b2)
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Div(%v) need %v, but %v", test.b1.String(), test.b2.String(), test.expected, actual)
		}
	}
}

func TestBigIntMod(t *testing.T) {
	tests := []struct {
		b1       BigInt
		b2       BigInt
		expected int
	}{
		{BigInt{0}, BigInt{12}, 0},
		{BigInt{2}, BigInt{2}, 0},
		{BigInt{10}, BigInt{2}, 0},
		{BigInt{11}, BigInt{2}, 1},
		{BigInt{0, 1}, BigInt{7}, 2},
		{BigInt{54}, BigInt{99}, 54},
		{BigInt{0, 1}, BigInt{99}, 1},
	}
	for _, test := range tests {
		num := test.b1.Clone()
		num.Mod(&test.b2)
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Mod(%v) need %v, but %v", test.b1.String(), test.b2.String(), test.expected, actual)
		}
	}
}

func TestBigIntSqrt(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected int
		d        int
	}{
		{BigInt{1}, 1, 0},
		{BigInt{2}, 1, 1},
		{BigInt{4}, 2, 0},
		{BigInt{64}, 8, 0},
		{BigInt{24, 10}, 32, 0},
		{BigInt{13}, 3, 4},
		{BigInt{0, 13}, 36, 4},
		{BigInt{0, 0, 13}, 360, 400},
	}
	for _, test := range tests {
		num := test.bigint.Clone()
		alD := num.Sqrt()
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigInt [%v] Sqrt() need %v, but %v", test.bigint.String(), test.expected, actual)
		}
		if alD.ToInt() != test.d {
			t.Fatalf("BigInt [%v] Sqrt() => need %v, but %v", test.bigint.String(), test.d, alD)
		}
	}
}

func TestBigIntToBin(t *testing.T) {
	tests := []struct {
		bigint   BigInt
		expected string
	}{
		{BigInt{0}, "0"},
		{BigInt{1}, "1"},
		{BigInt{2}, "10"},
		{BigInt{3}, "11"},
		{BigInt{12}, "1100"},
		{BigInt{34, 12}, "10011010010‬"},
	}
	for _, test := range tests {
		bins := test.bigint.ToBin()
		actual := ""
		flag := false
		for i, v := range bins {
			if v && test.expected[i] != '1' {
				flag = true
				break
			}
			if !v && test.expected[i] != '0' {
				flag = true
				break
			}
		}
		if flag {
			t.Fatalf("BigInt [%v] String() \nneed\t%v\nbut \t%v", test.bigint.String(), test.expected, actual)
		}
	}
}
