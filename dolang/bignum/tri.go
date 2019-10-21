package bignum

import "math"

func Cos(b *BigNum) *BigNum {
	return FloatToBigNum(cosf(b.ToFloat()))
}

func cosf(n float64) float64 {
	ni := int(math.Round(n))
	if ni > 90 || ni < 0 {
		return -1
	}
	return cosTable[ni-1]
}

func Sin(b *BigNum) *BigNum {
	return FloatToBigNum(sinf(b.ToFloat()))
}

func sinf(n float64) float64 {
	ni := int(math.Round(n))
	if ni > 90 || ni < 0 {
		return -1
	}
	return sinTable[ni-1]
}

func Tan(b *BigNum) *BigNum {
	return FloatToBigNum(tanf(b.ToFloat()))
}

func tanf(n float64) float64 {
	ni := int(math.Round(n))
	if ni > 90 || ni < 0 {
		return -1
	}
	return tanTable[ni-1]
}
