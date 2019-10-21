package bignum

import (
	"math"
	"strconv"
)

// int value
// int value*value
func squareTest(a int) (int, int) {
	lastV := 1
	for i, v := range squareTable {
		if v > a {
			return i, lastV
		}
		if v == a {
			return i + 1, v
		}
		lastV = v
	}
	return -1, -1
}

func decimalAilgn(a, b *BigInt) (*BigInt, *BigInt) {
	reta, retb := a.Clone(), b.Clone()
	alen := len(*a)
	blen := len(*b)
	t := alen - blen
	num := retb
	if alen < blen {
		t = blen - alen
		num = reta
	}
	for i := 0; i < t; i++ {
		num.RightShift100()
	}
	return reta, retb
}

func byteToNum(b byte) int {
	return int(b) % 100
}

func EvalInt(expr string) *BigInt {
	ret := BigInt{}
	length := len(expr)
	for i := 0; i < length; i += 2 {
		start := length - i - 2
		if start < 0 {
			start = 0
		}
		part := expr[start : length-i]
		num, _ := strconv.Atoi(part)
		ret = append(ret, byte(num))
	}
	return &ret
}

func EvalDec(expr string) *BigInt {
	ret := BigInt{}
	length := len(expr)
	for i := 0; i < length; i += 2 {
		end := i + 2
		if end > length {
			end = length
		}
		part := expr[i:end]
		num, _ := strconv.Atoi(part)
		if len(part) != 2 {
			num *= 10
		}
		ret = append([]byte{byte(num)}, ret...)
	}
	return &ret
}

type BigInt []byte

func intToBigInt(n int64) BigInt {
	ret := []byte{}
	for {
		ret = append(ret, byte(n%100))
		n = n / 100
		if n == 0 {
			break
		}
	}
	return BigInt(ret)
}

func (b *BigInt) Len() int {
	length := len(*b)
	ret := length
	if length == 1 {
		return ret
	}
	for i := 0; i < length; i++ {
		if (*b)[length-i-1] == 0 {
			ret--
		} else {
			break
		}
	}
	if length != 0 && ret == 0 {
		return 1
	}
	return ret
}

// 有多少十进制位
func (b *BigInt) Len10() int {
	length := b.Len()
	ret := length * 2
	if (*b)[length-1] < 10 {
		ret--
	}
	return ret
}

func (b *BigInt) IsZero() bool {
	if b == nil || *b == nil {
		return true
	}
	if len(*b) == 1 && (*b)[0] == 0 {
		return true
	}
	for _, v := range *b {
		if v != 0 {
			return false
		}
	}
	return true
}

func (b *BigInt) Clone() *BigInt {
	nb := BigInt(append([]byte{}, (*b)...))
	return &nb
}

func (b *BigInt) DecLess(a *BigInt) bool {
	na, nb := decimalAilgn(a, b)
	return nb.Less(na)
}

func (b *BigInt) Less(a *BigInt) bool {
	if a.IsZero() {
		return false
	} else if b.IsZero() {
		return true
	}
	blen := b.Len()
	alen := a.Len()
	if blen > alen {
		return false
	}
	if alen > blen {
		return true
	}
	for i := 1; i <= blen; i++ {
		if (*b)[blen-i] < (*a)[blen-i] {
			return true
		} else if (*b)[blen-i] > (*a)[blen-i] {
			return false
		}
	}
	return false
}

func (b *BigInt) DecGreat(a *BigInt) bool {
	na, nb := decimalAilgn(a, b)
	return nb.Great(na)
}

func (b *BigInt) Great(a *BigInt) bool {
	if b.IsZero() {
		return false
	} else if a.IsZero() {
		return true
	}
	blen := b.Len()
	alen := a.Len()
	if blen > alen {
		return true
	}
	if alen > blen {
		return false
	}
	for i := 1; i <= blen; i++ {
		if (*b)[blen-i] > (*a)[blen-i] {
			return true
		} else if (*b)[blen-i] < (*a)[blen-i] {
			break
		}
	}
	return false
}

