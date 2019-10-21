package bignum

import (
	"math"
	"strconv"
	"strings"
)

func isNegative(b *BigNum) bool {
	if len(*b) == 0 {
		return false
	}
	return ((*b)[0] & 0x80) != 0
}

// int 总共占几个byte
// int 表示length的值占了几个byte
func decimalLength(b *BigNum) (int, int) {
	length := []byte{}
	if len(*b) == 0 {
		return 1, 1
	}
	length = append(length, (*b)[0]&0x7F)
	count := 1
	for i, by := range *b {
		if i == 0 {
			continue
		}
		if by == byte(0) {
			break
		}
		count++
		length = append(length, by)
	}
	big := BigInt(length)
	return big.ToInt(), count
}

// bool 是否是负数
// bigint 小数部分
// bigint 整数部分
func bigNumSplit(b *BigNum) (bool, *BigInt, *BigInt) {
	neg := isNegative(b)
	declength, headLen := decimalLength(b)
	start := headLen + 1
	end := start + declength
	var dec BigInt
	if declength == 0 {
		dec = nil
	} else {
		dec = BigInt((*b)[start:end])
	}
	inter := BigInt((*b)[end:])
	return neg, &dec, &inter
}

func bigNumJoin(neg bool, dec, inter *BigInt) *BigNum {
	length := intToBigInt(int64(len(*dec)))
	if neg {
		length[0] |= 0x80
	}
	bs := append(length, byte(0))
	bs = append(bs, (*dec)...)
	bs = append(bs, (*inter)...)
	ret := BigNum(bs)
	return &ret
}

func IntToBigNum(n int64) *BigNum {
	inter := intToBigInt(n)
	return bigNumJoin(false, &BigInt{0}, &inter)
}

func FloatToBigNum(n float64) *BigNum {
	return Eval(strconv.FormatFloat(n, 'f', -1, 64))
}

type BigNum []byte

func Eval(expr string) *BigNum {
	var dec, inter *BigInt
	neg := expr[0] == '-'
	if expr[0] == '-' || expr[0] == '+' {
		expr = expr[1:]
	}
	if strings.Index(expr, ".") != -1 {
		nums := strings.Split(expr, ".")
		inter = EvalInt(nums[0])
		dec = EvalDec(nums[1])
	} else {
		inter = EvalInt(expr)
		dec = &BigInt{}
	}
	return bigNumJoin(neg, dec, inter)
}

func (b *BigNum) Clone() *BigNum {
	nb := BigNum(append([]byte{}, (*b)...))
	return &nb
}

func (b *BigNum) IsZero() bool {
	_, dec, inter := bigNumSplit(b)
	return inter.IsZero() && dec.IsZero()
}

func (b *BigNum) Less(a *BigNum) bool {
	bneg, bdec, binter := bigNumSplit(b)
	aneg, adec, ainter := bigNumSplit(a)
	if bneg != aneg {
		if !bneg && aneg {
			return false
		}
		if bneg && !aneg {
			return true
		}
	}
	if !bneg {
		return binter.Less(ainter) || bdec.DecLess(adec)
	}
	return binter.Great(ainter) || bdec.DecGreat(adec)
}

func (b *BigNum) ULess(a *BigNum) bool {
	_, bdec, binter := bigNumSplit(b)
	_, adec, ainter := bigNumSplit(a)
	return binter.Less(ainter) || bdec.DecLess(adec)
}

func (b *BigNum) Great(a *BigNum) bool {
	bneg, bdec, binter := bigNumSplit(b)
	aneg, adec, ainter := bigNumSplit(a)
	if bneg != aneg {
		if !bneg && aneg {
			return true
		}
		if bneg && !aneg {
			return false
		}
	}
	if bneg {
		return binter.Less(ainter) || bdec.DecLess(adec)
	}
	return binter.Great(ainter) || bdec.DecGreat(adec)
}

func (b *BigNum) UGreat(a *BigNum) bool {
	_, bdec, binter := bigNumSplit(b)
	_, adec, ainter := bigNumSplit(a)
	return binter.Great(ainter) || bdec.DecGreat(adec)
}

