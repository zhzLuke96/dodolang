package bignum

import (
	"math"
	"testing"
)

const (
	testPrecision = 0.005
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr     string
		expected BigNum
	}{
		{"1.02", BigNum{1, 0, 2, 1}},
		{"1.2", BigNum{1, 0, 20, 1}},
		{"1.20", BigNum{1, 0, 20, 1}},
		{"0.20", BigNum{1, 0, 20, 0}},
		{"0", BigNum{0, 0, 0}},
		{"1", BigNum{0, 0, 1}},
		{"10000", BigNum{0, 0, 0, 0, 1}},
		{"-1", BigNum{128, 0, 1}},
		{"-1", BigNum{128, 0, 1}},
		{"+1", BigNum{0, 0, 1}},
		{"+1.1", BigNum{1, 0, 10, 1}},
		{"-1.1", BigNum{129, 0, 10, 1}},
	}
	for _, test := range tests {
		actual := Eval(test.expr)
		if !actual.Equle(&test.expected) {
			t.Fatalf("Eval(%v) need %v, but %v", test.expr, test.expected.String(), actual.String())
		}
	}
}

func TestBigNumIsZero(t *testing.T) {
	tests := []struct {
		num      *BigNum
		expected bool
	}{
		{Eval("0"), true},
		{Eval("0.0"), true},
		{Eval("00.00"), true},
		{Eval("1"), false},
		{Eval("10"), false},
		{Eval("10000"), false},
		{Eval("1.000"), false},
		{Eval("0.5"), false},
	}
	for _, test := range tests {
		actual := test.num.IsZero()
		if actual != test.expected {
			t.Fatalf("BigNum [%v] IsZero() need %v, but %v", test.num.String(), test.expected, actual)
		}
	}
}

func TestBigNumLess(t *testing.T) {
	tests := []struct {
		num1     *BigNum
		num2     *BigNum
		expected bool
	}{
		{Eval("0"), Eval("0"), false},
		{Eval("1"), Eval("0"), false},
		{Eval("0"), Eval("1"), true},
		{Eval("1000"), Eval("10000"), true},
		{Eval("9999"), Eval("100000"), true},
		{Eval("1.1"), Eval("100000"), true},
		{Eval("9999"), Eval("100000"), true},
		{Eval("0.9"), Eval("1"), true},
		{Eval("0.000009"), Eval("0.001"), true},
		{Eval("-1"), Eval("0"), true},
		{Eval("0"), Eval("-1"), false},
		{Eval("-1.1"), Eval("-1"), true},
		{Eval("-1"), Eval("-2"), false},
	}
	for _, test := range tests {
		actual := test.num1.Less(test.num2)
		if actual != test.expected {
			t.Fatalf("BigNum [%v] Less(%v) need %v, but %v", test.num1.String(), test.num2.String(), test.expected, actual)
		}
	}
}

func TestBigNumPlus(t *testing.T) {
	tests := []struct {
		num1     *BigNum
		num2     *BigNum
		expected float64
	}{
		{Eval("0"), Eval("0"), 0},
		{Eval("1"), Eval("0"), 1},
		{Eval("0"), Eval("1"), 1},
		{Eval("1000"), Eval("1000"), 2000},
		{Eval("9999"), Eval("1"), 10000},
		{Eval("1.1"), Eval("10"), 11.1},
		{Eval("0.9"), Eval("1"), 1.9},
		{Eval("0.0009"), Eval("0.1"), 0.1009},
		{Eval("-1"), Eval("0"), -1},
		{Eval("0"), Eval("-1"), -1},
		{Eval("-1.1"), Eval("-1"), -2.1},
		{Eval("-1.1"), Eval("1"), -0.1},
		{Eval("-1"), Eval("-2"), -3},
	}
	for _, test := range tests {
		num := test.num1.Clone()
		num.Plus(test.num2)
		actual := num.ToFloat()
		if !(actual == test.expected || math.Abs(actual-test.expected) <= test.expected*testPrecision) {
			t.Fatalf("BigNum [%v] Plus(%v) need %v, but %v", test.num1.String(), test.num2.String(), test.expected, actual)
		}
	}
}

func TestBigNumDiv(t *testing.T) {
	tests := []struct {
		num1     *BigNum
		num2     *BigNum
		expected float64
	}{
		{Eval("1"), Eval("1"), 1},
		{Eval("10"), Eval("10"), 1},
		{Eval("100"), Eval("100"), 1},
		{Eval("1000"), Eval("1000"), 1},
		{Eval("10"), Eval("2"), 5},
		{Eval("12"), Eval("2"), 6},
		{Eval("1"), Eval("2"), 0.5},
		{Eval("2"), Eval("3"), 0.6666},
		{Eval("1024"), Eval("64"), 1024 / 64},
		{Eval("1024"), Eval("256"), 1024 / 256},
		{Eval("1024"), Eval("16"), 1024 / 16},
		{Eval("1024"), Eval("50"), 1024.0 / 50.0},
		{Eval("0.123"), Eval("2"), 0.123 / 2},
		{Eval("0.00123"), Eval("2"), 0.00123 / 2},
	}
	for _, test := range tests {
		num := test.num1.Clone()
		num.Div(test.num2)
		actual := num.ToFloat()
		if !(actual == test.expected || math.Abs(actual-test.expected) <= test.expected*testPrecision) {
			t.Fatalf("BigNum [%v] Div(%v) need %v, but %v", test.num1.String(), test.num2.String(), test.expected, actual)
		}
	}
}