func (b *BigInt) DecEqule(a *BigInt) bool {
	na, nb := decimalAilgn(a, b)
	return nb.Equle(na)
}

func (b *BigInt) Equle(a *BigInt) bool {
	blen := b.Len()
	alen := a.Len()
	if blen != alen {
		return false
	}
	for i := 0; i < blen; i++ {
		by := (*b)[i]
		if by != (*a)[i] {
			return false
		}
	}
	return true
}

// / 10
func (b *BigInt) LeftShift() {
	length := b.Len()
	flag := 0
	for i, v := range *b {
		if i+1 < length {
			flag = int((*b)[i+1]) % 10
		} else {
			flag = 0
		}
		nv := int(v)/10 + flag*10
		(*b)[i] = byte(nv)
	}
}

// / 100
func (b *BigInt) LeftShift100() {
	if b.IsZero() {
		return
	}
	*b = (*b)[1:]
}

// x 10
func (b *BigInt) RightShift() {
	flag := 0
	for i, v := range *b {
		num := int(v)*10 + flag
		nv := num % 100
		flag = num / 100
		(*b)[i] = byte(nv)
	}
	if flag != 0 {
		*b = append(*b, byte(flag))
	}
}

// x 100
func (b *BigInt) RightShift100() {
	if b.IsZero() {
		*b = append(*b, byte(0))
		return
	}
	*b = append([]byte{0}, (*b)...)
}

func (b *BigInt) Succ() {
	length := len(*b)
	if length == 0 {
		(*b) = []byte{1}
		return
	}
	i := 0
	for {
		if length == i {
			*b = append(*b, byte(1))
			return
		}
		if (*b)[i] == 99 {
			(*b)[i] = 0
		} else {
			(*b)[i]++
			return
		}
		i++
	}
}

func (b *BigInt) Pred() {
	length := len(*b)
	if length == 0 {
		(*b) = []byte{0}
		return
	}
	if b.IsZero() {
		return
	}
	i := 0
	for {
		if length == i {
			return
		}
		if (*b)[i] == 0 {
			if length == i+1 {
				return
			}
			(*b)[i] = 99
		} else {
			(*b)[i]--
			return
		}
		i++
	}
}

func (b *BigInt) DecPlus(a *BigInt) {
	na, nb := decimalAilgn(a, b)
	nb.Plus(na)
	*b = *nb
}

func (b *BigInt) Plus(a *BigInt) {
	if len(*a) == 1 && (*a)[0] == 1 {
		b.Succ()
		return
	}
	if a.IsZero() {
		return
	}
	if b.IsZero() {
		*b = *a.Clone()
		return
	}
	minlen := 0
	blen := b.Len()
	alen := a.Len()
	if blen > alen {
		minlen = alen
	} else {
		minlen = blen
	}
	flag := 0
	for i := 0; i < minlen; i++ {
		v := byteToNum((*b)[i]) + byteToNum((*a)[i]) + flag
		if v >= 100 {
			flag = 1
		} else {
			flag = 0
		}
		(*b)[i] = byte(v % 100)
	}
	if blen < alen {
		tail := (*a)[blen:]
		if flag == 1 {
			i := 0
			length := len(tail)
			for {
				if length == i {
					tail = append(tail, byte(1))
					break
				}
				tail[i]++
				if tail[i] != 100 {
					break
				} else {
					tail[i] = 0
				}
				i++
			}
		}
		*b = append(*b, tail...)
	} else {
		if flag == 1 {
			if len(*b) == alen {
				*b = append(*b, byte(1))
				return
			}
			i := alen
			length := len(*b)
			for {
				if length == i {
					*b = append(*b, byte(1))
					return
				}
				(*b)[i]++
				if (*b)[i] != 100 {
					break
				} else {
					(*b)[i] = 0
				}
				i++
			}
		}
	}
}