func (b *BigNum) Equle(a *BigNum) bool {
	bneg, bdec, binter := bigNumSplit(b)
	aneg, adec, ainter := bigNumSplit(a)
	if bneg != aneg {
		return false
	}
	return bdec.DecEqule(adec) && binter.Equle(ainter)
}

func (b *BigNum) UEqule(a *BigNum) bool {
	_, bdec, binter := bigNumSplit(b)
	_, adec, ainter := bigNumSplit(a)
	return bdec.DecEqule(adec) && binter.Equle(ainter)
}

func (b *BigNum) Neg() {
	neg, dec, inter := bigNumSplit(b)
	*b = *bigNumJoin(!neg, dec, inter)
}

func (b *BigNum) Succ() {
	neg, dec, inter := bigNumSplit(b)
	inter.Succ()
	*b = *bigNumJoin(neg, dec, inter)
}

func (b *BigNum) Pred() {
	neg, dec, inter := bigNumSplit(b)
	inter.Pred()
	*b = *bigNumJoin(neg, dec, inter)
}

func (b *BigNum) Plus(a *BigNum) {
	if b.IsZero() {
		*b = *a.Clone()
		return
	}
	if a.IsZero() {
		return
	}
	bneg, bdec, binter := bigNumSplit(b)
	aneg, adec, ainter := bigNumSplit(a)
	neg := false
	retdec, retint := &BigInt{}, &BigInt{}
	adec, bdec = decimalAilgn(adec, bdec)
	declength := adec.Len()

	taint := append(*adec, (*ainter)...)
	tbint := append(*bdec, (*binter)...)

	if bneg == aneg {
		neg = aneg
		taint.Plus(&tbint)
		*retdec = taint[:declength]
		*retint = taint[declength:]
	} else {
		if b.UGreat(a) {
			neg = bneg

			tbint.Minus(&taint)
			*retdec = tbint[:declength]
			*retint = tbint[declength:]
		} else {
			neg = aneg

			taint.Minus(&tbint)
			*retdec = taint[:declength]
			*retint = taint[declength:]
		}
	}

	*b = *bigNumJoin(neg, retdec, retint)
}

func (b *BigNum) Minus(a *BigNum) {
	aneg, adec, ainter := bigNumSplit(a)
	b.Plus(bigNumJoin(!aneg, adec, ainter))
}

func (b *BigNum) Mul(a *BigNum) {
	bneg, bdec, binter := bigNumSplit(b)
	blen, _ := decimalLength(b)
	aneg, adec, ainter := bigNumSplit(a)
	alen, _ := decimalLength(a)
	neg := bneg != aneg

	bnum := append(*bdec, (*binter)...)
	anum := append(*adec, (*ainter)...)
	bnum.Mul(&anum)

	length := blen + alen
	dec := bnum[:length]
	inter := bnum[length:]
	*b = *bigNumJoin(neg, &dec, &inter)
}

func (b *BigNum) DivMod(a *BigNum) (bool, *BigInt, int, *BigInt) {
	bneg, bdec, binter := bigNumSplit(b)
	blen, _ := decimalLength(b)
	aneg, adec, ainter := bigNumSplit(a)
	alen, _ := decimalLength(a)
	neg := bneg != aneg
	if blen < alen {
		times := alen - blen
		for i := 0; i < times; i++ {
			*bdec = append([]byte{0}, (*bdec)...)
		}
	}
	if blen > alen {
		times := blen - alen
		for i := 0; i < times; i++ {
			*adec = append([]byte{0}, (*adec)...)
		}
	}
	bnum := append(*bdec, (*binter)...)
	anum := append(*adec, (*ainter)...)

	length := 0
	maxlength := MAXLENGTH
	if bdec.Len() > maxlength {
		maxlength = bdec.Len()
	}
	var num, div, mod *BigInt
	for {
		num = bnum.Clone()
		for i := 0; i < length; i++ {
			*num = append([]byte{0}, (*num)...)
		}
		div, mod = num.DivMod(&anum)
		if mod.IsZero() {
			break
		}
		length++
		if length > maxlength {
			length--
			break
		}
	}
	return neg, div, length, mod
}