func TestBigNumMod(t *testing.T) {
	tests := []struct {
		num1     *BigNum
		num2     *BigNum
		expected int
	}{
		{Eval("1"), Eval("2"), 1},
		{Eval("2"), Eval("2"), 0},
		{Eval("3"), Eval("2"), 1},
		{Eval("4"), Eval("2"), 0},
		{Eval("9999"), Eval("2"), 9999 % 2},
		{Eval("9999"), Eval("3"), 9999 % 3},
		{Eval("9999"), Eval("5"), 9999 % 5},
		{Eval("9999"), Eval("7"), 9999 % 7},
		{Eval("9999"), Eval("11"), 9999 % 11},
		{Eval("9999"), Eval("27"), 9999 % 27},
		{Eval("10000"), Eval("27"), 10000 % 27},
	}
	for _, test := range tests {
		num := test.num1.Clone()
		num.Mod(test.num2)
		actual := num.ToInt()
		if actual != test.expected {
			t.Fatalf("BigNum [%v] Div(%v) need %v, but %v", test.num1.String(), test.num2.String(), test.expected, actual)
		}
	}
}

func TestBigNumMul(t *testing.T) {
	tests := []struct {
		num1     *BigNum
		num2     *BigNum
		expected float64
	}{
		{Eval("0"), Eval("0"), 0},
		{Eval("0"), Eval("1"), 0},
		{Eval("1"), Eval("0"), 0},
		{Eval("1"), Eval("1"), 1},
		{Eval("10"), Eval("10"), 100},
		{Eval("100"), Eval("100"), 10000},
		{Eval("123"), Eval("321"), 39483},
		{Eval("1.1"), Eval("1.1"), 1.21},
		{Eval("1.12"), Eval("2.11"), 2.3632},
		{Eval("0.12"), Eval("0.11"), 0.0132},
	}
	for _, test := range tests {
		num := test.num1.Clone()
		num.Mul(test.num2)
		actual := num.ToFloat()
		if actual != test.expected {
			t.Fatalf("BigNum [%v] Sqrt(%v) need %v, but %v", test.num1.String(), test.num2.String(), test.expected, actual)
		}
	}

}

func TestBigNumSqrt(t *testing.T) {
	tests := []struct {
		num      *BigNum
		expected float64
	}{
		{Eval("0"), 0},
		{Eval("1"), 1},
		{Eval("2"), 1.4142},
		{Eval("200"), 14.1421},
		{Eval("20"), 4.4721},
		{Eval("4"), 2},
		{Eval("5.5"), 2.3452},
		{Eval("2.12"), 1.456},
		{Eval("102"), 10.0995},
		{Eval("1.02"), 1.0099},
		{Eval("1.005"), 1.0024},
	}
	for _, test := range tests {
		num := test.num.Clone()
		num.Sqrt()
		actual := num.ToFloat()
		if !(actual == test.expected || math.Abs(actual-test.expected) <= test.expected*testPrecision) {
			t.Fatalf("BigNum [%v] Sqrt() need %v, but %v", test.num.String(), test.expected, actual)
		}
	}
}

func TestBigNumPower(t *testing.T) {
	tests := []struct {
		num1     *BigNum
		num2     *BigNum
		expected float64
	}{
		{Eval("0"), Eval("0"), 1.0},
		{Eval("1"), Eval("0"), 1.0},
		{Eval("2"), Eval("2"), 4.0},
		{Eval("3"), Eval("3"), 27.0},
		{Eval("7"), Eval("5"), 16807.0},
		{Eval("1.2"), Eval("3"), 1.728},
		{Eval("2.123"), Eval("3"), 9.568634867},
		{Eval("2.12"), Eval("0.1"), 1.078},
		{Eval("96"), Eval("0.123"), 1.75},
	}
	for _, test := range tests {
		num := test.num1.Clone()
		num.Power(test.num2)
		actual := num.ToFloat()
		if !(actual == test.expected || math.Abs(actual-test.expected) <= test.expected*testPrecision) {
			t.Fatalf("BigNum [%v] Power(%v) need %v, but %v", test.num1.String(), test.num2.String(), test.expected, actual)
		}
	}
}