func (b *BigInt) DecMinus(a *BigInt) {
	na, nb := decimalAilgn(a, b)
	nb.Minus(na)
	*b = *nb
}

func (b *BigInt) Minus(a *BigInt) {
	if len(*a) == 1 && (*a)[0] == 1 {
		b.Pred()
		return
	}
	if a.IsZero() {
		return
	}
	if b.IsZero() {
		*b = *a.Clone()
		return
	}
	if b.Great(a) {
		minlen := 0
		blen := b.Len()
		alen := a.Len()
		if blen > alen {
			minlen = alen
		} else {
			minlen = blen
		}
		flag := 0
		i := 0
		for ; i < minlen; i++ {
			bn := byteToNum((*b)[i]) - flag
			an := byteToNum((*a)[i])
			v := 0
			if bn < an {
				v = 100 + bn - flag - an
				flag = 1
			} else {
				v = bn - an
				flag = 0
			}
			(*b)[i] = byte(v)
		}
		if flag == 1 {
			for {
				if (*b)[i] == 0 {
					(*b)[i] = 99
				} else {
					(*b)[i]--
					break
				}
				i++
			}
		}
		if (*b)[blen-1] == byte(0) {
			(*b) = (*b)[:blen-1]
		}
		return
	}
	if b.Less(a) {
		nb := a.Clone()
		nb.Minus(b)
		*b = *nb
		return
	}
	*b = []byte{0}
}

func (b *BigInt) MulByte(a byte) *BigInt {
	if b.IsZero() || a == 0 {
		ret := BigInt{0}
		return &ret
	}
	ret := &BigInt{0}
	an := byteToNum(a)
	for i, by := range *b {
		bn := byteToNum(by)
		vn := intToBigInt(int64(bn * an))
		for j := 0; j < i; j++ {
			vn = append([]byte{0}, vn...)
		}
		ret.Plus(&vn)
	}
	return ret
}

func (b *BigInt) Mul(a *BigInt) {
	if b.IsZero() || a.IsZero() {
		*b = BigInt{0}
		return
	}
	tmp := &BigInt{0}
	for i, ay := range *a {
		vn := b.MulByte(ay)
		for j := 0; j < i; j++ {
			*vn = append([]byte{0}, (*vn)...)
		}
		tmp.Plus(vn)
	}
	*b = *tmp
}

func (b *BigInt) DivMod(a *BigInt) (*BigInt, *BigInt) {
	if b.IsZero() {
		return b, &BigInt{0}
	}
	if a.IsZero() {
		// [TODO] error catch
		return b, &BigInt{0}
	}
	if b.Less(a) {
		return &BigInt{0}, b
	}
	if b.Equle(a) {
		return &BigInt{1}, &BigInt{0}
	}

	base := BigInt{1}
	baseN := a.Clone()
	for {
		*baseN = append([]byte{0}, (*baseN)...)
		if baseN.Great(b) {
			*baseN = (*baseN)[1:]
			break
		}
		base = append([]byte{0}, base...)
	}

	tmp := b.Clone()
	count := BigInt{0}
	for {
		if !tmp.Less(baseN) {
			tmp.Minus(baseN)
			count.Plus(&base)
		} else if !baseN.Equle(a) {
			*baseN = (*baseN)[1:]
			base = base[1:]
		}
		if tmp.Less(a) {
			break
		}
	}

	return &count, tmp
}

func (b *BigInt) Div(a *BigInt) {
	div, _ := b.DivMod(a)
	*b = *div
}

func (b *BigInt) Mod(a *BigInt) {
	_, mod := b.DivMod(a)
	*b = *mod
}

func (b *BigInt) Sqrt() *BigInt {
	nb, d := b.Sqrtx(&BigInt{0})
	*b = *nb
	return d
}