func (b *BigNum) Div(a *BigNum) {
	neg, div, length, _ := b.DivMod(a)
	dec := (*div)[:length]
	inter := BigInt{}
	if length < len(*div) {
		inter = (*div)[length:]
	}
	*b = *bigNumJoin(neg, &dec, &inter)
}

func (b *BigNum) Mod(a *BigNum) {
	_, _, binter := bigNumSplit(b)
	_, _, ainter := bigNumSplit(a)
	binter.Mod(ainter)
	*b = *bigNumJoin(false, &BigInt{0}, binter)
}

func (b *BigNum) Square() {
	b.Mul(b)
}

func (b *BigNum) Cube() {
	b.Mul(b)
	b.Mul(b)
}

func (b *BigNum) Sqrt() {
	if b.IsZero() {
		return
	}
	_, dec, inter := bigNumSplit(b)
	preLen := MAXLENGTH - dec.Len()
	for i := 0; i < preLen; i++ {
		*dec = append([]byte{0}, (*dec)...)
	}
	tmp := &BigInt{}
	*tmp = append(*tmp, (*dec)...)
	*tmp = append(*tmp, (*inter)...)
	tmp.Sqrt()

	iilen := inter.Len()
	tmplen := len(*tmp)
	ilen := int(math.Round(float64(iilen) / 2.0))

	ninter := BigInt((*tmp)[tmplen-ilen:])

	if inter.Len10() < ninter.Len10() {
		tmp.LeftShift()
		ninter = BigInt((*tmp)[tmplen-ilen:])
	}

	ndec := BigInt((*tmp)[:tmplen-ilen])
	*b = *bigNumJoin(false, &ndec, &ninter)
}

func (b *BigNum) Power(a *BigNum) {
	if a.IsZero() {
		*b = *IntToBigNum(1)
		return
	}
	if b.IsZero() {
		return
	}
	_, dec, inter := bigNumSplit(a)
	tmp := b.PowerDecimal(dec)
	tmp.Mul(b.PowerInt(inter.ToInt()))
	*b = *tmp
}

func (b *BigNum) PowerDecimal(a *BigInt) *BigNum {
	if a.IsZero() {
		return IntToBigNum(1)
	}
	if b.IsZero() {
		return IntToBigNum(0)
	}
	ret := IntToBigNum(1)
	bins := a.ToDecBin()
	for i, bo := range bins {
		if bo {
			tmp := b.Clone()
			for j := 0; j < i+1; j++ {
				tmp.Sqrt()
			}
			ret.Mul(tmp)
		}
	}
	return ret
}

func (b *BigNum) PowerInt(n int) *BigNum {
	if n == 0 {
		return IntToBigNum(1)
	}
	if b.IsZero() {
		return IntToBigNum(0)
	}
	ret := b.Clone()
	for i := 1; i < n; i++ {
		ret.Mul(b)
	}
	return ret
}

func (b *BigNum) Factorial() {
	_, _, inter := bigNumSplit(b)
	inter.Factorial()
	*b = *bigNumJoin(false, &BigInt{}, inter)
}

func (b *BigNum) ToFloat() float64 {
	numString := b.String()
	f, _ := strconv.ParseFloat(numString, 64)
	return f
}

func (b *BigNum) ToInt() int {
	return int(b.ToFloat())
}

func (b *BigNum) StringPart() (string, string) {
	neg, dec, inter := bigNumSplit(b)
	negstr := ""
	if neg {
		negstr = "-"
	}
	return negstr + dec.ToDecString(), negstr + inter.String()
}

func (b *BigNum) String() string {
	neg, dec, inter := bigNumSplit(b)
	negstr := ""
	if neg {
		negstr = "-"
	}
	tail := ""
	if !dec.IsZero() {
		tail = "." + dec.ToDecString()
	}
	return negstr + inter.String() + tail
}