// return 余数
func (b *BigInt) Sqrtx(initd *BigInt) (*BigInt, *BigInt) {
	length := b.Len()

	// 查表
	if length <= 2 {
		v, s := squareTest(b.ToInt())
		vd := b.Clone()
		bs := intToBigInt(int64(s))
		vd.Minus(&bs)
		vret := intToBigInt(int64(v))
		return &vret, vd
	}

	r := &BigInt{0}
	d := initd
	tmp := []byte{}
	for i := 0; i < length; i++ {
		by := (*b)[length-i-1]

		r20 := r.MulByte(20)
		n := d.Clone()
		n.Mul(&BigInt{0, 1})
		n.Plus(&BigInt{by})
		x := 1

		tmpNum := &BigInt{0}
		for {
			tmp := r20.Clone()
			tmp.Plus(&BigInt{byte(x)})
			tmp = tmp.MulByte(byte(x))
			if tmp.Great(n) {
				x--
				break
			}
			tmpNum = tmp
			x++
		}
		r = r.MulByte(10)
		r.Plus(&BigInt{byte(x)})
		d = n.Clone()
		d.Minus(tmpNum)

		tmp = append([]byte{byte(x)}, tmp...)
	}
	ret := &BigInt{}
	bytmp := byte(0)
	onset := false
	for _, b := range tmp {
		if onset {
			*ret = append(*ret, bytmp+b*10)
			onset = false
		} else {
			bytmp = b
			onset = true
		}
	}
	if onset {
		*ret = append(*ret, bytmp)
	}
	return ret, d
}

func (b *BigInt) Factorial() {
	if b.IsZero() {
		*b = intToBigInt(1)
		return
	}
	temp := intToBigInt(1)
	count := intToBigInt(2)
	for {
		temp.Mul(&count)
		if b.Equle(&count) {
			break
		}
		count.Succ()
	}
	*b = temp
}

func (b *BigInt) String() string {
	ret := ""
	prefix := ""
	for _, by := range *b {
		num := byteToNum(by)
		ret = strconv.Itoa(num) + prefix + ret
		if num < 10 {
			prefix = "0"
		} else {
			prefix = ""
		}
	}
	return ret
}

func (b *BigInt) ToDecString() string {
	ret := ""
	for _, by := range *b {
		num := byteToNum(by)
		if num == 0 {
			ret = "00" + ret
		} else if num < 10 {
			ret = "0" + strconv.Itoa(num) + ret
		} else {
			ret = strconv.Itoa(num) + ret
		}
	}
	return ret
}

func (b *BigInt) ToInt() int {
	ret := 0
	for i, by := range *b {
		ret += int(by) % 100 * int(math.Pow(100, float64(i)))
	}
	return ret
}

func (b *BigInt) ToBin() []bool {
	tmp := b.Clone()
	ret := []bool{}
	num2 := &BigInt{2}
	for {
		div, mod := tmp.DivMod(num2)
		if mod.Equle(&BigInt{1}) {
			ret = append([]bool{true}, ret...)
		} else {
			ret = append([]bool{false}, ret...)
		}
		tmp = div
		if tmp.Less(num2) {
			if (*tmp)[0] == 1 {
				ret = append([]bool{true}, ret...)
			}
			break
		}
	}
	return ret
}

func (b *BigInt) ToDecBin() []bool {
	length := b.Len() * 2
	if length < MAXLENGTH {
		length = MAXLENGTH
	}
	tmp := b.Clone()
	ret := []bool{}
	num2 := &BigInt{2}
	for {
		plen := tmp.Len()
		tmp.Mul(num2)
		if tmp.Len() > plen {
			ret = append(ret, true)
			*tmp = (*tmp)[:tmp.Len()-1]
		} else {
			ret = append(ret, false)
		}
		if tmp.IsZero() {
			break
		}
		if len(ret) >= length {
			break
		}
	}
	return ret
}
